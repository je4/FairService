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
	"sync"
)

type Source struct {
	ID          int64
	Name        string
	DetailURL   string
	Description string
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
	Source    string      `json:"source"`
	Signature string      `json:"signature"`
	Metadata  myfair.Core `json:"metadata"`
	Set       []string    `json:"set"`
	Catalog   []string    `json:"catalog"`
	Access    DataAccess  `json:"access"`
	Deleted   bool        `json:"deleted"`
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

func (f *Fair) LoadSources() error {
	sqlstr := fmt.Sprintf("SELECT sourceid, name, detailurl, description FROM %s.source", f.dbschema)
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
		if err := rows.Scan(&src.ID, &src.Name, &src.DetailURL, &src.Description); err != nil {
			return errors.Wrap(err, "cannot scan values")
		}
		f.sources[src.ID] = src
	}
	return nil
}

func (f *Fair) GetSourceById(id int64) (*Source, error) {
	f.sourcesMutex.RLock()
	defer f.sourcesMutex.RUnlock()
	if s, ok := f.sources[id]; ok {
		return s, nil
	} else {
		return nil, errors.New(fmt.Sprintf("source #%v not found", id))
	}
}

func (f *Fair) GetSourceByName(name string) (*Source, error) {
	f.sourcesMutex.RLock()
	defer f.sourcesMutex.RUnlock()
	for _, src := range f.sources {
		if src.Name == name {
			return src, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("source %s not found", name))
}

func (f *Fair) GetItem(partitionName, uuidStr string) (*ItemData, error) {
	partition, ok := f.partitions[partitionName]
	if !ok {
		return nil, errors.New(fmt.Sprintf("partition %s not found", partitionName))
	}

	sqlstr := fmt.Sprintf("SELECT metadata, setspec, catalog, access, signature, source, deleted"+
		" FROM %s.oai"+
		" WHERE partition=$1 AND uuid=$2", f.dbschema)
	params := []interface{}{partition.Name, uuidStr}
	row := f.db.QueryRow(sqlstr, params...)

	var metaStr string
	var set, catalog []string
	var accessStr string
	var signature string
	var sourceStr string
	var deleted bool
	if err := row.Scan(&metaStr, pq.Array(&set), pq.Array(&catalog), &accessStr, &signature, &sourceStr, &deleted); err != nil {
		if err != sql.ErrNoRows {
			return nil, errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		return nil, nil
	}
	data := &ItemData{
		Source:    sourceStr,
		Signature: signature,
		Metadata:  myfair.Core{},
		Set:       set,
		Catalog:   catalog,
		Deleted:   deleted,
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

	src, err := f.GetSourceByName(data.Source)
	if err != nil {
		return "", errors.Wrapf(err, "cannot get source %s", data.Source)
	}

	sqlstr := fmt.Sprintf("SELECT uuid, metadata, setspec, catalog, access, deleted"+
		" FROM %s.oai "+
		" WHERE partition=$1 AND source=$2 AND signature=$3", f.dbschema)
	params := []interface{}{partition.Name, src.ID, data.Signature}
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
		sqlstr := fmt.Sprintf("INSERT INTO %s.oai"+
			" (uuid, partition, datestamp, setspec, metadata, dirty, signature, source, access, catalog, seq)"+
			" VALUES($1, $2, NOW(), $3, $4, false, $5, $6, $7, $8, NEXTVAL('lastchange'))", f.dbschema)
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
			return uuidStr, nil
		}

		dataMetaBytes, err := json.Marshal(data.Metadata)
		if err != nil {
			return "", errors.Wrapf(err, "[%s] cannot unmarshal data core", uuidStr)
		}

		sqlstr = fmt.Sprintf("UPDATE %s.oai"+
			" SET datestamp=NOW(), setspec=$1, metadata=$2, dirty=FALSE, access=$3, catalog=$4, seq=NEXTVAL('lastchange'), deleted=false"+
			" WHERE uuid=$5", f.dbschema)
		params := []interface{}{
			pq.Array(data.Set),
			string(dataMetaBytes),
			data.Access,
			pq.Array(data.Catalog),
			uuidStr}
		if _, err := f.db.Exec(sqlstr, params...); err != nil {
			return "", errors.Wrapf(err, "[%s] cannot update [%s] - [%v]: %v", uuidStr, sqlstr, params)
		}
		f.log.Infof("item [%s] updated")
		return uuidStr, nil
	}

}
