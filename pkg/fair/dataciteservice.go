package fair

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/FairService/v2/pkg/service/datacite"
	"github.com/je4/utils/v2/pkg/zLogger"
	"math/bits"
	"regexp"
)

type DataciteConfig struct {
	Api      string
	User     string
	Password string
	Prefix   string
}

func NewDataciteService(_fair *Fair, config DataciteConfig, logger zLogger.ZLogger) (*DataciteService, error) {
	client, err := datacite.NewClient(config.Api, config.User, config.Password, config.Prefix)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create datacite client")
	}
	if err := client.Heartbeat(); err != nil {
		logger.Error().Err(err).Msg("cannot connect to datacite")
	}
	return &DataciteService{fair: _fair, config: config, client: client, logger: logger}, nil
}

type DataciteService struct {
	fair   *Fair
	logger zLogger.ZLogger
	config DataciteConfig
	client *datacite.Client
}

func (srv *DataciteService) Resolve(pid string) (string, ResolveResultType, error) {
	//TODO implement me
	panic("implement me")
}

func (srv *DataciteService) CreatePID(fair *Fair, item *ItemData) (string, error) {
	part, err := fair.GetPartition(item.Partition)
	if err != nil {
		return "", errors.Wrap(err, "cannot get partition")
	}
	doiStr, err := srv.mint(fair)
	if err != nil {
		return "", errors.Wrap(err, "cannot mint doi")
	}
	if _, err = srv.client.RetrieveDOI(doiStr); err == nil {
		return "", errors.New(fmt.Sprintf("doi %s already exists", doiStr))
	}

	dataciteData := &dataciteModel.DataCite{}
	dataciteData.InitNamespace()
	dataciteData.FromCore(item.Metadata)

	_, suffix := srv.splitDOI(doiStr)

	api, err := srv.client.CreateDOI(dataciteData, suffix, part.RedirURL(item.UUID), datacite.DCEventDraft)
	if err != nil {
		return "", errors.Wrap(err, "cannot create doi")
	}
	_ = api

	return doiStr, nil
}

var doiRegexp = regexp.MustCompile(`(?i)^doi:(?P<prefix>10.[0-9]+)/(?P<suffix>[^;/?:@&=+$,!]+)$`)

func (srv *DataciteService) splitDOI(doi string) (prefix string, suffix string) {
	match := doiRegexp.FindStringSubmatch(doi)
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

var datacitechars = []rune("0123456789bcdfghjkmnpqrstvwxz")

var datacitecharslen = uint64(len(datacitechars))

func (srv *DataciteService) encode(nb uint64) string {
	var result string
	var i = 0
	for {
		if i > 0 && i%4 == 0 {
			result = "-" + result
		}
		result = string(datacitechars[nb%datacitecharslen]) + result
		nb /= datacitecharslen
		if nb == 0 {
			break
		}
		i++
	}
	return result
}

func (srv *DataciteService) Type() dataciteModel.RelatedIdentifierType {
	return dataciteModel.RelatedIdentifierTypeDOI
}

func (srv *DataciteService) Unify(doi string) (string, error) {
	prefix, suffix := srv.splitDOI(doi)
	if prefix == "" || suffix == "" {
		return "", errors.Wrapf(ErrInvalidIdentifier, "doi %s not valid", doi)
	}

	return fmt.Sprintf("doi:%s/%s", prefix, suffix), nil
}

func (srv *DataciteService) mint(fair *Fair) (string, error) {
	counter, err := fair.NextCounter("doiseq")
	if err != nil {
		return "", errors.Wrap(err, "cannot mint ark")
	}
	counter2 := bits.RotateLeft64(uint64(counter), -32)
	suffix := srv.encode(counter2)

	return fmt.Sprintf("doi:%s/%s", srv.config.Prefix, suffix), nil
}

var _ Resolver = (*DataciteService)(nil)
