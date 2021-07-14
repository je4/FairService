package oai

type ErrorCode string

const (
	ErrorCodeCannotDisseminateFormat ErrorCode = "cannotDisseminateFormat"
	ErrorCodeIdDoesNotExist          ErrorCode = "idDoesNotExist"
	ErrorCodeBadArgument             ErrorCode = "badArgument"
	ErrorCodeBadVerb                 ErrorCode = "badVerb"
	ErrorCodeNoMetadataFormats       ErrorCode = "noMetadataFormats"
	ErrorCodeNoRecordsMatch          ErrorCode = "noRecordsMatch"
	ErrorCodeBadResumptionToken      ErrorCode = "badResumptionToken"
	ErrorCodeNoSetHierarchy          ErrorCode = "noSetHierarchy"
)

type Error struct {
	Code ErrorCode `xml:"code"`
}
