package handle_net

import (
	"bytes"
	"emperror.dev/errors"
	"encoding/binary"
)

type OpCode uint32

const (
	OpCodeOCReserved           OpCode = 0
	OpCodeOCResolution                = 1   // 1 Handle query
	OpCodeOCGetSiteInfo               = 2   // 2 Get HS_SITE values
	OpCodeOCCreateHandle              = 100 // 100 Create new handle
	OpCodeOCDeleteHandle              = 101 // 101 Delete existing handle
	OpCodeOCAddValue                  = 102 // 102 Add handle value(s)
	OpCodeOCRemoveValue               = 103 // 103 Remove handle value(s)
	OpCodeOCModifyValue               = 104 // 104 Modify handle value(s)
	OpCodeOCListHandle                = 105 // 105 List handles
	OpCodeOCListNA                    = 106 // 106 List sub-naming authorities
	OpCodeOCChallengeResponse         = 200 // 200 Response to challenge
	OpCodeOCVerifyResponse            = 201 // 201 Verify challenge response
	OpCodeOCSessionSetup              = 400 // 400 Session setup request
	OpCodeOCSessionTerminate          = 401 // 401 Session termination request
	OpCodeOCSessionExchangeKey        = 402 // 402 Session key exchange
)

type ResponseCode uint32

const (
	RCReserved           ResponseCode = 0   // 0 Reserved for request
	RCSuccess                         = 1   // 1 Success response
	RCError                           = 2   // 2 General error
	RcServerBusy                      = 3   // 3 Server too busy to respond
	RcProtocolError                   = 4   // 4 Corrupted or unrecognizable message
	RcOperationDenied                 = 5   // 5 Unsupported operation
	RcRecurLimitExceeded              = 6   // 6 Too many recursions for the request
	RcHandleNotFound                  = 100 // 100 Handle not found
	RcHandleAlreadyExist              = 101 // 101 Handle already exists
	RcInvalidHandle                   = 102 // 102 Encoding (or syntax) error
	RcValueNotFound                   = 200 // 200 Value not found
	RcValueAlreadyExist               = 201 // 201 Value already exists
	RcValueInvalid                    = 202 // 202 Invalid handle value
	RcExpiredSiteInfo                 = 300 // 300 SITE_INFO out of date
	RcServerNotResp                   = 301 // 301 Server not responsible
	RcServiceReferral                 = 302 // 302 Server referral
	RcNaDelegate                      = 303 // 303 Naming authority delegation takes place.
	RcNotAuthorized                   = 400 // 400 Not authorized/permitted
	RcAccessDenied                    = 401 // 401 No access to data
	RcAuthenNeeded                    = 402 // 402 Authentication required
	RcAuthenFailed                    = 403 // 403 Failed to authenticate
	RcInvalidCredential               = 404 // 404 Invalid credential
	RcAuthenTimeout                   = 405 // 405 Authentication timed out
	RcUnableToAuthen                  = 406 // 406 Unable to authenticate
	RcSessionTimeout                  = 500 // 500 Session expired
	RcSessionFailed                   = 501 // 501 Unable to establish session
	RcNoSessionKey                    = 502 // 502 No session yet available
	RcSessionNoSupport                = 503 // 503 Session not supported
	RcSessionKeyInvalid               = 504 // 504 Invalid session key
	RcTrying                          = 900 // 900 Request under processing
	RcForwarded                       = 901 // 901 Request forwarded to another server
	RcQueued                          = 902 // 902 Request queued for later processing
)

type OpFlag struct {
	AT,
	CT,
	ENC,
	REC,
	CA,
	CN,
	KC,
	PO,
	RD bool
}

func (of *OpFlag) Size() uint {
	return 4
}

func (of *OpFlag) MarshalHandleBinary() ([]byte, error) {
	var result = make([]byte, 4)
	if of.AT {
		result[0] |= 1 << 7
	}
	if of.CT {
		result[0] |= 1 << 6
	}
	if of.ENC {
		result[0] |= 1 << 5
	}
	if of.REC {
		result[0] |= 1 << 4
	}
	if of.CA {
		result[0] |= 1 << 3
	}
	if of.CN {
		result[0] |= 1 << 2
	}
	if of.KC {
		result[0] |= 1 << 1
	}
	if of.PO {
		result[0] |= 1
	}
	if of.RD {
		result[1] |= 1 << 7
	}

	return result, nil
}

func (of *OpFlag) UnmarshalHandleBinary(data []byte) error {
	of.AT = (data[0] & (1 << 7)) > 0
	of.CT = (data[0] & (1 << 6)) > 0
	of.ENC = (data[0] & (1 << 5)) > 0
	of.REC = (data[0] & (1 << 4)) > 0
	of.CA = (data[0] & (1 << 3)) > 0
	of.CN = (data[0] & (1 << 2)) > 0
	of.KC = (data[0] & (1 << 1)) > 0
	of.PO = (data[0] & 1) > 0
	of.RD = (data[1] & (1 << 7)) > 0
	return nil
}

/*

AT |CT |ENC|REC|CA |CN |KC |PO |RD


AT   -  AuThoritative bit.  A request with the AT bit set (to 1)
               indicates that the request should be directed to the
               primary service site (instead of any mirroring sites).  A
               response message with the AT bit set (to 1) indicates
               that the message is returned from a primary server
               (within the primary service site).

       CT   -  CerTified bit.  A request with the CT bit set (to 1) asks
               the server to sign its response with its digital
               signature.  A response with the CT bit set (to 1)
               indicates that the message is signed.  The server must
               sign its response if the request has its CT bit set (to
               1).  If the server fails to provide a valid signature in
               its response, the client should discard the response and
               treat the request as failed.

       ENC  -  ENCryption bit.  A request with the ENC bit set (to 1)
               requires the server to encrypt its response using the
               pre-established session key.

Sun, et al.                  Informational                     [Page 14]
RFC 3652             Handle System Protocol (v2.1)         November 2003

       REC  -  RECursive bit.  A request with the REC bit set (to 1)
               asks the server to forward the query on behalf of the
               client if the request has to be processed by another
               handle server.  The server may honor the request by
               forwarding the request to the appropriate handle server
               and passing on any result back to the client.  The server
               may also deny any such request by sending a response
               with <ResponseCode> set to RC_SERVER_NOT_RESP.

       CA   -  Cache Authentication.  A request with the CA bit set (to
               1) asks the caching server (if any) to authenticate any
               server response (e.g., verifying the server's signature)
               on behalf of the client.  A response with the CA bit set
               (to 1) indicates that the response has been
               authenticated by the caching server.

       CN   -  ContiNuous bit.  A message with the CN bit set (to 1)
               tells the message recipient that more messages that are
               part of the same request (or response) will follow.  This
               happens if a request (or response) has data that is too
               large to fit into any single message and has to be
               fragmented into multiple messages.

       KC   -  Keep Connection bit.  A message with the KC bit set
               requires the message recipient to keep the TCP
               connection open (after the response is sent back).  This
               allows the same TCP connection to be used for multiple
               handle operations.

       PO   -  Public Only bit.  Used by query operations only.  A query
               request with the PO bit set (to 1) indicates that the
               client is only asking for handle values that have the
               PUB_READ permission.  A request with PO bit set to zero
               asks for all the handle values regardless of their read
               permission.  If any of the handle values require
               ADMIN_READ permission, the server must authenticate the
               client as the handle administrator.

       RD   -  Request-Digest bit.  A request with the RD bit set (to 1)
               asks the server to include in its response the message
               digest of the request.  A response message with the RD
               bit set (to 1) indicates that the first field in the
               Message MessageBody contains the message digest of the original
               request.  The message digest can be used to check the
               integrity of the server response.  Details of these are
               discussed later in this document.
*/

type MessageHeader struct {
	OpCode               OpCode
	ResponseCode         ResponseCode
	OpFlag               OpFlag
	SiteInfoSerialNumber uint16
	RecursionCount       uint8
	free                 byte
	ExpirationTime       uint32
	BodyLength           uint32
}

func (h *MessageHeader) Size() uint32 {
	return 24
}
func (h *MessageHeader) MarshalHandleBinary() ([]byte, error) {
	var buf = bytes.NewBuffer(nil)
	binary.Write(buf, binary.BigEndian, h.OpCode)
	binary.Write(buf, binary.BigEndian, h.ResponseCode)
	data, err := h.OpFlag.MarshalHandleBinary()
	if err != nil {
		return nil, errors.Wrap(err, "cannot marshal op flag")
	}
	buf.Write(data)
	binary.Write(buf, binary.BigEndian, h.SiteInfoSerialNumber)
	binary.Write(buf, binary.BigEndian, h.RecursionCount)
	binary.Write(buf, binary.BigEndian, h.free)
	binary.Write(buf, binary.BigEndian, h.ExpirationTime)
	binary.Write(buf, binary.BigEndian, h.BodyLength)

	return buf.Bytes(), nil
}

func (h *MessageHeader) UnmarshalHandleBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	if err := binary.Read(buf, binary.BigEndian, &h.OpCode); err != nil {
		return errors.Wrap(err, "cannot unmarshal opcode")
	}
	if err := binary.Read(buf, binary.BigEndian, &h.ResponseCode); err != nil {
		return errors.Wrap(err, "cannot unmarshal response code")
	}
	h.OpFlag = OpFlag{}
	if err := h.OpFlag.UnmarshalHandleBinary(buf.Next(4)); err != nil {
		return errors.Wrap(err, "cannot unmarshal op flag")
	}
	if err := binary.Read(buf, binary.BigEndian, &h.SiteInfoSerialNumber); err != nil {
		return errors.Wrap(err, "cannot unmarshal site info serial number")
	}
	if err := binary.Read(buf, binary.BigEndian, &h.RecursionCount); err != nil {
		return errors.Wrap(err, "cannot unmarshal recursion count")
	}
	if err := binary.Read(buf, binary.BigEndian, &h.free); err != nil {
		return errors.Wrap(err, "cannot unmarshal free")
	}
	if err := binary.Read(buf, binary.BigEndian, &h.ExpirationTime); err != nil {
		return errors.Wrap(err, "cannot unmarshal expiration time")
	}
	if err := binary.Read(buf, binary.BigEndian, &h.BodyLength); err != nil {
		return errors.Wrap(err, "cannot unmarshal body length")
	}
	return nil
}

var _ HandleBinaryMarshaler = (*OpFlag)(nil)
var _ HandleBinaryMarshaler = (*MessageHeader)(nil)
