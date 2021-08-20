package datacite

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

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

func (c *Client) GetPrefix() string {
	return c.prefix
}

func (c *Client) Heartbeat() error {
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
	var client http.Client

	urlStr := fmt.Sprintf("%s/dois/%s", c.api, doi)
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse url %s", urlStr)
	}
	req := &http.Request{
		Method: "GET",
		URL:    u,
		Header: map[string][]string{},
	}
	uPwd := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.user, c.password)))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", uPwd))
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot query %s", urlStr)
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
	dc := &API{}
	if err := json.Unmarshal(rData, dc); err != nil {
		return nil, errors.Wrapf(err, "cannot unmarshal result [%s]", string(rData))
	}
	return dc, nil
}

func (c *Client) CreateDOI(data *DataCite, doiSuffix, targetUrl string) (*API, error) {
	var client http.Client

	xmlBytes, err := xml.Marshal(data)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot marshal request data")
	}

	doiString := fmt.Sprintf("%s/%s", c.prefix, doiSuffix)
	a := API{Data: &APIDOIData{
		Id:   doiString,
		Type: "dois",
		Attributes: APIDOIDataAttributes{
			//			Event: "draft", // publish - register - hide
			DOI: doiString,
			Xml: Base64String(xmlBytes),
			Url: targetUrl,
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
