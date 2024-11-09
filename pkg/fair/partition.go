package fair

import (
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

func NewPartition(_fair *Fair, Name, AddrExt, Domain string, oai *OAIConfig, _ark *ARKConfig, _datacite *DataciteConfig, handle *HandleConfig, Description, JWTKey string, JWTAlg []string, logger zLogger.ZLogger) (*Partition, error) {
	p := &Partition{
		Name:        strings.ToLower(Name),
		AddrExt:     strings.TrimRight(AddrExt, "/"),
		Domain:      Domain,
		OAI:         oai,
		Handle:      handle,
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
	ARK          *ARKConfig
	Handle       *HandleConfig
	fair         *Fair
}

func (p *Partition) RedirURL(uuid string) string {
	return p.AddrExt + "/redir/" + uuid
}

func (p *Partition) GetFair() *Fair {
	return p.fair
}
