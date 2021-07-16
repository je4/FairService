package oai

type ListRecords struct {
	Record          []*Record        `xml:"record"`
	ResumptionToken *ResumptionToken `xml:"resumptionToken"`
}
