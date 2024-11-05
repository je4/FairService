package handle_net

import (
	"crypto/md5"
	"reflect"
	"testing"
	"time"
)

func TestMessage(t *testing.T) {

	md5Data := md5.Sum([]byte(""))
	md5Bytes := make([]byte, 16)
	copy(md5Bytes, md5Data[:])

	message := Message{
		Envelope: MessageEnvelope{
			MajorVersion: 1,
			MinorVersion: 0,
			MessageFlag: MessageFlag{
				CP: false,
				EC: false,
				TC: false,
			},
			SessionId:      100,
			RequestId:      200,
			SequenceNumber: 300,
			MessageLength:  0,
		},
		Header: MessageHeader{
			OpCode:       OpCodeOCResolution,
			ResponseCode: RcAuthenNeeded,
			OpFlag: OpFlag{
				AT:  false,
				CT:  false,
				ENC: false,
				REC: false,
				CA:  false,
				CN:  false,
				KC:  false,
				PO:  false,
				RD:  true,
			},
			SiteInfoSerialNumber: 400,
			RecursionCount:       2,
			ExpirationTime:       uint32(time.Now().Unix()),
			BodyLength:           0,
		},
		Body: MessageBody{
			Digest: &BodyDigest{
				Type:   BodyDigestTypeMD5,
				Digest: md5Bytes,
			},
			Data: nil,
			Request: &BodyQueryRequest{
				Handle: "test",
				Index:  IndexList{1, 2, 3},
				Type:   TypeList{"t1", "t2", "t3"},
			},
		},
		CredentialLength: 4,
		Credential: Credential{
			Version:       3,
			Reserved:      0,
			Options:       [2]byte{},
			Unimplemented: nil,
		},
	}

	message.SetSizes()

	data, err := message.MarshalHandleBinary()
	if err != nil {
		t.Error(err)
	}

	message2 := Message{}
	err = message2.UnmarshalHandleBinary(data)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(message.Envelope, message2.Envelope) {
		t.Error("Envelope are not equal")
	}

	if !reflect.DeepEqual(message.Header, message2.Header) {
		t.Error("Header are not equal")
	}

	if !reflect.DeepEqual(message.Body, message2.Body) {
		t.Error("Body are not equal")
	}

	if !reflect.DeepEqual(message.Credential, message2.Credential) {
		t.Error("Credential are not equal")
	}

	if !reflect.DeepEqual(message, message2) {
		t.Error("Messages are not equal")
	}
}
