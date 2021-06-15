package myfair

type CoreTitleType string

const (
	TitleTypeAlternativeTitle CoreTitleType = "AlternativeTitle"
	TitleTypeSubTitle         CoreTitleType = "Subtitle"
	TitleTypeTranslatedTitle  CoreTitleType = "TranslatedTitle"
	TitleTypeOther            CoreTitleType = "Other"
)

type NameType string

const (
	NameTypeDefault        NameType = ""
	NameTypeOrganizational NameType = "Organizational"
	NameTypePersonal       NameType = "Personal"
)

var NameTypeReverse = map[string]NameType{
	"":               NameTypeDefault,
	"Organizational": NameTypeOrganizational,
	"Personal":       NameTypePersonal,
}

type PersonType string

const (
	PersonTypeAuthor                PersonType = "Author"
	PersonTypeArtist                PersonType = "Artist"
	PersonTypeContactPerson         PersonType = "ContactPerson"
	PersonTypeDataCollector         PersonType = "DataCollector"
	PersonTypeDataCurator           PersonType = "DataCurator"
	PersonTypeDataManager           PersonType = "DataManager"
	PersonTypeDistributor           PersonType = "Distributor"
	PersonTypeEditor                PersonType = "Editor"
	PersonTypeHostingInstitution    PersonType = "HostingInstitution"
	PersonTypeOther                 PersonType = "Other"
	PersonTypeProducer              PersonType = "Producer"
	PersonTypeProjectLeader         PersonType = "ProjectLeader"
	PersonTypeProjectManager        PersonType = "ProjectManager"
	PersonTypeProjectMember         PersonType = "ProjectMember"
	PersonTypeRegistrationAgency    PersonType = "RegistrationAgency"
	PersonTypeRegistrationAuthority PersonType = "RegistrationAuthority"
	PersonTypeRelatedPerson         PersonType = "RelatedPerson"
	PersonTypeResearchGroup         PersonType = "ResearchGroup"
	PersonTypeRightsHolder          PersonType = "RightsHolder"
	PersonTypeResearcher            PersonType = "Researcher"
	PersonTypeSponsor               PersonType = "Sponsor"
	PersonTypeSupervisor            PersonType = "Supervisor"
	PersonTypeWorkPackageLeader     PersonType = "WorkPackageLeader"
)

var PersonTypeReverse = map[string]PersonType{
	string(PersonTypeAuthor):                PersonTypeAuthor,
	string(PersonTypeArtist):                PersonTypeArtist,
	string(PersonTypeContactPerson):         PersonTypeContactPerson,
	string(PersonTypeDataCollector):         PersonTypeDataCollector,
	string(PersonTypeDataCurator):           PersonTypeDataCurator,
	string(PersonTypeDataManager):           PersonTypeDataManager,
	string(PersonTypeDistributor):           PersonTypeDistributor,
	string(PersonTypeEditor):                PersonTypeEditor,
	string(PersonTypeHostingInstitution):    PersonTypeHostingInstitution,
	string(PersonTypeOther):                 PersonTypeOther,
	string(PersonTypeProducer):              PersonTypeProducer,
	string(PersonTypeProjectLeader):         PersonTypeProjectLeader,
	string(PersonTypeProjectManager):        PersonTypeProjectManager,
	string(PersonTypeProjectMember):         PersonTypeProjectMember,
	string(PersonTypeRegistrationAgency):    PersonTypeRegistrationAgency,
	string(PersonTypeRegistrationAuthority): PersonTypeRegistrationAuthority,
	string(PersonTypeRelatedPerson):         PersonTypeRelatedPerson,
	string(PersonTypeResearchGroup):         PersonTypeResearchGroup,
	string(PersonTypeRightsHolder):          PersonTypeRightsHolder,
	string(PersonTypeResearcher):            PersonTypeResearcher,
	string(PersonTypeSponsor):               PersonTypeSponsor,
	string(PersonTypeSupervisor):            PersonTypeSupervisor,
	string(PersonTypeWorkPackageLeader):     PersonTypeWorkPackageLeader,
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
	ResourceTypeBook                ResourceType = "book"
	ResourceTypeBookSection         ResourceType = "bookSection"
	ResourceTypeThesis              ResourceType = "thesis"
	ResourceTypeJournalArticle      ResourceType = "journalArticle"
	ResourceTypeMagazineArticle     ResourceType = "magazineArticle"
	ResourceTypeOnlineResource      ResourceType = "onlineResource"
	ResourceTypeReport              ResourceType = "report"
	ResourceTypeWebpage             ResourceType = "webpage"
	ResourceTypeConferencePaper     ResourceType = "conferencePaper"
	ResourceTypePatent              ResourceType = "patent"
	ResourceTypeNote                ResourceType = "note"
	ResourceTypeArtisticPerformance ResourceType = "artisticPerformance"
	ResourceTypeDataset             ResourceType = "dataset"
	ResourceTypePresentation        ResourceType = "presentation"
	ResourceTypePhysicalObject      ResourceType = "physicalObject"
	ResourceTypeComputerProgram     ResourceType = "computerProgram"
	ResourceTypeOther               ResourceType = "other"
	ResourceTypeArtwork             ResourceType = "artwork"
	ResourceTypeAttachment          ResourceType = "attachment"
	ResourceTypeAudioRecording      ResourceType = "audioRecording"
	ResourceTypeDocument            ResourceType = "document"
	ResourceTypeEmail               ResourceType = "email"
	ResourceTypeEncyclopediaArticle ResourceType = "encyclopediaArticle"
	ResourceTypeFilm                ResourceType = "film"
	ResourceTypeInstantMessage      ResourceType = "instantMessage"
	ResourceTypeInterview           ResourceType = "interview"
	ResourceTypeLetter              ResourceType = "letter"
	ResourceTypeManuscript          ResourceType = "manuscript"
	ResourceTypeMap                 ResourceType = "map"
	ResourceTypeNewspaperArticle    ResourceType = "newspaperArticle"
	ResourceTypePodcast             ResourceType = "podcast"
	ResourceTypeRadioBroadcast      ResourceType = "radioBroadcast"
	ResourceTypeTvBroadcast         ResourceType = "tvBroadcast"
	ResourceTypeVideoRecording      ResourceType = "videoRecording"
)

var ResourceTypeReverse = map[string]ResourceType{
	string(ResourceTypeBook):                ResourceTypeBook,
	string(ResourceTypeBookSection):         ResourceTypeBookSection,
	string(ResourceTypeThesis):              ResourceTypeThesis,
	string(ResourceTypeJournalArticle):      ResourceTypeJournalArticle,
	string(ResourceTypeMagazineArticle):     ResourceTypeMagazineArticle,
	string(ResourceTypeOnlineResource):      ResourceTypeOnlineResource,
	string(ResourceTypeReport):              ResourceTypeReport,
	string(ResourceTypeWebpage):             ResourceTypeWebpage,
	string(ResourceTypeConferencePaper):     ResourceTypeConferencePaper,
	string(ResourceTypePatent):              ResourceTypePatent,
	string(ResourceTypeNote):                ResourceTypeNote,
	string(ResourceTypeArtisticPerformance): ResourceTypeArtisticPerformance,
	string(ResourceTypeDataset):             ResourceTypeDataset,
	string(ResourceTypePresentation):        ResourceTypePresentation,
	string(ResourceTypePhysicalObject):      ResourceTypePhysicalObject,
	string(ResourceTypeComputerProgram):     ResourceTypeComputerProgram,
	string(ResourceTypeOther):               ResourceTypeOther,
	string(ResourceTypeArtwork):             ResourceTypeArtwork,
	string(ResourceTypeAttachment):          ResourceTypeAttachment,
	string(ResourceTypeAudioRecording):      ResourceTypeAudioRecording,
	string(ResourceTypeDocument):            ResourceTypeDocument,
	string(ResourceTypeEmail):               ResourceTypeEmail,
	string(ResourceTypeEncyclopediaArticle): ResourceTypeEncyclopediaArticle,
	string(ResourceTypeFilm):                ResourceTypeFilm,
	string(ResourceTypeInstantMessage):      ResourceTypeInstantMessage,
	string(ResourceTypeInterview):           ResourceTypeInterview,
	string(ResourceTypeLetter):              ResourceTypeLetter,
	string(ResourceTypeManuscript):          ResourceTypeManuscript,
	string(ResourceTypeMap):                 ResourceTypeMap,
	string(ResourceTypeNewspaperArticle):    ResourceTypeNewspaperArticle,
	string(ResourceTypePodcast):             ResourceTypePodcast,
	string(ResourceTypeRadioBroadcast):      ResourceTypeRadioBroadcast,
	string(ResourceTypeTvBroadcast):         ResourceTypeTvBroadcast,
	string(ResourceTypeVideoRecording):      ResourceTypeVideoRecording,
}
