package fairschema

type ZoteroItemType string

const (
	ZoteroItemTypeArtwork             ZoteroItemType = "Artwork"
	ZoteroItemTypeAttachment          ZoteroItemType = "Attachment"
	ZoteroItemTypeAudioRecording      ZoteroItemType = "Audio Recording"
	ZoteroItemTypeBill                ZoteroItemType = "Bill"
	ZoteroItemTypeBlogPost            ZoteroItemType = "Blog Post"
	ZoteroItemTypeBook                ZoteroItemType = "Book"
	ZoteroItemTypeBookSection         ZoteroItemType = "Book Section"
	ZoteroItemTypeCase                ZoteroItemType = "Case"
	ZoteroItemTypeComputerProgram     ZoteroItemType = "Computer Program"
	ZoteroItemTypeConferencePaper     ZoteroItemType = "Conference Paper"
	ZoteroItemTypeDictionaryEntry     ZoteroItemType = "Dictionary Entry"
	ZoteroItemTypeDocument            ZoteroItemType = "Document"
	ZoteroItemTypeEmail               ZoteroItemType = "E-mail"
	ZoteroItemTypeEncyclopediaArticle ZoteroItemType = "Encyclopedia Article"
	ZoteroItemTypeFilm                ZoteroItemType = "Film"
	ZoteroItemTypeForumPost           ZoteroItemType = "Forum Post"
	ZoteroItemTypeHearing             ZoteroItemType = "Hearing"
	ZoteroItemTypeInstantMessage      ZoteroItemType = "Instant Message"
	ZoteroItemTypeInterview           ZoteroItemType = "Interview"
	ZoteroItemTypeJournalArticle      ZoteroItemType = "Journal Article"
	ZoteroItemTypeLetter              ZoteroItemType = "Letter"
	ZoteroItemTypeMagazineArticle     ZoteroItemType = "Magazine Article"
	ZoteroItemTypeManuscript          ZoteroItemType = "Manuscript"
	ZoteroItemTypeMap                 ZoteroItemType = "Map"
	ZoteroItemTypeNewspaperArticle    ZoteroItemType = "Newspaper Article"
	ZoteroItemTypeNote                ZoteroItemType = "Note"
	ZoteroItemTypePatent              ZoteroItemType = "Patent"
	ZoteroItemTypePodcast             ZoteroItemType = "Podcast"
	ZoteroItemTypePresentation        ZoteroItemType = "Presentation"
	ZoteroItemTypeRadioBroadcast      ZoteroItemType = "Radio Broadcast"
	ZoteroItemTypeReport              ZoteroItemType = "Report"
	ZoteroItemTypeStatute             ZoteroItemType = "Statute"
	ZoteroItemTypeThesis              ZoteroItemType = "Thesis"
	ZoteroItemTypeTvBroadcast         ZoteroItemType = "TV Broadcast"
	ZoteroItemTypeVideoRecording      ZoteroItemType = "Video Recording"
	ZoteroItemTypeWebpage             ZoteroItemType = "Web Page"
)

var ZoteroItemTypeReverse = map[string]ZoteroItemType{
	"artwork":             ZoteroItemTypeArtwork,
	"attachment":          ZoteroItemTypeAttachment,
	"audioRecording":      ZoteroItemTypeAudioRecording,
	"bill":                ZoteroItemTypeBill,
	"blogPost":            ZoteroItemTypeBlogPost,
	"book":                ZoteroItemTypeBook,
	"bookSection":         ZoteroItemTypeBookSection,
	"case":                ZoteroItemTypeCase,
	"computerProgram":     ZoteroItemTypeComputerProgram,
	"conferencePaper":     ZoteroItemTypeConferencePaper,
	"dictionaryEntry":     ZoteroItemTypeDictionaryEntry,
	"document":            ZoteroItemTypeDocument,
	"email":               ZoteroItemTypeEmail,
	"encyclopediaArticle": ZoteroItemTypeEncyclopediaArticle,
	"film":                ZoteroItemTypeFilm,
	"forumPost":           ZoteroItemTypeForumPost,
	"hearing":             ZoteroItemTypeHearing,
	"instantMessage":      ZoteroItemTypeInstantMessage,
	"interview":           ZoteroItemTypeInterview,
	"journalArticle":      ZoteroItemTypeJournalArticle,
	"letter":              ZoteroItemTypeLetter,
	"magazineArticle":     ZoteroItemTypeMagazineArticle,
	"manuscript":          ZoteroItemTypeManuscript,
	"map":                 ZoteroItemTypeMap,
	"newspaperArticle":    ZoteroItemTypeNewspaperArticle,
	"note":                ZoteroItemTypeNote,
	"patent":              ZoteroItemTypePatent,
	"podcast":             ZoteroItemTypePodcast,
	"presentation":        ZoteroItemTypePresentation,
	"radioBroadcast":      ZoteroItemTypeRadioBroadcast,
	"report":              ZoteroItemTypeReport,
	"statute":             ZoteroItemTypeStatute,
	"thesis":              ZoteroItemTypeThesis,
	"tvBroadcast":         ZoteroItemTypeTvBroadcast,
	"videoRecording":      ZoteroItemTypeVideoRecording,
	"webpage":             ZoteroItemTypeWebpage,
}
