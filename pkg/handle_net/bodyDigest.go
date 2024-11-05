package handle_net

import "github.com/pkg/errors"

type BodyDigestType uint8

const (
	BodyDigestTypeNone BodyDigestType = 0
	BodyDigestTypeMD5  BodyDigestType = 1
	BodyDigestTypeSHA1 BodyDigestType = 2
)

type BodyDigest struct {
	Type   BodyDigestType
	Digest []byte
}

func (d *BodyDigest) Size() uint32 {
	switch d.Type {
	case BodyDigestTypeMD5:
		return 17
	case BodyDigestTypeSHA1:
		return 21
	default:
		return 0
	}
}

func (d *BodyDigest) MarshalHandleBinary() ([]byte, error) {
	switch d.Type {
	case BodyDigestTypeMD5:
		return append([]byte{byte(d.Type)}, d.Digest...), nil
	case BodyDigestTypeSHA1:
		return append([]byte{byte(d.Type)}, d.Digest...), nil
	}
	return nil, nil
}

func (d *BodyDigest) UnmarshalHandleBinary(data []byte) error {
	d.Type = BodyDigestType(data[0])
	switch d.Type {
	case BodyDigestTypeMD5:
		if len(data) < 17 {
			return errors.Errorf("invalid data length %d for MD5", len(data))
		}
		d.Digest = data[1:17]
	case BodyDigestTypeSHA1:
		if len(data) < 21 {
			return errors.Errorf("invalid data length %d for SHA-1", len(data))
		}
		d.Digest = data[1:21]
	default:
		return errors.Errorf("invalid digest type %d", d.Type)
	}
	return nil
}
