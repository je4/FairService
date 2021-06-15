package myfair

import (
	"github.com/je4/FairService/v2/pkg/model/zsearch"
	"github.com/je4/zsearch/pkg/search"
)

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

type Person struct {
	PersonType     PersonType     `json:"PersonType,attr"`
	PersonName     Name           `json:"PersonName"`
	GivenName      string         `json:"givenName,omitempty"`
	FamilyName     string         `json:"familyName,omitempty"`
	Affiliation    string         `json:"affiliation,omitempty"`
	NameIdentifier NameIdentifier `json:"nameIdentifier,omitempty"`
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

	// DataCite: #2 Person (with optional given name, family name, name identifier
	//              and affiliation sub-properties)
	Person []Person `json:"Person"`

	// DataCite: #3 Title (with optional type sub-properties
	Title []Title

	// DataCite: #4 Publisher
	Publisher string

	// DataCite: #5 Publicationyear
	PublicationYear string

	// DataCite: #10 ResourceType (with mandatory general type description subproperty)
	ResourceType ResourceType
}

func (core *Core) FromSearch(src *search.SourceData) error {
	if zit, err := zsearch.ItemTypeFromString(src.Type); err != nil {
		core.ResourceType = ResourceTypeOther
	} else {
		core.ResourceType = ZSearchItemTypeMap(zit)
	}

	for _, p := range src.Persons {
		ct, err := zsearch.PersonRoleFromString(p.Role)
		if err != nil {
			ct = zsearch.PersonRoleArtist
		}
		core.Person = append(core.Person, Person{
			PersonName: Name{Value: p.Name},
			PersonType: ct,
		})
	}

	return nil
}
