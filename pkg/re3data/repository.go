package re3data

import "encoding/xml"

type StringLang struct {
	Language string `xml:"lang,attr,omitempty"`
	Value    string `xml:",chardata"`
}

// RepositoryIdentifier An identifier provisioned for the website of the RDR (wrapper element).
type RepositoryIdentifier struct {
	Type  string `xml:"r3d:repositoryIdentifierType"`  // The type of the provider of the identifier for the RDR (e.g. DOI, URN, VIAF, DataCite).
	Value string `xml:"r3d:repositoryIdentifierValue"` // A globally unique identifier that refers to the RDR.
}

type RepositorySubject struct {
	Scheme string `xml:"subjectScheme,attr"`
	Id     string `xml:"r3d:subjectId"`
	Name   string `xml:"r3d:subjectName"`
}

type RepositoryInstitution struct {
	Name StringLang `xml:"r3d:institutionName"`
	// optional
	AdditionalName []StringLang `xml:"r3d:institutionAdditionalName"`
}

type RepositoryDatabaseAccess struct {
	DatabaseAccessType RE3DataAccessType `xml:"r3d:databaseAccessType"` // The type of access to the RDR
	// optional
	DatabaseAccessRestrictions []RE3DataAccessRestrictions `xml:"databaseAccessRestriction"` // All existing access restrictions to the RDR (required if restricted is chosen).
}

type RepositoryDataLicense struct {
	Name string `xml:"r3d:dataLicenseName"` // The name of the data license
	URL  string `xml:"r3d:dataLicenseUrl"`  // The data license URL
}

type RepositoryDataUpload struct {
	Type RE3DataAccessType `xml:"r3d:dataUploadType"`
	// optional
	Restriction []RE3DataAccessRestrictions `xml:"r3d:dataUploadRestriction,omitempty"` // All existing restrictions to the data upload (required if restricted is chosen).
}

type RepositoryIdentifiers struct {
	Re3Data string `xml:"r3d:re3data"` // A unique string to identify the RDR metadata entry. The internal identifier is assigned by re3data.org
	DOI     string `xml:"r3d:doi"`     // The DOI assigned to the re3data.org metadata entry of the RDR to make the metadata entries citable
}

type Repository struct {
	XMLName           xml.Name `xml:"r3d:repository"`
	R3DType           string   `xml:"xmlns:r3d,attr"`
	XsiType           string   `xml:"xmlns:xsi,attr"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr"`

	// mandatory fields
	Identifiers         RepositoryIdentifiers    `xml:"r3d:identifiers"`         // The identifiers provided by re3data.org (wrapper element).
	Name                StringLang               `xml:"r3d:repositoryName"`      // The full name of the RDR
	URL                 string                   `xml:"r3d:repositoryUrl"`       // The URL of the RDR
	Type                []RepositoryType         `xml:"r3d:type"`                // The type of the RDR
	Updated             string                   `xml:"r3d:updated"`             // The date of the last update of the RDR size
	Language            []string                 `xml:"r3d:repositoryLanguage"`  // The user interface language of the RDR
	Subject             []RepositorySubject      `xml:"r3d:subject"`             // The disciplinary focus of the RDR (wrapper element).
	ProviderType        []RE3DataProviderType    `xml:"r3d:providerType"`        // 1-2 The type of provider.
	Institution         []RepositoryInstitution  `xml:"r3d:institution"`         // All institutions being responsible for funding, creating and/or running the RDR (wrapper element).
	DataAccess          RepositoryDatabaseAccess `xml:"r3d:databaseAccess"`      // The access regulation to the RDR (wrapper element).
	DataLicense         RepositoryDataLicense    `xml:"r3d:dataLicense"`         // The license of the research data, existing in the RDR (wrapper element).
	DataUpload          []RepositoryDataUpload   `xml:"r3d:dataUpload"`          // The regulation for submitting research data to the RDR (wrapper element)
	Versioning          RE3DataYesNoUn           `xml:"r3d:versioning"`          // The RDR supports versioning of research data
	EnhancedPublication RE3DataYesNoUn           `xml:"r3d:enhancedPublication"` // The RDR offers the interlinking between publications and research data
	QualityManagement   RE3DataYesNoUn           `xml:"r3d:qualityManagement"`   // Any form of quality management concerning the research data or metadata of the RDR
	EntryDate           string                   `xml:"r3d:entryData"`           // The date the RDR was indexed in re3data.org
	LastUpdate          string                   `xml:"r3d:lastUpdate"`          // The date the metadata of the RDR was updated

	// optional fields
	AdditionalName        StringLang             `xml:"r3d:additionalName,omitempty"`       // The full name of the RDR
	RepositoryIdentifiers []RepositoryIdentifier `xml:"r3d:repositoryIdentifier,omitempty"` // An identifier provisioned for the website of the RDR (wrapper element).
	Description           StringLang             `xml:"r3d:description,omitempty"`          // A textual description providing additional information about the RDR (primary language is English).
	Contact               []StringLang           `xml:"r3d:repositoryContact,omitempty"`    // Email address of the contact or an URL of an online contact form of the RDR.
	Keyword               []string               `xml:"r3d:keyword,omitempty"`              // English keyword(s) describing the subject focus of the RDR
}

func (repository Repository) InitNamespace() {
	repository.R3DType = "http://www.re3data.org/schema/3-0"
	repository.XsiType = "http://www.w3.org/2001/XMLSchema-instance"
	repository.XsiSchemaLocation = "http://www.re3data.org/schema/3-0 http://schema.re3data.org/3-0/re3dataV3-0.xsd"
}
