package datacite

import "encoding/xml"

type Title struct {
	XMLName xml.Name  `xml:"title"`
	Value   string    `xml:",chardata"`
	Lang    string    `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	Type    TitleType `xml:"titleType,attr,omitempty"`
}

type Titles struct {
	Title []Title
}

type Name struct {
	Value string   `xml:",chardata"`
	Lang  string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	Type  NameType `xml:"nameType,attr,omitempty"`
}

type NameIdentifier struct {
	XMLName              xml.Name `xml:"nameIdentifier"`
	Value                string   `xml:",chardata"`
	Lang                 string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	SchemeURI            string   `xml:"schemeURI,attr,omitempty"`
	NameIdentifierScheme string   `xml:"nameIdentifierScheme,attr,omitempty"`
}

type Creator struct {
	XMLName        xml.Name       `xml:"creator"`
	CreatorName    Name           `xml:"creatorName"`
	GivenName      string         `xml:"givenName,omitempty"`
	FamilyName     string         `xml:"familyName,omitempty"`
	Affiliation    string         `xml:"affiliation,omitempty"`
	NameIdentifier NameIdentifier `xml:"nameIdentifier,omitempty"`
}

type Creators struct {
	Creator []Creator
}

type Contributor struct {
	ContributorType ContributorType `json:"contributorType"`
	ContributorName Name            `json:"contributorName"`
	GivenName       string          `json:"givenName,omitempty"`
	FamilyName      string          `json:"familyName,omitempty"`
	Affiliation     string          `json:"affiliation,omitempty"`
	NameIdentifier  NameIdentifier  `json:"nameIdentifier,omitempty"`
}

type Contributors struct {
	Contributor []Contributor
}

type Identifier struct {
	XMLName        xml.Name              `xml:"identifier"`
	Value          string                `xml:",chardata"`
	IdentifierType RelatedIdentifierType `xml:"identifierType,attr"`
}

type ResourceType struct {
	XMLName        xml.Name            `xml:"resourceType"`
	Value          string              `xml:",chardata"`
	IdentifierType ResourceTypeGeneral `xml:"resourceTypeGeneral,attr"`
}

type DataCite struct {
	XMLName           xml.Name `xml:"http://datacite.org/schema/kernel-4 resource"`
	XsiType           string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`

	// DataCite: #1 Identifier (with mandatory type sub-property)
	Identifier Identifier `xml:"identifier"`

	// DataCite: #2 Creator (with optional given name, family name, name identifier
	//              and affiliation sub-properties)
	Creators Creators `xml:"creators"`

	// DataCite: #3 Title (with optional type sub-properties
	Titles Titles `xml:"titles"`

	// DataCite: #4 Publisher
	Publisher string `xml:"publisher"`

	// DataCite: #5 PublicationYear
	PublicationYear int `xml:"publicationYear"`

	// DataCite: #10 ResourceType (with mandatory general type description subproperty)
	ResourceType ResourceType `xml:"resourceType"`

	// optional fields
	Contributors Contributors `xml:"contributors"`
}

func (datacite *DataCite) InitNamespace() {
	datacite.XsiType = "http://www.w3.org/2001/XMLSchema-instance"
	datacite.XsiSchemaLocation = "http://datacite.org/schema/kernel-4 http://schema.datacite.org/meta/kernel-4.1/metadata.xsd"
}
