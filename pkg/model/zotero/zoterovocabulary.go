package zotero

type ItemType string

const (
	ItemTypeArtwork             ItemType = "Artwork"
	ItemTypeAttachment          ItemType = "Attachment"
	ItemTypeAudioRecording      ItemType = "Audio Recording"
	ItemTypeBill                ItemType = "Bill"
	ItemTypeBlogPost            ItemType = "Blog Post"
	ItemTypeBook                ItemType = "Book"
	ItemTypeBookSection         ItemType = "Book Section"
	ItemTypeCase                ItemType = "Case"
	ItemTypeComputerProgram     ItemType = "Computer Program"
	ItemTypeConferencePaper     ItemType = "Conference Paper"
	ItemTypeDictionaryEntry     ItemType = "Dictionary Entry"
	ItemTypeDocument            ItemType = "Document"
	ItemTypeEmail               ItemType = "E-mail"
	ItemTypeEncyclopediaArticle ItemType = "Encyclopedia Article"
	ItemTypeFilm                ItemType = "Film"
	ItemTypeForumPost           ItemType = "Forum Post"
	ItemTypeHearing             ItemType = "Hearing"
	ItemTypeInstantMessage      ItemType = "Instant Message"
	ItemTypeInterview           ItemType = "Interview"
	ItemTypeJournalArticle      ItemType = "Journal Article"
	ItemTypeLetter              ItemType = "Letter"
	ItemTypeMagazineArticle     ItemType = "Magazine Article"
	ItemTypeManuscript          ItemType = "Manuscript"
	ItemTypeMap                 ItemType = "Map"
	ItemTypeNewspaperArticle    ItemType = "Newspaper Article"
	ItemTypeNote                ItemType = "Note"
	ItemTypePatent              ItemType = "Patent"
	ItemTypePodcast             ItemType = "Podcast"
	ItemTypePresentation        ItemType = "Presentation"
	ItemTypeRadioBroadcast      ItemType = "Radio Broadcast"
	ItemTypeReport              ItemType = "Report"
	ItemTypeStatute             ItemType = "Statute"
	ItemTypeThesis              ItemType = "Thesis"
	ItemTypeTvBroadcast         ItemType = "TV Broadcast"
	ItemTypeVideoRecording      ItemType = "Video Recording"
	ItemTypeWebpage             ItemType = "Web Page"
)

var ItemTypeReverse = map[string]ItemType{
	"artwork":             ItemTypeArtwork,
	"attachment":          ItemTypeAttachment,
	"audioRecording":      ItemTypeAudioRecording,
	"bill":                ItemTypeBill,
	"blogPost":            ItemTypeBlogPost,
	"book":                ItemTypeBook,
	"bookSection":         ItemTypeBookSection,
	"case":                ItemTypeCase,
	"computerProgram":     ItemTypeComputerProgram,
	"conferencePaper":     ItemTypeConferencePaper,
	"dictionaryEntry":     ItemTypeDictionaryEntry,
	"document":            ItemTypeDocument,
	"email":               ItemTypeEmail,
	"encyclopediaArticle": ItemTypeEncyclopediaArticle,
	"film":                ItemTypeFilm,
	"forumPost":           ItemTypeForumPost,
	"hearing":             ItemTypeHearing,
	"instantMessage":      ItemTypeInstantMessage,
	"interview":           ItemTypeInterview,
	"journalArticle":      ItemTypeJournalArticle,
	"letter":              ItemTypeLetter,
	"magazineArticle":     ItemTypeMagazineArticle,
	"manuscript":          ItemTypeManuscript,
	"map":                 ItemTypeMap,
	"newspaperArticle":    ItemTypeNewspaperArticle,
	"note":                ItemTypeNote,
	"patent":              ItemTypePatent,
	"podcast":             ItemTypePodcast,
	"presentation":        ItemTypePresentation,
	"radioBroadcast":      ItemTypeRadioBroadcast,
	"report":              ItemTypeReport,
	"statute":             ItemTypeStatute,
	"thesis":              ItemTypeThesis,
	"tvBroadcast":         ItemTypeTvBroadcast,
	"videoRecording":      ItemTypeVideoRecording,
	"webpage":             ItemTypeWebpage,
}
