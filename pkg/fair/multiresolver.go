package fair

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/pkg/errors"
	"strings"
)

func NewResolver(logger zLogger.ZLogger) (*MultiResolver, error) {
	return &MultiResolver{
		resolver: map[dataciteModel.RelatedIdentifierType]Resolver{},
		logger:   logger,
	}, nil
}

type MultiResolver struct {
	db       *pgxpool.Pool
	resolver map[dataciteModel.RelatedIdentifierType]Resolver
	fair     *Fair
	logger   zLogger.ZLogger
}

func (mr *MultiResolver) SetFair(fair *Fair) {
	mr.fair = fair
	mr.db = fair.GetDB()
}

func (mr *MultiResolver) AddResolver(resolver Resolver) {
	mr.resolver[resolver.Type()] = resolver
}

func (mr *MultiResolver) StorePID(uuid string, identifierType dataciteModel.RelatedIdentifierType, identifier string) error {
	sqlstr := `INSERT INTO pid (uuid, identifiertype, identifier) VALUES ($1,$2,$3)`
	if tag, err := mr.db.Exec(context.Background(), sqlstr, uuid, identifierType, identifier); err != nil {
		return errors.Wrapf(err, "cannot insert identifier %s", identifier)
	} else {
		if tag.RowsAffected() != 1 {
			mr.logger.Warn().Msgf("no row affected for identifier %s", identifier)
		}
	}
	return nil
}

func (mr *MultiResolver) CreatePID(uuid string, part *Partition, identifierType dataciteModel.RelatedIdentifierType) (string, error) {
	if _, ok := mr.resolver[identifierType]; !ok {
		return "", errors.Errorf("no resolver for identifier type %s", identifierType)
	}
	item, err := mr.fair.GetItem(part, uuid)
	if err != nil {
		return "", errors.Wrapf(err, "cannot load item %s/%s", part.Name, uuid)
	}
	identifier, err := mr.resolver[identifierType].CreatePID(mr.fair, item)
	if err != nil {
		return "", errors.Wrapf(err, "cannot mint identifier for %s", identifierType)
	}
	if err := mr.StorePID(uuid, identifierType, identifier); err != nil {
		return "", errors.Wrapf(err, "cannot store identifier %s for uuid %s", identifier, uuid)
	}
	return identifier, nil
}

func (mr *MultiResolver) InitPIDTable() error {
	sqlstr := `SELECT uuid, identifier FROM core`
	rows, err := mr.db.Query(context.Background(), sqlstr)
	if err != nil {
		return errors.Wrapf(err, "cannot query core table: %s", sqlstr)
	}
	defer rows.Close()
	for rows.Next() {
		var identifier []string
		var uuid string
		if err := rows.Scan(&uuid, &identifier); err != nil {
			return errors.Wrap(err, "cannot scan identifier")
		}
		for _, id := range identifier {
			id = strings.ToLower(id)
			if strings.HasPrefix(id, "ark:") {
				fmt.Println(dataciteModel.RelatedIdentifierTypeARK, id)
				mr.StorePID(uuid, dataciteModel.RelatedIdentifierTypeARK, id)
			} else if strings.HasPrefix(id, "doi:") {
				fmt.Println(dataciteModel.RelatedIdentifierTypeDOI, id)
				mr.StorePID(uuid, dataciteModel.RelatedIdentifierTypeDOI, id)
			} else if strings.HasPrefix(id, "handle:") {
				fmt.Println(dataciteModel.RelatedIdentifierTypeHandle, id)
				mr.StorePID(uuid, dataciteModel.RelatedIdentifierTypeHandle, id)
			} else {
				fmt.Println("ERROR:", id)
			}
		}
	}
	return nil
}

func (mr *MultiResolver) CreateAll(t dataciteModel.RelatedIdentifierType) error {
	resolver, ok := mr.resolver[t]
	if !ok {
		return errors.Errorf("no resolver for type %s", t)
	}
	var doRefresh = false
	defer func() {
		if doRefresh {
			mr.fair.RefreshSearch()
		}
	}()
	sqlStr := "SELECT coreview.uuid FROM coreview LEFT JOIN ark ON coreview.uuid = ark.uuid WHERE coreview.partition=$1 AND ark.uuid IS NULL LIMIT 1000"
	for _, p := range mr.fair.GetPartitions() {
		for {
			var uuids = make([]string, 0, 1000)
			rows, err := mr.db.Query(context.Background(), sqlStr, p.Name)
			if err != nil {
				return errors.Wrapf(err, "cannot execute %s", sqlStr)
			}
			for rows.Next() {
				var uuid string
				if err := rows.Scan(&uuid); err != nil {
					return errors.Wrapf(err, "cannot scan %s", sqlStr)
				}
				uuids = append(uuids, uuid)
			}
			rows.Close()
			if len(uuids) == 0 {
				break
			}

			for _, uuid := range uuids {
				data, err := mr.fair.GetItem(p, uuid)
				if err != nil {
					return errors.Wrapf(err, "error loading item")
				}
				if data == nil {
					return errors.New(fmt.Sprintf("item %s/%s not found", p.Name, uuid))
				}

				if data.Status != DataStatusActive {
					mr.logger.Error().Msgf("item %s/%s is not active", p.Name, uuid)
				}

				item, err := mr.fair.GetItem(p, uuid)
				if err != nil {
					return errors.Wrapf(err, "cannot load item %s/%s", p.Name, uuid)
				}

				if ark, err := resolver.CreatePID(mr.fair, item); err != nil {
					return errors.Wrapf(err, "cannot create ark for %s", uuid)
				} else {
					doRefresh = true
					mr.logger.Info().Msgf("created ark %s for %s", ark, uuid)

					// todo: write to ark

					ids := make([]string, 0, len(data.Identifier))
					for _, id := range data.Identifier {
						if strings.HasPrefix(id, fmt.Sprintf("ark:%s/%s", p.ARKNAAN, p.ARKShoulder)) {
							continue
						}
						if id == "" {
							continue
						}
						ids = append(ids, id)
					}
					data.Identifier = append(ids, ark)
					sqlstr := "UPDATE core" +
						" SET identifier=$1, datestamp=NOW(), seq=NEXTVAL('lastchange')" +
						" WHERE uuid=$2"
					params := []interface{}{
						pgtype.FlatArray[string](data.Identifier),
						data.UUID,
					}
					if _, err := mr.db.Exec(context.Background(), sqlstr, params...); err != nil {
						return errors.Wrapf(err, "[%s] cannot update [%s] - [%v]", uuid, sqlstr, params)
					}
				}
			}
		}
	}
	return nil

}
