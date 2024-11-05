package fair

import (
	"emperror.dev/errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/je4/FairService/v2/pkg/service/ark"
	"github.com/je4/FairService/v2/pkg/service/datacite"
	"github.com/je4/utils/v2/pkg/zLogger"
	"strings"
	"time"
)

type OAIConfig struct {
	RepositoryName         string
	AdminEmail             []string
	Pagesize               int64
	SampleIdentifier       string
	Delimiter              string
	Scheme                 string
	PageSize               int64
	ResumptionTokenTimeout time.Duration
}

type HandleConfig struct {
	ServiceName    string
	Addr           string
	JWTKey         string
	JWTAlg         string
	SkipCertVerify bool
	ID             string
	Prefix         string
}

func NewPartition(
	db *pgxpool.Pool,
	Name,
	AddrExt,
	Domain string,
	oai *OAIConfig,
	_ark *ark.Config,
	_datacite *datacite.Config,
	handle *HandleConfig,
	Description string,
	JWTKey string,
	JWTAlg []string,
	logger zLogger.ZLogger) (*Partition, error) {
	p := &Partition{
		Name:    strings.ToLower(Name),
		AddrExt: strings.TrimRight(AddrExt, "/"),
		Domain:  Domain,
		OAI:     oai,
		//ARK:     _ark,
		//DOI:         _datacite,
		Handle:      handle,
		Description: Description,
		JWTKey:      JWTKey,
		JWTAlg:      JWTAlg,
	}

	var dataciteClient *datacite.Client
	if _datacite != nil {
		var err error
		dataciteClient, err = datacite.NewClient(
			_datacite.Api,
			_datacite.User,
			_datacite.Password,
			_datacite.Prefix)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot create datacite client for partition %s", p.Name)
		}
		if err := dataciteClient.Heartbeat(); err != nil {
			return nil, errors.Wrapf(err, "cannot check datacite heartbeat for partition %s", p.Name)
		}
		p.datacite = dataciteClient

		var arkClient *ark.Service
		if _ark != nil {
			arkClient, err = ark.NewService(db, _ark, logger)
			if err != nil {
				return nil, errors.Wrapf(err, "cannot create ark client for partition %s", p.Name)
			}
			p.ark = arkClient
		}
	}

	return p, nil
}

type Partition struct {
	Name         string
	AddrExt      string
	Description  string
	JWTKey       string
	JWTAlg       []string
	Domain       string
	HandlePrefix string
	OAI          *OAIConfig
	ARK          *ark.Config
	Handle       *HandleConfig
	datacite     *datacite.Client
	ark          *ark.Service
}

func (p *Partition) RedirURL(uuid string) string {
	return p.AddrExt + "/redir/" + uuid
}
