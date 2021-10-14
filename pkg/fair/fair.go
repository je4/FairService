package fair

import (
	"database/sql"
	"emperror.dev/errors"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/je4/FairService/v2/pkg/datatable"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/FairService/v2/pkg/model/myfair"
	"github.com/je4/FairService/v2/pkg/service/datacite"
	//"github.com/je4/FairService/v2/pkg/service"
	"github.com/lib/pq"
	"github.com/op/go-logging"
	"net/url"
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
	DataAccessOpenAccess DataAccess = "open_access"
)

var DataAccessReverse = map[string]DataAccess{
	string(DataAccessPublic):     DataAccessPublic,
	string(DataAccessClosed):     DataAccessClosed,
	string(DataAccessClosedData): DataAccessClosedData,
	string(DataAccessOpenAccess): DataAccessOpenAccess,
}

type DataStatus string

const (
	DataStatusActive   DataStatus = "active"
	DataStatusDisabled DataStatus = "disabled"
	DataStatusDeleted  DataStatus = "deleted"
)

var DataStatusReverse = map[string]DataStatus{
	string(DataStatusActive):   DataStatusActive,
	string(DataStatusDisabled): DataStatusDisabled,
	string(DataStatusDeleted):  DataStatusDeleted,
}

type ItemData struct {
	Source     string      `json:"source"`
	Signature  string      `json:"signature"`
	Metadata   myfair.Core `json:"metadata"`
	Set        []string    `json:"set"`
	Catalog    []string    `json:"catalog"`
	Identifier []string    `json:"identifier"`
	Access     DataAccess  `json:"access"`
	Status     DataStatus  `json:"status"`
	Seq        int64       `json:"-"`
	UUID       string      `json:"uuid"`
	Datestamp  time.Time   `json:"datestamp"`
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
	dbSchema       string
	db             *sql.DB
	handle         *HandleServiceClient
	dataciteClient *datacite.Client
	sourcesMutex   sync.RWMutex
	sources        map[int64]*Source
	partitions     map[string]*Partition
	log            *logging.Logger
}

func NewFair(db *sql.DB, dbSchema string, handle *HandleServiceClient, dataciteClient *datacite.Client, log *logging.Logger) (*Fair, error) {
	f := &Fair{
		dbSchema:       dbSchema,
		db:             db,
		handle:         handle,
		dataciteClient: dataciteClient,
		sourcesMutex:   sync.RWMutex{},
		sources:        map[int64]*Source{},
		partitions:     map[string]*Partition{},
		log:            log,
	}
	if err := f.LoadSources(); err != nil {
		return nil, errors.Wrap(err, "cannot load sources")
	}
	return f, nil
}

func (f *Fair) nextCounter(name string) (int64, error) {
	sqlStr := fmt.Sprintf("SELECT NEXTVAL('%s.%s')", f.dbSchema, name)
	var next int64
	if err := f.db.QueryRow(sqlStr).Scan(&next); err != nil {
		return 0, errors.Wrapf(err, "cannot execute %s", sqlStr)
	}
	return next, nil
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

func (f *Fair) GetPartitions() map[string]*Partition {
	return f.partitions
}

func (f *Fair) LoadSources() error {
	sqlstr := fmt.Sprintf("SELECT sourceid, name, detailurl, description, oai_domain, partition FROM %s.source", f.dbSchema)
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
		" WHERE partition=$1", f.dbSchema)
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
			" WHERE %s", f.dbSchema, sqlWhere)
		if err := f.db.QueryRow(sqlstr, params...).Scan(completeListSize); err != nil {
			return errors.Wrapf(err, "cannot get number of result items [%s] - [%v]", sqlstr, params)
		}
	}
	sqlstr := fmt.Sprintf("SELECT uuid, metadata, setspec, catalog, access, signature, source, status, seq, datestamp, identifier"+
		" FROM %s.coreview"+
		" WHERE %s"+
		" ORDER BY seq ASC", f.dbSchema, sqlWhere)
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
		var source string
		var statusStr string
		var seq int64
		var identifier []string
		var datestamp time.Time
		if err := rows.Scan(&uuidStr, &metaStr, pq.Array(&set), pq.Array(&catalog), &accessStr, &signature, &source, &statusStr, &seq, &datestamp, pq.Array(&identifier)); err != nil {
			return errors.Wrapf(err, "cannot scan result of [%s] - [%v]", sqlstr, params)
		}
		data := &ItemData{
			UUID:       uuidStr,
			Source:     source,
			Signature:  signature,
			Metadata:   myfair.Core{},
			Set:        set,
			Catalog:    catalog,
			Seq:        seq,
			Datestamp:  datestamp,
			Identifier: identifier,
		}
		var ok bool
		data.Access, ok = DataAccessReverse[accessStr]
		if !ok {
			return errors.New(fmt.Sprintf("[%s] invalid access type %s", uuidStr, accessStr))
		}
		data.Status, ok = DataStatusReverse[statusStr]
		if !ok {
			return errors.New(fmt.Sprintf("[%s] invalid status type %s", uuidStr, accessStr))
		}
		if err := json.Unmarshal([]byte(metaStr), &data.Metadata); err != nil {
			return errors.Wrapf(err, "[%s] cannot unmarshal core [%s]", uuidStr, metaStr)
		}
		// add local identifiers
		for _, id := range data.Identifier {
			strs := strings.SplitN(id, ":", 2)
			if len(strs) != 2 {
				continue
			}

			idType, ok := myfair.RelatedIdentifierTypeReverse[strs[0]]
			if !ok {
				f.log.Warningf("[%s] unknown identifier type %s", uuidStr, id)
				continue
			}
			idStr := strs[1]
			found := false
			for _, di := range data.Metadata.Identifier {
				if di.IdentifierType == idType && di.Value == idStr {
					found = true
					break
				}
			}
			if !found {
				data.Metadata.Identifier = append(data.Metadata.Identifier, myfair.Identifier{
					Value:          idStr,
					IdentifierType: idType,
				})
			}
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

	sqlstr := fmt.Sprintf("SELECT metadata, setspec, catalog, access, signature, sourcename, status, seq, datestamp, identifier"+
		" FROM %s.coreview"+
		" WHERE partition=$1 AND uuid=$2", f.dbSchema)
	params := []interface{}{partition.Name, uuidStr}
	row := f.db.QueryRow(sqlstr, params...)

	var metaStr string
	var set, catalog []string
	var accessStr string
	var signature string
	var source string
	var statusStr string
	var seq int64
	var datestamp time.Time
	var identifier []string
	if err := row.Scan(&metaStr, pq.Array(&set), pq.Array(&catalog), &accessStr, &signature, &source, &statusStr, &seq, &datestamp, pq.Array(&identifier)); err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		return nil, nil
	}
	data := &ItemData{
		UUID:       uuidStr,
		Source:     source,
		Signature:  signature,
		Metadata:   myfair.Core{},
		Set:        set,
		Catalog:    catalog,
		Identifier: identifier,
		Seq:        seq,
		Datestamp:  datestamp,
	}
	data.Access, ok = DataAccessReverse[accessStr]
	if !ok {
		return nil, errors.New(fmt.Sprintf("[%s] invalid access type %s", uuidStr, accessStr))
	}
	if err := json.Unmarshal([]byte(metaStr), &data.Metadata); err != nil {
		return nil, errors.Wrapf(err, "[%s] cannot unmarshal core [%s]", uuidStr, metaStr)
	}
	data.Status, ok = DataStatusReverse[statusStr]
	if !ok {
		return nil, errors.New(fmt.Sprintf("[%s] invalid status type %s", uuidStr, accessStr))
	}
	// add local identifiers
	for _, id := range data.Identifier {
		strs := strings.SplitN(id, ":", 2)
		if len(strs) != 2 {
			continue
		}

		idType, ok := myfair.RelatedIdentifierTypeReverse[strs[0]]
		if !ok {
			f.log.Warningf("[%s] unknown identifier type %s", uuidStr, id)
			continue
		}
		idStr := strs[1]
		found := false
		for _, di := range data.Metadata.Identifier {
			if di.IdentifierType == idType && di.Value == idStr {
				found = true
				break
			}
		}
		if !found {
			data.Metadata.Identifier = append(data.Metadata.Identifier, myfair.Identifier{
				Value:          idStr,
				IdentifierType: idType,
			})
		}
	}

	return data, nil
}

func (f *Fair) DeleteItem(partitionName, uuidStr string) error {
	/*
		partition, ok := f.partitions[partitionName]
		if !ok {
			return errors.New(fmt.Sprintf("partition %s not found", partitionName))
		}
	*/
	data, err := f.GetItem(partitionName, uuidStr)
	if err != nil {
		return errors.Wrapf(err, "cannot get item %s/%s", partitionName, uuidStr)
	}
	// if already deleted, don't do anything
	if data.Status != DataStatusActive {
		return nil
	}
	return errors.New("function DeleteItem() not implemented")
	/*
		doiPrefix := fmt.Sprintf("%s:%s/", myfair.RelatedIdentifierTypeDOI, f.dataciteClient.GetPrefix())
		for _, id := range data.Identifier {
			if strings.HasPrefix(id, doiPrefix) {
				doi := strings.TrimPrefix(id, fmt.Sprintf("%s:", myfair.RelatedIdentifierTypeDOI))
				if _, err := f.dataciteClient.Delete(doi); err != nil {
					if _, err := f.dataciteClient.SetEvent(doi, datacite.DCEventHide); err != nil {
						return errors.Wrapf(err, "cannot hide doi %s", doi)
					}
				} else {
					// todo: remove doi from metadata and store it

				}
			}
		}
		return nil
	*/
}

func (f *Fair) RefreshSearch() error {
	f.log.Info("refreshing meterialized view searchable")
	sqlstr := fmt.Sprintf("SELECT %s.refresh()", f.dbSchema)
	_, err := f.db.Exec(sqlstr)
	if err != nil {
		f.log.Errorf("cannot refresh materialized view: %v - %v", sqlstr, err)
		return errors.Wrapf(err, "error refreshing search view %s", sqlstr)
	}
	return nil
}

func (f *Fair) CreateItem(partitionName string, data *ItemData) (string, error) {
	partition, ok := f.partitions[partitionName]
	if !ok {
		return "", errors.New(fmt.Sprintf("partition %s not found", partitionName))
	}

	sort.Strings(data.Catalog)
	sort.Strings(data.Set)

	src, err := f.GetSourceByName(data.Source, partitionName)
	if err != nil {
		return "", errors.Wrapf(err, "cannot get source %s", data.Source)
	}

	sqlstr := fmt.Sprintf("SELECT uuid, metadata, setspec, catalog, access, status, identifier"+
		" FROM %s.coreview "+
		" WHERE source=$1 AND signature=$2", f.dbSchema)
	params := []interface{}{src.ID, data.Signature}
	row := f.db.QueryRow(sqlstr, params...)

	var metaStr string
	var uuidStr string
	var set, catalog, identifiers []string
	var accessStr string
	var statusStr string
	if err := row.Scan(&uuidStr, &metaStr, pq.Array(&set), pq.Array(&catalog), &accessStr, &statusStr, pq.Array(&identifiers)); err != nil {
		if err != sql.ErrNoRows {
			return "", errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}

		//
		// Create new Entry
		//

		uuidVal, err := uuid.NewUUID()
		if err != nil {
			return "", errors.Wrap(err, "cannot generate uuid")
		}
		uuidStr := uuidVal.String()
		identifiers = []string{}
		var sqlHandle = sql.NullString{}
		if f.handle != nil {
			//
			// Create Handle and add to local identifier list
			//
			next, err := f.nextCounter("handle")
			if err != nil {
				return "", errors.Wrap(err, "cannot get next handle value")
			}
			newHandle := fmt.Sprintf("%s/%s/%v", partition.HandleID, partition.HandlePrefix, next)
			newURL, err := url.Parse(fmt.Sprintf("%s/redir/%s", partition.AddrExt, uuidStr))
			if err != nil {
				return "", errors.Wrapf(err, "cannot parse url %s", fmt.Sprintf("%s/redir/%s", partition.AddrExt, uuidStr))
			}
			if err := f.handle.Create(newHandle, newURL); err != nil {
				return "", errors.Wrapf(err, "cannot create handle %s for %s", newHandle, newURL.String())
			}
			sqlHandle.String = newHandle
			sqlHandle.Valid = true
			data.Metadata.Identifier = append(data.Metadata.Identifier, myfair.Identifier{
				Value:          newHandle,
				IdentifierType: myfair.RelatedIdentifierTypeHandle,
			})
			identifiers = append(identifiers, fmt.Sprintf("%s:%s", myfair.RelatedIdentifierTypeHandle, newHandle))
		}

		coreBytes, err := json.Marshal(data.Metadata)
		if err != nil {
			return "", errors.Wrapf(err, "cannot marshal core data [%v]", data.Metadata)
		}
		sqlstr := fmt.Sprintf("INSERT INTO %s.core"+
			" (uuid, datestamp, setspec, metadata, signature, source, access, catalog, seq, status, identifier)"+
			" VALUES($1, NOW(), $2, $3, $4, $5, $6, $7, NEXTVAL('lastchange'), $8, $9)", f.dbSchema)
		params := []interface{}{
			uuidStr, // uuid
			// datestamp
			pq.Array(data.Set),     // setspec
			string(coreBytes),      // metadata
			data.Signature,         // signature
			src.ID,                 // source
			data.Access,            // access
			pq.Array(data.Catalog), // catalog
			// seq
			DataStatusActive,
			pq.Array(identifiers), // identifier
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
			" VALUES($1)", f.dbSchema)
		params = []interface{}{
			uuidStr, // uuid
		}
		if _, err = f.db.Exec(sqlstr, params...); err != nil {
			return "", errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		f.log.Infof("new item [%s] inserted", uuidStr)

		return uuidStr, nil

	} else {
		var ok bool
		access, ok := DataAccessReverse[accessStr]
		if !ok {
			return "", errors.Wrapf(err, "[%s] invalid access type %s", uuidStr, accessStr)
		}
		status, ok := DataStatusReverse[statusStr]
		if !ok {
			return "", errors.Wrapf(err, "[%s] invalid status type %s", uuidStr, accessStr)
		}
		sort.Strings(catalog)
		sort.Strings(set)
		sort.Strings(identifiers)

		// add local identifiers
		for _, id := range identifiers {
			strs := strings.SplitN(id, ":", 2)
			if len(strs) != 2 {
				return "", errors.New(fmt.Sprintf("[%s] invalid identifier format %s", uuidStr, id))
			}

			idType, ok := myfair.RelatedIdentifierTypeReverse[strs[0]]
			if !ok {
				f.log.Warningf("[%s] unknown identifier type %s", uuidStr, id)
				continue
			}
			idStr := strs[1]
			found := false
			for _, di := range data.Metadata.Identifier {
				if di.IdentifierType == idType && di.Value == idStr {
					found = true
					break
				}
			}
			if !found {
				data.Metadata.Identifier = append(data.Metadata.Identifier, myfair.Identifier{
					Value:          idStr,
					IdentifierType: idType,
				})
			}
		}

		/*
			identifiers = []string{}
			for _, id := range data.Metadata.Identifier {
				identifiers = append(identifiers, fmt.Sprintf("%s:%s", id.IdentifierType, id.Value))
			}
		*/

		// do update here
		var oldMeta myfair.Core
		if err := json.Unmarshal([]byte(metaStr), &oldMeta); err != nil {
			return "", errors.Wrapf(err, "[%s] cannot unmarshal core [%s]", uuidStr, metaStr)
		}
		if status == "active" &&
			reflect.DeepEqual(oldMeta, data.Metadata) &&
			equalStrings(set, data.Set) &&
			equalStrings(catalog, data.Catalog) &&
			access == data.Access {
			f.log.Infof("no update needed for item [%v]", uuidStr)
			sqlstr = fmt.Sprintf("INSERT INTO %s.core_dirty"+
				" (uuid)"+
				" VALUES($1)", f.dbSchema)
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
			" SET setspec=$1, metadata=$2, access=$3, catalog=$4, status=$5, datestamp=NOW(), seq=NEXTVAL('lastchange')"+
			" WHERE uuid=$6", f.dbSchema)
		params := []interface{}{
			pq.Array(data.Set),
			string(dataMetaBytes),
			data.Access,
			pq.Array(data.Catalog),
			DataStatusActive,
			uuidStr}
		if _, err := f.db.Exec(sqlstr, params...); err != nil {
			return "", errors.Wrapf(err, "[%s] cannot update [%s] - [%v]", uuidStr, sqlstr, params)
		}
		sqlstr = fmt.Sprintf("INSERT INTO %s.core_dirty"+
			" (uuid)"+
			" VALUES($1)", f.dbSchema)
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
		" LEFT JOIN %s.set s ON s.setspec=setspecx", f.dbSchema, f.dbSchema)
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
		" WHERE uuid IN (SELECT uuid FROM %s.core WHERE source=$1)", f.dbSchema, f.dbSchema)
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
		sqlstr := fmt.Sprintf("UPDATE %s.core SET dirty=FALSE WHERE deleted=FALSE AND source=$1", f.dbSchema)
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
		" SET status=$1"+
		" WHERE source=$2 AND uuid NOT IN (SELECT uuid FROM %s.core_dirty)", f.dbSchema, f.dbSchema)
	params := []interface{}{DataStatusDisabled, src.ID}
	if _, err := f.db.Exec(sqlstr, params...); err != nil {
		return errors.Wrapf(err, "cannot execute dirty update - %s - %v", sqlstr, params)
	}
	f.RefreshSearch()
	return f.StartUpdate(partitionName, source)
}

func (f *Fair) CreateDOI(partitionName, uuidStr, targetUrl string) (*datacite.API, error) {
	part, err := f.GetPartition(partitionName)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get partition %s", partitionName)
	}

	data, err := f.GetItem(partitionName, uuidStr)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading item")
	}
	if data == nil {
		return nil, errors.New(fmt.Sprintf("item %s/%s not found", partitionName, uuidStr))
	}

	if data.Status != DataStatusActive {
		return nil, errors.New(fmt.Sprintf("item %s/%s is not active", partitionName, uuidStr))
	}

	var hasDOI string
	var doiPrefix = fmt.Sprintf("%s:", myfair.RelatedIdentifierTypeDOI)
	for _, id := range data.Identifier {
		if strings.HasPrefix(id, doiPrefix) {
			hasDOI = strings.TrimPrefix(id, doiPrefix)
			break
		}
	}
	if hasDOI != "" {
		return nil, errors.New(fmt.Sprintf("doi %s for uuid %s already exists", hasDOI, uuidStr))
	}

	next, err := f.nextCounter("doi")
	if err != nil {
		return nil, errors.Wrap(err, "cannot get next doi sequence value")
	}

	doiSuffix := fmt.Sprintf("%s/%v", part.HandlePrefix, next)
	doiStr := fmt.Sprintf("%s/%s", f.dataciteClient.GetPrefix(), doiSuffix)
	_, err = f.dataciteClient.RetrieveDOI(doiStr)
	if err == nil {
		return nil, errors.New(fmt.Sprintf("doi %s already exists", doiStr))
	}

	dataciteData := &dataciteModel.DataCite{}
	dataciteData.InitNamespace()
	dataciteData.FromCore(data.Metadata)

	api, err := f.dataciteClient.CreateDOI(dataciteData, doiSuffix, targetUrl)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create doi")
	}

	data.Identifier = append(data.Identifier, fmt.Sprintf("%s:%s", myfair.RelatedIdentifierTypeDOI, doiStr))
	sqlstr := fmt.Sprintf("UPDATE %s.core"+
		" SET identifiers=$1, datestamp=NOW(), seq=NEXTVAL('lastchange')"+
		" WHERE uuid=$2", f.dbSchema)
	params := []interface{}{
		pq.Array(data.Identifier),
		data.UUID,
	}
	if _, err := f.db.Exec(sqlstr, params...); err != nil {
		return nil, errors.Wrapf(err, "[%s] cannot update [%s] - [%v]", uuidStr, sqlstr, params)
	}
	f.RefreshSearch()
	return api, nil
}

var fieldDef = map[string]string{
	"partition":       "src.partition AS partition",
	"sourcename":      "src.name AS sourcename",
	"uuid":            "s.uuid AS uuid",
	"persons":         "array_to_string(s.persons, '; ') AS persons",
	"titles":          "array_to_string(s.titles, '; ') AS titles",
	"sets":            "array_to_string(s.setspec, '; ') AS sets",
	"catalogs":        "array_to_string(s.catalog, '; ') AS catalogs",
	"signature":       "s.signature AS signature",
	"access":          "s.access AS access",
	"status":          "s.status::TEXT AS status",
	"identifiers":     "array_to_string(s.identifier, '; ') AS identifiers",
	"handle":          "s.handle AS handle",
	"doi":             "s.doi AS doi",
	"resourcetype":    "s.resourcetype AS resourcetype",
	"publicationyear": "s.publicationyear AS publicationyear",
}

func (f *Fair) Search(partitionName string, dtr *datatable.Request) ([]map[string]string, int64, int64, error) {
	part, err := f.GetPartition(partitionName)
	if err != nil {
		return nil, 0, 0, errors.Wrapf(err, "cannot get partition %s", partitionName)
	}

	var fields []string
	var resultData = map[string]*string{}
	var resultVals []interface{}

	for key, col := range dtr.Columns {
		str, ok := fieldDef[col.Data]
		if !ok {
			return nil, 0, 0, errors.New(fmt.Sprintf("invalid field name for column %v: %s", key, col.Data))
		}
		fields = append(fields, str)
		var h string
		resultData[col.Data] = &h
		resultVals = append(resultVals, &h)
	}
	sqlFields := strings.Join(fields, ", ")

	var orderList = []string{}
	for key, order := range dtr.Order {
		col, ok := dtr.Columns[order.Column]
		if !ok {
			return nil, 0, 0, errors.New(fmt.Sprintf("invalid column nuber %v in order %v", order.Column, key))
		}
		switch order.Dir {
		case "asc":
			orderList = append(orderList, fmt.Sprintf("%s %s", col.Data, order.Dir))
		case "desc":
			orderList = append(orderList, fmt.Sprintf("%s %s", col.Data, order.Dir))
		}
	}

	var params = []interface{}{part.Name}
	sqlWhere := "s.source=src.sourceid AND src.partition=$1 "

	sqlstr := fmt.Sprintf("SELECT COUNT(*) AS num"+
		" FROM %s.searchable s, %s.source src "+
		" WHERE %s", f.dbSchema, f.dbSchema, sqlWhere)
	var num, total int64
	if err := f.db.QueryRow(sqlstr, params...).Scan(&total); err != nil {
		return nil, 0, 0, errors.Wrapf(err, "cannot execute query %s - %v", sqlstr, params)
	}

	if dtr.Search.Value != "" {
		params = append(params, dtr.Search.Value+":*")
		sqlWhere += " AND s.fulltext @@ to_tsquery($2) "
		sqlstr = fmt.Sprintf("SELECT COUNT(*) AS num"+
			" FROM %s.searchable s, %s.source src "+
			" WHERE %s", f.dbSchema, f.dbSchema, sqlWhere)
	}
	if err := f.db.QueryRow(sqlstr, params...).Scan(&num); err != nil {
		return nil, 0, 0, errors.Wrapf(err, "cannot execute query %s - %v", sqlstr, params)
	}

	if dtr.Start >= num {
		return []map[string]string{}, num, total, nil
	}

	sqlOrder := ""
	if len(orderList) > 0 {
		sqlOrder = fmt.Sprintf("ORDER BY %s", strings.Join(orderList, ", "))
	}

	sqlstr = fmt.Sprintf("SELECT %s "+
		" FROM %s.searchable s, %s.source src "+
		" WHERE %s "+
		" %s "+
		" LIMIT %v OFFSET %v", sqlFields, f.dbSchema, f.dbSchema, sqlWhere, sqlOrder, dtr.Length, dtr.Start)
	f.log.Infof("%s - %v", sqlstr, params)
	rows, err := f.db.Query(sqlstr, params...)
	if err != nil {
		return nil, 0, 0, errors.Wrapf(err, "cannot execute query %s - %v", sqlstr, params)
	}
	var result []map[string]string
	for rows.Next() {
		if err := rows.Scan(resultVals...); err != nil {
			return nil, 0, 0, errors.Wrapf(err, "cannot scan row %s - %v", sqlstr, params)
		}
		var rLine = map[string]string{}
		for key, val := range resultData {
			rLine[key] = *val
		}
		result = append(result, rLine)
	}
	return result, num, total, nil
}
