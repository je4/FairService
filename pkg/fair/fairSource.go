package fair

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

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

func (f *Fair) GetSourceById(partition *Partition, id int64) (*Source, error) {
	f.sourcesMutex.RLock()
	defer f.sourcesMutex.RUnlock()
	if s, ok := f.sources[id]; ok {
		if s.Partition == partition.Name {
			return s, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("source #%v for partition %s not found", id, partition.Name))
}

func (f *Fair) GetSourceByName(partition *Partition, name string) (*Source, error) {
	f.sourcesMutex.RLock()
	defer f.sourcesMutex.RUnlock()
	for _, src := range f.sources {
		if src.Name == name && src.Partition == partition.Name {
			return src, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("source %s for partition %s not found", name, partition.Name))
}

func (f *Fair) SetSource(src *Source) error {
	sqlstr := fmt.Sprintf("SELECT sourceid FROM %s.source WHERE name=$1", f.dbSchema)
	var sourceId int64
	if err := f.db.QueryRow(sqlstr, src.Name).Scan(&sourceId); err != nil {
		if err != sql.ErrNoRows {
			return errors.Wrapf(err, "cannot query database - %s [%v]", sqlstr, src.ID)
		}
	}
	if sourceId > 0 {
		sqlstr = fmt.Sprintf("UPDATE %s.source "+
			"SET name=$1, detailurl=$2, description=$3, oai_domain=$4, partition=$5 WHERE sourceid=$6 ", f.dbSchema)
		values := []interface{}{src.Name, src.DetailURL, src.Description, src.OAIDomain, src.Partition, sourceId}
		if _, err := f.db.Exec(sqlstr, values...); err != nil {
			return errors.Wrapf(err, "cannot insert into source database - %s [%v]", sqlstr, values)
		}
	} else {
		sqlstr = fmt.Sprintf("INSERT INTO %s.source (name, detailurl, description, oai_domain, partition) "+
			"VALUES($1, $2, $3, $4, $5)", f.dbSchema)
		values := []interface{}{src.Name, src.DetailURL, src.Description, src.OAIDomain, src.Partition}
		if _, err := f.db.Exec(sqlstr, values...); err != nil {
			return errors.Wrapf(err, "cannot update source database - %s [%v]", sqlstr, values)
		}
	}
	if err := f.LoadSources(); err != nil {
		return errors.Wrap(err, "cannot load sources")
	}
	return nil
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
