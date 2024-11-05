package handle_net

import (
	"emperror.dev/errors"
	"encoding/binary"
)

type TypeList []UTF8String

func (t *TypeList) Size() uint32 {
	var size uint32 = 4
	for _, s := range *t {
		size += s.Size()
	}
	return size
}

func (t *TypeList) MarshalHandleBinary() ([]byte, error) {
	var result = make([]byte, t.Size())
	binary.BigEndian.PutUint32(result, uint32(len(*t)))
	offset := 4
	for _, s := range *t {
		data, err := s.MarshalHandleBinary()
		if err != nil {
			return nil, errors.Wrap(err, "cannot marshal type")
		}
		copy(result[offset:], data)
		offset += len(data)
	}
	return result, nil
}

func (t *TypeList) UnmarshalHandleBinary(data []byte) error {
	l := binary.BigEndian.Uint32(data)
	offset := 4
	*t = make([]UTF8String, l)
	for i := uint32(0); i < l; i++ {
		if err := (*t)[i].UnmarshalHandleBinary(data[offset:]); err != nil {
			return errors.Wrapf(err, "cannot unmarshal type %d", i)
		}
		offset += int((*t)[i].Size())
	}
	return nil
}
