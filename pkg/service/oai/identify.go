package oai

type Identify struct {
	RepositoryName    string   `xml:"repositoryName"`
	BaseURL           string   `xml:"baseURL"`
	ProtocolVersion   string   `xml:"protocolVersion"`
	EarliestDatestamp string   `xml:"earliestDatestamp"`
	AdminEmail        []string `xml:"adminEmail"`
	DeletedRecord     string   `xml:"deletedRecord"`
	Granularity       string   `xml:"granularity"`
	Compression       []string `xml:"compression,omitempty"`
}
