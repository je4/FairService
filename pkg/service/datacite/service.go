package datacite

import (
	"emperror.dev/errors"
	"fmt"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/utils/v2/pkg/zLogger"
	"math/bits"
	"regexp"
)

type Config struct {
	Api      string
	User     string
	Password string
	Prefix   string
}

func NewService(config Config, logger zLogger.ZLogger) (*Service, error) {
	client, err := NewClient(config.Api, config.User, config.Password, config.Prefix)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create datacite client")
	}
	if err := client.Heartbeat(); err != nil {
		logger.Error().Err(err).Msg("cannot connect to datacite")
	}
	return &Service{config: config, client: client, logger: logger}, nil
}

type PartitionParam struct {
	prefix string
}

type Service struct {
	logger zLogger.ZLogger
	config Config
	client *Client
}

func (srv *Service) CreatePID(fair *fair.Fair, item *fair.ItemData) (string, error) {
	part, err := fair.GetPartition(item.Partition)
	if err != nil {
		return "", errors.Wrap(err, "cannot get partition")
	}
	doiStr, err := srv.mint(fair)
	if err != nil {
		return "", errors.Wrap(err, "cannot mint doi")
	}
	_, err = srv.client.RetrieveDOI(doiStr)
	if err == nil {
		return "", errors.New(fmt.Sprintf("doi %s already exists", doiStr))
	}

	dataciteData := &dataciteModel.DataCite{}
	dataciteData.InitNamespace()
	dataciteData.FromCore(item.Metadata)

	_, suffix := splitDOI(doiStr)

	api, err := srv.client.CreateDOI(dataciteData, suffix, part.RedirURL(item.UUID), DCEventDraft)
	if err != nil {
		return "", errors.Wrap(err, "cannot create doi")
	}
	_ = api

	return doiStr, nil
}

var doiRegexp = regexp.MustCompile(`(?i)^doi:(?P<prefix>10.[0-9]+)/(?P<suffix>[^;/?:@&=+$,!]+)$`)

func splitDOI(doi string) (prefix string, suffix string) {
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

var chars = []rune("0123456789bcdfghjkmnpqrstvwxz")

var l = uint64(len(chars))

func encode(nb uint64) string {
	var result string
	var i = 0
	for {
		if i > 0 && i%4 == 0 {
			result = "-" + result
		}
		result = string(chars[nb%l]) + result
		nb /= l
		if nb == 0 {
			break
		}
		i++
	}
	return result
}

func (srv *Service) Type() dataciteModel.RelatedIdentifierType {
	return dataciteModel.RelatedIdentifierTypeDOI
}

func (srv *Service) Unify(doi string) (string, error) {
	prefix, suffix := splitDOI(doi)
	if prefix == "" || suffix == "" {
		return "", errors.Wrapf(fair.ErrInvalidIdentifier, "doi %s not valid", doi)
	}

	return fmt.Sprintf("doi:%s/%s", prefix, suffix), nil
}

func (srv *Service) mint(fair *fair.Fair) (string, error) {
	counter, err := fair.NextCounter("doiseq")
	if err != nil {
		return "", errors.Wrap(err, "cannot mint ark")
	}
	counter2 := bits.RotateLeft64(uint64(counter), -32)
	suffix := encode(counter2)

	return fmt.Sprintf("doi:%s/%s", srv.config.Prefix, suffix), nil
}

var _ fair.Resolver = (*Service)(nil)
