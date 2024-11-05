package handle_net

import "encoding/binary"

type UTF8String string

func (s *UTF8String) Size() uint32 {
	return uint32(4 + len([]byte(*s)))
}

func (s *UTF8String) MarshalHandleBinary() ([]byte, error) {
	var l = uint32(len(*s))
	data := []byte(*s)
	var result = make([]byte, 4+l)
	binary.BigEndian.PutUint32(result, l)
	copy(result[4:], data)
	return result, nil
}

func (s *UTF8String) UnmarshalHandleBinary(data []byte) error {
	l := binary.BigEndian.Uint32(data)
	*s = UTF8String(data[4 : 4+l])
	return nil
}
