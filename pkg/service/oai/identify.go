package oai

import "encoding/xml"

type Identifier struct {
	XMLName              xml.Name `xml:"oai-identifier"`
	NS                   string   `xml:"xmlns,attr"`
	XsiType              string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation    string   `xml:"xsi:schemaLocation,attr"`
	Scheme               string   `xml:"scheme"`
	RepositoryIdentifier string   `xml:"repositoryIdentifier"`
	Delimiter            string   `xml:"delimiter"`
	SampleIdentifier     string   `xml:"sampleIdentifier"`
}

func (ident *Identifier) InitNamespace() {
	ident.NS = "http://www.openarchives.org/OAI/2.0/oai-identifier"
	ident.XsiType = "http://www.w3.org/2001/XMLSchema-instance"
	ident.XsiSchemaLocation = "http://www.openarchives.org/OAI/2.0/oai-identifier http://www.openarchives.org/OAI/2.0/oai-identifier.xsd"
}

/*
<oai-identifier xmlns="http://www.openarchives.org/OAI/2.0/oai-identifier" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://www.openarchives.org/OAI/2.0/oai-identifier http://www.openarchives.org/OAI/2.0/oai-identifier.xsd">
	<scheme>oai</scheme>
	<repositoryIdentifier>www.research-collection.ethz.ch</repositoryIdentifier>
	<delimiter>:</delimiter>
	<sampleIdentifier>oai:www.research-collection.ethz.ch:20.500.11850/1234</sampleIdentifier>
</oai-identifier>
*/
type Description struct {
	Identifier Identifier `xml:"oai-identifier"`
}

type Identify struct {
	RepositoryName    string      `xml:"repositoryName"`
	BaseURL           string      `xml:"baseURL"`
	ProtocolVersion   string      `xml:"protocolVersion"`
	EarliestDatestamp string      `xml:"earliestDatestamp"`
	AdminEmail        []string    `xml:"adminEmail"`
	DeletedRecord     string      `xml:"deletedRecord"`
	Granularity       string      `xml:"granularity"`
	Compression       []string    `xml:"compression,omitempty"`
	Description       Description `xml:"description,omitempty"`
}
