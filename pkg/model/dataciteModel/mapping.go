package dataciteModel

import (
	"github.com/araddon/dateparse"
	"github.com/je4/FairService/v2/pkg/model/myfair"
)

var MyfairResourceTypeGeneral = map[myfair.ResourceType]ResourceTypeGeneral{
	myfair.ResourceTypeBook:                ResourceTypeGeneralBook,
	myfair.ResourceTypeBookSection:         ResourceTypeGeneralBookChapter,
	myfair.ResourceTypeThesis:              ResourceTypeGeneralDissertation,
	myfair.ResourceTypeJournalArticle:      ResourceTypeGeneralJournalArticle,
	myfair.ResourceTypeMagazineArticle:     ResourceTypeGeneralJournalArticle,
	myfair.ResourceTypeOnlineResource:      ResourceTypeGeneralOther,
	myfair.ResourceTypeReport:              ResourceTypeGeneralReport,
	myfair.ResourceTypeWebpage:             ResourceTypeGeneralOther,
	myfair.ResourceTypeConferencePaper:     ResourceTypeGeneralConferencePaper,
	myfair.ResourceTypePatent:              ResourceTypeGeneralOther,
	myfair.ResourceTypeNote:                ResourceTypeGeneralOther,
	myfair.ResourceTypeArtisticPerformance: ResourceTypeGeneralOther,
	myfair.ResourceTypeDataset:             ResourceTypeGeneralDataset,
	myfair.ResourceTypePresentation:        ResourceTypeGeneralOther,
	myfair.ResourceTypePhysicalObject:      ResourceTypeGeneralPhysicalObject,
	myfair.ResourceTypeComputerProgram:     ResourceTypeGeneralSoftware,
	myfair.ResourceTypeOther:               ResourceTypeGeneralOther,
	myfair.ResourceTypeArtwork:             ResourceTypeGeneralOther,
	myfair.ResourceTypeAttachment:          ResourceTypeGeneralOther,
	myfair.ResourceTypeAudioRecording:      ResourceTypeGeneralSound,
	myfair.ResourceTypeDocument:            ResourceTypeGeneralOther,
	myfair.ResourceTypeEmail:               ResourceTypeGeneralOther,
	myfair.ResourceTypeEncyclopediaArticle: ResourceTypeGeneralOther,
	myfair.ResourceTypeFilm:                ResourceTypeGeneralOther,
	myfair.ResourceTypeInstantMessage:      ResourceTypeGeneralOther,
	myfair.ResourceTypeInterview:           ResourceTypeGeneralOther,
	myfair.ResourceTypeLetter:              ResourceTypeGeneralOther,
	myfair.ResourceTypeManuscript:          ResourceTypeGeneralOther,
	myfair.ResourceTypeMap:                 ResourceTypeGeneralOther,
	myfair.ResourceTypeNewspaperArticle:    ResourceTypeGeneralOther,
	myfair.ResourceTypePodcast:             ResourceTypeGeneralOther,
	myfair.ResourceTypeRadioBroadcast:      ResourceTypeGeneralOther,
	myfair.ResourceTypeTvBroadcast:         ResourceTypeGeneralOther,
	myfair.ResourceTypeVideoRecording:      ResourceTypeGeneralOther,
}

func resourceTypeFromCore(rt myfair.ResourceType) ResourceType {
	t, ok := MyfairResourceTypeGeneral[rt]
	if !ok {
		t = ResourceTypeGeneralOther
	}
	return ResourceType{
		Value:          string(rt),
		IdentifierType: t,
	}
}

var MyfairTitleType = map[myfair.CoreTitleType]TitleType{
	myfair.TitleTypeMain:             TitleTypeDefaultTitle,
	myfair.TitleTypeAlternativeTitle: TitleTypeAlternativeTitle,
	myfair.TitleTypeSubTitle:         TitleTypeSubTitle,
	myfair.TitleTypeTranslatedTitle:  TitleTypeTranslatedTitle,
	myfair.TitleTypeOther:            TitleTypeOther,
}

func titleTypeFromCore(tt myfair.CoreTitleType) TitleType {
	t, ok := MyfairTitleType[tt]
	if !ok {
		t = TitleTypeOther
	}
	return t
}

var MyfairNameType = map[myfair.NameType]NameType{
	myfair.NameTypeDefault:        NameTypeDefault,
	myfair.NameTypePersonal:       NameTypePersonal,
	myfair.NameTypeOrganizational: NameTypeOrganizational,
}

func nameTypeFromCore(nt myfair.NameType) NameType {
	t, ok := MyfairNameType[nt]
	if !ok {
		t = NameTypeDefault
	}
	return t
}

var MyfairRelatedIdentifierType = map[myfair.RelatedIdentifierType]RelatedIdentifierType{
	myfair.RelatedIdentifierTypeARK:     RelatedIdentifierTypeARK,
	myfair.RelatedIdentifierTypeArXiv:   RelatedIdentifierTypeArXiv,
	myfair.RelatedIdentifierTypeBibcode: RelatedIdentifierTypeBibcode,
	myfair.RelatedIdentifierTypeDOI:     RelatedIdentifierTypeDOI,
	myfair.RelatedIdentifierTypeEAN13:   RelatedIdentifierTypeEAN13,
	myfair.RelatedIdentifierTypeEISSN:   RelatedIdentifierTypeEISSN,
	myfair.RelatedIdentifierTypeHandle:  RelatedIdentifierTypeHandle,
	myfair.RelatedIdentifierTypeIGSN:    RelatedIdentifierTypeIGSN,
	myfair.RelatedIdentifierTypeISBN:    RelatedIdentifierTypeISBN,
	myfair.RelatedIdentifierTypeISSN:    RelatedIdentifierTypeISSN,
	myfair.RelatedIdentifierTypeISTC:    RelatedIdentifierTypeISTC,
	myfair.RelatedIdentifierTypeLISSN:   RelatedIdentifierTypeLISSN,
	myfair.RelatedIdentifierTypeLSID:    RelatedIdentifierTypeLSID,
	myfair.RelatedIdentifierTypePMID:    RelatedIdentifierTypePMID,
	myfair.RelatedIdentifierTypePURL:    RelatedIdentifierTypePURL,
	myfair.RelatedIdentifierTypeUPC:     RelatedIdentifierTypeUPC,
	myfair.RelatedIdentifierTypeURL:     RelatedIdentifierTypeURL,
	myfair.RelatedIdentifierTypeURN:     RelatedIdentifierTypeURN,
	myfair.RelatedIdentifierTypeW3id:    RelatedIdentifierTypeW3id,
}

func relatedIdentifierFromCore(rit myfair.RelatedIdentifierType) RelatedIdentifierType {
	t, ok := MyfairRelatedIdentifierType[rit]
	if !ok {
		t = ""
	}
	return t
}
func (datacite *DataCite) FromCore(core myfair.Core) error {

	datacite.ResourceType = resourceTypeFromCore(core.ResourceType)

	// Title
	datacite.Titles = Titles{Title: []Title{}}

	for _, t := range core.Title {
		datacite.Titles.Title = append(datacite.Titles.Title, Title{
			Value: t.Data,
			Type:  titleTypeFromCore(t.Type),
		})
	}

	// Persons
	datacite.Creators = Creators{Creator: []Creator{}}
	datacite.Contributors = Contributors{Contributor: []Contributor{}}
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
			datacite.Creators.Creator = append(datacite.Creators.Creator, Creator{
				CreatorName: Name{
					Value: p.PersonName.Value,
					Type:  nameTypeFromCore(p.PersonName.Type),
				},
				GivenName:   p.GivenName,
				FamilyName:  p.FamilyName,
				Affiliation: []string{p.Affiliation},
				NameIdentifier: []NameIdentifier{NameIdentifier{
					Value:                p.NameIdentifier.Value,
					Lang:                 p.NameIdentifier.Lang,
					SchemeURI:            p.NameIdentifier.SchemeURI,
					NameIdentifierScheme: p.NameIdentifier.NameIdentifierScheme,
				}},
			})
		} else {
			datacite.Contributors.Contributor = append(datacite.Contributors.Contributor, Contributor{
				ContributorName: Name{
					Value: p.PersonName.Value,
					Type:  nameTypeFromCore(p.PersonName.Type),
				},
				GivenName:   p.GivenName,
				FamilyName:  p.FamilyName,
				Affiliation: p.Affiliation,
				NameIdentifier: NameIdentifier{
					Value:                p.NameIdentifier.Value,
					Lang:                 p.NameIdentifier.Lang,
					SchemeURI:            p.NameIdentifier.SchemeURI,
					NameIdentifierScheme: p.NameIdentifier.NameIdentifierScheme,
				},
			})
		}
	}

	datacite.Publisher = core.Publisher
	t, err := dateparse.ParseAny(core.PublicationYear)
	if err == nil {
		datacite.PublicationYear = t.Year()
	}

	datacite.AlternateIdentifiers = AlternateIdentifiers{AlternateIdentifier: []AlternateIdentifier{}}
	if len(core.Identifier) > 0 {
		for _, id := range core.Identifier {
			if id.IdentifierType != myfair.RelatedIdentifierTypeDOI {
				idType := relatedIdentifierFromCore(id.IdentifierType)
				if idType != "" {
					datacite.AlternateIdentifiers.AlternateIdentifier = append(datacite.AlternateIdentifiers.AlternateIdentifier, AlternateIdentifier{
						Value:                   id.Value,
						AlternateIdentifierType: idType,
					})
				}
			} else {
				datacite.Identifier = Identifier{
					Value:          id.Value,
					IdentifierType: relatedIdentifierFromCore(id.IdentifierType),
				}
			}
		}
	}

	return nil
}
