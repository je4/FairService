package handle_net

import (
	"emperror.dev/errors"
	"encoding/binary"
)

type IndexList []byte

func (i *IndexList) Size() uint32 {
	return uint32(4 + len(*i))
}

func (i *IndexList) MarshalHandleBinary() ([]byte, error) {
	var l = uint32(len(*i))
	var result = make([]byte, i.Size())
	binary.BigEndian.PutUint32(result, l)
	copy(result[4:], *i)
	return result, nil
}

func (i *IndexList) UnmarshalHandleBinary(data []byte) error {
	l := binary.BigEndian.Uint32(data)
	if len(data[4:]) < int(l) {
		return errors.New("not enough data")
	}
	*i = data[4 : 4+l]
	return nil
}
