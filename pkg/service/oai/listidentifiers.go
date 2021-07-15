package oai

type ResumptionToken struct {
	ExpirationDate   string `xml:"expirationDate,attr,omitempty"`
	CompleteListSize int64  `xml:"completeListSize,attr,omitempty"`
	Cursor           int64  `xml:"cursor,attr,omitempty"`
	Value            string `xml:",chardata"`
}

type ListIdentifiers struct {
	Header          []*RecordHeader  `xml:"header"`
	ResumptionToken *ResumptionToken `xml:"resumptionToken,omitempty"`
}
