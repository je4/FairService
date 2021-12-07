package fair

import (
	"fmt"
	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (f *Fair) AddArchive(part *Partition, name, description string) error {
	sqlstr := fmt.Sprintf("INSERT INTO %s.archive(partition, name, description) "+
		" VALUES($1, $2, $3) "+
		" ON CONFLICT(partition, name) "+
		" DO UPDATE SET description=$3 ", f.dbSchema)
	result, err := f.db.Exec(sqlstr, part.Name, name, description)
	if err != nil {
		return errors.Wrapf(err, "cannot query database - %s [%v, %v]", sqlstr, name, description)
	}
	if rows, err := result.RowsAffected(); err != nil {
		return errors.Wrap(err, "cannot get affected rows")
	} else {
		if rows != 1 {
			return errors.New(fmt.Sprintf("invalid number of affected rows: %v", rows))
		}
	}
	return nil
}

func (f *Fair) AddArchiveItem(part *Partition, archive string, item *ItemData) error {
	sqlstr := fmt.Sprintf("INSERT INTO %s.core_archive(partition, archive_name, core_uuid, lastchange, files) "+
		" VALUES($1, $2, $3, $5) "+
		" ON CONFLICT(partition, archive_name, core_uuid) "+
		" DO NOTHING",
		// "     UPDATE SET lastchange=$4, files=$5",
		f.dbSchema)

	var medias = []string{}
	/*
		for _, media := range item.Metadata.Media {
			medias = append(medias, media.Uri)
		}
	*/
	var params = []interface{}{part.Name, archive, item.UUID, item.Datestamp, pq.Array(medias)}
	_, err := f.db.Exec(sqlstr, params...)
	if err != nil {
		return errors.Wrapf(err, "cannot exec insert query - %s %v", sqlstr, params)
	}
	return nil
}
