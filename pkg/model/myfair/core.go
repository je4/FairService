package myfair

type CoreTitleType string

const TitleTypeAlternativeTitle CoreTitleType = "AlternativeTitle"
const TitleTypeSubTitle CoreTitleType = "Subtitle"
const TitleTypeTranslatedTitle CoreTitleType = "TranslatedTitle"
const TitleTypeOther CoreTitleType = "Other"

type Name struct {
	Value string   `json:"value"`
	Lang  string   `json:"lang,omitempty"`
	Type  NameType `json:"type,omitempty"`
}

type NameIdentifier struct {
	Value                string `json:"value"`
	Lang                 string `json:"lang,omitempty"`
	SchemeURI            string `json:"schemeURI,omitempty"`
	NameIdentifierScheme string `json:"nameIdentifierScheme,omitempty"`
}

type Creator struct {
	CreatorName    Name           `json:"creatorName"`
	GivenName      string         `json:"givenName,omitempty"`
	FamilyName     string         `json:"familyName,omitempty"`
	Affiliation    string         `json:"affiliation,omitempty"`
	NameIdentifier NameIdentifier `json:"nameIdentifier,omitempty"`
}

type Contributor struct {
	ContributorType ContributorType `xml:"contributorType,attr"`
	ContributorName Name            `xml:"contributorName"`
	GivenName       string          `xml:"givenName,omitempty"`
	FamilyName      string          `xml:"familyName,omitempty"`
	Affiliation     string          `xml:"affiliation,omitempty"`
	NameIdentifier  NameIdentifier  `xml:"nameIdentifier,omitempty"`
}

type Identifier struct {
	Value          string                `json:"value"`
	IdentifierType RelatedIdentifierType `json:"identifierType"`
}

type Title struct {
	Data string
	Type CoreTitleType
}

type Core struct {
	// DataCite: #1 Identifier (with mandatory type sub-property)
	Identifier []Identifier

	// DataCite: #2 Creator (with optional given name, family name, name identifier
	//              and affiliation sub-properties)
	Creator []Creator `json:"creator"`

	Contributor []Contributor `json:"contributor,omitempty"`

	// DataCite: #3 Title (with optional type sub-properties
	Title []Title

	// DataCite: #4 Publisher
	Publisher string

	// DataCite: #5 Publicationyear
	PublicationYear string

	// DataCite: #10 ResourceType (with mandatory general type description subproperty)
	ResourceType ResourceType
}
