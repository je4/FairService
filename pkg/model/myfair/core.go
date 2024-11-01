package myfair

type Name struct {
	Value string   `json:"value"`
	Lang  string   `json:"lang,omitempty"`
	Type  NameType `json:"type,omitempty"`
}

type NameIdentifier struct {
	Value                string `json:"value,omitempty"`
	Lang                 string `json:"lang,omitempty"`
	SchemeURI            string `json:"schemeURI,omitempty"`
	NameIdentifierScheme string `json:"nameIdentifierScheme,omitempty"`
}

type Person struct { // Todo: json anpassen
	PersonType     PersonType     `json:"personType"`
	PersonName     Name           `json:"personName"`
	GivenName      string         `json:"givenName,omitempty"`
	FamilyName     string         `json:"familyName,omitempty"`
	Affiliation    string         `json:"affiliation,omitempty"`
	NameIdentifier NameIdentifier `json:"nameIdentifier,omitempty"`
}

type Identifier struct {
	Value          string                `json:"value"` // Todo: value durch identifier ersetzen
	IdentifierType RelatedIdentifierType `json:"identifierType"`
}

type Title struct {
	Lang string        `json:"lang"`
	Data string        `json:"value"`
	Type CoreTitleType `json:"type"`
}

type Media struct {
	Name        string `json:"name"`
	Mimetype    string `json:"mimetype"`
	Type        string `json:"type"`
	Uri         string `json:"uri"`
	Width       int64  `json:"width,omitempty"`
	Height      int64  `json:"height,omitempty"`
	Orientation int64  `json:"orientation,omitempty"`
	Duration    int64  `json:"duration,omitempty"`
	Fulltext    string `json:"fulltext,omitempty"`
}

type Core struct {
	// DataCite: #1 Identifier (with mandatory type sub-property)
	Identifier []Identifier `json:"identifier"` // Todo: identifier druch identifiers ersetzen

	// DataCite: #2 Person (with optional given name, family name, name identifier
	//              and affiliation sub-properties)
	Person []Person `json:"person"` // Todo: person durch creators ersetzen

	// DataCite: #3 Title (with optional type sub-properties
	Title []Title `json:"title"` // todo: json anpassen

	// DataCite: #4 Publisher
	Publisher string `json:"publisher"` // Todo: create struct

	// DataCite: #5 Publicationyear
	PublicationYear string `json:"publicationYear"`

	// DataCite: #10 ResourceType (with mandatory general type description subproperty)
	ResourceType ResourceType `json:"resourceType"`

	Rights  string `json:"rights,omitempty"`
	License string `json:"license,omitempty"`

	Media  []*Media `json:"media"`
	Poster *Media   `json:"poster"`
}
