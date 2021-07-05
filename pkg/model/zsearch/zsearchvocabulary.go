package zsearch

type PersonRole string

const (
	PersonRoleArtist      PersonRole = "artist"
	PersonRoleContributor PersonRole = "contributor"
)

var CreatorTypeReverse = map[string]PersonRole{
	string(PersonRoleArtist):      PersonRoleArtist,
	string(PersonRoleContributor): PersonRoleContributor,
}

type ItemType string

const (
	ItemTypeBook                ItemType = "book"
	ItemTypeBookSection         ItemType = "bookSection"
	ItemTypeThesis              ItemType = "thesis"
	ItemTypeJournalArticle      ItemType = "journalArticle"
	ItemTypeMagazineArticle     ItemType = "magazineArticle"
	ItemTypeReport              ItemType = "report"
	ItemTypeWebpage             ItemType = "webpage"
	ItemTypeConferencePaper     ItemType = "conferencePaper"
	ItemTypePatent              ItemType = "patent"
	ItemTypeNote                ItemType = "note"
	ItemTypePresentation        ItemType = "presentation"
	ItemTypeComputerProgram     ItemType = "computerProgram"
	ItemTypeArtwork             ItemType = "artwork"
	ItemTypePerformance         ItemType = "performance"
	ItemTypeAttachment          ItemType = "attachment"
	ItemTypeAudioRecording      ItemType = "audioRecording"
	ItemTypeBill                ItemType = "bill"
	ItemTypeBlogPost            ItemType = "blogPost"
	ItemTypeCase                ItemType = "case"
	ItemTypeDictionaryEntry     ItemType = "dictionaryEntry"
	ItemTypeDocument            ItemType = "document"
	ItemTypeEmail               ItemType = "email"
	ItemTypeEncyclopediaArticle ItemType = "encyclopediaArticle"
	ItemTypeFilm                ItemType = "film"
	ItemTypeForumPost           ItemType = "forumPost"
	ItemTypeHearing             ItemType = "hearing"
	ItemTypeInstantMessage      ItemType = "instantMessage"
	ItemTypeInterview           ItemType = "interview"
	ItemTypeLetter              ItemType = "letter"
	ItemTypeManuscript          ItemType = "manuscript"
	ItemTypeMap                 ItemType = "map"
	ItemTypeNewspaperArticle    ItemType = "newspaperArticle"
	ItemTypePodcast             ItemType = "podcast"
	ItemTypeRadioBroadcast      ItemType = "radioBroadcast"
	ItemTypeStatute             ItemType = "statute"
	ItemTypeTvBroadcast         ItemType = "tvBroadcast"
	ItemTypeVideoRecording      ItemType = "videoRecording"
	ItemTypeOther               ItemType = "other"
)

var ItemTypeReverse = map[string]ItemType{
	string(ItemTypeBook):                ItemTypeBook,
	string(ItemTypeBookSection):         ItemTypeBookSection,
	string(ItemTypeThesis):              ItemTypeThesis,
	string(ItemTypeJournalArticle):      ItemTypeJournalArticle,
	string(ItemTypeMagazineArticle):     ItemTypeMagazineArticle,
	string(ItemTypeReport):              ItemTypeReport,
	string(ItemTypeWebpage):             ItemTypeWebpage,
	string(ItemTypeConferencePaper):     ItemTypeConferencePaper,
	string(ItemTypePatent):              ItemTypePatent,
	string(ItemTypeNote):                ItemTypeNote,
	string(ItemTypePresentation):        ItemTypePresentation,
	string(ItemTypeComputerProgram):     ItemTypeComputerProgram,
	string(ItemTypeArtwork):             ItemTypeArtwork,
	string(ItemTypePerformance):         ItemTypePerformance,
	string(ItemTypeAttachment):          ItemTypeAttachment,
	string(ItemTypeAudioRecording):      ItemTypeAudioRecording,
	string(ItemTypeBill):                ItemTypeBill,
	string(ItemTypeBlogPost):            ItemTypeBlogPost,
	string(ItemTypeCase):                ItemTypeCase,
	string(ItemTypeDictionaryEntry):     ItemTypeDictionaryEntry,
	string(ItemTypeDocument):            ItemTypeDocument,
	string(ItemTypeEmail):               ItemTypeEmail,
	string(ItemTypeEncyclopediaArticle): ItemTypeEncyclopediaArticle,
	string(ItemTypeFilm):                ItemTypeFilm,
	string(ItemTypeForumPost):           ItemTypeForumPost,
	string(ItemTypeHearing):             ItemTypeHearing,
	string(ItemTypeInstantMessage):      ItemTypeInstantMessage,
	string(ItemTypeInterview):           ItemTypeInterview,
	string(ItemTypeLetter):              ItemTypeLetter,
	string(ItemTypeManuscript):          ItemTypeManuscript,
	string(ItemTypeMap):                 ItemTypeMap,
	string(ItemTypeNewspaperArticle):    ItemTypeNewspaperArticle,
	string(ItemTypePodcast):             ItemTypePodcast,
	string(ItemTypeRadioBroadcast):      ItemTypeRadioBroadcast,
	string(ItemTypeStatute):             ItemTypeStatute,
	string(ItemTypeTvBroadcast):         ItemTypeTvBroadcast,
	string(ItemTypeVideoRecording):      ItemTypeVideoRecording,
	string(ItemTypeOther):               ItemTypeOther,
}
