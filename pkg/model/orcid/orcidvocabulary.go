package orcid

type ORCIDWorkType string

const (
	ORCIDWorkTypeBook                         ORCIDWorkType = "book"
	ORCIDWorkTypeBookChapter                  ORCIDWorkType = "book-chapter"
	ORCIDWorkTypeBookReview                   ORCIDWorkType = "book-review"
	ORCIDWorkTypeDictionaryEntry              ORCIDWorkType = "dictionary-entry"
	ORCIDWorkTypeDissertation                 ORCIDWorkType = "dissertation"
	ORCIDWorkTypeDissertationThesis           ORCIDWorkType = "dissertation-thesis"
	ORCIDWorkTypeEncyclopediaEntry            ORCIDWorkType = "encyclopedia-entry"
	ORCIDWorkTypeEditedBook                   ORCIDWorkType = "edited-book"
	ORCIDWorkTypeJournalArticle               ORCIDWorkType = "journal-article"
	ORCIDWorkTypeJournalIssue                 ORCIDWorkType = "journal-issue"
	ORCIDWorkTypeMagazineArticle              ORCIDWorkType = "magazine-article"
	ORCIDWorkTypeManual                       ORCIDWorkType = "manual"
	ORCIDWorkTypeOnlineResource               ORCIDWorkType = "online-resource"
	ORCIDWorkTypeNewsletterArticle            ORCIDWorkType = "newsletter-article"
	ORCIDWorkTypeNewspaperArticle             ORCIDWorkType = "newspaper-article"
	ORCIDWorkTypePreprint                     ORCIDWorkType = "preprint"
	ORCIDWorkTypeReport                       ORCIDWorkType = "report"
	ORCIDWorkTypeResearchTool                 ORCIDWorkType = "research-tool"
	ORCIDWorkTypeSupervisedStudentPublication ORCIDWorkType = "supervised-student-publication"
	ORCIDWorkTypeTest                         ORCIDWorkType = "test"
	ORCIDWorkTypeTranslation                  ORCIDWorkType = "translation"
	ORCIDWorkTypeWebsite                      ORCIDWorkType = "website"
	ORCIDWorkTypeWorkingPaper                 ORCIDWorkType = "working-paper"
	ORCIDWorkTypeConference                   ORCIDWorkType = "Conference"
	ORCIDWorkTypeConferenceAbstract           ORCIDWorkType = "conference-abstract"
	ORCIDWorkTypeConferencePaper              ORCIDWorkType = "conference-paper"
	ORCIDWorkTypeConferencePoster             ORCIDWorkType = "conference-poster"
	ORCIDWorkTypeIntellectualProperty         ORCIDWorkType = "Intellectual Property"
	ORCIDWorkTypeDisclosure                   ORCIDWorkType = "disclosure"
	ORCIDWorkTypeLicense                      ORCIDWorkType = "license"
	ORCIDWorkTypePatent                       ORCIDWorkType = "patent"
	ORCIDWorkTypeRegisteredCopyright          ORCIDWorkType = "registered-copyright"
	ORCIDWorkTypeTrademark                    ORCIDWorkType = "trademark"
	ORCIDWorkTypeAnnotation                   ORCIDWorkType = "annotation"
	ORCIDWorkTypeArtisticPerformance          ORCIDWorkType = "artistic-performance"
	ORCIDWorkTypeDataManagementPlan           ORCIDWorkType = "data-management-plan"
	ORCIDWorkTypeDataSet                      ORCIDWorkType = "data-set"
	ORCIDWorkTypeInvention                    ORCIDWorkType = "invention"
	ORCIDWorkTypeLectureSpeech                ORCIDWorkType = "lecture-speech"
	ORCIDWorkTypePhysicalObject               ORCIDWorkType = "physical-object"
	ORCIDWorkTypeResearchTechnique            ORCIDWorkType = "research-technique"
	ORCIDWorkTypeSoftware                     ORCIDWorkType = "software"
	ORCIDWorkTypeSpinOffCompany               ORCIDWorkType = "spin-off-company"
	ORCIDWorkTypeStandardsAndPolicy           ORCIDWorkType = "standards-and-policy"
	ORCIDWorkTypeTechnicalStandard            ORCIDWorkType = "technical-standard"
	ORCIDWorkTypeOther                        ORCIDWorkType = "other"
)

var ORCIDWorkTypeReverse = map[string]ORCIDWorkType{
	"book":                           ORCIDWorkTypeBook,
	"book-chapter":                   ORCIDWorkTypeBookChapter,
	"book-review":                    ORCIDWorkTypeBookReview,
	"dictionary-entry":               ORCIDWorkTypeDictionaryEntry,
	"dissertation":                   ORCIDWorkTypeDissertation,
	"dissertation-thesis":            ORCIDWorkTypeDissertationThesis,
	"encyclopedia-entry":             ORCIDWorkTypeEncyclopediaEntry,
	"edited-book":                    ORCIDWorkTypeEditedBook,
	"journal-article":                ORCIDWorkTypeJournalArticle,
	"journal-issue":                  ORCIDWorkTypeJournalIssue,
	"magazine-article":               ORCIDWorkTypeMagazineArticle,
	"manual":                         ORCIDWorkTypeManual,
	"online-resource":                ORCIDWorkTypeOnlineResource,
	"newsletter-article":             ORCIDWorkTypeNewsletterArticle,
	"newspaper-article":              ORCIDWorkTypeNewspaperArticle,
	"preprint":                       ORCIDWorkTypePreprint,
	"report":                         ORCIDWorkTypeReport,
	"research-tool":                  ORCIDWorkTypeResearchTool,
	"supervised-student-publication": ORCIDWorkTypeSupervisedStudentPublication,
	"test":                           ORCIDWorkTypeTest,
	"translation":                    ORCIDWorkTypeTranslation,
	"website":                        ORCIDWorkTypeWebsite,
	"working-paper":                  ORCIDWorkTypeWorkingPaper,
	"Conference":                     ORCIDWorkTypeConference,
	"conference-abstract":            ORCIDWorkTypeConferenceAbstract,
	"conference-paper":               ORCIDWorkTypeConferencePaper,
	"conference-poster":              ORCIDWorkTypeConferencePoster,
	"Intellectual Property":          ORCIDWorkTypeIntellectualProperty,
	"disclosure":                     ORCIDWorkTypeDisclosure,
	"license":                        ORCIDWorkTypeLicense,
	"patent":                         ORCIDWorkTypePatent,
	"registered-copyright":           ORCIDWorkTypeRegisteredCopyright,
	"trademark":                      ORCIDWorkTypeTrademark,
	"annotation":                     ORCIDWorkTypeAnnotation,
	"artistic-performance":           ORCIDWorkTypeArtisticPerformance,
	"data-management-plan":           ORCIDWorkTypeDataManagementPlan,
	"data-set":                       ORCIDWorkTypeDataSet,
	"invention":                      ORCIDWorkTypeInvention,
	"lecture-speech":                 ORCIDWorkTypeLectureSpeech,
	"physical-object":                ORCIDWorkTypePhysicalObject,
	"research-technique":             ORCIDWorkTypeResearchTechnique,
	"software":                       ORCIDWorkTypeSoftware,
	"spin-off-company":               ORCIDWorkTypeSpinOffCompany,
	"standards-and-policy":           ORCIDWorkTypeStandardsAndPolicy,
	"technical-standard":             ORCIDWorkTypeTechnicalStandard,
	"other":                          ORCIDWorkTypeOther,
}
