package ark

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/pkg/errors"
	"math/bits"
	"regexp"
	"strings"
)

func NewService(db *pgxpool.Pool, logger zLogger.ZLogger) (*Service, error) {
	return &Service{db: db, logger: logger}, nil
}

type Service struct {
	db     *pgxpool.Pool
	logger zLogger.ZLogger
}

var arkRegexp = regexp.MustCompile(`(?i)^ark:(?P<naan>[^/]+)/(?P<qualifier>[^./]+)(/(?P<component>[^.]+))?(\.(?P<variant>[^?]+))?(\?.*)?$`)

func (srv *Service) ResolveUUID(ark string) (uuid, components, variants string, err error) {
	match := arkRegexp.FindStringSubmatch(ark)
	if match == nil {
		return "", "", variants, errors.Errorf("ark %s not valid", ark)
	}
	result := make(map[string]string)
	for i, name := range arkRegexp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	naan, _ := result["naan"]
	qualifier, _ := result["qualifier"]
	components, _ = result["component"]
	variants, _ = result["variant"]
	// hyphen is removed
	ark = "ark:" + strings.ReplaceAll(strings.Join([]string{naan, qualifier}, "/"), "-", "")
	sqlStr := "SELECT ark.uuid FROM ark WHERE ark.ark=$1"
	if err = srv.db.QueryRow(context.Background(), sqlStr, ark).Scan(&uuid); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", variants, errors.Errorf("ark %s not found", ark)
		}
		return "", "", variants, errors.Wrapf(err, "cannot execute %s [%s]", sqlStr, ark)
	}
	return

	/*
		var detailurl, signature string
		var url zeronull.Text
		sqlStr := "SELECT source.detailurl, core.signature, core.url FROM ark, source, core WHERE source.name=core.source AND core.uuid=ark.uuid AND ark.ark=$1"
		var detailurl, signature string
		var url zeronull.Text
		if err := srv.db.QueryRow(context.Background(), sqlStr, ark).Scan(&detailurl, &signature, &url); err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return "", "", errors.Errorf("ark %s not found", ark)
			}
			return "", "", errors.Wrapf(err, "cannot execute %s [%s]", sqlStr, ark)
		}
		var resultURL string
		if url != "" {
			resultURL = string(url)
		} else {
			resultURL = strings.ReplaceAll(detailurl, "{signature}", signature)
		}
		partsVariants := strings.Join(parts[2:], "/")
		return resultURL + partsVariants, nil

	*/
}

var chars = []rune("0123456789bcdfghjkmnpqrstvwxz")

var l = uint64(len(chars))

func encode(nb uint64) string {
	var result string
	for {
		result = string(chars[nb%l]) + result
		nb /= l
		if nb == 0 {
			break
		}
	}
	return result
}

func (srv *Service) Mint(naan, shoulder, prefix string) (string, error) {
	counter, err := srv.nextCounter()
	if err != nil {
		return "", errors.Wrap(err, "cannot mint ark")
	}
	counter2 := bits.RotateLeft64(uint64(counter), -32)
	b := encode(counter2)

	return fmt.Sprintf("ark:%s/%s%s%s", naan, shoulder, prefix, b), nil
}

func (srv *Service) nextCounter() (int64, error) {
	sqlStr := "SELECT NEXTVAL('arkseq')"
	var next int64
	if err := srv.db.QueryRow(context.Background(), sqlStr).Scan(&next); err != nil {
		return 0, errors.Wrapf(err, "cannot execute %s", sqlStr)
	}
	return next, nil
}

func (srv *Service) CreateNew(uuid, naan, shoulder, prefix string) (string, error) {
	ark, err := srv.Mint(naan, shoulder, prefix)
	if err != nil {
		return "", errors.Wrap(err, "cannot mint ark")
	}
	sqlStr := "INSERT INTO ark(uuid, ark) VALUES ($1, $2)"
	if _, err := srv.db.Exec(context.Background(), sqlStr, uuid, ark); err != nil {
		return "", errors.Wrapf(err, "cannot execute %s [%s, %s]", sqlStr, uuid, ark)
	}
	return ark, nil
}

func (srv *Service) Close() error {
	return nil
}
