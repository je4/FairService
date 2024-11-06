package handle

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/pkg/errors"
	"math/bits"
	"regexp"
	"strings"
)

type Config struct {
	Prefix string
}

func NewService(db *pgxpool.Pool, config Config, logger zLogger.ZLogger) (*Service, error) {
	return &Service{db: db, config: config, logger: logger}, nil
}

type Service struct {
	db     *pgxpool.Pool
	logger zLogger.ZLogger
	config Config
}

func (srv *Service) Resolve(pid string) (string, fair.ResolveResultType, error) {
	//TODO implement me
	panic("implement me")
}

func (srv *Service) CreatePID(fair *fair.Fair, item *fair.ItemData) (string, error) {

	return srv.mint(fair, item.UUID)
}

var handleRegexp = regexp.MustCompile(`(?i)^handle:(?P<prefix>[^/]+)/(?P<suffix>[^?]+)$`)

func (srv *Service) ResolveUUID(ark string) (uuid, components, variants string, err error) {
	match := handleRegexp.FindStringSubmatch(ark)
	if match == nil {
		return "", "", variants, errors.Errorf("ark %s not valid", ark)
	}
	result := make(map[string]string)
	for i, name := range handleRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	prefix, _ := result["prefix"]
	suffix, _ := result["suffix"]
	// hyphen is removed
	handle := "handle:" + strings.Join([]string{prefix, suffix}, "/")
	sqlStr := "SELECT ark.uuid FROM ark WHERE ark.ark=$1"
	if err = srv.db.QueryRow(context.Background(), sqlStr, ark).Scan(&uuid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", variants, errors.Errorf("ark %s not found", ark)
		}
		return "", "", variants, errors.Wrapf(err, "cannot execute %s [%s]", sqlStr, ark)
	}
	return

	/*
		var detailurl, signature string
		var url zeronull.Text
		sqlStr := "SELECT source.detailurl, core.signature, core.url FROM ark, source, core WHERE source.name=core.source AND core.uuid=ark.uuid AND ark.ark=$1"
		var detailurl, signature string
		var url zeronull.Text
		if err := srv.db.QueryRow(context.Background(), sqlStr, ark).Scan(&detailurl, &signature, &url); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", "", errors.Errorf("ark %s not found", ark)
			}
			return "", "", errors.Wrapf(err, "cannot execute %s [%s]", sqlStr, ark)
		}
		var resultURL string
		if url != "" {
			resultURL = string(url)
		} else {
			resultURL = strings.ReplaceAll(detailurl, "{signature}", signature)
		}
		partsVariants := strings.Join(parts[2:], "/")
		return resultURL + partsVariants, nil

	*/
}

var chars = []rune("0123456789bcdfghjkmnpqrstvwxz")

var l = uint64(len(chars))

func encode(nb uint64) string {
	var result string
	for {
		result = string(chars[nb%l]) + result
		nb /= l
		if nb == 0 {
			break
		}
	}
	return result
}

func (srv *Service) Type() dataciteModel.RelatedIdentifierType {
	return dataciteModel.RelatedIdentifierTypeARK
}

func (srv *Service) Unify(handle string) (string, error) {
	match := handleRegexp.FindStringSubmatch(handle)
	if match == nil {
		return "", errors.Wrapf(fair.ErrInvalidIdentifier, "handle %s not valid", handle)
	}
	var prefix, suffix string
	for i, name := range handleRegexp.SubexpNames() {
		if i != 0 {
			switch name {
			case "prefix":
				prefix = match[i]
			case "suffix":
				suffix = match[i]
			}
		}
	}
	if prefix == "" || suffix == "" {
		return "", errors.Wrapf(fair.ErrInvalidIdentifier, "handle %s not valid", handle)
	}
	return fmt.Sprintf("handle:%s/%s", prefix, suffix), nil
}

func (srv *Service) mint(fair *fair.Fair, uuid string) (string, error) {
	counter, err := fair.NextCounter("handleseq")
	if err != nil {
		return "", errors.Wrap(err, "cannot mint handle")
	}
	counter2 := bits.RotateLeft64(uint64(counter), -32)
	b := encode(counter2)

	return fmt.Sprintf("handle:%s/%s", srv.config.Prefix, b), nil
}

var _ fair.Resolver = (*Service)(nil)
