package fairclient

import (
	"bytes"
	"crypto/sha512"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/service"
	"github.com/je4/utils/v2/pkg/JWTInterceptor"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"time"
)

type FairClient struct {
	service        string
	address        string
	certSkipVerify bool
	jwtKey         string
	jwtAlg         string
	client         *http.Client
}

func NewFairService(service, address string, certSkipVerify bool, jwtKey string, jwtAlg string, jwtLifetime time.Duration) (*FairClient, error) {
	// create transport with authorization bearer
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: certSkipVerify}
	tr, err := JWTInterceptor.NewJWTTransport(service, nil, sha512.New(), jwtKey, jwtAlg, jwtLifetime)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create jwt transport")
	}

	fs := &FairClient{
		service:        service,
		address:        address,
		certSkipVerify: certSkipVerify,
		jwtKey:         jwtKey,
		jwtAlg:         jwtAlg,
		client:         &http.Client{Transport: tr},
	}
	return fs, nil
}

func (fs *FairClient) StartUpdate(source string) error {
	srcData := fair.SourceData{Source: source}
	data, err := json.Marshal(srcData)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal [%v]", srcData)
	}

	response, err := fs.client.Post(fs.address+"/startupdate", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrapf(err, "cannot post to %s", fs.address)
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body")
	}
	result := service.CreateResultStatus{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return errors.Wrapf(err, "cannot decode result %s", string(bodyBytes))
	}
	if result.Status != "ok" {
		return errors.New(fmt.Sprintf("error starting update: %s", result.Message))
	}
	return nil

}
func (fs *FairClient) EndUpdate(source string) error {
	srcData := fair.SourceData{Source: source}
	data, err := json.Marshal(srcData)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal [%v]", srcData)
	}

	response, err := fs.client.Post(fs.address+"/endupdate", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrapf(err, "cannot post to %s", fs.address)
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body")
	}
	result := service.CreateResultStatus{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return errors.Wrapf(err, "cannot decode result %s", string(bodyBytes))
	}
	if result.Status != "ok" {
		return errors.New(fmt.Sprintf("error starting update: %s", result.Message))
	}
	return nil

}
func (fs *FairClient) AbortUpdate(source string) error {
	srcData := fair.SourceData{Source: source}
	data, err := json.Marshal(srcData)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal [%v]", srcData)
	}

	response, err := fs.client.Post(fs.address+"/abortupdate", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrapf(err, "cannot post to %s", fs.address)
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body")
	}
	result := service.CreateResultStatus{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return errors.Wrapf(err, "cannot decode result %s", string(bodyBytes))
	}
	if result.Status != "ok" {
		return errors.New(fmt.Sprintf("error starting update: %s", result.Message))
	}
	return nil

}

func (fs *FairClient) Create(item *fair.ItemData) (*fair.ItemData, error) {
	data, err := json.Marshal(item)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot marshal [%v]", item)
	}

	response, err := fs.client.Post(fs.address+"/item", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot post to %s", fs.address)
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}
	result := service.CreateResultStatus{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrapf(err, "cannot decode result %s", string(bodyBytes))
	}
	if result.Status != "ok" {
		return nil, errors.New(fmt.Sprintf("error creating item: %s", result.Message))
	}
	if result.Item == nil {
		return nil, errors.New(fmt.Sprintf("no item in result: %s", result.Message))
	}
	return result.Item, nil
}
