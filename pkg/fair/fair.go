package fair

import (
	"database/sql"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/je4/FairService/v2/pkg/model/myfair"
	"github.com/lib/pq"
	"github.com/op/go-logging"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

type Source struct {
	ID          int64
	Name        string
	DetailURL   string
	Description string
	OAIDomain   string
	Partition   string
}

type DataAccess string

const (
	DataAccessPublic     DataAccess = "public"
	DataAccessClosed     DataAccess = "closed"
	DataAccessClosedData DataAccess = "closed_data"
)

var DataAccessReverse = map[string]DataAccess{
	string(DataAccessPublic):     DataAccessPublic,
	string(DataAccessClosed):     DataAccessClosed,
	string(DataAccessClosedData): DataAccessClosedData,
}

type ItemData struct {
	Source    int64       `json:"source"`
	Signature string      `json:"signature"`
	Metadata  myfair.Core `json:"metadata"`
	Set       []string    `json:"set"`
	Catalog   []string    `json:"catalog"`
	Access    DataAccess  `json:"access"`
	Deleted   bool        `json:"deleted"`
	Seq       int64
	UUID      string
	Datestamp time.Time
}

type SourceData struct {
	Source string `json:"source"`
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

type Fair struct {
	dbschema     string
	db           *sql.DB
	sourcesMutex sync.RWMutex
	sources      map[int64]*Source
	partitions   map[string]*Partition
	log          *logging.Logger
}

func NewFair(db *sql.DB, dbschema string, log *logging.Logger) (*Fair, error) {
	f := &Fair{
		dbschema:     dbschema,
		db:           db,
		sourcesMutex: sync.RWMutex{},
		sources:      map[int64]*Source{},
		partitions:   map[string]*Partition{},
		log:          log,
	}
	if err := f.LoadSources(); err != nil {
		return nil, errors.Wrap(err, "cannot load sources")
	}
	return f, nil
}

func (f *Fair) AddPartition(p *Partition) {
	f.partitions[p.Name] = p
}

func (f *Fair) GetPartition(name string) (*Partition, error) {
	p, ok := f.partitions[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("cannot find partition %s", name))
	}
	return p, nil
}

func (f *Fair) LoadSources() error {
	sqlstr := fmt.Sprintf("SELECT sourceid, name, detailurl, description, oai_domain, partition FROM %s.source", f.dbschema)
	rows, err := f.db.Query(sqlstr)
	if err != nil {
		return errors.Wrapf(err, "cannot execute %s", sqlstr)
	}
	defer rows.Close()
	f.sourcesMutex.Lock()
	defer f.sourcesMutex.Unlock()
	f.sources = make(map[int64]*Source)
	for rows.Next() {
		src := &Source{}
		if err := rows.Scan(&src.ID, &src.Name, &src.DetailURL, &src.Description, &src.OAIDomain, &src.Partition); err != nil {
			return errors.Wrap(err, "cannot scan values")
		}
		f.sources[src.ID] = src
	}
	return nil
}

func (f *Fair) GetSourceById(id int64, partitionName string) (*Source, error) {
	f.sourcesMutex.RLock()
	defer f.sourcesMutex.RUnlock()
	if s, ok := f.sources[id]; ok {
		if s.Partition == partitionName {
			return s, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("source #%v for partition %s not found", id, partitionName))
}

func (f *Fair) GetSourceByName(name string, partitionName string) (*Source, error) {
	f.sourcesMutex.RLock()
	defer f.sourcesMutex.RUnlock()
	for _, src := range f.sources {
		if src.Name == name && src.Partition == partitionName {
			return src, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("source %s for partition %s not found", name, partitionName))
}

func (f *Fair) GetSourceByOAIDomain(name string) (*Source, error) {
	f.sourcesMutex.RLock()
	defer f.sourcesMutex.RUnlock()
	for _, src := range f.sources {
		if src.OAIDomain == name {
			return src, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("source %s not found", name))
}

func (f *Fair) GetMinimumDatestamp(partitionName string) (time.Time, error) {
	partition, ok := f.partitions[partitionName]
	if !ok {
		return time.Time{}, errors.New(fmt.Sprintf("partition %s not found", partitionName))
	}
	sqlstr := fmt.Sprintf("SELECT MIN(datestamp) AS mindate"+
		" FROM %s.coreview"+
		" WHERE partition=$1", f.dbschema)
	params := []interface{}{partition.Name}
	var datestamp time.Time
	if err := f.db.QueryRow(sqlstr, params...).Scan(&datestamp); err != nil {
		if err != sql.ErrNoRows {
			return time.Time{}, errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		return time.Now(), nil
	}
	return datestamp, nil
}

func (f *Fair) getItems(sqlWhere string, params []interface{}, limit, offset int64, completeListSize *int64, fn func(item *ItemData) error) error {
	if completeListSize != nil {
		sqlstr := fmt.Sprintf("SELECT COUNT(*) AS num"+
			" FROM %s.coreview"+
			" WHERE %s", f.dbschema, sqlWhere)
		if err := f.db.QueryRow(sqlstr, params...).Scan(completeListSize); err != nil {
			return errors.Wrapf(err, "cannot get number of result items [%s] - [%v]", sqlstr, params)
		}
	}
	sqlstr := fmt.Sprintf("SELECT uuid, metadata, setspec, catalog, access, signature, source, deleted, seq, datestamp"+
		" FROM %s.coreview"+
		" WHERE %s"+
		" ORDER BY seq ASC", f.dbschema, sqlWhere)
	if limit > 0 {
		sqlstr += fmt.Sprintf(" LIMIT %v", limit)
	}
	if offset > 0 {
		sqlstr += fmt.Sprintf(" OFFSET %v", offset)
	}
	rows, err := f.db.Query(sqlstr, params...)
	if err != nil {
		if err != sql.ErrNoRows {
			return errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
	}
	defer rows.Close()
	for rows.Next() {
		var uuidStr string
		var metaStr string
		var set, catalog []string
		var accessStr string
		var signature string
		var source int64
		var deleted bool
		var seq int64
		var datestamp time.Time
		if err := rows.Scan(&uuidStr, &metaStr, pq.Array(&set), pq.Array(&catalog), &accessStr, &signature, &source, &deleted, &seq, &datestamp); err != nil {
			return errors.Wrapf(err, "cannot scan result of [%s] - [%v]", sqlstr, params)
		}
		data := &ItemData{
			UUID:      uuidStr,
			Source:    source,
			Signature: signature,
			Metadata:  myfair.Core{},
			Set:       set,
			Catalog:   catalog,
			Deleted:   deleted,
			Seq:       seq,
			Datestamp: datestamp,
		}
		var ok bool
		data.Access, ok = DataAccessReverse[accessStr]
		if !ok {
			return errors.New(fmt.Sprintf("[%s] invalid access type %s", uuidStr, accessStr))
		}
		if err := json.Unmarshal([]byte(metaStr), &data.Metadata); err != nil {
			return errors.Wrapf(err, "[%s] cannot unmarshal core [%s]", uuidStr, metaStr)
		}
		if err := fn(data); err != nil {
			return errors.Wrap(err, "error calling fn")
		}
	}
	return nil

}

func (f *Fair) GetItemsDatestamp(partitionName string, datestamp, until time.Time, access []DataAccess, limit, offset int64, completeListSize *int64, fn func(item *ItemData) error) error {
	partition, ok := f.partitions[partitionName]
	if !ok {
		return errors.New(fmt.Sprintf("partition %s not found", partitionName))
	}

	sqlWhere := "partition=$1 AND datestamp>=$2"
	params := []interface{}{partition.Name, datestamp}
	if len(access) > 0 {
		var accessList []string
		for key, acc := range access {
			accessList = append(accessList, fmt.Sprintf("access=$%v", key+3))
			params = append(params, acc)
		}
		sqlWhere += fmt.Sprintf(" AND (%s)", strings.Join(accessList, " OR "))
	}
	if !until.Equal(time.Time{}) {
		params = append(params, until)
		sqlWhere += fmt.Sprintf(" AND datestamp<=$%v", len(params))
	}
	return f.getItems(sqlWhere, params, limit, offset, completeListSize, fn)
}

func (f *Fair) GetItemsSeq(partitionName string, seq int64, until time.Time, access []DataAccess, limit, offset int64, completeListSize *int64, fn func(item *ItemData) error) error {
	partition, ok := f.partitions[partitionName]
	if !ok {
		return errors.New(fmt.Sprintf("partition %s not found", partitionName))
	}

	sqlWhere := "partition=$1 AND seq>=$2"
	params := []interface{}{partition.Name, seq}
	if len(access) > 0 {
		var accessList []string
		for key, acc := range access {
			accessList = append(accessList, fmt.Sprintf("access=$%v", key+3))
			params = append(params, acc)
		}
		sqlWhere += fmt.Sprintf(" AND (%s)", strings.Join(accessList, " OR "))
	}
	if !until.Equal(time.Time{}) {
		params = append(params, until)
		sqlWhere += fmt.Sprintf(" AND datestamp<=$%v", len(params))
	}
	return f.getItems(sqlWhere, params, limit, offset, completeListSize, fn)
}

func (f *Fair) GetItem(partitionName, uuidStr string) (*ItemData, error) {
	partition, ok := f.partitions[partitionName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("partition %s not found", partitionName))
	}

	sqlstr := fmt.Sprintf("SELECT metadata, setspec, catalog, access, signature, source, deleted, seq, datestamp"+
		" FROM %s.coreview"+
		" WHERE partition=$1 AND uuid=$2", f.dbschema)
	params := []interface{}{partition.Name, uuidStr}
	row := f.db.QueryRow(sqlstr, params...)

	var metaStr string
	var set, catalog []string
	var accessStr string
	var signature string
	var source int64
	var deleted bool
	var seq int64
	var datestamp time.Time
	if err := row.Scan(&metaStr, pq.Array(&set), pq.Array(&catalog), &accessStr, &signature, &source, &deleted, &seq, &datestamp); err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		return nil, nil
	}
	data := &ItemData{
		UUID:      uuidStr,
		Source:    source,
		Signature: signature,
		Metadata:  myfair.Core{},
		Set:       set,
		Catalog:   catalog,
		Deleted:   deleted,
		Seq:       seq,
		Datestamp: datestamp,
	}
	data.Access, ok = DataAccessReverse[accessStr]
	if !ok {
		return nil, errors.New(fmt.Sprintf("[%s] invalid access type %s", uuidStr, accessStr))
	}
	if err := json.Unmarshal([]byte(metaStr), &data.Metadata); err != nil {
		return nil, errors.Wrapf(err, "[%s] cannot unmarshal core [%s]", uuidStr, metaStr)
	}
	return data, nil
}

func (f *Fair) CreateItem(partitionName string, data ItemData) (string, error) {
	partition, ok := f.partitions[partitionName]
	if !ok {
		return "", errors.New(fmt.Sprintf("partition %s not found", partitionName))
	}

	sort.Strings(data.Catalog)
	sort.Strings(data.Set)

	src, err := f.GetSourceById(data.Source, partitionName)
	if err != nil {
		return "", errors.Wrapf(err, "cannot get source %s", data.Source)
	}

	sqlstr := fmt.Sprintf("SELECT uuid, metadata, setspec, catalog, access, deleted"+
		" FROM %s.coreview "+
		" WHERE source=$1 AND signature=$2", f.dbschema)
	params := []interface{}{src.ID, data.Signature}
	row := f.db.QueryRow(sqlstr, params...)

	var metaStr string
	var uuidStr string
	var set, catalog []string
	var accessStr string
	var deleted bool
	if err := row.Scan(&uuidStr, &metaStr, pq.Array(&set), pq.Array(&catalog), &accessStr, &deleted); err != nil {
		if err != sql.ErrNoRows {
			return "", errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		// do insert here

		uuidVal, err := uuid.NewUUID()
		if err != nil {
			return "", errors.Wrap(err, "cannot generate uuid")
		}
		uuidStr := uuidVal.String()
		coreBytes, err := json.Marshal(data.Metadata)
		if err != nil {
			return "", errors.Wrapf(err, "cannot marshal core data [%v]", data.Metadata)
		}
		sqlstr := fmt.Sprintf("INSERT INTO %s.core"+
			" (uuid, datestamp, setspec, metadata, signature, source, access, catalog, seq, deleted)"+
			" VALUES($1, NOW(), $2, $3, $4, $5, $6, $7, NEXTVAL('lastchange')), false", f.dbschema)
		params := []interface{}{
			uuidStr,        // uuid
			partition.Name, // partition
			// datestamp
			pq.Array(data.Set), // setspec
			string(coreBytes),  // metadata
			// dirty
			data.Signature,         // signature
			src.ID,                 // source
			data.Access,            // access
			pq.Array(data.Catalog), // catalog
			// seq
		}
		ret, err := f.db.Exec(sqlstr, params...)
		if err != nil {
			return "", errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		num, err := ret.RowsAffected()
		if err != nil {
			return "", errors.Wrap(err, "cannot get affected rows")
		}
		if num == 0 {
			return "", errors.Wrap(err, "no affected rows")
		}
		sqlstr = fmt.Sprintf("INSERT INTO %s.core_dirty"+
			" (uuid)"+
			" VALUES($1)", f.dbschema)
		params = []interface{}{
			uuidStr, // uuid
		}
		if _, err = f.db.Exec(sqlstr, params...); err != nil {
			return "", errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		f.log.Infof("new item [%s] inserted", uuidStr)
		return uuidStr, nil

	} else {
		access, ok := DataAccessReverse[accessStr]
		if !ok {
			return "", errors.Wrapf(err, "[%s] invalid access type %s", uuidStr, accessStr)
		}
		sort.Strings(catalog)
		sort.Strings(set)
		// do update here
		var oldMeta myfair.Core
		if err := json.Unmarshal([]byte(metaStr), &oldMeta); err != nil {
			return "", errors.Wrapf(err, "[%s] cannot unmarshal core [%s]", uuidStr, metaStr)
		}
		if !deleted &&
			reflect.DeepEqual(oldMeta, data.Metadata) &&
			equalStrings(set, data.Set) &&
			equalStrings(catalog, data.Catalog) &&
			access == data.Access {
			f.log.Infof("no update needed for item [%v]", uuidStr)
			sqlstr = fmt.Sprintf("INSERT INTO %s.core_dirty"+
				" (uuid)"+
				" VALUES($1)", f.dbschema)
			params = []interface{}{
				uuidStr, // uuid
			}
			if _, err = f.db.Exec(sqlstr, params...); err != nil {
				return "", errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
			}
			return uuidStr, nil
		}

		dataMetaBytes, err := json.Marshal(data.Metadata)
		if err != nil {
			return "", errors.Wrapf(err, "[%s] cannot unmarshal data core", uuidStr)
		}

		sqlstr = fmt.Sprintf("UPDATE %s.core"+
			" SET setspec=$1, metadata=$2, access=$3, catalog=$4, deleted=false"+
			" WHERE uuid=$5", f.dbschema)
		params := []interface{}{
			pq.Array(data.Set),
			string(dataMetaBytes),
			data.Access,
			pq.Array(data.Catalog),
			uuidStr}
		if _, err := f.db.Exec(sqlstr, params...); err != nil {
			return "", errors.Wrapf(err, "[%s] cannot update [%s] - [%v]", uuidStr, sqlstr, params)
		}
		sqlstr = fmt.Sprintf("INSERT INTO %s.core_dirty"+
			" (uuid)"+
			" VALUES($1)", f.dbschema)
		params = []interface{}{
			uuidStr, // uuid
		}
		if _, err = f.db.Exec(sqlstr, params...); err != nil {
			return "", errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		f.log.Infof("item [%s] updated", uuidStr)
		return uuidStr, nil
	}

}

func (f *Fair) GetSets(partitionName string) (map[string]string, error) {
	sqlstr := fmt.Sprintf("SELECT s.setspec, s.setname"+
		" FROM (SELECT DISTINCT unnest(setspec) AS setspecx FROM %s.coreview WHERE partition=$1) specs"+
		" LEFT JOIN %s.set s ON s.setspec=setspecx", f.dbschema, f.dbschema)
	rows, err := f.db.Query(sqlstr, partitionName)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot query sets %s", sqlstr)
	}
	var sets map[string]string = make(map[string]string)
	defer rows.Close()
	for rows.Next() {
		var setspec, setname string
		if err := rows.Scan(&setspec, &setname); err != nil {
			return nil, errors.Wrapf(err, "cannot scan sets query result - %s", sqlstr)
		}
		sets[setspec] = setname
	}
	return sets, nil
}

func (f *Fair) StartUpdate(partitionName string, source string) error {
	src, err := f.GetSourceByName(source, partitionName)
	if err != nil {
		return errors.Wrapf(err, "cannot get source %s", source)
	}
	sqlstr := fmt.Sprintf("DELETE FROM %s.core_dirty"+
		" WHERE uuid IN (SELECT uuid FROM %s.core WHERE source=$1)", f.dbschema, f.dbschema)
	params := []interface{}{src.ID}
	if _, err := f.db.Exec(sqlstr, params...); err != nil {
		return errors.Wrapf(err, "cannot execute dirty update - %s - %v", sqlstr, params)
	}
	return nil
}

func (f *Fair) AbortUpdate(partitionName string, source string) error {
	return f.StartUpdate(partitionName, source)
	/*
		src, err := f.GetSourceByName(source)
		if err != nil {
			return errors.Wrapf(err, "cannot get source %s", source)
		}
		sqlstr := fmt.Sprintf("UPDATE %s.core SET dirty=FALSE WHERE deleted=FALSE AND source=$1", f.dbschema)
		params := []interface{}{src.ID}
		if _, err := f.db.Exec(sqlstr, params...); err != nil {
			return errors.Wrapf(err, "cannot execute dirty reset update - %s - %v", sqlstr, params)
		}
		return nil
	*/
}

func (f *Fair) EndUpdate(partitionName string, source string) error {
	src, err := f.GetSourceByName(source, partitionName)
	if err != nil {
		return errors.Wrapf(err, "cannot get source %s", source)
	}
	sqlstr := fmt.Sprintf("UPDATE %s.core"+
		" SET deleted=TRUE"+
		" WHERE source=$1 AND uuid NOT IN (SELECT uuid FROM %s.core_dirty)", f.dbschema, f.dbschema)
	params := []interface{}{src.ID}
	if _, err := f.db.Exec(sqlstr, params...); err != nil {
		return errors.Wrapf(err, "cannot execute dirty update - %s - %v", sqlstr, params)
	}
	return f.StartUpdate(partitionName, source)
}
