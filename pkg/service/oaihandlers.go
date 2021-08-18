package service

import (
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/datacite"
	"github.com/je4/FairService/v2/pkg/model/dcmi"
	"github.com/je4/FairService/v2/pkg/service/oai"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// const STYLESHEET = "../static/oai2.xsl"
const STYLESHEET = "../static/dspace/oai.xsl"

func sendError(w http.ResponseWriter, code oai.ErrorCodeType, message, verb, identifier, metadataPrefix, baseURL string) {
	w.Header().Set("Content-type", "text/xml")
	pmh := &oai.OAIPMH{}
	pmh.InitNamespace()
	pmh.Error = &oai.Error{
		Code:  code,
		Value: message,
	}
	pmh.ResponseDate = time.Now().Format("2006-01-02T15:04:05Z")
	pmh.Request = &oai.Request{
		Verb:           verb,
		Identifier:     identifier,
		MetadataPrefix: metadataPrefix,
		Value:          baseURL,
	}
	w.Write([]byte(fmt.Sprintf("<?xml version=\"1.0\" ?><?xml-stylesheet type=\"text/xsl\" href=\"%s\"?>", STYLESHEET)))
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
	context := vars["context"]

	partition, err := s.fair.GetPartition(pName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("invalid partition: %s", pName)))
		return
	}
	if context != "request" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("invalid context %s for partition: %s", context, pName)))
		return
	}
	values := req.URL.Query()
	verb, ok := getVar("verb", values)
	if !ok {
		verb = "Identify"
	}
	switch verb {
	case "ListIdentifiers":
		var fromStr, untilStr, set, resumptionToken, metadataPrefix string
		for key, vals := range values {
			if len(vals) < 1 {
				continue
			}
			switch key {
			case "from":
				fromStr = vals[0]
			case "until":
				untilStr = vals[0]
			case "set":
				set = vals[0]
			case "resumptionToken":
				resumptionToken = vals[0]
			case "metadataPrefix":
				metadataPrefix = vals[0]
			case "verb":
			default:
				sendError(w, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %s", key), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		var from, until time.Time
		if fromStr != "" {
			from, err = time.Parse("2006-01-02T15:04:05Z", fromStr)
			if err != nil {
				sendError(w, oai.ErrorCodeBadArgument, fmt.Sprintf("error parsing from [%s]: %v", fromStr, err), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		if untilStr != "" {
			until, err = time.Parse("2006-01-02T15:04:05Z", untilStr)
			if err != nil {
				sendError(w, oai.ErrorCodeBadArgument, fmt.Sprintf("error parsing until [%s]: %v", untilStr, err), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		s.oaiHandlerListIdentifiers(w, req, partition, context, from, until, set, resumptionToken, metadataPrefix)
	case "ListRecords":
		var fromStr, untilStr, set, resumptionToken, metadataPrefix string
		for key, vals := range values {
			if len(vals) < 1 {
				continue
			}
			switch key {
			case "from":
				fromStr = vals[0]
			case "until":
				untilStr = vals[0]
			case "set":
				set = vals[0]
			case "resumptionToken":
				resumptionToken = vals[0]
			case "metadataPrefix":
				metadataPrefix = vals[0]
			case "verb":
			default:
				sendError(w, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %s", key), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		var from, until time.Time
		if fromStr != "" {
			from, err = time.Parse("2006-01-02T15:04:05Z", fromStr)
			if err != nil {
				sendError(w, oai.ErrorCodeBadArgument, fmt.Sprintf("error parsing from [%s]: %v", fromStr, err), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		if untilStr != "" {
			until, err = time.Parse("2006-01-02T15:04:05Z", untilStr)
			if err != nil {
				sendError(w, oai.ErrorCodeBadArgument, fmt.Sprintf("error parsing until [%s]: %v", untilStr, err), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		s.oaiHandlerListRecords(w, req, partition, context, from, until, set, resumptionToken, metadataPrefix)
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
				sendError(w, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %s", key), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		if identifier == "" || metadataPrefix == "" {
			sendError(w, oai.ErrorCodeBadArgument, "identifier AND metadataPrefix needed", verb, identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		s.oaiHandlerGetRecord(w, req, partition, context, identifier, metadataPrefix)
	case "Identify":
		for key, _ := range values {
			if key == "verb" {
				continue
			}
			sendError(w, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %s", key), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		s.oaiHandlerIdentify(w, req, partition, context)
	case "ListSets":
		for key, _ := range values {
			if key == "verb" {
				continue
			}
			sendError(w, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %s", key), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		s.oaiHandlerListSets(w, req, partition, context)
	case "ListMetadataFormats":
		for key, _ := range values {
			if key == "verb" || key == "identifier" {
				continue
			}
			sendError(w, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %s", key), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		s.oaiHandlerListMetadataFormats(w, req, partition, context)
	default:
		sendError(w, oai.ErrorCodeBadVerb, fmt.Sprintf("unknown verb: %s", verb), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
	}
	return
}

func (s *Server) oaiHandlerListMetadataFormats(w http.ResponseWriter, req *http.Request, partition *fair.Partition, context string) {
	pmh := &oai.OAIPMH{}
	pmh.InitNamespace()
	pmh.ResponseDate = time.Now().Format("2006-01-02T15:04:05Z")
	pmh.Request = &oai.Request{
		Verb:  "ListMetadataFormats",
		Value: fmt.Sprintf("%s/%s/%s", partition.AddrExt, oai.APIPATH, context),
	}
	pmh.ListMetadataFormats = &oai.ListMetadataFormats{
		MetadataFormat: []*oai.MetadataFormat{{
			MetadataPrefix:    "oai_dc",
			Schema:            "https://www.openarchives.org/OAI/2.0/oai_dc.xsd",
			MetadataNamespace: "https://www.openarchives.org/OAI/2.0/oai_dc/",
		},
			{
				MetadataPrefix:    "oai_datacite",
				Schema:            "http://schema.datacite.org/meta/kernel-4.1/metadata.xsd",
				MetadataNamespace: "http://datacite.org/schema/kernel-4",
			},
		},
	}
	/*
		pmh.ListMetadataFormats.ResumptionToken = &oai.ResumptionToken{
			ExpirationDate:   "",
			Value:            "",
			Cursor:           0,
			CompleteListSize: 2,
		}
	*/

	w.Write([]byte(fmt.Sprintf("<?xml version=\"1.0\" ?><?xml-stylesheet type=\"text/xsl\" href=\"%s\"?>", STYLESHEET)))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(pmh); err != nil {
		s.log.Error("cannot encode pmh - %v: %v", pmh, err)
	}
}

func (s *Server) oaiHandlerIdentify(w http.ResponseWriter, req *http.Request, partition *fair.Partition, context string) {
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
		Value: fmt.Sprintf("%s/%s/%s", partition.AddrExt, oai.APIPATH, context),
	}
	pmh.Identify = &oai.Identify{
		RepositoryName:    partition.OAIRepositoryName,
		BaseURL:           partition.AddrExt + "/" + oai.APIPATH,
		ProtocolVersion:   "2.0",
		EarliestDatestamp: earliestDatestamp.Format("2006-01-02T15:04:05Z"),
		AdminEmail:        partition.OAIAdminEmail,
		DeletedRecord:     "transient",
		Granularity:       "YYYY-MM-DDThh:mm:ssZ",
		Compression:       []string{"gzip", "deflate"},
		Description: oai.Description{Identifier: oai.Identifier{
			Scheme:               partition.OAIScheme,
			RepositoryIdentifier: partition.Domain,
			Delimiter:            partition.OAIDelimiter,
			SampleIdentifier:     partition.OAISampleIdentifier,
		}},
	}
	pmh.Identify.Description.Identifier.InitNamespace()
	w.Write([]byte(fmt.Sprintf("<?xml version=\"1.0\" ?><?xml-stylesheet type=\"text/xsl\" href=\"%s\"?>", STYLESHEET)))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(pmh); err != nil {
		s.log.Error("cannot encode pmh - %v: %v", pmh, err)
	}
}

func (s *Server) oaiHandlerGetRecord(w http.ResponseWriter, req *http.Request, partition *fair.Partition, context, identifier, metadataPrefix string) {
	uuidStr := strings.TrimPrefix(identifier, fmt.Sprintf("%s:%s:", partition.OAIScheme, partition.Domain))
	if uuidStr == identifier {
		s.log.Infof("invalid identifier for partition %s: %s", partition.Name, identifier)
		sendError(w, oai.ErrorCodeIdDoesNotExist, "", "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	data, err := s.fair.GetItem(partition.Name, uuidStr)
	if err != nil {
		s.log.Infof("cannot get item %s: %v", uuidStr, err)
		sendError(w, oai.ErrorCodeIdDoesNotExist, "", "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	if data == nil {
		s.log.Infof("item %s not found", uuidStr)
		sendError(w, oai.ErrorCodeIdDoesNotExist, "", "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	metadata := &oai.Metadata{}
	if metadataPrefix == "" {
		metadataPrefix = "oai_dc"
	}
	switch metadataPrefix {
	case "oai_dc":
		dcmiData := &dcmi.DCMI{}
		dcmiData.InitNamespace()
		dcmiData.FromCore(data.Metadata)
		metadata.OAIDC = dcmiData
	case "oai_datacite":
		dataciteData := &datacite.DataCite{}
		dataciteData.InitNamespace()
		dataciteData.FromCore(data.Metadata)
		metadata.Datacite = dataciteData
	default:
		s.log.Infof("invalid metadataPrefix %s", metadataPrefix)
		sendError(w, oai.ErrorCodeCannotDisseminateFormat, fmt.Sprintf("invalid metadataPrefix %s", metadataPrefix), "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	pmh := &oai.OAIPMH{}
	pmh.InitNamespace()
	pmh.ResponseDate = time.Now().Format("2006-01-02T15:04:05Z")
	pmh.Request = &oai.Request{
		Verb:           "GetRecord",
		Identifier:     identifier,
		MetadataPrefix: metadataPrefix,
		Value:          fmt.Sprintf("%s/%s/%s", partition.AddrExt, oai.APIPATH, context),
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
	w.Write([]byte(fmt.Sprintf("<?xml version=\"1.0\" ?><?xml-stylesheet type=\"text/xsl\" href=\"%s\"?>", STYLESHEET)))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(pmh); err != nil {
		s.log.Error("cannot encode pmh - %v: %v", pmh, err)
	}
}

type resumptionData struct {
	from, until      time.Time
	partitionName    string
	set              string
	metadataPrefix   string
	seq              int64
	cursor           int64
	completeListSize int64
	lastToken        string
}

func (s *Server) oaiHandlerListIdentifiers(w http.ResponseWriter, req *http.Request, partition *fair.Partition, context string, from, until time.Time, set, resumptionToken, metadataPrefix string) {
	var rData *resumptionData
	if resumptionToken != "" {
		data, err := s.resumptionTokenCache.Get(resumptionToken)
		if err != nil || data == nil {
			sendError(w, oai.ErrorCodeBadResumptionToken, fmt.Sprintf("cannot load resumption data for %s: %v", resumptionToken, err), "ListenIdentifier", "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		var ok bool
		rData, ok = data.(*resumptionData)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("invalid resumption data: %v", data)))
			s.log.Errorf("invalid resumption data: %v", data)
			s.resumptionTokenCache.Remove(resumptionToken)
			return
		}
		//s.resumptionTokenCache.Remove(resumptionToken)
	} else {
		rData = &resumptionData{
			from:             from,
			until:            until,
			partitionName:    partition.Name,
			set:              set,
			metadataPrefix:   metadataPrefix,
			cursor:           0,
			completeListSize: 0,
			lastToken:        "",
		}
	}

	listIdentifiers := &oai.ListIdentifiers{
		Header:          []*oai.RecordHeader{},
		ResumptionToken: nil,
	}

	var count int64

	itemFunc := func(item *fair.ItemData) error {
		if count < partition.OAIPagesize {
			var status oai.RecordHeaderStatusType = oai.RecordHeaderStatusOK
			if item.Deleted {
				status = oai.RecordHeaderStatusDeleted
			}
			header := &oai.RecordHeader{
				Identifier: fmt.Sprintf("%s:%s:%s", partition.OAIScheme, partition.Domain, item.UUID),
				Datestamp:  item.Datestamp.Format("2006-01-02T15:04:05Z"),
				SetSpec:    item.Set,
				Status:     status,
			}
			listIdentifiers.Header = append(listIdentifiers.Header, header)
		} else {
			rData.seq = item.Seq
			uuidData, err := uuid.NewUUID()
			if err != nil {
				return errors.Wrapf(err, "cannot create uuid")
			}
			uuidStr := uuidData.String()
			listIdentifiers.ResumptionToken = &oai.ResumptionToken{
				ExpirationDate:   time.Now().Add(partition.ResumptionTokenTimeout).Format("2006-01-02T15:04:05Z"),
				Value:            uuidStr,
				Cursor:           rData.cursor,
				CompleteListSize: rData.completeListSize,
			}
			rData.cursor++
			lastToken := rData.lastToken
			rData.lastToken = resumptionToken
			if err := s.resumptionTokenCache.SetWithExpire(uuidStr, rData, partition.ResumptionTokenTimeout); err != nil {
				return errors.Wrapf(err, "cannot store resumption data")
			}
			s.resumptionTokenCache.Remove(lastToken)
		}
		count++
		if listIdentifiers.ResumptionToken == nil {
			listIdentifiers.ResumptionToken = &oai.ResumptionToken{
				ExpirationDate:   "",
				Value:            "",
				Cursor:           rData.cursor,
				CompleteListSize: rData.completeListSize,
			}
		}
		return nil
	}

	if rData.seq > 0 {
		if err := s.fair.GetItemsSeq(partition.Name,
			rData.seq,
			until,
			[]fair.DataAccess{fair.DataAccessPublic, fair.DataAccessPublic},
			partition.OAIPagesize+1,
			0,
			nil,
			itemFunc); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("cannot read content from database: %v", err)))
			return
		}

	} else {
		if err := s.fair.GetItemsDatestamp(partition.Name,
			rData.from,
			rData.until,
			[]fair.DataAccess{fair.DataAccessPublic, fair.DataAccessClosedData},
			partition.OAIPagesize+1,
			0,
			&rData.completeListSize,
			itemFunc); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("cannot read content from database: %v", err)))
			return
		}
	}

	if len(listIdentifiers.Header) == 0 {
		sendError(w, oai.ErrorCodeNoRecordsMatch, "no records match", "ListenIdentifier", "", metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return

	}

	pmh := &oai.OAIPMH{}
	pmh.InitNamespace()
	pmh.ResponseDate = time.Now().Format("2006-01-02T15:04:05Z")
	pmh.Request = &oai.Request{
		Verb:           "ListIdentifiers",
		MetadataPrefix: metadataPrefix,
		Value:          fmt.Sprintf("%s/%s/%s", partition.AddrExt, oai.APIPATH, context),
	}
	pmh.ListIdentifiers = listIdentifiers
	w.Write([]byte(fmt.Sprintf("<?xml version=\"1.0\" ?><?xml-stylesheet type=\"text/xsl\" href=\"%s\"?>", STYLESHEET)))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(pmh); err != nil {
		s.log.Error("cannot encode pmh - %v: %v", pmh, err)
	}

}

func (s *Server) oaiHandlerListRecords(w http.ResponseWriter, req *http.Request, partition *fair.Partition, context string, from, until time.Time, set, resumptionToken, metadataPrefix string) {
	if metadataPrefix == "" {
		metadataPrefix = "oai_dc"
	}
	switch metadataPrefix {
	case "oai_dc":
	case "oai_datacite":
	default:
		s.log.Infof("invalid metadataPrefix %s", metadataPrefix)
		sendError(w, oai.ErrorCodeCannotDisseminateFormat, fmt.Sprintf("invalid metadataPrefix %s", metadataPrefix), "ListRecords", "", metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}

	var rData *resumptionData
	if resumptionToken != "" {
		data, err := s.resumptionTokenCache.Get(resumptionToken)
		if err != nil || data == nil {
			sendError(w, oai.ErrorCodeBadResumptionToken, fmt.Sprintf("cannot load resumption data for %s: %v", resumptionToken, err), "ListRecords", "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		var ok bool
		rData, ok = data.(*resumptionData)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("invalid resumption data: %v", data)))
			s.log.Errorf("invalid resumption data: %v", data)
			s.resumptionTokenCache.Remove(resumptionToken)
			return
		}
		//s.resumptionTokenCache.Remove(resumptionToken)
	} else {
		rData = &resumptionData{
			from:             from,
			until:            until,
			partitionName:    partition.Name,
			set:              set,
			metadataPrefix:   metadataPrefix,
			cursor:           0,
			completeListSize: 0,
			lastToken:        "",
		}
	}

	listRecords := &oai.ListRecords{
		Record:          []*oai.Record{},
		ResumptionToken: nil,
	}

	var count int64

	itemFunc := func(item *fair.ItemData) error {
		if count < partition.OAIPagesize {
			var status oai.RecordHeaderStatusType = oai.RecordHeaderStatusOK
			if item.Deleted {
				status = oai.RecordHeaderStatusDeleted
			}
			record := &oai.Record{
				Header: &oai.RecordHeader{
					Identifier: fmt.Sprintf("%s:%s:%s", partition.OAIScheme, partition.Domain, item.UUID),
					Datestamp:  item.Datestamp.Format("2006-01-02T15:04:05Z"),
					SetSpec:    item.Set,
					Status:     status,
				},
			}
			if !item.Deleted {
				metadata := &oai.Metadata{}
				switch metadataPrefix {
				case "oai_dc":
					dcmiData := &dcmi.DCMI{}
					dcmiData.InitNamespace()
					dcmiData.FromCore(item.Metadata)
					metadata.OAIDC = dcmiData
				case "oai_datacite":
					dataciteData := &datacite.DataCite{}
					dataciteData.InitNamespace()
					dataciteData.FromCore(item.Metadata)
					metadata.Datacite = dataciteData
				default:
					return errors.New(fmt.Sprintf("invalid metadataPrefix %s", metadataPrefix))
				}
				record.Metadata = metadata
			}
			listRecords.Record = append(listRecords.Record, record)
		} else {
			rData.seq = item.Seq
			uuidData, err := uuid.NewUUID()
			if err != nil {
				return errors.Wrapf(err, "cannot create uuid")
			}
			uuidStr := uuidData.String()
			listRecords.ResumptionToken = &oai.ResumptionToken{
				ExpirationDate:   time.Now().Add(partition.ResumptionTokenTimeout).Format("2006-01-02T15:04:05Z"),
				Value:            uuidStr,
				Cursor:           rData.cursor,
				CompleteListSize: rData.completeListSize,
			}
			rData.cursor++
			lastToken := rData.lastToken
			rData.lastToken = resumptionToken
			if err := s.resumptionTokenCache.SetWithExpire(uuidStr, rData, partition.ResumptionTokenTimeout); err != nil {
				return errors.Wrapf(err, "cannot store resumption data")
			}
			s.resumptionTokenCache.Remove(lastToken)
		}
		count++
		if listRecords.ResumptionToken == nil {
			listRecords.ResumptionToken = &oai.ResumptionToken{
				ExpirationDate:   "",
				Value:            "",
				Cursor:           rData.cursor,
				CompleteListSize: rData.completeListSize,
			}
		}
		return nil
	}

	if rData.seq > 0 {
		if err := s.fair.GetItemsSeq(partition.Name,
			rData.seq,
			until,
			[]fair.DataAccess{fair.DataAccessPublic, fair.DataAccessPublic},
			partition.OAIPagesize+1,
			0,
			nil,
			itemFunc); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("cannot read content from database: %v", err)))
			return
		}

	} else {
		if err := s.fair.GetItemsDatestamp(partition.Name,
			rData.from,
			rData.until,
			[]fair.DataAccess{fair.DataAccessPublic, fair.DataAccessClosedData},
			partition.OAIPagesize+1,
			0,
			&rData.completeListSize,
			itemFunc); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("cannot read content from database: %v", err)))
			return
		}
	}

	if len(listRecords.Record) == 0 {
		sendError(w, oai.ErrorCodeNoRecordsMatch, "no records match", "ListRecords", "", metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return

	}

	pmh := &oai.OAIPMH{}
	pmh.InitNamespace()
	pmh.ResponseDate = time.Now().Format("2006-01-02T15:04:05Z")
	pmh.Request = &oai.Request{
		Verb:           "ListRecords",
		MetadataPrefix: metadataPrefix,
		Value:          fmt.Sprintf("%s/%s/%s", partition.AddrExt, oai.APIPATH, context),
	}
	pmh.ListRecords = listRecords
	w.Write([]byte(fmt.Sprintf("<?xml version=\"1.0\" ?><?xml-stylesheet type=\"text/xsl\" href=\"%s\"?>", STYLESHEET)))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(pmh); err != nil {
		s.log.Error("cannot encode pmh - %v: %v", pmh, err)
	}

}

// todo: add paging with resumption token
func (s *Server) oaiHandlerListSets(w http.ResponseWriter, req *http.Request, partition *fair.Partition, context string) {
	sets, err := s.fair.GetSets(partition.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("cannot read sets from database: %v", err)))
		return
	}
	listSets := &oai.ListSets{
		Set:             []*oai.Set{},
		ResumptionToken: nil,
	}
	for setspec, setname := range sets {
		set := &oai.Set{
			SetSpec: setspec,
			SetName: setname,
		}
		listSets.Set = append(listSets.Set, set)
	}
	pmh := &oai.OAIPMH{}
	pmh.InitNamespace()
	pmh.ResponseDate = time.Now().Format("2006-01-02T15:04:05Z")
	pmh.Request = &oai.Request{
		Verb:  "ListSets",
		Value: fmt.Sprintf("%s/%s/%s", partition.AddrExt, oai.APIPATH, context),
	}
	pmh.ListSets = listSets
	w.Write([]byte(fmt.Sprintf("<?xml version=\"1.0\" ?><?xml-stylesheet type=\"text/xsl\" href=\"%s\"?>", STYLESHEET)))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(pmh); err != nil {
		s.log.Error("cannot encode pmh - %v: %v", pmh, err)
	}
}
