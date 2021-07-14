package fair

import (
	"strings"
)

type Partition struct {
	Name               string
	AddrExt            string
	JWTKey             string
	JWTAlg             []string
	OAISignatureDomain string
	OAIDescrption      string
	OAIAdminEmail      []string
	OAIRepositoryName  string
}

func NewPartition(Name, AddrExt, OAISignatureDomain, OAIRepositoryName string, OAIAdminEmail []string, OAIDescrption string, JWTKey string, JWTAlg []string) (*Partition, error) {
	p := &Partition{
		Name:               strings.ToLower(Name),
		AddrExt:            strings.TrimRight(AddrExt, "/"),
		OAISignatureDomain: OAISignatureDomain,
		OAIRepositoryName:  OAIRepositoryName,
		OAIAdminEmail:      OAIAdminEmail,
		OAIDescrption:      OAIDescrption,
		JWTKey:             JWTKey,
		JWTAlg:             JWTAlg,
	}
	return p, nil
}
