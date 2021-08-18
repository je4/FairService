package oai

type ErrorCodeType string

const (
	ErrorCodeCannotDisseminateFormat ErrorCodeType = "cannotDisseminateFormat"
	ErrorCodeIdDoesNotExist          ErrorCodeType = "idDoesNotExist"
	ErrorCodeBadArgument             ErrorCodeType = "badArgument"
	ErrorCodeBadVerb                 ErrorCodeType = "badVerb"
	ErrorCodeNoMetadataFormats       ErrorCodeType = "noMetadataFormats"
	ErrorCodeNoRecordsMatch          ErrorCodeType = "noRecordsMatch"
	ErrorCodeBadResumptionToken      ErrorCodeType = "badResumptionToken"
	ErrorCodeNoSetHierarchy          ErrorCodeType = "noSetHierarchy"
)

type Error struct {
	Code  ErrorCodeType `xml:"code,attr"`
	Value string        `xml:",chardata"`
}

/*
type ErrorCode struct {
	Code  ErrorCodeType `xml:"code,attr"`
	Value string        `xml:",chardata"`
}

type Error struct {
	Code ErrorCode `xml:"code"`
}

*/
