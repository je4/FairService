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
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type FairClient struct {
	service        string
	address        string
	certSkipVerify bool
	jwtKey         string
	jwtAlg         string
	jwtLifetime    time.Duration
}

func postHelper(client http.Client, urlstr string, data []byte) error {
	response, err := client.Post(urlstr, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrapf(err, "cannot post to %s", urlstr)
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body")
	}
	result := service.FairResultStatus{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return errors.Wrapf(err, "cannot decode result %s", string(bodyBytes))
	}
	if strings.ToLower(result.Status) != "ok" {
		return errors.New(fmt.Sprintf("error on POST::%s: %s", urlstr, result.Message))
	}
	return nil
}

func getHelper(client http.Client, urlstr string) (*service.FairResultStatus, error) {
	response, err := client.Get(urlstr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot post to %s", urlstr)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}
	result := &service.FairResultStatus{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrapf(err, "cannot decode result %s", string(bodyBytes))
	}
	if strings.ToLower(result.Status) != "ok" {
		return nil, errors.New(fmt.Sprintf("error on POST::%s: %s", urlstr, result.Message))
	}
	return result, nil
}

func NewFairService(service, address string, certSkipVerify bool, jwtKey string, jwtAlg string, jwtLifetime time.Duration) (*FairClient, error) {
	// create transport with authorization bearer
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: certSkipVerify}

	fs := &FairClient{
		service:        service,
		address:        address,
		certSkipVerify: certSkipVerify,
		jwtKey:         jwtKey,
		jwtAlg:         jwtAlg,
		jwtLifetime:    jwtLifetime,
	}
	return fs, nil
}

func (fs *FairClient) Ping() error {
	response, err := http.Get(fs.address + "/ping")
	if err != nil {
		return errors.Wrapf(err, "cannot post to %s", fs.address)
	}
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "cannot read response body")
	}
	result := service.FairResultStatus{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return errors.Wrapf(err, "cannot decode result %s", string(bodyBytes))
	}
	if strings.ToLower(result.Status) != "ok" {
		return errors.New(fmt.Sprintf("ping error: %s", result.Message))
	}
	return nil
}

func (fs *FairClient) StartUpdate(source string) error {
	tr, err := JWTInterceptor.NewJWTTransport(
		fs.service,
		"StartUpdate",
		JWTInterceptor.Secure,
		nil,
		sha512.New(),
		fs.jwtKey,
		fs.jwtAlg,
		fs.jwtLifetime)
	if err != nil {
		return errors.Wrapf(err, "cannot create jwt transport")
	}
	client := http.Client{Transport: tr}

	srcData := fair.SourceData{Source: source}
	data, err := json.Marshal(srcData)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal [%v]", srcData)
	}

	urlstr := fs.address + "/startupdate"
	return postHelper(client, urlstr, data)
}
func (fs *FairClient) EndUpdate(source string) error {
	tr, err := JWTInterceptor.NewJWTTransport(
		fs.service,
		"EndUpdate",
		JWTInterceptor.Secure,
		nil,
		sha512.New(),
		fs.jwtKey,
		fs.jwtAlg,
		fs.jwtLifetime)
	if err != nil {
		return errors.Wrapf(err, "cannot create jwt transport")
	}
	client := http.Client{Transport: tr}

	srcData := fair.SourceData{Source: source}
	data, err := json.Marshal(srcData)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal [%v]", srcData)
	}

	urlstr := fs.address + "/endupdate"
	return postHelper(client, urlstr, data)
}
func (fs *FairClient) AbortUpdate(source string) error {
	tr, err := JWTInterceptor.NewJWTTransport(
		fs.service,
		"AbortUpdate",
		JWTInterceptor.Secure,
		nil,
		sha512.New(),
		fs.jwtKey,
		fs.jwtAlg,
		fs.jwtLifetime)
	if err != nil {
		return errors.Wrapf(err, "cannot create jwt transport")
	}
	client := http.Client{Transport: tr}

	srcData := fair.SourceData{Source: source}
	data, err := json.Marshal(srcData)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal [%v]", srcData)
	}

	urlstr := fs.address + "/abortupdate"
	return postHelper(client, urlstr, data)
}

func (fs *FairClient) Create(item *fair.ItemData) (*fair.ItemData, error) {
	tr, err := JWTInterceptor.NewJWTTransport(
		fs.service,
		"CreateItem",
		JWTInterceptor.Secure,
		nil,
		sha512.New(),
		fs.jwtKey,
		fs.jwtAlg,
		fs.jwtLifetime)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create jwt transport")
	}
	client := http.Client{Transport: tr}

	data, err := json.Marshal(item)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot marshal [%v]", item)
	}

	response, err := client.Post(fs.address+"/item", "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot post to %s", fs.address)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}
	result := service.FairResultStatus{}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, errors.Wrapf(err, "cannot decode result %s", string(bodyBytes))
	}
	if strings.ToLower(result.Status) != "ok" {
		return nil, errors.New(fmt.Sprintf("error creating item: %s", result.Message))
	}
	if result.Item == nil {
		return nil, errors.New(fmt.Sprintf("no item in result: %s", result.Message))
	}
	return result.Item, nil
}

func (fs *FairClient) SetSource(src *fair.Source) error {
	tr, err := JWTInterceptor.NewJWTTransport(
		fs.service,
		"setSource",
		JWTInterceptor.Secure,
		nil,
		sha512.New(),
		fs.jwtKey,
		fs.jwtAlg,
		fs.jwtLifetime)
	if err != nil {
		return errors.Wrapf(err, "cannot create jwt transport")
	}
	client := http.Client{Transport: tr}

	data, err := json.Marshal(src)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal [%v]", src)
	}

	urlstr := fs.address + "/source"
	return postHelper(client, urlstr, data)
}

func (fs *FairClient) WriteOriginalData(item *fair.ItemData, data []byte) error {
	tr, err := JWTInterceptor.NewJWTTransport(
		fs.service,
		"originalDataWrite",
		JWTInterceptor.Secure,
		nil,
		sha512.New(),
		fs.jwtKey,
		fs.jwtAlg,
		fs.jwtLifetime)
	if err != nil {
		return errors.Wrapf(err, "cannot create jwt transport")
	}
	client := http.Client{Transport: tr}

	urlString := fmt.Sprintf("%s/item/%s/originaldata", fs.address, item.UUID)
	return postHelper(client, urlString, data)

}

func (fs *FairClient) ReadOriginalData(item *fair.ItemData) ([]byte, error) {
	tr, err := JWTInterceptor.NewJWTTransport(
		fs.service,
		"originalDataRead",
		JWTInterceptor.Secure,
		nil,
		sha512.New(),
		fs.jwtKey,
		fs.jwtAlg,
		fs.jwtLifetime)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create jwt transport")
	}
	client := http.Client{Transport: tr}

	urlString := fmt.Sprintf("%s/item/%s/originaldata", fs.address, item.UUID)
	response, err := client.Get(urlString)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot post to %s", fs.address)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "cannot read response body")
	}
	if response.StatusCode != http.StatusOK {
		result := service.FairResultStatus{}
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			return nil, errors.Wrapf(err, "cannot decode result %s", string(bodyBytes))
		}
		return nil, errors.New(fmt.Sprintf("error reading original data of item %s: %s", item.UUID, result.Message))
	}
	return bodyBytes, nil
}

func (fs *FairClient) AddArchive(name, description string) error {
	tr, err := JWTInterceptor.NewJWTTransport(
		fs.service,
		"AddArchive",
		JWTInterceptor.Secure,
		nil,
		sha512.New(),
		fs.jwtKey,
		fs.jwtAlg,
		fs.jwtLifetime)
	if err != nil {
		return errors.Wrapf(err, "cannot create jwt transport")
	}
	client := http.Client{Transport: tr}

	ar := service.Archive{
		Name:        name,
		Description: description,
	}

	data, err := json.Marshal(ar)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal [%v]", ar)
	}

	urlstr := fs.address + "/archive"
	return postHelper(client, urlstr, data)
}

func (fs *FairClient) AddArchiveItem(archive, uuid string) error {
	tr, err := JWTInterceptor.NewJWTTransport(
		fs.service,
		"AddArchiveItem",
		JWTInterceptor.Secure,
		nil,
		sha512.New(),
		fs.jwtKey,
		fs.jwtAlg,
		fs.jwtLifetime)
	if err != nil {
		return errors.Wrapf(err, "cannot create jwt transport")
	}
	client := http.Client{Transport: tr}

	data, err := json.Marshal(uuid)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal [%v]", uuid)
	}

	urlstr := fmt.Sprintf("%s/archive/%s", fs.address, archive)
	return postHelper(client, urlstr, data)
}

func (fs *FairClient) GetArchiveItem(archive, uuid string) error {
	tr, err := JWTInterceptor.NewJWTTransport(
		fs.service,
		"GetArchiveItem",
		JWTInterceptor.Secure,
		nil,
		sha512.New(),
		fs.jwtKey,
		fs.jwtAlg,
		fs.jwtLifetime)
	if err != nil {
		return errors.Wrapf(err, "cannot create jwt transport")
	}
	client := http.Client{Transport: tr}

	data, err := json.Marshal(uuid)
	if err != nil {
		return errors.Wrapf(err, "cannot marshal [%v]", uuid)
	}

	urlstr := fmt.Sprintf("%s/archive/%s", fs.address, archive)
	return postHelper(client, urlstr, data)
}
