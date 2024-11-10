package fair

import (
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/utils/v2/pkg/zLogger"
	"strings"
	"time"
)

type OAIConfig struct {
	RepositoryName         string
	AdminEmail             []string
	SampleIdentifier       string
	Delimiter              string
	Scheme                 string
	PageSize               int64
	ResumptionTokenTimeout time.Duration
}

func NewPartition(_fair *Fair, Name, AddrExt, Domain string, oai *OAIConfig, Description, JWTKey string, JWTAlg []string, logger zLogger.ZLogger) (*Partition, error) {
	p := &Partition{
		Name:        strings.ToLower(Name),
		AddrExt:     strings.TrimRight(AddrExt, "/"),
		Domain:      Domain,
		OAI:         oai,
		Description: Description,
		JWTKey:      JWTKey,
		JWTAlg:      JWTAlg,
		fair:        _fair,
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
	fair         *Fair
	mr           *MultiResolver
}

func (p *Partition) Resolve(pid string) (string, ResolveResultType, error) {
	return p.mr.Resolve(pid)
}

func (p *Partition) CreatePID(uuid string, identifierType dataciteModel.RelatedIdentifierType) (string, error) {
	return p.mr.CreatePID(uuid, p, identifierType)
}

func (p *Partition) RedirURL(uuid string) string {
	return p.AddrExt + "/redir/" + uuid
}

func (p *Partition) GetFair() *Fair {
	return p.fair
}

func (p *Partition) AddResolver(mr *MultiResolver) {
	p.mr = mr
}

func (p *Partition) GetMultiResolver() *MultiResolver {
	return p.mr

}
