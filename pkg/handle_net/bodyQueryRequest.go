package handle_net

import "emperror.dev/errors"

type BodyQueryRequest struct {
	Handle UTF8String
	Index  IndexList
	Type   TypeList
}

func (b *BodyQueryRequest) Size() uint32 {
	return b.Handle.Size() + b.Index.Size() + b.Type.Size()
}

func (b *BodyQueryRequest) MarshalHandleBinary() ([]byte, error) {
	var result = make([]byte, b.Size())
	offset := 0
	data, err := b.Handle.MarshalHandleBinary()
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal handle")
	}
	copy(result[offset:], data)
	offset += len(data)
	data, err = b.Index.MarshalHandleBinary()
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal index")
	}
	copy(result[offset:], data)
	offset += len(data)
	data, err = b.Type.MarshalHandleBinary()
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal type")
	}
	copy(result[offset:], data)
	return result, nil
}

func (b *BodyQueryRequest) UnmarshalHandleBinary(data []byte) error {
	offset := 0
	if err := b.Handle.UnmarshalHandleBinary(data[offset:]); err != nil {
		return errors.Wrap(err, "cannot unmarshal handle")
	}
	offset += int(b.Handle.Size())
	if err := b.Index.UnmarshalHandleBinary(data[offset:]); err != nil {
		return errors.Wrap(err, "cannot unmarshal index")
	}
	offset += int(b.Index.Size())
	if err := b.Type.UnmarshalHandleBinary(data[offset:]); err != nil {
		return errors.Wrap(err, "cannot unmarshal type")
	}
	return nil
}
