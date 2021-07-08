package fair

import (
	"strings"
)

type Partition struct {
	Name    string
	AddrExt string
	JWTKey  string
	JWTAlg  []string
}

func NewPartition(Name, AddrExt, JWTKey string, JWTAlg []string) (*Partition, error) {
	p := &Partition{
		Name:    strings.ToLower(Name),
		AddrExt: strings.TrimRight(AddrExt, "/"),
		JWTKey:  JWTKey,
		JWTAlg:  JWTAlg,
	}
	return p, nil
}
