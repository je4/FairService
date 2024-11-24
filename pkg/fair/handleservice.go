package fair

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
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

func NewHandleService(mr *MultiResolver, config *HandleConfig, logger zLogger.ZLogger) (*HandleService, error) {
	handleClient, err := hcClient.NewHandleCreatorClient(config.ServiceName, config.Addr, string(config.JWTKey), config.JWTAlg, config.SkipCertVerify, logger)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create handle client")
	}
	srv := &HandleService{mr: mr, handleClient: handleClient, config: config, logger: logger}
	mr.AddResolver(srv)
	return srv, nil
}

type HandleService struct {
	handleClient *hcClient.HandleCreatorClient
	logger       zLogger.ZLogger
	config       *HandleConfig
	mr           *MultiResolver
}

func (srv *HandleService) Resolve(pid string) (string, ResolveResultType, error) {
	part := srv.mr.GetPartition()
	fair := part.GetFair()
	db := fair.GetDB()
	prefix, suffix := srv.splitHandle(pid)
	if prefix == "" || suffix == "" {
		return "", ResolveResultTypeUnknown, errors.Wrapf(ErrInvalidIdentifier, "handle %s not valid", pid)
	}
	_pid := fmt.Sprintf("handle:%s/%s", prefix, suffix)
	sqlStr := "SELECT pid.uuid FROM pid WHERE pid.identifier=$1"
	var uuid string
	if err := db.QueryRow(context.Background(), sqlStr, _pid).Scan(&uuid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ResolveResultTypeUnknown, errors.Errorf("handle %s not found", _pid)
		}
		return "", ResolveResultTypeUnknown, errors.Wrapf(err, "cannot execute %s [%s]", sqlStr, _pid)
	}
	item, err := fair.GetItem(part, uuid)
	if err != nil {
		return "", ResolveResultTypeUnknown, errors.Wrapf(err, "cannot get item %s/%s", part.Name, uuid)
	}
	return item.URL, ResolveResultTypeRedirect, nil
}

func (srv *HandleService) CreatePID(fair *Fair, item *ItemData) (string, error) {
	handle, err := srv.mint()
	if err != nil {
		return "", errors.Wrap(err, "cannot mint handle")
	}
	part := srv.mr.GetPartition()
	urlStr := part.RedirURL(item.UUID)
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", errors.Wrap(err, "cannot parse url")
	}
	return handle, errors.Wrapf(srv.handleClient.Create(handle, u), "cannot create handle %s", handle)
}

var handleRegexp = regexp.MustCompile(`(?i)^handle:(?P<prefix>[^/]+)/(?P<suffix>[^?]+)$`)

func (srv *HandleService) splitHandle(doi string) (prefix string, suffix string) {
	match := handleRegexp.FindStringSubmatch(doi)
	if match == nil {
		return "", ""
	}
	for i, name := range doiRegexp.SubexpNames() {
		if i != 0 {
			switch name {
			case "prefix":
				prefix = match[i]
			case "suffix":
				suffix = match[i]
			}
		}
	}
	return
}

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
	return dataciteModel.RelatedIdentifierTypeHandle
}

func (srv *HandleService) mint() (string, error) {
	counter, err := srv.mr.GetPartition().GetFair().NextCounter("handle")
	if err != nil {
		return "", errors.Wrap(err, "cannot mint handle")
	}
	counter2 := bits.RotateLeft64(uint64(counter), -32)
	b := srv.encode(counter2)

	return fmt.Sprintf("handle:%s/%s", srv.config.Prefix, b), nil
}

var _ Resolver = (*HandleService)(nil)
