package fair

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgtype/zeronull"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/FairService/v2/pkg/model/myfair"
	"github.com/pkg/errors"
	"reflect"
	"sort"
	"strings"
	"time"
)

func itemDataFromRow(row interface{}, lastCols ...interface{}) (*ItemData, error) {
	var uuidStr string
	var metaStr string
	var accessStr string
	var signature string
	var sourceName string
	var partition string
	var statusStr string
	var seq int64
	var datestamp time.Time
	var setArr, catalogArr, identifierArr pgtype.FlatArray[string]
	var url zeronull.Text
	cols := []interface{}{}
	cols = append(cols,
		&uuidStr,
		&metaStr,
		&setArr,
		&catalogArr,
		&accessStr,
		&signature,
		&sourceName,
		&partition,
		&statusStr,
		&seq,
		&datestamp,
		&identifierArr,
		&url,
	)
	cols = append(cols, lastCols...)
	pgxrow, ok := row.(pgx.Row)
	if !ok {
		return nil, errors.New(fmt.Sprintf("invalid type %T for data row", row))
	}
	if err := pgxrow.Scan(cols...); err != nil {
		return nil, errors.Wrapf(err, "cannot scan result")
	}
	/*	switch r := row.(type) {
		case *sql.Row:
			if err := r.Scan(cols...); err != nil {
				return nil, errors.Wrapf(err, "cannot scan result")
			}
		case *sql.Rows:
			if err := r.Scan(cols...); err != nil {
				return nil, errors.Wrapf(err, "cannot scan result")
			}
		default:
			return nil, errors.New(fmt.Sprintf("invalid type %T for data row", r))
		}
	*/
	data := &ItemData{
		UUID:       uuidStr,
		Source:     sourceName,
		Partition:  partition,
		Signature:  signature,
		Metadata:   myfair.Core{},
		Set:        setArr,
		Catalog:    catalogArr,
		Seq:        seq,
		Datestamp:  datestamp,
		Identifier: identifierArr,
		URL:        string(url),
	}
	data.Access, ok = DataAccessReverse[accessStr]
	if !ok {
		return nil, errors.New(fmt.Sprintf("[%s] invalid access type %s", uuidStr, accessStr))
	}
	data.Status, ok = DataStatusReverse[statusStr]
	if !ok {
		return nil, errors.New(fmt.Sprintf("[%s] invalid status type %s", uuidStr, statusStr))
	}
	if err := json.Unmarshal([]byte(metaStr), &data.Metadata); err != nil {
		return nil, errors.Wrapf(err, "[%s] cannot unmarshal core [%s]", uuidStr, metaStr)
	}
	// add local identifiers
	for _, id := range data.Identifier {
		strs := strings.SplitN(id, ":", 2)
		if len(strs) != 2 {
			continue
		}

		idType, ok := myfair.RelatedIdentifierTypeReverse[strings.ToLower(strs[0])]
		if !ok {
			//f.log.Warningf("[%s] unknown identifier type %s", uuidStr, id)
			continue
		}
		idStr := strings.TrimPrefix(strs[1], "/")
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

func (f *Fair) getItems(sqlWhere string, params []interface{}, limit, offset int64, completeListSize *int64, fn func(item *ItemData) error) error {
	if completeListSize != nil {
		/*		sqlstr := fmt.Sprintf("SELECT COUNT(*) AS num"+
				" FROM %s.coreview"+
				" WHERE %s", f.dbSchema, sqlWhere)
		*/
		sqlstr := fmt.Sprintf("SELECT COUNT(*) AS num"+
			" FROM coreview_new"+
			" WHERE %s ", sqlWhere)
		if err := f.db.QueryRow(context.Background(), sqlstr, params...).Scan(completeListSize); err != nil {
			return errors.Wrapf(err, "cannot get number of result items [%s] - [%v]", sqlstr, params)
		}
	}
	/*	sqlstr := fmt.Sprintf("SELECT uuid, metadata, setspec, catalog, access, signature, sourcename, status, seq, datestamp, identifier"+
		" FROM %s.coreview"+
		" WHERE %s"+
		" ORDER BY seq ASC", f.dbSchema, sqlWhere)
	*/
	sqlstr := fmt.Sprintf("SELECT uuid, metadata, setspec, catalog, access, signature, sourcename, partition, status, seq, datestamp, identifier, url"+
		" FROM coreview_new"+
		" WHERE %s"+
		" ORDER BY seq ASC", sqlWhere)
	if limit > 0 {
		sqlstr += fmt.Sprintf(" LIMIT %v", limit)
	}
	if offset > 0 {
		sqlstr += fmt.Sprintf(" OFFSET %v", offset)
	}
	rows, err := f.db.Query(context.Background(), sqlstr, params...)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
	}
	defer rows.Close()
	for rows.Next() {
		data, err := itemDataFromRow(rows)
		if err != nil {
			return errors.Wrap(err, "cannot get data from result row")
		}
		if err := fn(data); err != nil {
			return errors.Wrap(err, "error cal	ling fn")
		}
	}
	return nil

}

func (f *Fair) GetItemsDatestamp(partition *Partition, datestamp, until time.Time, set string, access []DataAccess, limit, offset int64, completeListSize *int64, fn func(item *ItemData) error) error {
	sqlWhere := "partition=$1 AND datestamp>=$2"
	params := []interface{}{partition.Name, datestamp}
	key := 3
	if set != "" {
		sqlWhere += fmt.Sprintf(" AND $%d = ANY(setspec)", key)
		params = append(params, set)
		key++
	}
	if len(access) > 0 {
		var accessList []string
		for _, acc := range access {
			accessList = append(accessList, fmt.Sprintf("access=$%v", key))
			params = append(params, acc)
			key++
		}
		sqlWhere += fmt.Sprintf(" AND (%s)", strings.Join(accessList, " OR "))
	}
	if !until.Equal(time.Time{}) {
		params = append(params, until)
		sqlWhere += fmt.Sprintf(" AND datestamp<=$%v", key)
		key++
	}
	return f.getItems(sqlWhere, params, limit, offset, completeListSize, fn)
}

func (f *Fair) GetItemsSeq(partition *Partition, seq int64, until time.Time, access []DataAccess, limit, offset int64, completeListSize *int64, fn func(item *ItemData) error) error {
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

func (f *Fair) GetItem(partition *Partition, uuidStr string) (*ItemData, error) {

	var item *ItemData

	sqlWhere := "uuid=$1"
	params := []interface{}{uuidStr}
	if partition != nil {
		sqlWhere += " AND partition=$2"
		params = append(params, partition.Name)
	}
	var completeListSize int64
	if err := f.getItems(sqlWhere, params, 1, 0, &completeListSize, func(found *ItemData) error {
		item = found
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "error querying item cannot query item %v", uuidStr)
	}

	return item, nil
}

func (f *Fair) GetItemSource(partition *Partition, sourceid int64, signature string) (*ItemData, error) {
	var item *ItemData

	sqlWhere := "partition=$1 AND source=$2 AND signature=$3"
	params := []interface{}{partition.Name, sourceid, signature}
	var completeListSize int64
	if err := f.getItems(sqlWhere, params, 1, 0, &completeListSize, func(found *ItemData) error {
		item = found
		return nil
	}); err != nil {
		return nil, errors.Wrapf(err, "error querying item cannot query item %v.%s", sourceid, signature)
	}

	return item, nil
}

func (f *Fair) DeleteItem(partition *Partition, uuidStr string) error {
	data, err := f.GetItem(partition, uuidStr)
	if err != nil {
		return errors.Wrapf(err, "cannot get item %s/%s", partition.Name, uuidStr)
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

func (f *Fair) CreateItem(partition *Partition, data *ItemData) (*ItemData, error) {
	sort.Strings(data.Catalog)
	sort.Strings(data.Set)

	src, err := f.GetSourceByName(partition.Name, data.Source)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get source %s", data.Source)
	}

	item, err := f.GetItemSource(partition, src.ID, data.Signature)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get item #%v.%v from partition %s", src.ID, data.Signature, partition.Name)
	}

	// not found
	if item == nil {
		//
		// Create new Entry
		//

		item = data
		//		item.Identifier = []string{}

		uuidVal, err := uuid.NewUUID()
		if err != nil {
			return nil, errors.Wrap(err, "cannot generate uuid")
		}
		item.UUID = uuidVal.String()

		/*
			if f.handle != nil {
				//
				// Create Handle and add to local identifier list
				//
				next, err := f.NextCounter("handleseq")
				if err != nil {
					return nil, errors.Wrap(err, "cannot get next handle value")
				}
				newHandle := fmt.Sprintf("%s/%s/%v", partition.HandleID, partition.HandlePrefix, next)
				newURL, err := url.Parse(fmt.Sprintf("%s/redir/%s", partition.AddrExt, item.UUID))
				if err != nil {
					return nil, errors.Wrapf(err, "cannot parse url %s", fmt.Sprintf("%s/redir/%s", partition.AddrExt, item.UUID))
				}
				if err := f.handle.Create(newHandle, newURL); err != nil {
					return nil, errors.Wrapf(err, "cannot create handle %s for %s", newHandle, newURL.String())
				}
				sqlHandle.String = newHandle
				sqlHandle.Valid = true
				item.Metadata.Identifier = append(item.Metadata.Identifier, myfair.Identifier{
					Value:          newHandle,
					IdentifierType: myfair.RelatedIdentifierTypeHandle,
				})
				item.Identifier = append(item.Identifier, fmt.Sprintf("%s:%s", myfair.RelatedIdentifierTypeHandle, newHandle))
			}
		*/
		sort.Strings(item.Identifier)

		coreBytes, err := json.Marshal(item.Metadata)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot marshal core data [%v]", data.Metadata)
		}
		/*		sqlstr := fmt.Sprintf("INSERT INTO %s.core"+
				" (uuid, datestamp, setspec, metadata, signature, source, access, catalog, seq, status, identifier)"+
				" VALUES($1, NOW(), $2, $3, $4, $5, $6, $7, NEXTVAL('lastchange'), $8, $9)", f.dbSchema)
		*/
		sqlstr := `INSERT INTO core 
    				(uuid, datestamp, setspec, metadata, signature, identifier, source, access, catalog, seq, status) 
					VALUES($1, NOW(), $2, $3, $4, $5, $6, $7, $8, NEXTVAL('lastchange'), $9)`
		params := []interface{}{
			item.UUID, // uuid
			// datestamp
			pgtype.FlatArray[string](item.Set),        // setspec
			string(coreBytes),                         // metadata
			item.Signature,                            // signature
			pgtype.FlatArray[string](item.Identifier), // identifier
			src.ID,                                 // source
			item.Access,                            // access
			pgtype.FlatArray[string](item.Catalog), // catalog
			// seq
			DataStatusActive,
		}
		ret, err := f.db.Exec(context.Background(), sqlstr, params...)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		if num := ret.RowsAffected(); num == 0 {
			return nil, errors.Wrap(err, "no affected rows")
		}
		/*		sqlstr = fmt.Sprintf("INSERT INTO %s.core_dirty"+
				" (uuid)"+
				" VALUES($1)", f.dbSchema)
		*/
		ark, err := partition.CreatePID(item.UUID, dataciteModel.RelatedIdentifierTypeHandle)
		if err != nil {
			f.log.Error().Err(err).Msgf("cannot create handle for %s", item.UUID)
		}
		handle, err := partition.CreatePID(item.UUID, dataciteModel.RelatedIdentifierTypeARK)
		if err != nil {
			f.log.Error().Err(err).Msgf("cannot create ark for %s", item.UUID)
		}
		item.Identifier = append(item.Identifier, ark, handle)

		sqlstr = "INSERT INTO core_dirty (uuid) VALUES($1)"
		params = []interface{}{
			item.UUID, // uuid
		}

		if _, err = f.db.Exec(context.Background(), sqlstr, params...); err != nil {
			return nil, errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		f.log.Info().Msgf("new item [%s] inserted", item.UUID)

		return item, nil

	} else {
		sort.Strings(data.Catalog)
		sort.Strings(data.Set)

		// add local identifiers
		for _, id := range item.Identifier {
			strs := strings.SplitN(id, ":", 2)
			if len(strs) != 2 {
				return nil, errors.New(fmt.Sprintf("[%s] invalid identifier format %s", item.UUID, id))
			}

			idType, ok := myfair.RelatedIdentifierTypeReverse[strings.ToLower(strs[0])]
			if !ok {
				f.log.Warn().Msgf("[%s] unknown identifier type %s", item.UUID, id)
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

		// do update here
		if item.Status == DataStatusActive &&
			reflect.DeepEqual(item.Metadata, data.Metadata) &&
			equalStrings(item.Set, data.Set) &&
			equalStrings(item.Catalog, data.Catalog) &&
			equalStrings(item.Identifier, data.Identifier) &&
			item.Access == data.Access {
			f.log.Info().Msgf("no update needed for item [%v]", item.UUID)
			/*			sqlstr := fmt.Sprintf("INSERT INTO %s.core_dirty"+
						" (uuid)"+
						" VALUES($1)", f.dbSchema)
			*/
			sqlstr := "INSERT INTO core_dirty (uuid) VALUES($1)"
			params := []interface{}{
				item.UUID, // uuid
			}
			if _, err = f.db.Exec(context.Background(), sqlstr, params...); err != nil {
				return nil, errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
			}
			return item, nil
		}

		dataMetaBytes, err := json.Marshal(data.Metadata)
		if err != nil {
			return nil, errors.Wrapf(err, "[%s] cannot unmarshal data core", item.UUID)
		}

		identifiers := []string{}
		identifiers = append(identifiers, item.Identifier...)
		identifiers = append(identifiers, data.Identifier...)
		sort.Strings(identifiers)
		identifiers = UniqString(identifiers)

		sqlstr := fmt.Sprintf("UPDATE %s.core"+
			" SET setspec=$1, identifier=$7, metadata=$2, access=$3, catalog=$4, status=$5, datestamp=NOW(), seq=NEXTVAL('lastchange')"+
			" WHERE uuid=$6", f.dbSchema)
		params := []any{
			pgtype.FlatArray[string](data.Set),
			string(dataMetaBytes),
			data.Access,
			pgtype.FlatArray[string](data.Catalog),
			DataStatusActive,
			item.UUID,
			pgtype.FlatArray[string](identifiers),
		}
		if _, err := f.db.Exec(context.Background(), sqlstr, params...); err != nil {
			return nil, errors.Wrapf(err, "[%s] cannot update [%s] - [%v]", item.UUID, sqlstr, params)
		}
		if len(identifiers) > 0 {
			sqlstr = `WITH DATA(uuid, identifier, identifierType)  AS (
   VALUES
`
			params = []any{}
			counter := 1
			for _, id := range identifiers {
				parts := strings.SplitN(id, ":", 2)
				if len(parts) != 2 {
					f.log.Warn().Msgf("[%s] invalid identifier format %s", item.UUID, id)
					continue
				}
				idType := myfair.RelatedIdentifierTypeReverse[strings.ToLower(parts[0])]
				sqlstr += fmt.Sprintf("	($%d, $%d, $%d ::\"IdentifierType\"),\n", counter, counter+1, counter+2)
				params = append(params, item.UUID, id, idType)
				counter += 3
			}
			sqlstr = strings.TrimSuffix(sqlstr, ",\n") + `
) 
insert into pid (uuid, identifier, identifierType) 
select d.uuid, d.identifier, d.identifierType
from data d
where not exists (select 1
                  from pid p2
                  where p2.uuid = d.uuid and p2.identifier = d.identifier and p2.identifierType = d.identifierType);
`
			if _, err := f.db.Exec(context.Background(), sqlstr, params...); err != nil {
				return nil, errors.Wrapf(err, "[%s] cannot update pid [%s] - [%v]", item.UUID, sqlstr, params)
			}
		}
		/*		sqlstr = fmt.Sprintf("INSERT INTO %s.core_dirty"+
				" (uuid)"+
				" VALUES($1)", f.dbSchema)
		*/
		sqlstr = "INSERT INTO core_dirty (uuid) VALUES($1)"
		params = []interface{}{
			item.UUID, // uuid
		}
		if _, err = f.db.Exec(context.Background(), sqlstr, params...); err != nil {
			return nil, errors.Wrapf(err, "cannot execute query [%s] - [%v]", sqlstr, params)
		}
		f.log.Info().Msgf("item [%s] updated", item.UUID)
		item.Set = data.Set
		item.Metadata = data.Metadata
		item.Access = data.Access
		item.Catalog = data.Catalog
		item.Status = DataStatusActive
		item.Metadata = data.Metadata
		return item, nil
	}

}
