package datacite

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/je4/FairService/v2/pkg/model/myfair"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

/*
Object	data
String	data.id
String	data.type
Object	data.attributes
String	data.attributes.doi
String	data.attributes.prefix
String	data.attributes.suffix
String	data.attributes.event
Can be set to trigger a DOI state change.

[Object]	data.attributes.identifiers
[Object]	data.attributes.creators
[Object]	data.attributes.titles
String	data.attributes.publisher
Object	data.attributes.container
Number	data.attributes.publicationYear
[Object]	data.attributes.subjects
[Object]	data.attributes.contributors
[Object]	data.attributes.dates
String	data.attributes.language
Object	data.attributes.types
[Object]	data.attributes.relatedIdentifiers
[String]	data.attributes.sizes
[String]	data.attributes.formats
String	data.attributes.version
[Object]	data.attributes.rightsList
[Object]	data.attributes.descriptions
[Object]	data.attributes.geoLocations
[Object]	data.attributes.fundingReferences
String	data.attributes.url
[String]	data.attributes.contentUrl
Number	data.attributes.metadataVersion
String	data.attributes.schemaVersion
String	data.attributes.source
Boolean	data.attributes.isActive
String	data.attributes.state
String	data.attributes.reason
Object	data.attributes.landingPage
Data describing the landing page, used by link checking.

String	data.attributes.created
String	data.attributes.registered
String	data.attributes.updated

*/

type Client struct {
	api      string
	user     string
	password string
	prefix   string
}

func NewClient(api, user, password string, prefix string) (*Client, error) {
	client := &Client{
		api:      strings.TrimRight(api, "/"),
		user:     user,
		password: password,
		prefix:   prefix,
	}
	return client, nil
}

func (c *Client) Hearbeat() error {
	urlStr := fmt.Sprintf("%s/heartbeat", c.api)
	resp, err := http.Get(urlStr)
	if err != nil {
		return errors.Wrapf(err, "cannot query %s", urlStr)
	}
	rData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "cannot get result data of %s", urlStr)
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("status: %s - %s", resp.Status, string(rData)))
	}
	return nil
}

func (c *Client) RetrieveDOI(doi string) (*API, error) {
	urlStr := fmt.Sprintf("%s/dois/%s", c.api, doi)
	resp, err := http.Get(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot query %s", urlStr)
	}
	rData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get result data of %s", urlStr)
	}

	dc := &API{}
	if err := json.Unmarshal(rData, dc); err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal result [%s]", string(rData))
	}
	return dc, nil
}

func (c *Client) CreateDOI(data *myfair.Core) (*API, error) {
	var client http.Client

	a := API{Data: &APIDOIData{
		Type: "dois",
		Attributes: APIDOIDataAttributes{
			DOI: fmt.Sprintf("%s/%s", c.prefix, "6fzw-t035"),
		},
	}}
	aJson, err := json.Marshal(a)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot marshal request data")
	}

	urlStr := fmt.Sprintf("%s/dois", c.api)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse url %s", urlStr)
	}
	req := &http.Request{
		Method: "POST",
		URL:    u,
		Body:   ioutil.NopCloser(bytes.NewBuffer(aJson)),
		Header: map[string][]string{},
	}
	upwd := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.user, c.password)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", upwd))
	req.Header.Set("Content-type", "application/vnd.api+json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot request url %s", urlStr)
	}
	rData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get result data of %s", urlStr)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		rErr := &APIErrorResult{}
		if err := json.Unmarshal(rData, rErr); err != nil {
			return nil, errors.Wrapf(err, "cannot unmarshal result [%s]", string(rData))
		}
		errs := []string{}
		for _, e := range rErr.Errors {
			errs = append(errs, fmt.Sprintf("%v:%s", e.Status, e.Title))
		}
		return nil, errors.New(fmt.Sprintf("%s - %s", resp.Status, strings.Join(errs, " / ")))
	}

	dca := &API{}
	if err := json.Unmarshal(rData, dca); err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal result [%s]", string(rData))
	}
	return dca, nil

}
