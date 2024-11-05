package handle_net

import (
	"emperror.dev/errors"
	"encoding/binary"
)

type HandleBinaryMarshaler interface {
	MarshalHandleBinary() ([]byte, error)
	UnmarshalHandleBinary([]byte) error
}

type Message struct {
	Envelope         MessageEnvelope
	Header           MessageHeader
	Body             MessageBody
	CredentialLength uint32
	Credential       Credential
}

func (m *Message) SetSizes() {
	m.Header.BodyLength = m.Body.Size()
	m.CredentialLength = m.Credential.GetLength()
	m.Envelope.MessageLength = uint32(20 + m.Header.Size() + m.Header.BodyLength + 4 + m.CredentialLength)
}

func (m *Message) MarshalHandleBinary() ([]byte, error) {
	var result = []byte{}
	data, err := m.Envelope.MarshalHandleBinary()
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal envelope")
	}
	result = append(result, data...)
	data, err = m.Header.MarshalHandleBinary()
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal header")
	}
	result = append(result, data...)
	data, err = m.Body.GetData(m.Header.OpFlag.RD)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get body data")
	}
	result = append(result, data...)
	data = make([]byte, 4)
	binary.BigEndian.PutUint32(data, m.Credential.GetLength())
	result = append(result, data...)
	data, err = m.Credential.MarshalHandleBinary()
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal credential")
	}
	result = append(result, data...)
	return result, nil
}

func (m *Message) UnmarshalHandleBinary(data []byte) error {
	if err := m.Envelope.UnmarshalHandleBinary(data); err != nil {
		return errors.Wrap(err, "cannot unmarshal envelope")
	}
	data = data[20:]
	if err := m.Header.UnmarshalHandleBinary(data); err != nil {
		return errors.Wrap(err, "cannot unmarshal header")
	}
	data = data[24:]
	m.Body.SetData(data[0:m.Header.BodyLength], m.Header.OpFlag.RD, m.Header.OpCode)
	data = data[m.Header.BodyLength:]
	m.CredentialLength = binary.BigEndian.Uint32(data)
	data = data[4:]
	if err := m.Credential.UnmarshalHandleBinary(data); err != nil {
		return errors.Wrap(err, "cannot unmarshal credential")
	}
	return nil
}

var _ HandleBinaryMarshaler = (*Message)(nil)
