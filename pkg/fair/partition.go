package fair

import (
	"strings"
	"time"
)

type Partition struct {
	Name                    string
	AddrExt                 string
	Description             string
	JWTKey                  string
	JWTAlg                  []string
	OAIAdminEmail           []string
	OAIRepositoryName       string
	OAIPagesize             int64
	OAIRepositoryIdentifier string
	OAISampleIdentifier     string
	OAIDelimiter            string
	OAIScheme               string

	ResumptionTokenTimeout time.Duration
}

func NewPartition(
	Name,
	AddrExt,
	OAISignatureDomain,
	OAIRepositoryName string,
	OAIAdminEmail []string,
	OAIRepositoryIdentifier,
	OAISampleIdentifier,
	OAIDelimiter string,
	OAIScheme string,
	Description string,
	pagesize int64,
	resumptionTokenTimeout time.Duration,
	JWTKey string,
	JWTAlg []string) (*Partition, error) {
	p := &Partition{
		Name:                    strings.ToLower(Name),
		AddrExt:                 strings.TrimRight(AddrExt, "/"),
		OAIRepositoryName:       OAIRepositoryName,
		OAIAdminEmail:           OAIAdminEmail,
		Description:             Description,
		OAIPagesize:             pagesize,
		OAIRepositoryIdentifier: OAIRepositoryIdentifier,
		OAISampleIdentifier:     OAISampleIdentifier,
		OAIDelimiter:            OAIDelimiter,
		OAIScheme:               OAIScheme,
		ResumptionTokenTimeout:  resumptionTokenTimeout,
		JWTKey:                  JWTKey,
		JWTAlg:                  JWTAlg,
	}
	return p, nil
}
