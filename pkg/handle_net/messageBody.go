package handle_net

import "github.com/pkg/errors"

type MessageBody struct {
	Digest  *BodyDigest
	Data    []byte
	Request *BodyQueryRequest
}

func (b *MessageBody) Size() uint32 {
	var size uint32
	if b.Digest != nil {
		size += b.Digest.Size()
	}
	if b.Request != nil {
		size += b.Request.Size()
	} else if b.Data != nil {
		size += uint32(len(b.Data))
	}
	return size
}

func (b *MessageBody) SetData(data []byte, hasDigest bool, op OpCode) error {
	if hasDigest {
		b.Digest = &BodyDigest{}
		if err := b.Digest.UnmarshalHandleBinary(data); err != nil {
			return errors.Wrap(err, "cannot unmarshal digest")
		}
		data = data[b.Digest.Size():]
	} else {
		b.Digest = nil
	}
	b.Data = nil
	b.Request = nil
	if len(data) > 0 {
		switch op {
		case OpCodeOCResolution:
			b.Request = &BodyQueryRequest{}
			if err := b.Request.UnmarshalHandleBinary(data); err != nil {
				return errors.Wrap(err, "cannot unmarshal request")
			}
		default:
			b.Data = data
		}
	}
	return nil
}

func (b *MessageBody) GetData(hasDigest bool) ([]byte, error) {
	var result = []byte{}
	if b.Digest != nil {
		data, err := b.Digest.MarshalHandleBinary()
		if err != nil {
			return nil, errors.Wrap(err, "cannot marshal digest")
		}
		result = append(result, data...)
	}
	if b.Request != nil {
		data, err := b.Request.MarshalHandleBinary()
		if err != nil {
			return nil, errors.Wrap(err, "cannot marshal request")
		}
		result = append(result, data...)
	} else if b.Data != nil {
		result = append(result, b.Data...)
	}
	return result, nil
}
