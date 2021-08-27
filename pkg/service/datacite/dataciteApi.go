package datacite

import (
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/pkg/errors"
)

type Base64String string

func (base64str *Base64String) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	base64byte, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return errors.Wrapf(err, "cannot decode base64 string - %s", str)
	}
	*base64str = Base64String(base64byte)
	return err
}

func (base64str *Base64String) MarshalJSON() ([]byte, error) {
	base64bytes := base64.StdEncoding.EncodeToString([]byte(*base64str))
	return json.Marshal(base64bytes)
}

type APIDOIDataAttributes struct {
	XMLName           xml.Name `xml:"http://datacite.org/schema/kernel-4 resource" json:"-"`
	XsiType           string   `xml:"xmlns:xsi,attr" json:"-"`
	XsiSchemaLocation string   `xml:"xsi:schemaLocation,attr" json:"-"`

	DOI                string                      `json:"doi,omitempty"`
	Prefix             string                      `json:"prefix,omitempty"`
	Suffix             string                      `json:"suffix,omitempty"`
	Event              string                      `json:"event,omitempty"` // Can be set to trigger a DOI state change.
	Identifiers        []dataciteModel.Identifier  `json:"identifiers,omitempty"`
	Creators           []dataciteModel.Creator     `json:"creators,omitempty"`
	Titles             []dataciteModel.Title       `json:"titles,omitempty"`
	Publisher          string                      `json:"publisher,omitempty"`
	Container          interface{}                 `json:"container,omitempty"`
	PublicationYear    int64                       `json:"publicationYear,omitempty"`
	Subjects           interface{}                 `json:"subjects,omitempty"`
	Contributors       []dataciteModel.Contributor `json:"contributors,omitempty"`
	Dates              []interface{}               `json:"dates,omitempty"`
	Language           string                      `json:"language,omitempty"`
	Types              map[string]string           `json:"types,omitempty"`
	RelatedIdentifiers []interface{}               `json:"related_identifiers,omitempty"`
	Sizes              []string                    `json:"sizes,omitempty"`
	Formats            []string                    `json:"formats,omitempty"`
	Version            string                      `json:"version,omitempty"`
	RightsList         []interface{}               `json:"rightsList,omitempty"`
	Descriptions       []interface{}               `json:"descriptions,omitempty"`
	GeoLocations       []interface{}               `json:"geoLocations,omitempty"`
	FundingReferences  []interface{}               `json:"fundingReferences,omitempty"`
	Xml                Base64String                `json:"xml,omitempty"`
	Url                string                      `json:"url,omitempty"`
	ContentUrl         []string                    `json:"contentUrl,omitempty"`
	MetadataVersion    int                         `json:"metadataVersion,omitempty"`
	SchemaVersion      string                      `json:"schemaVersion,omitempty"`
	Source             string                      `json:"source,omitempty"`
	IsActive           bool                        `json:"isActive,omitempty"`
	State              string                      `json:"state,omitempty"`
	Reason             string                      `json:"reason,omitempty"`
	LandingPage        interface{}                 `json:"landingPage,omitempty"` // Data describing the landing page, used by link checking.
	Created            string                      `json:"created,omitempty"`
	Registered         string                      `json:"registered,omitempty"`
	Updated            string                      `json:"updated,omitempty"`
}

func (dda *APIDOIDataAttributes) InitNamespace() {
	dda.XsiType = "http://www.w3.org/2001/XMLSchema-instance"
	dda.XsiSchemaLocation = "http://datacite.org/schema/kernel-4 http://schema.datacite.org/meta/kernel-4.1/metadata.xsd"
}

type APIDOIData struct {
	Id   string `json:"id,omitempty"`
	Type string `json:"type"`

	Attributes APIDOIDataAttributes `json:"attributes"`
}

type API struct {
	Data *APIDOIData `json:"data"`
}

type APIError struct {
	Status string `json:"status"`
	Title  string `json:"title"`
}

type APIErrorResult struct {
	Errors []*APIError `json:"errors"`
}
