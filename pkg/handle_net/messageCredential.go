package handle_net

import "emperror.dev/errors"

type Credential struct {
	Version       uint8
	Reserved      uint8
	Options       [2]byte
	Unimplemented []byte
}

func (c *Credential) GetLength() uint32 {
	return uint32(4 + len(c.Unimplemented))
}

func (c *Credential) MarshalHandleBinary() ([]byte, error) {
	var result = []byte{}
	result = append(result, c.Version)
	result = append(result, c.Reserved)
	result = append(result, c.Options[:]...)
	result = append(result, c.Unimplemented...)
	return result, nil
}

func (c *Credential) UnmarshalHandleBinary(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if len(data) < 4 {
		return errors.New("data is too short")
	}
	c.Version = data[0]
	c.Reserved = data[1]
	c.Options = [2]byte{data[2], data[3]}
	if len(data[4:]) > 0 {
		c.Unimplemented = data[4:]
	} else {
		c.Unimplemented = nil
	}
	return nil
}

var _ HandleBinaryMarshaler = (*Credential)(nil)
