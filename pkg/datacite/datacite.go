package datacite

import "encoding/xml"

type DataCiteTitle struct {
	XMLName xml.Name          `xml:"title"`
	Value   string            `xml:",chardata"`
	Lang    string            `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	Type    DataCiteTitleType `xml:"titleType,attr,omitempty"`
}

type Titles struct {
	Title []DataCiteTitle
}

type DataCiteName struct {
	Value string           `xml:",chardata"`
	Lang  string           `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	Type  DataCiteNameType `xml:"nameType,attr,omitempty"`
}

type DataCiteNameIdentifier struct {
	XMLName              xml.Name `xml:"nameIdentifier"`
	Value                string   `xml:",chardata"`
	Lang                 string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty"`
	SchemeURI            string   `xml:"schemeURI,attr,omitempty"`
	NameIdentifierScheme string   `xml:"nameIdentifierScheme,attr,omitempty"`
}

type DataCiteCreator struct {
	XMLName        xml.Name               `xml:"creator"`
	CreatorName    DataCiteName           `xml:"creatorName"`
	GivenName      string                 `xml:"givenName,omitempty"`
	FamilyName     string                 `xml:"familyName,omitempty"`
	Affiliation    string                 `xml:"affiliation,omitempty"`
	NameIdentifier DataCiteNameIdentifier `xml:"nameIdentifier,omitempty"`
}

type Creators struct {
	Creator []DataCiteCreator
}

type DataCiteContributor struct {
	XMLName         xml.Name                `xml:"contributor"`
	ContributorType DataCiteContributorType `xml:"contributorType,attr"`
	ContributorName DataCiteName            `xml:"contributorName"`
	GivenName       string                  `xml:"givenName,omitempty"`
	FamilyName      string                  `xml:"familyName,omitempty"`
	Affiliation     string                  `xml:"affiliation,omitempty"`
	NameIdentifier  DataCiteNameIdentifier  `xml:"nameIdentifier,omitempty"`
}

type Contributors struct {
	Contributor []DataCiteContributor
}

type DataCiteIdentifier struct {
	XMLName        xml.Name                      `xml:"identifier"`
	Value          string                        `xml:",chardata"`
	IdentifierType DataCiteRelatedIdentifierType `xml:"identifierType,attr"`
}

type DataCiteResourceType struct {
	XMLName        xml.Name                    `xml:"resourceType"`
	Value          string                      `xml:",chardata"`
	IdentifierType DataCiteResourceTypeGeneral `xml:"resourceTypeGeneral,attr"`
}

type DataCite struct {
	XMLName           xml.Name `xml:"http://datacite.org/schema/kernel-4 resource"`
	XsiType           string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`

	// DataCite: #1 Identifier (with mandatory type sub-property)
	Identifier DataCiteIdentifier `xml:"identifier"`

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
	ResourceType DataCiteResourceType `xml:"resourceType"`

	// optional fields
	Contributors Contributors `xml:"contributors"`
}

func (datacite *DataCite) InitNamespace() {
	datacite.XsiType = "http://www.w3.org/2001/XMLSchema-instance"
	datacite.XsiSchemaLocation = "http://datacite.org/schema/kernel-4 http://schema.datacite.org/meta/kernel-4.1/metadata.xsd"
}
