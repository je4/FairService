package oai

import "encoding/xml"

const APIPATH string = "oai"

type Request struct {
	Verb           string `xml:"verb,attr"`
	Identifier     string `xml:"identifier,attr,omitempty"`
	MetadataPrefix string `xml:"metadataPrefix,attr,omitempty"`
	Value          string `xml:",chardata"`
}

type OAIPMH struct {
	XMLName             xml.Name             `xml:"OAI-PMH"`
	NS                  string               `xml:"xmlns,attr"`
	XsiType             string               `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation   string               `xml:"xsi:schemaLocation,attr"`
	ResponseDate        string               `xml:"responseDate"`
	Request             *Request             `xml:"request"`
	Error               *Error               `xml:"error,omitempty"`
	GetRecord           *GetRecord           `xml:"GetRecord,omitempty"`
	Identify            *Identify            `xml:"Identify,omitempty"`
	ListIdentifiers     *ListIdentifiers     `xml:"ListIdentifiers,omitempty"`
	ListRecords         *ListRecords         `xml:"ListRecords,omitempty"`
	ListSets            *ListSets            `xml:"ListSets,omitempty"`
	ListMetadataFormats *ListMetadataFormats `xml:"ListMetadataFormats,omitempty"`
}

func (pmh *OAIPMH) InitNamespace() {
	pmh.NS = "http://www.openarchives.org/OAI/2.0/"
	pmh.XsiType = "http://www.w3.org/2001/XMLSchema-instance"
	pmh.XsiSchemaLocation = "http://www.openarchives.org/OAI/2.0/ http://www.openarchives.org/OAI/2.0/OAI-PMH.xsd"
}
