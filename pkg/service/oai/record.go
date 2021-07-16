package oai

type RecordHeaderStatusType string

const (
	RecordHeaderStatusOK      RecordHeaderStatusType = ""
	RecordHeaderStatusDeleted RecordHeaderStatusType = "deleted"
)

type RecordHeader struct {
	Status     RecordHeaderStatusType `xml:"status,attr,omitempty"`
	Identifier string                 `xml:"identifier"`
	Datestamp  string                 `xml:"datestamp"`
	SetSpec    []string               `xml:"setSpec"`
}

type RecordAboutProvenanceOriginDescription struct {
	BaseURL           string      `xml:"baseURL"`
	Identifier        string      `xml:"identifier"`
	Datestamp         string      `xml:"datestamp"`
	MetadataNamespace string      `xml:"metadataNamespace"`
	HarvestDate       string      `xml:"harvestDate"`
	Altered           bool        `xml:"altered"`
	RepositoryId      string      `xml:"repositoryId,omitempty"`
	RepositoryName    string      `xml:"repositoryName,omitempty"`
	OriginDescription interface{} `xml:"originDescription,omitempty"`
}

type RecordAboutProvenance struct {
	OriginDescription []*RecordAboutProvenanceOriginDescription `xml:"originDescription"`
}

type RecordAbout struct {
	Provenance *RecordAboutProvenance `xml:"provenance"`
}

type Metadata struct {
	OAIDC interface{} `xml:"oai_dc:dc,omitempty"`
}

type Record struct {
	Header   *RecordHeader `xml:"header"`
	Metadata *Metadata     `xml:"metadata"`
	About    *RecordAbout  `xml:"about,omitempty"`
}
