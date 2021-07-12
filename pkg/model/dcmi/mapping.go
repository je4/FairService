package dcmi

import (
	"fmt"
	"github.com/je4/FairService/v2/pkg/model/myfair"
)

var MyfairResource = map[myfair.ResourceType]Type{
	myfair.ResourceTypeBook:                TypeText,
	myfair.ResourceTypeBookSection:         TypeText,
	myfair.ResourceTypeThesis:              TypeText,
	myfair.ResourceTypeJournalArticle:      TypeText,
	myfair.ResourceTypeMagazineArticle:     TypeText,
	myfair.ResourceTypeOnlineResource:      TypeText,
	myfair.ResourceTypeReport:              TypeText,
	myfair.ResourceTypeWebpage:             TypeText,
	myfair.ResourceTypeConferencePaper:     TypeText,
	myfair.ResourceTypePatent:              TypeText,
	myfair.ResourceTypeNote:                TypeText,
	myfair.ResourceTypeArtisticPerformance: TypeEvent,
	myfair.ResourceTypeDataset:             TypeDataset,
	myfair.ResourceTypePresentation:        TypeEvent,
	myfair.ResourceTypePhysicalObject:      TypePhysicalObject,
	myfair.ResourceTypeComputerProgram:     TypeSoftware,
	myfair.ResourceTypeOther:               TypeText,
	myfair.ResourceTypeArtwork:             TypeText,
	myfair.ResourceTypeAttachment:          TypeText,
	myfair.ResourceTypeAudioRecording:      TypeSound,
	myfair.ResourceTypeDocument:            TypeText,
	myfair.ResourceTypeEmail:               TypeText,
	myfair.ResourceTypeEncyclopediaArticle: TypeText,
	myfair.ResourceTypeFilm:                TypeMovingImage,
	myfair.ResourceTypeInstantMessage:      TypeText,
	myfair.ResourceTypeInterview:           TypeText,
	myfair.ResourceTypeLetter:              TypeText,
	myfair.ResourceTypeManuscript:          TypeText,
	myfair.ResourceTypeMap:                 TypeText,
	myfair.ResourceTypeNewspaperArticle:    TypeText,
	myfair.ResourceTypePodcast:             TypeSound,
	myfair.ResourceTypeRadioBroadcast:      TypeSound,
	myfair.ResourceTypeTvBroadcast:         TypeMovingImage,
	myfair.ResourceTypeVideoRecording:      TypeMovingImage,
}

func resourceTypeFromCore(rt myfair.ResourceType) Type {
	t, ok := MyfairResource[rt]
	if !ok {
		return TypeText
	}
	return t
}

func (dcmi *DCMI) FromCore(core myfair.Core) error {

	dcmi.Type = resourceTypeFromCore(core.ResourceType)
	// Title
	dcmi.Title = []string{}
	for _, t := range core.Title {
		if t.Type == myfair.TitleTypeMain {
			dcmi.Title = append(dcmi.Title, t.Data)
		} else {
			dcmi.AlternativeTitle = append(dcmi.AlternativeTitle, t.Data)
		}
	}

	// Persons
	dcmi.Creator = []string{}
	dcmi.Contributor = []string{}
	for _, p := range core.Person {
		pTypeCreator := true
		switch p.PersonType {
		case myfair.PersonTypeAuthor:
			pTypeCreator = true
		case myfair.PersonTypeContactPerson:
			pTypeCreator = false
		case myfair.PersonTypeDataCollector:
			pTypeCreator = false
		case myfair.PersonTypeArtist:
			pTypeCreator = true
		case myfair.PersonTypeDataCurator:
			pTypeCreator = true
		}
		if pTypeCreator {
			dcmi.Creator = append(dcmi.Creator, p.PersonName.Value)
		} else {
			dcmi.Contributor = append(dcmi.Contributor, p.PersonName.Value)
		}
	}

	if core.Publisher != "" {
		dcmi.Publisher = []string{core.Publisher}
	}
	dcmi.Date = []string{core.PublicationYear}
	if len(core.Identifier) > 0 {
		dcmi.Identifier = []string{}
		for _, id := range core.Identifier {
			dcmi.Identifier = append(dcmi.Identifier, fmt.Sprintf("%v:%s", id.IdentifierType, id.Value))
		}
	}
	if core.Rights != "" {
		dcmi.Rights = []string{core.Rights}
	}

	return nil
}
