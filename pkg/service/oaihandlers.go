package service

import (
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/dcmi"
	"github.com/je4/FairService/v2/pkg/service/oai"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func sendError(w http.ResponseWriter, code oai.ErrorCode, verb, identifier, metadataPrefix, baseURL string) {
	w.Header().Set("Content-type", "text/xml")
	pmh := &oai.OAIPMH{}
	pmh.InitNamespace()
	pmh.Error = &oai.Error{Code: code}
	pmh.ResponseDate = time.Now().Format("2006-01-02T15:04:05Z")
	pmh.Request = &oai.Request{
		Verb:           verb,
		Identifier:     identifier,
		MetadataPrefix: metadataPrefix,
		Value:          baseURL,
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	enc.Encode(pmh)
}

func getVar(key string, values url.Values) (string, bool) {
	vals, ok := values[key]
	if !ok {
		return "", false
	}
	if len(vals) < 1 {
		return "", false
	}
	return vals[0], true
}

func (s *Server) oaiHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]
	partition, err := s.fair.GetPartition(pName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("invalid partition: %s", pName)))
		return
	}
	values := req.URL.Query()
	verb, ok := getVar("verb", values)
	if !ok {
		sendError(w, oai.ErrorCodeBadVerb, "", "", "", partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	switch verb {
	case "GetRecord":
		var identifier, metadataPrefix string
		for key, vals := range values {
			if len(vals) < 1 {
				continue
			}
			switch key {
			case "identifier":
				identifier = vals[0]
			case "metadataPrefix":
				metadataPrefix = vals[0]
			case "verb":
			default:
				sendError(w, oai.ErrorCodeBadArgument, "GetRecord", "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		if identifier == "" || metadataPrefix == "" {
			sendError(w, oai.ErrorCodeBadArgument, verb, identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		s.oaiHandlerGetRecord(w, req, partition, identifier, metadataPrefix)
	case "Identify":
		if len(values) > 1 {
			sendError(w, oai.ErrorCodeBadArgument, verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		s.oaiHandlerIdentify(w, req, partition)
	default:
		sendError(w, oai.ErrorCodeBadVerb, verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
	}
	return
}

func (s *Server) oaiHandlerIdentify(w http.ResponseWriter, req *http.Request, partition *fair.Partition) {
	earliestDatestamp, err := s.fair.GetMinimumDatestamp(partition.Name)
	if err != nil {
		s.log.Errorf("cannot get earliest datestamp: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("cannot get earliest datestamp")))
		return
	}
	pmh := &oai.OAIPMH{}
	pmh.InitNamespace()
	pmh.ResponseDate = time.Now().Format("2006-01-02T15:04:05Z")
	pmh.Request = &oai.Request{
		Verb:  "Identify",
		Value: partition.AddrExt + "/" + oai.APIPATH,
	}
	pmh.Identify = &oai.Identify{
		RepositoryName:    partition.OAIRepositoryName,
		BaseURL:           partition.AddrExt + "/" + oai.APIPATH,
		ProtocolVersion:   "2.0",
		EarliestDatestamp: earliestDatestamp.Format("2006-01-02T15:04:05Z"),
		AdminEmail:        partition.OAIAdminEmail,
		DeletedRecord:     "transient",
		Granularity:       "YYYY-MM-DDThh:mm:ssZ",
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(pmh); err != nil {
		s.log.Error("cannot encode pmh - %v: %v", pmh, err)
	}
}

func (s *Server) oaiHandlerGetRecord(w http.ResponseWriter, req *http.Request, partition *fair.Partition, identifier, metadataPrefix string) {
	uuidStr := strings.TrimPrefix(identifier, fmt.Sprintf("oai:%s:", partition.OAISignatureDomain))
	if uuidStr == identifier {
		s.log.Infof("invalid identifier for partition %s: %s", partition.Name, identifier)
		sendError(w, oai.ErrorCodeIdDoesNotExist, "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	data, err := s.fair.GetItem(partition.Name, uuidStr)
	if err != nil {
		s.log.Infof("cannot get item %s: %v", uuidStr, err)
		sendError(w, oai.ErrorCodeIdDoesNotExist, "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	if data == nil {
		s.log.Infof("item %s not found", uuidStr)
		sendError(w, oai.ErrorCodeIdDoesNotExist, "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	metadata := &oai.Metadata{}
	switch metadataPrefix {
	case "oai_dc":
		dcmiData := &dcmi.DCMI{}
		dcmiData.InitNamespace()
		dcmiData.FromCore(data.Metadata)
		metadata.OAIDC = dcmiData
	default:
		s.log.Infof("invalid metadataPrefix %s", metadataPrefix)
		sendError(w, oai.ErrorCodeCannotDisseminateFormat, "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	pmh := &oai.OAIPMH{}
	pmh.InitNamespace()
	pmh.ResponseDate = time.Now().Format("2006-01-02T15:04:05Z")
	pmh.Request = &oai.Request{
		Verb:           "GetRecord",
		Identifier:     identifier,
		MetadataPrefix: metadataPrefix,
		Value:          partition.AddrExt + "/" + oai.APIPATH,
	}
	pmh.GetRecord = &oai.GetRecord{Record: []*oai.Record{{
		Header: &oai.RecordHeader{
			Identifier: identifier,
			Datestamp:  data.Datestamp.Format("2006-01-02T15:04:05Z"),
			SetSpec:    data.Set,
		},
		Metadata: metadata,
		About:    nil,
	}}}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(pmh); err != nil {
		s.log.Error("cannot encode pmh - %v: %v", pmh, err)
	}
}
