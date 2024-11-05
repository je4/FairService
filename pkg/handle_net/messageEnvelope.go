package handle_net

import (
	"emperror.dev/errors"
	"encoding/binary"
)

type MessageFlag struct {
	CP,
	EC,
	TC bool
}

func (mf *MessageFlag) MarshalHandleBinary() ([]byte, error) {
	var result = make([]byte, 2)
	if mf.CP {
		result[0] |= 1 << 7
	}
	if mf.EC {
		result[0] |= 1 << 6
	}
	if mf.TC {
		result[0] |= 1 << 5
	}

	return result, nil
}

func (mf *MessageFlag) UnmarshalHandleBinary(data []byte) error {
	mf.CP = (data[0] & (1 << 7)) > 0
	mf.EC = (data[0] & (1 << 6)) > 0
	mf.TC = (data[0] & (1 << 5)) > 0
	return nil
}

type MessageEnvelope struct {
	MajorVersion   uint8
	MinorVersion   uint8
	MessageFlag    MessageFlag
	SessionId      uint32
	RequestId      uint32
	SequenceNumber uint32
	MessageLength  uint32
}

func (env *MessageEnvelope) MarshalHandleBinary() ([]byte, error) {
	var result = make([]byte, 5*4)
	result[0] = env.MajorVersion
	result[1] = env.MinorVersion
	data, err := env.MessageFlag.MarshalHandleBinary()
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal message flag")
	}
	copy(result[2:4], data)
	binary.BigEndian.PutUint32(result[4:8], env.SessionId)
	binary.BigEndian.PutUint32(result[8:12], env.RequestId)
	binary.BigEndian.PutUint32(result[12:16], env.SequenceNumber)
	binary.BigEndian.PutUint32(result[16:20], env.MessageLength)

	return result, nil
}

func (env *MessageEnvelope) UnmarshalHandleBinary(data []byte) error {
	env.MajorVersion = data[0]
	env.MinorVersion = data[1]
	env.MessageFlag = MessageFlag{}
	if err := env.MessageFlag.UnmarshalHandleBinary(data[2:4]); err != nil {
		return errors.Wrap(err, "cannot unmarshal message flag")
	}
	env.SessionId = binary.BigEndian.Uint32(data[4:8])
	env.RequestId = binary.BigEndian.Uint32(data[8:12])
	env.SequenceNumber = binary.BigEndian.Uint32(data[12:16])
	env.MessageLength = binary.BigEndian.Uint32(data[16:20])
	return nil
}

var _ HandleBinaryMarshaler = (*MessageFlag)(nil)
var _ HandleBinaryMarshaler = (*MessageEnvelope)(nil)
