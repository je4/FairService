package myfair

type NameType string

const NameTypeDefault NameType = ""
const NameTypeOrganizational NameType = "Organizational"
const NameTypePersonal NameType = "Personal"

var NameTypeReverse = map[string]NameType{
	"":               NameTypeDefault,
	"Organizational": NameTypeOrganizational,
	"Personal":       NameTypePersonal,
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

type ResourceType string

const (
	TypeBook                ResourceType = "book"
	TypeBookSection         ResourceType = "bookSection"
	TypeThesis              ResourceType = "thesis"
	TypeJournalArticle      ResourceType = "journalArticle"
	TypeMagazineArticle     ResourceType = "magazineArticle"
	TypeOnlineResource      ResourceType = "onlineResource"
	TypeReport              ResourceType = "report"
	TypeWebpage             ResourceType = "webpage"
	TypeConferencePaper     ResourceType = "conferencePaper"
	TypePatent              ResourceType = "patent"
	TypeNote                ResourceType = "note"
	TypeArtisticPerformance ResourceType = "artisticPerformance"
	TypeDataset             ResourceType = "dataset"
	TypePresentation        ResourceType = "presentation"
	TypePhysicalObject      ResourceType = "physicalObject"
	TypeComputerProgram     ResourceType = "computerProgram"
	TypeOther               ResourceType = "other"
	TypeArtwork             ResourceType = "artwork"
	TypeAttachment          ResourceType = "attachment"
	TypeAudioRecording      ResourceType = "audioRecording"
	TypeDocument            ResourceType = "document"
	TypeEmail               ResourceType = "email"
	TypeEncyclopediaArticle ResourceType = "encyclopediaArticle"
	TypeFilm                ResourceType = "film"
	TypeInstantMessage      ResourceType = "instantMessage"
	TypeInterview           ResourceType = "interview"
	TypeLetter              ResourceType = "letter"
	TypeManuscript          ResourceType = "manuscript"
	TypeMap                 ResourceType = "map"
	TypeNewspaperArticle    ResourceType = "newspaperArticle"
	TypePodcast             ResourceType = "podcast"
	TypeRadioBroadcast      ResourceType = "radioBroadcast"
	TypeTvBroadcast         ResourceType = "tvBroadcast"
	TypeVideoRecording      ResourceType = "videoRecording"
)

var TypeReverse = map[string]ResourceType{
	string(TypeBook):                TypeBook,
	string(TypeBookSection):         TypeBookSection,
	string(TypeThesis):              TypeThesis,
	string(TypeJournalArticle):      TypeJournalArticle,
	string(TypeMagazineArticle):     TypeMagazineArticle,
	string(TypeOnlineResource):      TypeOnlineResource,
	string(TypeReport):              TypeReport,
	string(TypeWebpage):             TypeWebpage,
	string(TypeConferencePaper):     TypeConferencePaper,
	string(TypePatent):              TypePatent,
	string(TypeNote):                TypeNote,
	string(TypeArtisticPerformance): TypeArtisticPerformance,
	string(TypeDataset):             TypeDataset,
	string(TypePresentation):        TypePresentation,
	string(TypePhysicalObject):      TypePhysicalObject,
	string(TypeComputerProgram):     TypeComputerProgram,
	string(TypeOther):               TypeOther,
	string(TypeArtwork):             TypeArtwork,
	string(TypeAttachment):          TypeAttachment,
	string(TypeAudioRecording):      TypeAudioRecording,
	string(TypeDocument):            TypeDocument,
	string(TypeEmail):               TypeEmail,
	string(TypeEncyclopediaArticle): TypeEncyclopediaArticle,
	string(TypeFilm):                TypeFilm,
	string(TypeInstantMessage):      TypeInstantMessage,
	string(TypeInterview):           TypeInterview,
	string(TypeLetter):              TypeLetter,
	string(TypeManuscript):          TypeManuscript,
	string(TypeMap):                 TypeMap,
	string(TypeNewspaperArticle):    TypeNewspaperArticle,
	string(TypePodcast):             TypePodcast,
	string(TypeRadioBroadcast):      TypeRadioBroadcast,
	string(TypeTvBroadcast):         TypeTvBroadcast,
	string(TypeVideoRecording):      TypeVideoRecording,
}
