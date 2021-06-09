package datacite

type DataCiteTitleType string

const DataCiteTitleTypeDefaultTitle DataCiteTitleType = ""
const DataCiteTitleTypeAlternativeTitle DataCiteTitleType = "AlternativeTitle"
const DataCiteTitleTypeSubTitle DataCiteTitleType = "Subtitle"
const DataCiteTitleTypeTranslatedTitle DataCiteTitleType = "TranslatedTitle"
const DataCiteTitleTypeOther DataCiteTitleType = "Other"

var DataCiteTitleTypeReverse = map[string]DataCiteTitleType{
	"":                 DataCiteTitleTypeDefaultTitle,
	"AlternativeTitle": DataCiteTitleTypeAlternativeTitle,
	"Subtitle":         DataCiteTitleTypeSubTitle,
	"TranslatedTitle":  DataCiteTitleTypeTranslatedTitle,
	"Other":            DataCiteTitleTypeOther,
}

type DataCiteNameType string

const DataCiteNameTypeDefault DataCiteNameType = ""
const DataCiteNameTypeOrganizational DataCiteNameType = "Organizational"
const DataCiteNameTypePersonal DataCiteNameType = "Personal"

var DataCiteNameTypeReverse = map[string]DataCiteNameType{
	"":               DataCiteNameTypeDefault,
	"Organizational": DataCiteNameTypeOrganizational,
	"Personal":       DataCiteNameTypePersonal,
}

type DataCiteRelatedIdentifierType string

const (
	DataCiteRelatedIdentifierTypeARK     DataCiteRelatedIdentifierType = "ARK"
	DataCiteRelatedIdentifierTypeArXiv   DataCiteRelatedIdentifierType = "arXiv"
	DataCiteRelatedIdentifierTypeBibcode DataCiteRelatedIdentifierType = "bibcode"
	DataCiteRelatedIdentifierTypeDOI     DataCiteRelatedIdentifierType = "DOI"
	DataCiteRelatedIdentifierTypeEAN13   DataCiteRelatedIdentifierType = "EAN13"
	DataCiteRelatedIdentifierTypeEISSN   DataCiteRelatedIdentifierType = "EISSN"
	DataCiteRelatedIdentifierTypeHandle  DataCiteRelatedIdentifierType = "Handle"
	DataCiteRelatedIdentifierTypeIGSN    DataCiteRelatedIdentifierType = "IGSN"
	DataCiteRelatedIdentifierTypeISBN    DataCiteRelatedIdentifierType = "ISBN"
	DataCiteRelatedIdentifierTypeISSN    DataCiteRelatedIdentifierType = "ISSN"
	DataCiteRelatedIdentifierTypeISTC    DataCiteRelatedIdentifierType = "ISTC"
	DataCiteRelatedIdentifierTypeLISSN   DataCiteRelatedIdentifierType = "LISSN"
	DataCiteRelatedIdentifierTypeLSID    DataCiteRelatedIdentifierType = "LSID"
	DataCiteRelatedIdentifierTypePMID    DataCiteRelatedIdentifierType = "PMID"
	DataCiteRelatedIdentifierTypePURL    DataCiteRelatedIdentifierType = "PURL"
	DataCiteRelatedIdentifierTypeUPC     DataCiteRelatedIdentifierType = "UPC"
	DataCiteRelatedIdentifierTypeURL     DataCiteRelatedIdentifierType = "URL"
	DataCiteRelatedIdentifierTypeURN     DataCiteRelatedIdentifierType = "URN"
	DataCiteRelatedIdentifierTypeW3id    DataCiteRelatedIdentifierType = "w3id"
)

var DataCiteRelatedIdentifierTypeReverse = map[string]DataCiteRelatedIdentifierType{
	"ARK":     DataCiteRelatedIdentifierTypeARK,
	"arXiv":   DataCiteRelatedIdentifierTypeArXiv,
	"bibcode": DataCiteRelatedIdentifierTypeBibcode,
	"DOI":     DataCiteRelatedIdentifierTypeDOI,
	"EAN13":   DataCiteRelatedIdentifierTypeEAN13,
	"EISSN":   DataCiteRelatedIdentifierTypeEISSN,
	"Handle":  DataCiteRelatedIdentifierTypeHandle,
	"IGSN":    DataCiteRelatedIdentifierTypeIGSN,
	"ISBN":    DataCiteRelatedIdentifierTypeISBN,
	"ISSN":    DataCiteRelatedIdentifierTypeISSN,
	"ISTC":    DataCiteRelatedIdentifierTypeISTC,
	"LISSN":   DataCiteRelatedIdentifierTypeLISSN,
	"LSID":    DataCiteRelatedIdentifierTypeLSID,
	"PMID":    DataCiteRelatedIdentifierTypePMID,
	"PURL":    DataCiteRelatedIdentifierTypePURL,
	"UPC":     DataCiteRelatedIdentifierTypeUPC,
	"URL":     DataCiteRelatedIdentifierTypeURL,
	"URN":     DataCiteRelatedIdentifierTypeURN,
	"w3id":    DataCiteRelatedIdentifierTypeW3id,
}

type DataCiteResourceTypeGeneral string

const (
	DataCiteResourceTypeGeneralAudiovisual           DataCiteResourceTypeGeneral = "Audiovisual"
	DataCiteResourceTypeGeneralBook                  DataCiteResourceTypeGeneral = "Book"
	DataCiteResourceTypeGeneralBookChapter           DataCiteResourceTypeGeneral = "BookChapter"
	DataCiteResourceTypeGeneralCollection            DataCiteResourceTypeGeneral = "Collection"
	DataCiteResourceTypeGeneralComputationalNotebook DataCiteResourceTypeGeneral = "ComputationalNotebook"
	DataCiteResourceTypeGeneralConferencePaper       DataCiteResourceTypeGeneral = "ConferencePaper"
	DataCiteResourceTypeGeneralConferenceProceeding  DataCiteResourceTypeGeneral = "ConferenceProceeding"
	DataCiteResourceTypeGeneralDataPaper             DataCiteResourceTypeGeneral = "DataPaper"
	DataCiteResourceTypeGeneralDataset               DataCiteResourceTypeGeneral = "Dataset"
	DataCiteResourceTypeGeneralDissertation          DataCiteResourceTypeGeneral = "Dissertation"
	DataCiteResourceTypeGeneralEvent                 DataCiteResourceTypeGeneral = "Event"
	DataCiteResourceTypeGeneralImage                 DataCiteResourceTypeGeneral = "Image"
	DataCiteResourceTypeGeneralInteractiveResource   DataCiteResourceTypeGeneral = "InteractiveResource"
	DataCiteResourceTypeGeneralJournal               DataCiteResourceTypeGeneral = "Journal"
	DataCiteResourceTypeGeneralJournalArticle        DataCiteResourceTypeGeneral = "JournalArticle"
	DataCiteResourceTypeGeneralModel                 DataCiteResourceTypeGeneral = "Model"
	DataCiteResourceTypeGeneralOutputManagementPlan  DataCiteResourceTypeGeneral = "OutputManagementPlan"
	DataCiteResourceTypeGeneralPeerReview            DataCiteResourceTypeGeneral = "PeerReview"
	DataCiteResourceTypeGeneralPhysicalObject        DataCiteResourceTypeGeneral = "PhysicalObject"
	DataCiteResourceTypeGeneralPreprint              DataCiteResourceTypeGeneral = "Preprint"
	DataCiteResourceTypeGeneralReport                DataCiteResourceTypeGeneral = "Report"
	DataCiteResourceTypeGeneralService               DataCiteResourceTypeGeneral = "Service"
	DataCiteResourceTypeGeneralSoftware              DataCiteResourceTypeGeneral = "Software"
	DataCiteResourceTypeGeneralSound                 DataCiteResourceTypeGeneral = "Sound"
	DataCiteResourceTypeGeneralStandard              DataCiteResourceTypeGeneral = "Standard"
	DataCiteResourceTypeGeneralText                  DataCiteResourceTypeGeneral = "Text"
	DataCiteResourceTypeGeneralWorkflow              DataCiteResourceTypeGeneral = "Workflow"
	DataCiteResourceTypeGeneralOther                 DataCiteResourceTypeGeneral = "Other"
)

var DataCiteResourceTypeGeneralReverse = map[string]DataCiteResourceTypeGeneral{
	"Audiovisual":           DataCiteResourceTypeGeneralAudiovisual,
	"Book":                  DataCiteResourceTypeGeneralBook,
	"BookChapter":           DataCiteResourceTypeGeneralBookChapter,
	"Collection":            DataCiteResourceTypeGeneralCollection,
	"ComputationalNotebook": DataCiteResourceTypeGeneralComputationalNotebook,
	"ConferencePaper":       DataCiteResourceTypeGeneralConferencePaper,
	"ConferenceProceeding":  DataCiteResourceTypeGeneralConferenceProceeding,
	"DataPaper":             DataCiteResourceTypeGeneralDataPaper,
	"Dataset":               DataCiteResourceTypeGeneralDataset,
	"Dissertation":          DataCiteResourceTypeGeneralDissertation,
	"Event":                 DataCiteResourceTypeGeneralEvent,
	"Image":                 DataCiteResourceTypeGeneralImage,
	"InteractiveResource":   DataCiteResourceTypeGeneralInteractiveResource,
	"Journal":               DataCiteResourceTypeGeneralJournal,
	"JournalArticle":        DataCiteResourceTypeGeneralJournalArticle,
	"Model":                 DataCiteResourceTypeGeneralModel,
	"OutputManagementPlan":  DataCiteResourceTypeGeneralOutputManagementPlan,
	"PeerReview":            DataCiteResourceTypeGeneralPeerReview,
	"PhysicalObject":        DataCiteResourceTypeGeneralPhysicalObject,
	"Preprint":              DataCiteResourceTypeGeneralPreprint,
	"Report":                DataCiteResourceTypeGeneralReport,
	"Service":               DataCiteResourceTypeGeneralService,
	"Software":              DataCiteResourceTypeGeneralSoftware,
	"Sound":                 DataCiteResourceTypeGeneralSound,
	"Standard":              DataCiteResourceTypeGeneralStandard,
	"Text":                  DataCiteResourceTypeGeneralText,
	"Workflow":              DataCiteResourceTypeGeneralWorkflow,
	"Other":                 DataCiteResourceTypeGeneralOther,
}

type DataCiteContributorType string

const (
	DataCiteContributorTypeContactPerson         DataCiteContributorType = "ContactPerson"
	DataCiteContributorTypeDataCollector         DataCiteContributorType = "DataCollector"
	DataCiteContributorTypeDataCurator           DataCiteContributorType = "DataCurator"
	DataCiteContributorTypeDataManager           DataCiteContributorType = "DataManager"
	DataCiteContributorTypeDistributor           DataCiteContributorType = "Distributor"
	DataCiteContributorTypeEditor                DataCiteContributorType = "Editor"
	DataCiteContributorTypeHostingInstitution    DataCiteContributorType = "HostingInstitution"
	DataCiteContributorTypeOther                 DataCiteContributorType = "Other"
	DataCiteContributorTypeProducer              DataCiteContributorType = "Producer"
	DataCiteContributorTypeProjectLeader         DataCiteContributorType = "ProjectLeader"
	DataCiteContributorTypeProjectManager        DataCiteContributorType = "ProjectManager"
	DataCiteContributorTypeProjectMember         DataCiteContributorType = "ProjectMember"
	DataCiteContributorTypeRegistrationAgency    DataCiteContributorType = "RegistrationAgency"
	DataCiteContributorTypeRegistrationAuthority DataCiteContributorType = "RegistrationAuthority"
	DataCiteContributorTypeRelatedPerson         DataCiteContributorType = "RelatedPerson"
	DataCiteContributorTypeResearchGroup         DataCiteContributorType = "ResearchGroup"
	DataCiteContributorTypeRightsHolder          DataCiteContributorType = "RightsHolder"
	DataCiteContributorTypeResearcher            DataCiteContributorType = "Researcher"
	DataCiteContributorTypeSponsor               DataCiteContributorType = "Sponsor"
	DataCiteContributorTypeSupervisor            DataCiteContributorType = "Supervisor"
	DataCiteContributorTypeWorkPackageLeader     DataCiteContributorType = "WorkPackageLeader"
)

var DataCiteContributorTypeReverse = map[string]DataCiteContributorType{
	"ContactPerson":         DataCiteContributorTypeContactPerson,
	"DataCollector":         DataCiteContributorTypeDataCollector,
	"DataCurator":           DataCiteContributorTypeDataCurator,
	"DataManager":           DataCiteContributorTypeDataManager,
	"Distributor":           DataCiteContributorTypeDistributor,
	"Editor":                DataCiteContributorTypeEditor,
	"HostingInstitution":    DataCiteContributorTypeHostingInstitution,
	"Other":                 DataCiteContributorTypeOther,
	"Producer":              DataCiteContributorTypeProducer,
	"ProjectLeader":         DataCiteContributorTypeProjectLeader,
	"ProjectManager":        DataCiteContributorTypeProjectManager,
	"ProjectMember":         DataCiteContributorTypeProjectMember,
	"RegistrationAgency":    DataCiteContributorTypeRegistrationAgency,
	"RegistrationAuthority": DataCiteContributorTypeRegistrationAuthority,
	"RelatedPerson":         DataCiteContributorTypeRelatedPerson,
	"ResearchGroup":         DataCiteContributorTypeResearchGroup,
	"RightsHolder":          DataCiteContributorTypeRightsHolder,
	"Researcher":            DataCiteContributorTypeResearcher,
	"Sponsor":               DataCiteContributorTypeSponsor,
	"Supervisor":            DataCiteContributorTypeSupervisor,
	"WorkPackageLeader":     DataCiteContributorTypeWorkPackageLeader,
}
