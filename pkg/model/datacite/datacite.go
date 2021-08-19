package datacite

import (
	"encoding/json"
	"encoding/xml"
)

type Title struct {
	XMLName xml.Name  `xml:"title" json:"-"`
	Value   string    `xml:",chardata" json:"title"`
	Lang    string    `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty" json:"lang,omitempty"`
	Type    TitleType `xml:"titleType,attr,omitempty" json:"type,omitempty"`
}

type Titles struct {
	Title []Title
}

func (titles Titles) UnmarshalJSON(data []byte) error {
	titles.Title = []Title{}
	if err := json.Unmarshal(data, &titles.Title); err != nil {
		return err
	}
	return nil
}

func (titles Titles) MarshalJSON() ([]byte, error) {
	return json.Marshal(titles.Title)
}

type Name struct {
	Value string   `xml:",chardata" json:"value"`
	Lang  string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty" json:"lang,omitempty"`
	Type  NameType `xml:"nameType,attr,omitempty" json:"type,omitempty"`
}

type NameIdentifier struct {
	XMLName              xml.Name `xml:"nameIdentifier" json:"-"`
	Value                string   `xml:",chardata" json:"nameIdentifier"`
	Lang                 string   `xml:"http://www.w3.org/XML/1998/namespace lang,attr,omitempty" json:"lang,omitempty"`
	SchemeURI            string   `xml:"schemeURI,attr,omitempty" json:"schemeUri,omitempty"`
	NameIdentifierScheme string   `xml:"nameIdentifierScheme,attr,omitempty" json:"name_identifierScheme,omitempty"`
}

type Creator struct {
	XMLName        xml.Name         `xml:"creator" json:"-"`
	CreatorName    Name             `xml:"creatorName" json:"creatorName,omitempty"`
	GivenName      string           `xml:"givenName,omitempty" json:"givenName,omitempty"`
	FamilyName     string           `xml:"familyName,omitempty" json:"familyName,omitempty"`
	Affiliation    []string         `xml:"affiliation,omitempty" json:"affiliation,omitempty"`
	NameIdentifier []NameIdentifier `xml:"nameIdentifier,omitempty" json:"nameIdentifier,omitempty"`
}

type Creators struct {
	Creator []Creator
}

func (creators Creators) UnmarshalJSON(data []byte) error {
	creators.Creator = []Creator{}
	if err := json.Unmarshal(data, &creators.Creator); err != nil {
		return err
	}
	return nil
}

func (creators Creators) MarshalJSON() ([]byte, error) {
	return json.Marshal(creators.Creator)
}

type Contributor struct {
	ContributorType ContributorType `xml:"contributorType" json:"contributorType"`
	ContributorName Name            `xml:"contributorName" json:"contributorName"`
	GivenName       string          `xml:"givenName,omitempty" json:"givenName,omitempty"`
	FamilyName      string          `xml:"familyName,omitempty" json:"familyName,omitempty"`
	Affiliation     string          `xml:"affiliation,omitempty" json:"affiliation,omitempty"`
	NameIdentifier  NameIdentifier  `xml:"nameIdentifier,omitempty" json:"nameIdentifier,omitempty"`
}

type Contributors struct {
	Contributor []Contributor
}

func (contributors Contributors) UnmarshalJSON(data []byte) error {
	contributors.Contributor = []Contributor{}
	if err := json.Unmarshal(data, &contributors.Contributor); err != nil {
		return err
	}
	return nil
}

func (contributors Contributors) MarshalJSON() ([]byte, error) {
	return json.Marshal(contributors.Contributor)
}

type Identifier struct {
	XMLName        xml.Name              `xml:"identifier" json:"-"`
	Value          string                `xml:",chardata" json:"identifier"`
	IdentifierType RelatedIdentifierType `xml:"identifierType,attr" json:"identifierType"`
}

type AlternateIdentifier struct {
	XMLName                 xml.Name              `xml:"identifier" json:"-"`
	Value                   string                `xml:",chardata" json:"alternateIdentifier"`
	AlternateIdentifierType RelatedIdentifierType `xml:"alternateIdentifierType,attr" json:"alternateIdentifierType"`
}

type ResourceType struct {
	XMLName        xml.Name            `xml:"resourceType" json:"-"`
	Value          string              `xml:",chardata" json:"resourceType"`
	IdentifierType ResourceTypeGeneral `xml:"resourceTypeGeneral,attr" json:"identifierType"`
}

type DataCite struct {
	XMLName           xml.Name `xml:"http://datacite.org/schema/kernel-4 resource" json:"-"`
	XsiType           string   `xml:"xmlns:xsi,attr" json:"-"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr" json:"-"`

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
