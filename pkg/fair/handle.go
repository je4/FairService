package fair

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
)

type HandleServiceClient struct {
	client http.Client
	addr   string
	logger *logging.Logger
}

func NewHandleServiceClient(addr string, tr http.RoundTripper, logger *logging.Logger) (*HandleServiceClient, error) {
	hsc := &HandleServiceClient{
		addr: addr,
		client: http.Client{
			Transport: tr,
		},
		logger: logger,
	}
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
	hsc.logger.Infof("creating handle %s", handle)
	resp, err := hsc.client.Post(u,
		"application/json",
		bytes.NewBuffer(data))
	if err != nil {
		hsc.logger.Errorf("cannot query %s: %v", u, err)
		return errors.Wrapf(err, "cannot query %s", u)
	}
	result, err := io.ReadAll(resp.Body)
	if err != nil {
		hsc.logger.Errorf("cannot read result: %v", err)
		return errors.Wrap(err, "cannot read result")
	}
	if resp.StatusCode != http.StatusCreated {
		hsc.logger.Errorf("handle %s not created: %s", handle, string(result))
		return errors.New(fmt.Sprintf("handle %s not created: %s", handle, string(result)))
	}
	return nil
}
