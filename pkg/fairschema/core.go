package fairschema

type CoreTitleType string

const TitleTypeAlternativeTitle CoreTitleType = "AlternativeTitle"
const TitleTypeSubTitle CoreTitleType = "Subtitle"
const TitleTypeTranslatedTitle CoreTitleType = "TranslatedTitle"
const TitleTypeOther CoreTitleType = "Other"

type CoreTitle struct {
	Data string
	Type CoreTitleType
}

type Core struct {
	// DataCite: #1 Identifier (with mandatory type sub-property)
	Identifier string

	// DataCite: #2 Creator (with optional given name, family name, name identifier
	//              and affiliation sub-properties)
	Creator string

	// DataCite: #3 Title (with optional type sub-properties
	Title CoreTitle

	// DataCite: #4 Publisher
	Publisher string

	// DataCite: #5 Publicationyear
	Publicationyear string

	// DataCite: #10 ResourceType (with mandatory general type description subproperty)
	ResourceType string
}
