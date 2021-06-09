package datacite

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func TestDataciteExample(t *testing.T) {
	dataCiteCore := &DataCite{
		Identifier: DataCiteIdentifier{
			Value:          "10.5072/example-full",
			IdentifierType: DataCiteRelatedIdentifierTypeDOI,
		},
		Creators: Creators{
			Creator: []DataCiteCreator{
				{
					CreatorName: DataCiteName{
						Value: "Miller, Elizabeth",
						Type:  DataCiteNameTypePersonal,
					},
					GivenName:   "Elizabeth",
					FamilyName:  "Miller",
					Affiliation: "DataCite",
					NameIdentifier: DataCiteNameIdentifier{
						Value:                "0000-0001-5000-0007",
						SchemeURI:            "http://orcid.org/",
						NameIdentifierScheme: "ORCID",
					},
				},
			},
		},
		Contributors: Contributors{
			Contributor: []DataCiteContributor{
				{
					ContributorName: DataCiteName{
						Value: "Miller, Elizabeth",
						Type:  DataCiteNameTypePersonal,
					},
					GivenName:   "Elizabeth",
					FamilyName:  "Miller",
					Affiliation: "DataCite",
					NameIdentifier: DataCiteNameIdentifier{
						Value:                "0000-0001-5000-0007",
						SchemeURI:            "http://orcid.org/",
						NameIdentifierScheme: "ORCID",
					},
				},
			},
		},
		Titles: Titles{
			Title: []DataCiteTitle{
				{
					Value: "Full DataCite XML Example",
					Type:  DataCiteTitleTypeDefaultTitle,
					Lang:  "en-US",
				},
				{
					Value: "Demonstration of DataCite Properties.",
					Type:  DataCiteTitleTypeSubTitle,
					Lang:  "en-US",
				},
			},
		},
		Publisher:       "DataCite",
		PublicationYear: 2014,
		ResourceType: DataCiteResourceType{
			Value:          "XML",
			IdentifierType: DataCiteResourceTypeGeneralSoftware,
		},
	}
	dataCiteCore.InitNamespace()

	xmlBytes, err := xml.Marshal(dataCiteCore)
	if err != nil {
		t.Fatalf("cannot marshal xml: %v", err)
	}
	xmlstr := string(xmlBytes) // fmt.Sprintf("%s\n%s", "<?xml version=\"1.0\" encoding=\"UTF-8\"?>", string(xmlBytes))
	fmt.Print(xmlstr, "\n")
	xmlTestStr := `<resource xmlns="http://datacite.org/schema/kernel-4" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:schemaLocation="http://datacite.org/schema/kernel-4 http://schema.datacite.org/meta/kernel-4.1/metadata.xsd"><identifier identifierType="DOI">10.5072/example-full</identifier><creators><creator><creatorName nameType="Personal">Miller, Elizabeth</creatorName><givenName>Elizabeth</givenName><familyName>Miller</familyName><affiliation>DataCite</affiliation><nameIdentifier schemeURI="http://orcid.org/" nameIdentifierScheme="ORCID">0000-0001-5000-0007</nameIdentifier></creator></creators><titles><title xml:lang="en-US">Full DataCite XML Example</title><title xml:lang="en-US" titleType="Subtitle">Demonstration of DataCite Properties.</title></titles><publisher>DataCite</publisher><publicationYear>2014</publicationYear><resourceType resourceTypeGeneral="Software">XML</resourceType><contributors><contributor contributorType=""><contributorName nameType="Personal">Miller, Elizabeth</contributorName><givenName>Elizabeth</givenName><familyName>Miller</familyName><affiliation>DataCite</affiliation><nameIdentifier schemeURI="http://orcid.org/" nameIdentifierScheme="ORCID">0000-0001-5000-0007</nameIdentifier></contributor></contributors></resource>`
	if xmlstr != xmlTestStr {
		t.Errorf("invalid result from xml marshalling")
	}
}
