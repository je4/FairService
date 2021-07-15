package fair

import (
	"strings"
	"time"
)

type Partition struct {
	Name                   string
	AddrExt                string
	JWTKey                 string
	JWTAlg                 []string
	OAISignatureDomain     string
	OAIDescrption          string
	OAIAdminEmail          []string
	OAIRepositoryName      string
	OAIPagesize            int64
	ResumptionTokenTimeout time.Duration
}

func NewPartition(Name, AddrExt, OAISignatureDomain, OAIRepositoryName string, OAIAdminEmail []string, OAIDescrption string, pagesize int64, resumptionTokenTimeout time.Duration, JWTKey string, JWTAlg []string) (*Partition, error) {
	p := &Partition{
		Name:                   strings.ToLower(Name),
		AddrExt:                strings.TrimRight(AddrExt, "/"),
		OAISignatureDomain:     OAISignatureDomain,
		OAIRepositoryName:      OAIRepositoryName,
		OAIAdminEmail:          OAIAdminEmail,
		OAIDescrption:          OAIDescrption,
		OAIPagesize:            pagesize,
		ResumptionTokenTimeout: resumptionTokenTimeout,
		JWTKey:                 JWTKey,
		JWTAlg:                 JWTAlg,
	}
	return p, nil
}
