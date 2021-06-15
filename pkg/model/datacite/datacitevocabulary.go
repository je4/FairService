package datacite

type TitleType string

const TitleTypeDefaultTitle TitleType = ""
const TitleTypeAlternativeTitle TitleType = "AlternativeTitle"
const TitleTypeSubTitle TitleType = "Subtitle"
const TitleTypeTranslatedTitle TitleType = "TranslatedTitle"
const TitleTypeOther TitleType = "Other"

var TitleTypeReverse = map[string]TitleType{
	"":                 TitleTypeDefaultTitle,
	"AlternativeTitle": TitleTypeAlternativeTitle,
	"Subtitle":         TitleTypeSubTitle,
	"TranslatedTitle":  TitleTypeTranslatedTitle,
	"Other":            TitleTypeOther,
}

type NameType string

const NameTypeDefault NameType = ""
const NameTypeOrganizational NameType = "Organizational"
const NameTypePersonal NameType = "Personal"

var NameTypeReverse = map[string]NameType{
	"":               NameTypeDefault,
	"Organizational": NameTypeOrganizational,
	"Personal":       NameTypePersonal,
}

type RelatedIdentifierType string

const (
	RelatedIdentifierTypeARK     RelatedIdentifierType = "ARK"
	RelatedIdentifierTypeArXiv   RelatedIdentifierType = "arXiv"
	RelatedIdentifierTypeBibcode RelatedIdentifierType = "bibcode"
	RelatedIdentifierTypeDOI     RelatedIdentifierType = "DOI"
	RelatedIdentifierTypeEAN13   RelatedIdentifierType = "EAN13"
	RelatedIdentifierTypeEISSN   RelatedIdentifierType = "EISSN"
	RelatedIdentifierTypeHandle  RelatedIdentifierType = "Handle"
	RelatedIdentifierTypeIGSN    RelatedIdentifierType = "IGSN"
	RelatedIdentifierTypeISBN    RelatedIdentifierType = "ISBN"
	RelatedIdentifierTypeISSN    RelatedIdentifierType = "ISSN"
	RelatedIdentifierTypeISTC    RelatedIdentifierType = "ISTC"
	RelatedIdentifierTypeLISSN   RelatedIdentifierType = "LISSN"
	RelatedIdentifierTypeLSID    RelatedIdentifierType = "LSID"
	RelatedIdentifierTypePMID    RelatedIdentifierType = "PMID"
	RelatedIdentifierTypePURL    RelatedIdentifierType = "PURL"
	RelatedIdentifierTypeUPC     RelatedIdentifierType = "UPC"
	RelatedIdentifierTypeURL     RelatedIdentifierType = "URL"
	RelatedIdentifierTypeURN     RelatedIdentifierType = "URN"
	RelatedIdentifierTypeW3id    RelatedIdentifierType = "w3id"
)

var RelatedIdentifierTypeReverse = map[string]RelatedIdentifierType{
	"ARK":     RelatedIdentifierTypeARK,
	"arXiv":   RelatedIdentifierTypeArXiv,
	"bibcode": RelatedIdentifierTypeBibcode,
	"DOI":     RelatedIdentifierTypeDOI,
	"EAN13":   RelatedIdentifierTypeEAN13,
	"EISSN":   RelatedIdentifierTypeEISSN,
	"Handle":  RelatedIdentifierTypeHandle,
	"IGSN":    RelatedIdentifierTypeIGSN,
	"ISBN":    RelatedIdentifierTypeISBN,
	"ISSN":    RelatedIdentifierTypeISSN,
	"ISTC":    RelatedIdentifierTypeISTC,
	"LISSN":   RelatedIdentifierTypeLISSN,
	"LSID":    RelatedIdentifierTypeLSID,
	"PMID":    RelatedIdentifierTypePMID,
	"PURL":    RelatedIdentifierTypePURL,
	"UPC":     RelatedIdentifierTypeUPC,
	"URL":     RelatedIdentifierTypeURL,
	"URN":     RelatedIdentifierTypeURN,
	"w3id":    RelatedIdentifierTypeW3id,
}

type ResourceTypeGeneral string

const (
	ResourceTypeGeneralAudiovisual           ResourceTypeGeneral = "Audiovisual"
	ResourceTypeGeneralBook                  ResourceTypeGeneral = "Book"
	ResourceTypeGeneralBookChapter           ResourceTypeGeneral = "BookChapter"
	ResourceTypeGeneralCollection            ResourceTypeGeneral = "Collection"
	ResourceTypeGeneralComputationalNotebook ResourceTypeGeneral = "ComputationalNotebook"
	ResourceTypeGeneralConferencePaper       ResourceTypeGeneral = "ConferencePaper"
	ResourceTypeGeneralConferenceProceeding  ResourceTypeGeneral = "ConferenceProceeding"
	ResourceTypeGeneralDataPaper             ResourceTypeGeneral = "DataPaper"
	ResourceTypeGeneralDataset               ResourceTypeGeneral = "Dataset"
	ResourceTypeGeneralDissertation          ResourceTypeGeneral = "Dissertation"
	ResourceTypeGeneralEvent                 ResourceTypeGeneral = "Event"
	ResourceTypeGeneralImage                 ResourceTypeGeneral = "Image"
	ResourceTypeGeneralInteractiveResource   ResourceTypeGeneral = "InteractiveResource"
	ResourceTypeGeneralJournal               ResourceTypeGeneral = "Journal"
	ResourceTypeGeneralJournalArticle        ResourceTypeGeneral = "JournalArticle"
	ResourceTypeGeneralModel                 ResourceTypeGeneral = "Model"
	ResourceTypeGeneralOutputManagementPlan  ResourceTypeGeneral = "OutputManagementPlan"
	ResourceTypeGeneralPeerReview            ResourceTypeGeneral = "PeerReview"
	ResourceTypeGeneralPhysicalObject        ResourceTypeGeneral = "PhysicalObject"
	ResourceTypeGeneralPreprint              ResourceTypeGeneral = "Preprint"
	ResourceTypeGeneralReport                ResourceTypeGeneral = "Report"
	ResourceTypeGeneralService               ResourceTypeGeneral = "Service"
	ResourceTypeGeneralSoftware              ResourceTypeGeneral = "Software"
	ResourceTypeGeneralSound                 ResourceTypeGeneral = "Sound"
	ResourceTypeGeneralStandard              ResourceTypeGeneral = "Standard"
	ResourceTypeGeneralText                  ResourceTypeGeneral = "Text"
	ResourceTypeGeneralWorkflow              ResourceTypeGeneral = "Workflow"
	ResourceTypeGeneralOther                 ResourceTypeGeneral = "Other"
)

var ResourceTypeGeneralReverse = map[string]ResourceTypeGeneral{
	"Audiovisual":           ResourceTypeGeneralAudiovisual,
	"Book":                  ResourceTypeGeneralBook,
	"BookChapter":           ResourceTypeGeneralBookChapter,
	"Collection":            ResourceTypeGeneralCollection,
	"ComputationalNotebook": ResourceTypeGeneralComputationalNotebook,
	"ConferencePaper":       ResourceTypeGeneralConferencePaper,
	"ConferenceProceeding":  ResourceTypeGeneralConferenceProceeding,
	"DataPaper":             ResourceTypeGeneralDataPaper,
	"Dataset":               ResourceTypeGeneralDataset,
	"Dissertation":          ResourceTypeGeneralDissertation,
	"Event":                 ResourceTypeGeneralEvent,
	"Image":                 ResourceTypeGeneralImage,
	"InteractiveResource":   ResourceTypeGeneralInteractiveResource,
	"Journal":               ResourceTypeGeneralJournal,
	"JournalArticle":        ResourceTypeGeneralJournalArticle,
	"Model":                 ResourceTypeGeneralModel,
	"OutputManagementPlan":  ResourceTypeGeneralOutputManagementPlan,
	"PeerReview":            ResourceTypeGeneralPeerReview,
	"PhysicalObject":        ResourceTypeGeneralPhysicalObject,
	"Preprint":              ResourceTypeGeneralPreprint,
	"Report":                ResourceTypeGeneralReport,
	"Service":               ResourceTypeGeneralService,
	"Software":              ResourceTypeGeneralSoftware,
	"Sound":                 ResourceTypeGeneralSound,
	"Standard":              ResourceTypeGeneralStandard,
	"Text":                  ResourceTypeGeneralText,
	"Workflow":              ResourceTypeGeneralWorkflow,
	"Other":                 ResourceTypeGeneralOther,
}

type ContributorType string

const (
	ContributorTypeContactPerson         ContributorType = "ContactPerson"
	ContributorTypeDataCollector         ContributorType = "DataCollector"
	ContributorTypeDataCurator           ContributorType = "DataCurator"
	ContributorTypeDataManager           ContributorType = "DataManager"
	ContributorTypeDistributor           ContributorType = "Distributor"
	ContributorTypeEditor                ContributorType = "Editor"
	ContributorTypeHostingInstitution    ContributorType = "HostingInstitution"
	ContributorTypeOther                 ContributorType = "Other"
	ContributorTypeProducer              ContributorType = "Producer"
	ContributorTypeProjectLeader         ContributorType = "ProjectLeader"
	ContributorTypeProjectManager        ContributorType = "ProjectManager"
	ContributorTypeProjectMember         ContributorType = "ProjectMember"
	ContributorTypeRegistrationAgency    ContributorType = "RegistrationAgency"
	ContributorTypeRegistrationAuthority ContributorType = "RegistrationAuthority"
	ContributorTypeRelatedPerson         ContributorType = "RelatedPerson"
	ContributorTypeResearchGroup         ContributorType = "ResearchGroup"
	ContributorTypeRightsHolder          ContributorType = "RightsHolder"
	ContributorTypeResearcher            ContributorType = "Researcher"
	ContributorTypeSponsor               ContributorType = "Sponsor"
	ContributorTypeSupervisor            ContributorType = "Supervisor"
	ContributorTypeWorkPackageLeader     ContributorType = "WorkPackageLeader"
)

var ContributorTypeReverse = map[string]ContributorType{
	"ContactPerson":         ContributorTypeContactPerson,
	"DataCollector":         ContributorTypeDataCollector,
	"DataCurator":           ContributorTypeDataCurator,
	"DataManager":           ContributorTypeDataManager,
	"Distributor":           ContributorTypeDistributor,
	"Editor":                ContributorTypeEditor,
	"HostingInstitution":    ContributorTypeHostingInstitution,
	"Other":                 ContributorTypeOther,
	"Producer":              ContributorTypeProducer,
	"ProjectLeader":         ContributorTypeProjectLeader,
	"ProjectManager":        ContributorTypeProjectManager,
	"ProjectMember":         ContributorTypeProjectMember,
	"RegistrationAgency":    ContributorTypeRegistrationAgency,
	"RegistrationAuthority": ContributorTypeRegistrationAuthority,
	"RelatedPerson":         ContributorTypeRelatedPerson,
	"ResearchGroup":         ContributorTypeResearchGroup,
	"RightsHolder":          ContributorTypeRightsHolder,
	"Researcher":            ContributorTypeResearcher,
	"Sponsor":               ContributorTypeSponsor,
	"Supervisor":            ContributorTypeSupervisor,
	"WorkPackageLeader":     ContributorTypeWorkPackageLeader,
}
