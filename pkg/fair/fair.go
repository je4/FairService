package fair

import (
	"bytes"
	"compress/gzip"
	"context"
	"emperror.dev/errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/FairService/v2/pkg/model/myfair"
	"github.com/je4/FairService/v2/pkg/service/datacite"
	hcClient "github.com/je4/HandleCreator/v2/pkg/client"
	"github.com/je4/utils/v2/pkg/datatable"
	"github.com/je4/utils/v2/pkg/zLogger"
	"io"

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
	DataStatusActive      DataStatus = "active"
	DataStatusDisabled    DataStatus = "disabled"
	DataStatusDeleted     DataStatus = "deleted"
	DataStatusDeletedMeta DataStatus = "deleted_meta"
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

type ArchiveItem struct {
	ItemData
	NewFiles []string
}

type Archive struct {
	Name         string
	CreationDate time.Time
	LastVersion  int64
	Description  string
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
	db             *pgxpool.Pool
	handle         *hcClient.HandleCreatorClient
	dataciteClient *datacite.Client
	sourcesMutex   sync.RWMutex
	sources        map[int64]*Source
	partitions     map[string]*Partition
	log            zLogger.ZLogger
}

func NewFair(db *pgxpool.Pool, dbSchema string, handle *hcClient.HandleCreatorClient, dataciteClient *datacite.Client, log zLogger.ZLogger) (*Fair, error) {
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
	//	sqlStr := fmt.Sprintf("SELECT NEXTVAL('%s.%s')", f.dbSchema, name)
	sqlStr := fmt.Sprintf("SELECT NEXTVAL('%s')", name)
	var next int64
	if err := f.db.QueryRow(context.Background(), sqlStr).Scan(&next); err != nil {
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

func (f *Fair) GetMinimumDatestamp(partition *Partition) (time.Time, error) {
	/*	sqlstr := fmt.Sprintf("SELECT MIN(datestamp) AS mindate"+
		" FROM %s.coreview"+
		" WHERE partition=$1", f.dbSchema)
	*/sqlstr := "SELECT MIN(datestamp) AS mindate" +
		" FROM coreview" +
		" WHERE partition=$1"
	params := []interface{}{partition.Name}
	var datestamp time.Time
	if err := f.db.QueryRow(context.Background(), sqlstr, params...).Scan(&datestamp); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return time.Time{}, errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		return time.Now(), nil
	}
	return datestamp, nil
}

func (f *Fair) RefreshSearch() error {
	f.log.Info().Msg("refreshing meterialized view searchable")
	//	sqlstr := fmt.Sprintf("SELECT %s.refresh()", f.dbSchema)
	sqlstr := "SELECT refresh()"
	_, err := f.db.Exec(context.Background(), sqlstr)
	if err != nil {
		f.log.Error().Msgf("cannot refresh materialized view: %v - %v", sqlstr, err)
		return errors.Wrapf(err, "error refreshing search view %s", sqlstr)
	}
	return nil
}

func (f *Fair) SetOriginalData(p *Partition, uuid string, format string, data []byte) error {
	formatOK := false
	for _, f := range []string{"Other", "XML", "JSON"} {
		if f == format {
			formatOK = true
			break
		}
	}
	if !formatOK {
		return errors.Errorf("invalid format. only \"Other\", \"XML\", \"JSON\" allowed: %s", format)
	}

	var compressed = false
	var target *bytes.Buffer
	if len(data) < 4*1024 {
		target = bytes.NewBuffer(data)
	} else {
		target = bytes.NewBuffer(nil)
		w := gzip.NewWriter(target)
		if _, err := w.Write(data); err != nil {
			return errors.Wrapf(err, "cannot write compressed metadata")
		}
		w.Close()
		compressed = true
	}
	//	sqlstr := fmt.Sprintf("SELECT COUNT(*) FROM %s.originaldata WHERE uuid=$1", f.dbSchema)
	sqlstr := "SELECT COUNT(*) FROM originaldata WHERE uuid=$1"
	var num int64
	if err := f.db.QueryRow(context.Background(), sqlstr, uuid).Scan(&num); err != nil {
		return errors.Wrapf(err, "cannot query %s - [%v]", sqlstr, uuid)
	}
	if num > 0 {
		//sqlstr = fmt.Sprintf("UPDATE %s.originaldata SET type=$1, data=$2, compressed=$3 WHERE uuid=$4", f.dbSchema)
		sqlstr = "UPDATE originaldata SET type=$1, data=$2, compressed=$3 WHERE uuid=$4"
	} else {
		// sqlstr = fmt.Sprintf("INSERT INTO %s.originaldata (type, data, compressed, uuid) VALUES( $1, $2, $3, $4)", f.dbSchema)
		sqlstr = "INSERT INTO originaldata (type, data, compressed, uuid) VALUES( $1, $2, $3, $4)"
	}
	if _, err := f.db.Exec(context.Background(), sqlstr, format, target.Bytes(), compressed, uuid); err != nil {
		return errors.Wrapf(err, "cannot query %s - [%v, %v, %v, %v]", sqlstr, format, target, compressed, uuid)
	}
	return nil
}

func (f *Fair) GetOriginalData(p *Partition, uuid string) ([]byte, string, error) {

	// sqlstr := fmt.Sprintf("SELECT type, data, compressed FROM %s.originaldata WHERE uuid=$1", f.dbSchema)
	sqlstr := "SELECT type, data, compressed FROM originaldata WHERE uuid=$1"
	var t string
	var compressed bool
	var data []byte
	if err := f.db.QueryRow(context.Background(), sqlstr, uuid).Scan(&t, &compressed, &data); err != nil {
		return nil, "", errors.Wrapf(err, "cannot query %s - [%v]", sqlstr, uuid)
	}
	if compressed {
		buf := bytes.NewBuffer(data)
		reader, err := gzip.NewReader(buf)
		if err != nil {
			return nil, "", errors.Wrapf(err, "cannot create gzip reader")
		}
		d2, err := io.ReadAll(reader)
		if err != nil {
			return nil, "", errors.Wrapf(err, "cannot read zipped data")
		}
		return d2, t, nil

	}
	return data, t, nil
}

func (f *Fair) GetSets(p *Partition) (map[string]string, error) {
	/*	sqlstr := fmt.Sprintf("SELECT specs.setspecx, s.setname"+
		" FROM (SELECT DISTINCT unnest(setspec) AS setspecx FROM %s.coreview WHERE partition=$1) specs"+
		" LEFT JOIN %s.set s ON s.setspec=setspecx", f.dbSchema, f.dbSchema)
	*/
	sqlstr := "SELECT specs.setspecx, s.setname" +
		" FROM (SELECT DISTINCT unnest(setspec) AS setspecx FROM coreview WHERE partition=$1) specs" +
		" LEFT JOIN set s ON s.setspec=setspecx"
	rows, err := f.db.Query(context.Background(), sqlstr, p.Name)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot query sets %s", sqlstr)
	}
	var sets map[string]string = make(map[string]string)
	defer rows.Close()
	for rows.Next() {
		var setspec string
		var setname zeronull.Text
		if err := rows.Scan(&setspec, &setname); err != nil {
			return nil, errors.Wrapf(err, "cannot scan sets query result - %s", sqlstr)
		}
		sets[setspec] = setspec
	}
	return sets, nil
}

func (f *Fair) StartUpdate(p *Partition, source string) error {
	src, err := f.GetSourceByName(p, source)
	if err != nil {
		return errors.Wrapf(err, "cannot get source %s", source)
	}
	/*	sqlstr := fmt.Sprintf("DELETE FROM %s.core_dirty"+
		" WHERE uuid IN (SELECT uuid FROM %s.core WHERE source=$1)", f.dbSchema, f.dbSchema)
	*/
	sqlstr := "DELETE FROM core_dirty WHERE uuid IN (SELECT uuid FROM core WHERE source=$1)"
	params := []interface{}{src.ID}
	if _, err := f.db.Exec(context.Background(), sqlstr, params...); err != nil {
		return errors.Wrapf(err, "cannot execute dirty update - %s - %v", sqlstr, params)
	}
	return nil
}

func (f *Fair) AbortUpdate(p *Partition, source string) error {
	return f.StartUpdate(p, source)
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

func (f *Fair) EndUpdate(p *Partition, source string) error {
	src, err := f.GetSourceByName(p, source)
	if err != nil {
		return errors.Wrapf(err, "cannot get source %s", source)
	}
	/*	sqlstr := fmt.Sprintf("UPDATE %s.core"+
		" SET status=$1"+
		" WHERE source=$2 AND uuid NOT IN (SELECT uuid FROM %s.core_dirty)", f.dbSchema, f.dbSchema)
	*/
	sqlstr := "UPDATE core SET status=$1 WHERE source=$2 AND uuid NOT IN (SELECT uuid FROM core_dirty)"
	params := []interface{}{DataStatusDeletedMeta, src.ID}
	if _, err := f.db.Exec(context.Background(), sqlstr, params...); err != nil {
		return errors.Wrapf(err, "cannot execute dirty update - %s - %v", sqlstr, params)
	}
	f.RefreshSearch()
	return f.StartUpdate(p, source)
}

func (f *Fair) CreateDOI(p *Partition, uuidStr, targetUrl string) (*datacite.API, error) {
	data, err := f.GetItem(p, uuidStr)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading item")
	}
	if data == nil {
		return nil, errors.New(fmt.Sprintf("item %s/%s not found", p.Name, uuidStr))
	}

	if data.Status != DataStatusActive {
		return nil, errors.New(fmt.Sprintf("item %s/%s is not active", p.Name, uuidStr))
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

	doiSuffix := fmt.Sprintf("%s/%v", p.HandlePrefix, next)
	doiStr := fmt.Sprintf("%s/%s", f.dataciteClient.GetPrefix(), doiSuffix)
	_, err = f.dataciteClient.RetrieveDOI(doiStr)
	if err == nil {
		return nil, errors.New(fmt.Sprintf("doi %s already exists", doiStr))
	}

	dataciteData := &dataciteModel.DataCite{}
	dataciteData.InitNamespace()
	dataciteData.FromCore(data.Metadata)

	api, err := f.dataciteClient.CreateDOI(dataciteData, doiSuffix, targetUrl, datacite.DCEventDraft)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create doi")
	}

	data.Identifier = append(data.Identifier, fmt.Sprintf("%s:%s", myfair.RelatedIdentifierTypeDOI, doiStr))
	/*	sqlstr := fmt.Sprintf("UPDATE %s.core"+
		" SET identifier=$1, datestamp=NOW(), seq=NEXTVAL('lastchange')"+
		" WHERE uuid=$2", f.dbSchema)
	*/
	sqlstr := "UPDATE core" +
		" SET identifier=$1, datestamp=NOW(), seq=NEXTVAL('lastchange')" +
		" WHERE uuid=$2"
	params := []interface{}{
		pgtype.FlatArray[string](data.Identifier),
		data.UUID,
	}
	if _, err := f.db.Exec(context.Background(), sqlstr, params...); err != nil {
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

func (f *Fair) Search(p *Partition, dtr *datatable.Request) ([]map[string]string, int64, int64, error) {
	part, err := f.GetPartition(p.Name)
	if err != nil {
		return nil, 0, 0, errors.Wrapf(err, "cannot get partition %s", p.Name)
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

	/*	sqlstr := fmt.Sprintf("SELECT COUNT(*) AS num"+
		" FROM %s.searchable s, %s.source src "+
		" WHERE %s", f.dbSchema, f.dbSchema, sqlWhere)
	*/
	sqlstr := fmt.Sprintf("SELECT COUNT(*) AS num"+
		" FROM searchable s, source src "+
		" WHERE %s", sqlWhere)
	var num, total int64
	if err := f.db.QueryRow(context.Background(), sqlstr, params...).Scan(&total); err != nil {
		return nil, 0, 0, errors.Wrapf(err, "cannot execute query %s - %v", sqlstr, params)
	}

	if dtr.Search.Value != "" {
		tsquery := ""
		for _, fld := range strings.Fields(dtr.Search.Value) {
			if tsquery != "" {
				tsquery += " & "
			}
			tsquery += fld + ":*"
		}
		params = append(params, "simple", tsquery)
		sqlWhere += " AND s.fulltext @@ to_tsquery($2, $3) "
		/*		sqlstr = fmt.Sprintf("SELECT COUNT(*) AS num"+
				" FROM %s.searchable s, %s.source src "+
				" WHERE %s", f.dbSchema, f.dbSchema, sqlWhere)
		*/
		sqlstr = fmt.Sprintf("SELECT COUNT(*) AS num"+
			" FROM searchable s, source src "+
			" WHERE %s", sqlWhere)
		if err := f.db.QueryRow(context.Background(), sqlstr, params...).Scan(&num); err != nil {
			return nil, 0, 0, errors.Wrapf(err, "cannot execute query %s - %v", sqlstr, params)
		}
	}

	if dtr.Start >= num {
		return []map[string]string{}, num, total, nil
	}

	sqlOrder := ""
	if len(orderList) > 0 {
		sqlOrder = fmt.Sprintf("ORDER BY %s", strings.Join(orderList, ", "))
	}

	/*	sqlstr = fmt.Sprintf("SELECT %s "+
		" FROM %s.searchable s, %s.source src "+
		" WHERE %s "+
		" %s "+
		" LIMIT %v OFFSET %v", sqlFields, f.dbSchema, f.dbSchema, sqlWhere, sqlOrder, dtr.Length, dtr.Start)
	*/
	sqlstr = fmt.Sprintf("SELECT %s "+
		" FROM searchable s, source src "+
		" WHERE %s "+
		" %s "+
		" LIMIT %v OFFSET %v", sqlFields, sqlWhere, sqlOrder, dtr.Length, dtr.Start)
	f.log.Info().Msgf("%s - %v", sqlstr, params)
	rows, err := f.db.Query(context.Background(), sqlstr, params...)
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
