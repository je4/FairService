package fair

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

type HandleServiceClient struct {
	client http.Client
	addr   string
}

func NewHandleServiceClient(addr string, tr http.RoundTripper) (*HandleServiceClient, error) {
	hsc := &HandleServiceClient{
		addr: addr,
		client: http.Client{
			Transport: tr,
		}}
	return hsc, nil
}

func (hsc *HandleServiceClient) Create(handle string, URL *url.URL) error {
	createStruct := struct {
		Handle string `json:"handle"`
		Url    string `json:"url"`
	}{
		Handle: handle,
		Url:    URL.String(),
	}
	data, err := json.Marshal(createStruct)
	if err != nil {
		return errors.Wrapf(err, fmt.Sprintf("cannot marshal create struct %v", createStruct))
	}
	u := fmt.Sprintf("%s/create", hsc.addr)
	resp, err := hsc.client.Post(u,
		"application/json",
		bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrapf(err, "cannot query %s", u)
	}
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read result")
	}
	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("handle %s not created: %s", handle, string(result)))
	}
	return nil
}
