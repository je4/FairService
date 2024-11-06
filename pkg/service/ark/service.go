package ark

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/pkg/errors"
	"math/bits"
	"regexp"
	"strings"
)

type Config struct {
	NAAN     string
	Shoulder string
	Prefix   string
}

func NewService(_fair *fair.Fair, config *Config, logger zLogger.ZLogger) (*Service, error) {
	return &Service{fair: _fair, config: config, logger: logger}, nil
}

type Service struct {
	fair   *fair.Fair
	logger zLogger.ZLogger
	config *Config
}

func (srv *Service) Resolve(pid string) (string, fair.ResolveResultType, error) {
	//TODO implement me
	panic("implement me")
}

func (srv *Service) CreatePID(fair *fair.Fair, item *fair.ItemData) (string, error) {
	return srv.mint(fair, item.UUID)
}

var arkRegexp = regexp.MustCompile(`(?i)^ark:(?P<naan>[^/]+)/(?P<qualifier>[^./]+)(/(?P<component>[^.]+))?(\.(?P<variant>[^?]+))?(\?.*)?$`)

func (srv *Service) ResolveUUID(ark string) (uuid, components, variants string, err error) {
	db := srv.fair.GetDB()
	match := arkRegexp.FindStringSubmatch(ark)
	if match == nil {
		return "", "", variants, errors.Errorf("ark %s not valid", ark)
	}
	result := make(map[string]string)
	for i, name := range arkRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	naan, _ := result["naan"]
	qualifier, _ := result["qualifier"]
	components, _ = result["component"]
	variants, _ = result["variant"]
	// hyphen is removed
	ark = "ark:" + strings.ReplaceAll(strings.Join([]string{naan, qualifier}, "/"), "-", "")
	sqlStr := "SELECT ark.uuid FROM ark WHERE ark.ark=$1"
	if err = db.QueryRow(context.Background(), sqlStr, ark).Scan(&uuid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", variants, errors.Errorf("ark %s not found", ark)
		}
		return "", "", variants, errors.Wrapf(err, "cannot execute %s [%s]", sqlStr, ark)
	}
	return
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

func (srv *Service) Unify(ark string) (string, error) {
	match := arkRegexp.FindStringSubmatch(ark)
	if match == nil {
		return "", errors.Wrapf(fair.ErrInvalidIdentifier, "ark %s not valid", ark)
	}
	var naan, qualifier string
	for i, name := range arkRegexp.SubexpNames() {
		if i != 0 {
			switch name {
			case "naan":
				naan = match[i]
			case "qualifier":
				qualifier = match[i]
			}
		}
	}
	if naan == "" || qualifier == "" {
		return "", errors.Wrapf(fair.ErrInvalidIdentifier, "ark %s not valid", ark)
	}

	return fmt.Sprintf("ark:%s/%s", naan, qualifier), nil
}

func (srv *Service) mint(fair *fair.Fair, uuid string) (string, error) {
	counter, err := fair.NextCounter("arkseq")
	if err != nil {
		return "", errors.Wrap(err, "cannot mint ark")
	}
	counter2 := bits.RotateLeft64(uint64(counter), -32)
	b := encode(counter2)

	return fmt.Sprintf("ark:%s/%s%s%s", srv.config.NAAN, srv.config.Shoulder, srv.config.Prefix, b), nil
}

var _ fair.Resolver = (*Service)(nil)
