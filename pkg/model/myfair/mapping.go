package myfair

import "github.com/je4/FairService/v2/pkg/model/zsearch"

var zsarchItemTypeMap = map[zsearch.ItemType]ResourceType{
	zsearch.ItemTypeBook:                ResourceTypeBook,
	zsearch.ItemTypeBookSection:         ResourceTypeBookSection,
	zsearch.ItemTypeThesis:              ResourceTypeThesis,
	zsearch.ItemTypeJournalArticle:      ResourceTypeJournalArticle,
	zsearch.ItemTypeMagazineArticle:     ResourceTypeMagazineArticle,
	zsearch.ItemTypeReport:              ResourceTypeReport,
	zsearch.ItemTypeWebpage:             ResourceTypeWebpage,
	zsearch.ItemTypeConferencePaper:     ResourceTypeConferencePaper,
	zsearch.ItemTypePatent:              ResourceTypePatent,
	zsearch.ItemTypeNote:                ResourceTypeNote,
	zsearch.ItemTypePresentation:        ResourceTypePresentation,
	zsearch.ItemTypeComputerProgram:     ResourceTypeComputerProgram,
	zsearch.ItemTypeOther:               ResourceTypeOther,
	zsearch.ItemTypeArtwork:             ResourceTypeArtwork,
	zsearch.ItemTypeAttachment:          ResourceTypeAttachment,
	zsearch.ItemTypeAudioRecording:      ResourceTypeAudioRecording,
	zsearch.ItemTypeBill:                ResourceTypeOther,
	zsearch.ItemTypeBlogPost:            ResourceTypeOther,
	zsearch.ItemTypeCase:                ResourceTypeOther,
	zsearch.ItemTypeDictionaryEntry:     ResourceTypeOther,
	zsearch.ItemTypeDocument:            ResourceTypeDocument,
	zsearch.ItemTypeEmail:               ResourceTypeEmail,
	zsearch.ItemTypeEncyclopediaArticle: ResourceTypeEncyclopediaArticle,
	zsearch.ItemTypeFilm:                ResourceTypeFilm,
	zsearch.ItemTypeForumPost:           ResourceTypeOther,
	zsearch.ItemTypeHearing:             ResourceTypeOther,
	zsearch.ItemTypeInstantMessage:      ResourceTypeInstantMessage,
	zsearch.ItemTypeInterview:           ResourceTypeInterview,
	zsearch.ItemTypeLetter:              ResourceTypeLetter,
	zsearch.ItemTypeManuscript:          ResourceTypeManuscript,
	zsearch.ItemTypeMap:                 ResourceTypeMap,
	zsearch.ItemTypeNewspaperArticle:    ResourceTypeNewspaperArticle,
	zsearch.ItemTypePodcast:             ResourceTypePodcast,
	zsearch.ItemTypeRadioBroadcast:      ResourceTypeRadioBroadcast,
	zsearch.ItemTypeStatute:             ResourceTypeOther,
	zsearch.ItemTypeTvBroadcast:         ResourceTypeTvBroadcast,
	zsearch.ItemTypeVideoRecording:      ResourceTypeVideoRecording,
}

func ZSearchItemTypeMap(zit zsearch.ItemType) ResourceType {
	if rt, ok := zsarchItemTypeMap[zit]; ok {
		return rt
	} else {
		return ResourceTypeOther
	}
}
