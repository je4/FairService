package oai

type Set struct {
	SetSpec string `xml:"setSpec"`
	SetName string `xml:"setName"`
}

type ListSets struct {
	Set             []*Set           `xml:"set"`
	ResumptionToken *ResumptionToken `xml:"resumptionToken"`
}
