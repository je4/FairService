package fair

import (
	"fmt"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	hcClient "github.com/je4/HandleCreator/v2/pkg/client"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/pkg/errors"
	"math/bits"
	"net/url"
	"regexp"
)

type HandleConfig struct {
	ServiceName    string
	Addr           string
	JWTKey         string
	JWTAlg         string
	SkipCertVerify bool
	ID             string
	Prefix         string
}

func NewHandleService(_fair *Fair, config *HandleConfig, logger zLogger.ZLogger) (*HandleService, error) {
	handleClient, err := hcClient.NewHandleCreatorClient(config.ServiceName, config.Addr, string(config.JWTKey), config.JWTAlg, config.SkipCertVerify, logger)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create handle client")
	}
	return &HandleService{fair: _fair, handleClient: handleClient, config: config, logger: logger}, nil
}

type HandleService struct {
	handleClient *hcClient.HandleCreatorClient
	logger       zLogger.ZLogger
	config       *HandleConfig
	fair         *Fair
}

func (srv *HandleService) Resolve(pid string) (string, ResolveResultType, error) {
	//TODO implement me
	panic("implement me")
}

func (srv *HandleService) CreatePID(fair *Fair, item *ItemData) (string, error) {
	handle, err := srv.mint(fair, item.UUID)
	if err != nil {
		return "", errors.Wrap(err, "cannot mint handle")
	}
	url, err := url.Parse(fmt.Sprintf("%s/%s", srv.config.ID, item.UUID))
	if err != nil {
		return "", errors.Wrapf(err, "cannot parse url %s/%s", srv.config.ID, item.UUID)
	}
	return handle, errors.Wrapf(srv.handleClient.Create(handle, url), "cannot create handle %s", handle)
}

var handleRegexp = regexp.MustCompile(`(?i)^handle:(?P<prefix>[^/]+)/(?P<suffix>[^?]+)$`)

/*
func (srv *HandleService) ResolveUUID(ark string) (uuid, components, variants string, err error) {
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
}
*/

var handlechars = []rune("0123456789bcdfghjkmnpqrstvwxz")

var handlecharslen = uint64(len(handlechars))

func (srv *HandleService) encode(nb uint64) string {
	var result string
	for {
		result = string(handlechars[nb%handlecharslen]) + result
		nb /= handlecharslen
		if nb == 0 {
			break
		}
	}
	return result
}

func (srv *HandleService) Type() dataciteModel.RelatedIdentifierType {
	return dataciteModel.RelatedIdentifierTypeARK
}

func (srv *HandleService) Unify(handle string) (string, error) {
	match := handleRegexp.FindStringSubmatch(handle)
	if match == nil {
		return "", errors.Wrapf(ErrInvalidIdentifier, "handle %s not valid", handle)
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
		return "", errors.Wrapf(ErrInvalidIdentifier, "handle %s not valid", handle)
	}
	return fmt.Sprintf("handle:%s/%s", prefix, suffix), nil
}

func (srv *HandleService) mint(fair *Fair, uuid string) (string, error) {
	counter, err := fair.NextCounter("handleseq")
	if err != nil {
		return "", errors.Wrap(err, "cannot mint handle")
	}
	counter2 := bits.RotateLeft64(uint64(counter), -32)
	b := srv.encode(counter2)

	return fmt.Sprintf("handle:%s/%s", srv.config.Prefix, b), nil
}

var _ Resolver = (*HandleService)(nil)
