package fair

import (
	"strings"
	"time"
)

type Partition struct {
	Name                string
	AddrExt             string
	Description         string
	JWTKey              string
	JWTAlg              []string
	Domain              string
	HandlePrefix        string
	OAIAdminEmail       []string
	OAIRepositoryName   string
	OAIPagesize         int64
	OAISampleIdentifier string
	OAIDelimiter        string
	OAIScheme           string
	ARKNAAN             string
	ARKShoulder         string
	ARKPrefix           string

	ResumptionTokenTimeout time.Duration
	HandleID               string
}

func NewPartition(
	Name,
	AddrExt,
	Domain,
	HandlePrefix,
	OAIRepositoryName string,
	OAIAdminEmail []string,
	OAISampleIdentifier,
	OAIDelimiter string,
	OAIScheme string,
	ARKNAAN string,
	ARKShoulder string,
	ARKPrefix string,
	HandleID string,
	Description string,
	pagesize int64,
	resumptionTokenTimeout time.Duration,
	JWTKey string,
	JWTAlg []string) (*Partition, error) {
	p := &Partition{
		Name:                   strings.ToLower(Name),
		AddrExt:                strings.TrimRight(AddrExt, "/"),
		Domain:                 Domain,
		HandlePrefix:           HandlePrefix,
		OAIRepositoryName:      OAIRepositoryName,
		OAIAdminEmail:          OAIAdminEmail,
		Description:            Description,
		OAIPagesize:            pagesize,
		OAISampleIdentifier:    OAISampleIdentifier,
		OAIDelimiter:           OAIDelimiter,
		OAIScheme:              OAIScheme,
		ARKNAAN:                ARKNAAN,
		ARKShoulder:            ARKShoulder,
		ARKPrefix:              ARKPrefix,
		HandleID:               HandleID,
		ResumptionTokenTimeout: resumptionTokenTimeout,
		JWTKey:                 JWTKey,
		JWTAlg:                 JWTAlg,
	}
	return p, nil
}
