package dcmi

import (
	"fmt"
	"github.com/je4/FairService/v2/pkg/model/myfair"
)

func (dcmi *DCMI) FromCore(core myfair.Core) error {

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
