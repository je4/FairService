package dcmi

type DCMI struct {
	OAI_DCNS          string `xml:"xmlns:oai_dc,attr"`
	DCNS              string `xml:"xmlns:dc"`
	XsiType           string `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string `xml:"xsi:schemaLocation,attr"`
	Type              Type   `xml:"dc:type"`

	Title       []string `xml:"dc:title"`
	Creator     []string `xml:"dc:creator"`
	Subject     []string `xml:"dc:subject"`
	Description []string `xml:"dc:description,omitempty"`
	Publisher   []string `xml:"dc:publisher"`
	Contributor []string `xml:"dc:contributor"`
	Date        []string `xml:"dc:date"`
	Format      []string `xml:"dc:format,omitempty"`
	Identifier  []string `xml:"dc:identifier"`
	Source      []string `xml:"dc:source"`
	Language    []string `xml:"dc:language,omitempty"`
	Relation    []string `xml:"dc:relation,omitempty"`
	Coverage    []string `xml:"dc:coverage,omitempty"`
	Rights      []string `xml:"dc:rights,omitempty"`

	AlternativeTitle      []string `xml:"dc:alternativeTitle,omitempty"`
	TableOfContents       []string `xml:"dc:tableOfContents,omitempty"`
	Abstract              []string `xml:"dc:abstract,omitempty"`
	Created               []string `xml:"dc:created,omitempty"`
	Valid                 []string `xml:"dc:valid,omitempty"`
	Available             []string `xml:"dc:available,omitempty"`
	Issued                []string `xml:"dc:issued,omitempty"`
	Modified              []string `xml:"dc:modified,omitempty"`
	DateAccepted          []string `xml:"dc:dateAccepted,omitempty"`
	DateCopyrighted       []string `xml:"dc:dateCopyrighted,omitempty"`
	DateSubmitted         []string `xml:"dc:dateSubmitted,omitempty"`
	Extent                []string `xml:"dc:extent,omitempty"`
	Medium                []string `xml:"dc:medium,omitempty"`
	IsVersionOf           []string `xml:"dc:isVersionOf,omitempty"`
	HasVersion            []string `xml:"dc:hasVersion,omitempty"`
	IsReplacedBy          []string `xml:"dc:isReplacedBy,omitempty"`
	Replaces              []string `xml:"dc:replaces,omitempty"`
	IsRequiredBy          []string `xml:"dc:isRequiredBy,omitempty"`
	Requires              []string `xml:"dc:requires,omitempty"`
	IsPartOf              []string `xml:"dc:isPartOf,omitempty"`
	HasPart               []string `xml:"dc:hasPart,omitempty"`
	IsReferencedBy        []string `xml:"dc:isReferencedBy,omitempty"`
	References            []string `xml:"dc:references,omitempty"`
	IsFormatOf            []string `xml:"dc:isFormatOf,omitempty"`
	HasFormat             []string `xml:"dc:hasFormat,omitempty"`
	ConformsTo            []string `xml:"dc:conformsTo,omitempty"`
	Spatial               []string `xml:"dc:spatial,omitempty"`
	Temporal              []string `xml:"dc:temporal,omitempty"`
	Audience              []string `xml:"dc:audience,omitempty"`
	AccrualMethod         []string `xml:"dc:accrualMethod,omitempty"`
	AccrualPeriodicity    []string `xml:"dc:accrualPeriodicity,omitempty"`
	AccrualPolicy         []string `xml:"dc:accrualPolicy,omitempty"`
	InstructionalMethod   []string `xml:"dc:instructionalMethod,omitempty"`
	Provenance            []string `xml:"dc:provenance,omitempty"`
	RightsHolder          []string `xml:"dc:rightsHolder,omitempty"`
	Mediator              []string `xml:"dc:mediator,omitempty"`
	EducationLevel        []string `xml:"dc:educationLevel,omitempty"`
	AccessRights          []string `xml:"dc:accessRights,omitempty"`
	License               []string `xml:"dc:license,omitempty"`
	BibliographicCitation []string `xml:"dc:bibliographicCitation,omitempty"`
}

func (dcmi *DCMI) InitNamespace() {
	dcmi.OAI_DCNS = "http://www.openarchives.org/OAI/2.0/oai_dc/"
	dcmi.DCNS = "http://purl.org/dc/elements/1.1/"
	dcmi.XsiType = "http://www.w3.org/2001/XMLSchema-instance"
	dcmi.XsiSchemaLocation = "http://www.openarchives.org/OAI/2.0/oai_dc/ http://www.openarchives.org/OAI/2.0/oai_dc.xsd"
}

/*
<?xml version="1.0" encoding="UTF-8"?>

<metadata
xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
xmlns:dc="http://purl.org/dc/elements/1.1/">
<dc:title>title 1</dc:title>
<dc:title>title 2</dc:title>
<dc:creator>Enge, Jürgen</dc:creator>
<dc:subject>testing</dc:subject>
<dc:description>lorem ipsum dolor sit amet
</dc:description>
<dc:relation>sadfsdfsd</dc:relation>

</metadata>
*/
