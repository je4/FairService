package fair

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/pkg/errors"
	"math/bits"
	"regexp"
	"slices"
	"strings"
)

type ARKConfig struct {
	NAAN     string
	Shoulder string
	Prefix   string
}

var arkRegexp = regexp.MustCompile(`(?i)^ark:/(?P<naan>[^/]+)/(?P<qualifier>[^./]+)(/(?P<component>[^.]+))?(\.(?P<variant>[^?]+))?(?P<inflection>\?.*)?$`)

func ArkParts(pid string) (naan, qualifier, component, variant, inflection string, err error) {
	match := arkRegexp.FindStringSubmatch(pid)
	if match == nil {
		return "", "", "", "", "", errors.Errorf("ark %s not valid", pid)
	}
	result := make(map[string]string)
	for i, name := range arkRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	naan, _ = result["naan"]
	qualifier, _ = result["qualifier"]
	component, _ = result["component"]
	variant, _ = result["variant"]
	inflection, _ = result["inflection"]
	return
}

func NewARKService(mr *MultiResolver, config *ARKConfig, logger zLogger.ZLogger) (*ARKService, error) {
	srv := &ARKService{mr: mr, config: config, logger: logger}
	mr.AddResolver(srv)
	return srv, nil
}

type ARKService struct {
	logger zLogger.ZLogger
	config *ARKConfig
	mr     *MultiResolver
}

func (srv *ARKService) Resolve(pid string) (string, ResolveResultType, error) {
	part := srv.mr.GetPartition()
	fair := part.GetFair()
	db := fair.GetDB()
	naan, qualifier, components, variants, inflection, err := ArkParts(pid)

	// hyphen is removed
	_pid := "ark:/" + strings.ReplaceAll(strings.Join([]string{naan, qualifier}, "/"), "-", "")
	sqlStr := "SELECT pid.uuid FROM pid WHERE pid.identifier=$1"
	var uuid string
	if err = db.QueryRow(context.Background(), sqlStr, _pid).Scan(&uuid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ResolveResultTypeUnknown, errors.Errorf("ark %s not found", _pid)
		}
		return "", ResolveResultTypeUnknown, errors.Wrapf(err, "cannot execute %s [%s]", sqlStr, _pid)
	}
	item, err := fair.GetItem(part, uuid)
	if err != nil {
		return "", ResolveResultTypeUnknown, errors.Wrapf(err, "cannot get item %s", uuid)
	}
	if slices.Contains([]string{"?info", "?", "??"}, inflection) {
		data := "erc\n"
		for _, creator := range item.Metadata.Person {
			name := creator.FamilyName
			if name != "" && creator.GivenName != "" {
				name += ", "
			}
			name += creator.GivenName
			data += fmt.Sprintf("who: %s\n", name)
		}
		data += fmt.Sprintf("what: %s\n", item.Metadata.Title)
		if item.Metadata.PublicationYear != "" {
			data += fmt.Sprintf("when: %s\n", item.Metadata.PublicationYear)
		}
		data += fmt.Sprintf("where: %s/resolver/%s\n", part.AddrExt, pid)
		if slices.Contains([]string{"??", "?info"}, inflection) {
			data += "erc-support\n"
			if item.Metadata.Publisher != "" {
				data += fmt.Sprintf("who: %s\n", part.Description)
			}
			data += fmt.Sprintf("where: %s/resolver/%s/\n", part.AddrExt, naan)
		}
		return data, ResolveResultTypeTextPlain, nil
	}
	redirURL := item.URL
	if redirURL == "" {
		src, err := fair.GetSourceByName(part.Name, item.Source)
		if err != nil {
			return "", ResolveResultTypeUnknown, errors.Wrapf(err, "cannot get source %s", item.Source)
		}
		redirURL = strings.ReplaceAll(src.DetailURL, "{signature}", item.Signature)
	}
	if redirURL == "" {
		return "", ResolveResultTypeUnknown, errors.Errorf("no URL found for item %s", uuid)
	}
	if components != "" {
		redirURL += "/" + components
	}
	if variants != "" {
		redirURL += "." + variants
	}
	if inflection != "" {
		redirURL += inflection
	}
	return redirURL, ResolveResultTypeRedirect, nil
}

func (srv *ARKService) CreatePID(fair *Fair, item *ItemData) (string, error) {
	return srv.mint(fair, item.UUID)
}

/*
func (srv *ARKService) ResolveUUID(ark string) (uuid, components, variants string, err error) {
	db := srv.part.GetFair().GetDB()
	var naan, qualifier string
	naan, qualifier, components, variants, err = ArkParts(ark)

	// hyphen is removed
	ark = "ark:" + strings.ReplaceAll(strings.Join([]string{naan, qualifier}, "/"), "-", "")
	sqlStr := "SELECT ark.uuid FROM ark WHERE ark.ark=$1"
	if err = db.QueryRow(context.Background(), sqlStr, ark).Scan(&uuid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", variants, errors.Errorf("ark %s not found", ark)
		}
		return "", "", variants, errors.Wrapf(err, "cannot execute %s [%s]", sqlStr, ark)
	}
	return
}
*/

var arkchars = []rune("0123456789bcdfghjkmnpqrstvwxz")
var arkcharlen = uint64(len(arkchars))

func (srv *ARKService) encode(nb uint64) string {
	var result string
	for {
		result = string(arkchars[nb%arkcharlen]) + result
		nb /= arkcharlen
		if nb == 0 {
			break
		}
	}
	return result
}

func (srv *ARKService) Type() dataciteModel.RelatedIdentifierType {
	return dataciteModel.RelatedIdentifierTypeARK
}

func (srv *ARKService) Unify(ark string) (string, error) {
	match := arkRegexp.FindStringSubmatch(ark)
	if match == nil {
		return "", errors.Wrapf(ErrInvalidIdentifier, "ark %s not valid", ark)
	}
	var naan, qualifier string
	for i, name := range arkRegexp.SubexpNames() {
		if i != 0 {
			switch name {
			case "naan":
				naan = match[i]
			case "qualifier":
				qualifier = match[i]
			}
		}
	}
	if naan == "" || qualifier == "" {
		return "", errors.Wrapf(ErrInvalidIdentifier, "ark %s not valid", ark)
	}

	return fmt.Sprintf("ark:%s/%s", naan, qualifier), nil
}

func (srv *ARKService) mint(fair *Fair, uuid string) (string, error) {
	counter, err := fair.NextCounter("arkseq")
	if err != nil {
		return "", errors.Wrap(err, "cannot mint ark")
	}
	counter2 := bits.RotateLeft64(uint64(counter), -32)
	b := srv.encode(counter2)

	return fmt.Sprintf("ark:/%s/%s%s%s", srv.config.NAAN, srv.config.Shoulder, srv.config.Prefix, b), nil
}

var _ Resolver = (*ARKService)(nil)
