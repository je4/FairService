package fair

import (
	"emperror.dev/errors"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
)

var ErrInvalidIdentifier = errors.New("invalid identifier")

type ResolveResultType uint32

const (
	ResolveResultTypeUnknown ResolveResultType = iota
	ResolveResultTypeRedirect
	ResolveResultTypeContent
)

type ResolverResolve interface {
	Resolve(pid string) (string, ResolveResultType, error)
}

type Resolver interface {
	ResolverResolve
	Type() dataciteModel.RelatedIdentifierType
	Unify(ark string) (string, error)
	CreatePID(fair *Fair, item *ItemData) (string, error)
}
