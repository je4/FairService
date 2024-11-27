package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/FairService/v2/pkg/model/dcmi"
	"github.com/je4/FairService/v2/pkg/service/oai"
	"github.com/pkg/errors"
	"golang.org/x/exp/maps"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
)

// const STYLESHEET = "../static/oai2.xsl"
const STYLESHEET = "../static/dspace/oai.xsl"

var xmlHeader = fmt.Sprintf("<?xml version=\"1.0\" ?><?xml-stylesheet type=\"text/xsl\" href=\"%s\"?>", STYLESHEET)

func sendOAIError(ctx *gin.Context, code oai.ErrorCodeType, message, verb, identifier, metadataPrefix, baseURL string) {
	//	w.Header().Set("Content-type", "text/xml")
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
	//	ctx.Writer.Write([]byte(xmlHeader))
	ctx.XML(http.StatusOK, pmh)
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

func additional(arr1, arr2 []string) []string {
	var add = make([]string, 0)
	for _, a := range arr1 {
		if !slices.Contains(arr2, a) {
			add = append(add, a)
		}
	}
	return add
}

func (s *Server) oaiHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")
	context := ctx.Param("context")

	partition, err := s.fair.GetPartition(pName)
	if err != nil {
		NewResultMessage(ctx, http.StatusNotFound, errors.Wrapf(err, "partition [%s] not found", pName))
		return
	}
	if context != "request" {
		NewResultMessage(ctx, http.StatusNotFound, errors.Wrapf(err, "invalid context %s for partition: %s", context, pName))
		return
	}

	values := ctx.Request.URL.Query()
	params := maps.Keys(values)

	verb := ctx.DefaultQuery("verb", "Identify")
	switch verb {
	case "ListIdentifiers":
		if add := additional(params, []string{"verb", "from", "until", "set", "resumptionToken", "metadataPrefix"}); len(add) > 0 {
			sendOAIError(ctx, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %v", add), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		fromStr := ctx.Query("from")
		untilStr := ctx.Query("until")
		set := ctx.Query("set")
		resumptionToken := ctx.Query("resumptionToken")
		metadataPrefix := ctx.Query("metadataPrefix")
		var from, until time.Time
		if fromStr != "" {
			from, err = time.Parse("2006-01-02T15:04:05Z", fromStr)
			if err != nil {
				sendOAIError(ctx, oai.ErrorCodeBadArgument, fmt.Sprintf("error parsing from [%s]: %v", fromStr, err), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		if untilStr != "" {
			until, err = time.Parse("2006-01-02T15:04:05Z", untilStr)
			if err != nil {
				sendOAIError(ctx, oai.ErrorCodeBadArgument, fmt.Sprintf("error parsing until [%s]: %v", untilStr, err), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		s.oaiHandlerListIdentifiers(ctx, partition, context, from, until, set, resumptionToken, metadataPrefix)
	case "ListRecords":
		if add := additional(params, []string{"verb", "from", "until", "set", "resumptionToken", "metadataPrefix"}); len(add) > 0 {
			sendOAIError(ctx, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %v", add), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		fromStr := ctx.Query("from")
		untilStr := ctx.Query("until")
		set := ctx.Query("set")
		resumptionToken := ctx.Query("resumptionToken")
		metadataPrefix := ctx.Query("metadataPrefix")
		var from, until time.Time
		if fromStr != "" {
			from, err = time.Parse("2006-01-02T15:04:05Z", fromStr)
			if err != nil {
				sendOAIError(ctx, oai.ErrorCodeBadArgument, fmt.Sprintf("error parsing from [%s]: %v", fromStr, err), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		if untilStr != "" {
			until, err = time.Parse("2006-01-02T15:04:05Z", untilStr)
			if err != nil {
				sendOAIError(ctx, oai.ErrorCodeBadArgument, fmt.Sprintf("error parsing until [%s]: %v", untilStr, err), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
				return
			}
		}
		s.oaiHandlerListRecords(ctx, partition, context, from, until, set, resumptionToken, metadataPrefix)
	case "GetRecord":
		if add := additional(params, []string{"verb", "identifier", "metadataPrefix"}); len(add) > 0 {
			sendOAIError(ctx, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %v", add), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		identifier := ctx.Query("identifier")
		metadataPrefix := ctx.Query("metadataPrefix")
		if identifier == "" || metadataPrefix == "" {
			sendOAIError(ctx, oai.ErrorCodeBadArgument, "identifier AND metadataPrefix needed", verb, identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		s.oaiHandlerGetRecord(ctx, partition, context, identifier, metadataPrefix)
	case "Identify":
		if add := additional(params, []string{"verb"}); len(add) > 0 {
			sendOAIError(ctx, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %v", add), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		s.oaiHandlerIdentify(ctx, partition, context)
	case "ListSets":
		if add := additional(params, []string{"verb"}); len(add) > 0 {
			sendOAIError(ctx, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %v", add), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		s.oaiHandlerListSets(ctx, partition, context)
	case "ListMetadataFormats":
		if add := additional(params, []string{"verb", "identifier"}); len(add) > 0 {
			sendOAIError(ctx, oai.ErrorCodeBadArgument, fmt.Sprintf("unknown parameter %v", add), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		s.oaiHandlerListMetadataFormats(ctx, partition, context)
	default:
		sendOAIError(ctx, oai.ErrorCodeBadVerb, fmt.Sprintf("unknown verb: %s", verb), verb, "", "", partition.AddrExt+"/"+oai.APIPATH)
	}
	return
}

func (s *Server) oaiHandlerListMetadataFormats(ctx *gin.Context, partition *fair.Partition, context string) {
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
	ctx.XML(http.StatusOK, pmh)
}

func (s *Server) oaiHandlerIdentify(ctx *gin.Context, partition *fair.Partition, context string) {
	earliestDatestamp, err := s.fair.GetMinimumDatestamp(partition)
	if err != nil {
		NewResultMessage(ctx, http.StatusInternalServerError, errors.Wrapf(err, "cannot get earliest datestamp"))
		s.log.Error().Msgf("cannot get earliest datestamp: %v", err)
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
		RepositoryName:    partition.OAI.RepositoryName,
		BaseURL:           partition.AddrExt + "/" + oai.APIPATH,
		ProtocolVersion:   "2.0",
		EarliestDatestamp: earliestDatestamp.Format("2006-01-02T15:04:05Z"),
		AdminEmail:        partition.OAI.AdminEmail,
		DeletedRecord:     "transient",
		Granularity:       "YYYY-MM-DDThh:mm:ssZ",
		Compression:       []string{"gzip", "deflate"},
		Description: oai.Description{Identifier: oai.Identifier{
			Scheme:               partition.OAI.Scheme,
			RepositoryIdentifier: partition.Domain,
			Delimiter:            partition.OAI.Delimiter,
			SampleIdentifier:     partition.OAI.SampleIdentifier,
		}},
	}
	pmh.Identify.Description.Identifier.InitNamespace()
	ctx.XML(http.StatusOK, pmh)
}

func (s *Server) oaiHandlerGetRecord(ctx *gin.Context, partition *fair.Partition, context, identifier, metadataPrefix string) {
	uuidStr := strings.TrimPrefix(identifier, fmt.Sprintf("%s:%s:", partition.OAI.Scheme, partition.Domain))
	if uuidStr == identifier {
		s.log.Info().Msgf("invalid identifier for partition %s: %s", partition.Name, identifier)
		sendOAIError(ctx, oai.ErrorCodeIdDoesNotExist, "", "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	data, err := s.fair.GetItem(partition, uuidStr)
	if err != nil {
		s.log.Info().Msgf("cannot get item %s: %v", uuidStr, err)
		sendOAIError(ctx, oai.ErrorCodeIdDoesNotExist, "", "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}
	if data == nil {
		s.log.Info().Msgf("item %s not found", uuidStr)
		sendOAIError(ctx, oai.ErrorCodeIdDoesNotExist, "", "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
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
		dataciteData := &dataciteModel.DataCite{}
		dataciteData.InitNamespace()
		dataciteData.FromCore(data.Metadata)
		metadata.Datacite = dataciteData
	default:
		s.log.Info().Msgf("invalid metadataPrefix %s", metadataPrefix)
		sendOAIError(ctx, oai.ErrorCodeCannotDisseminateFormat, fmt.Sprintf("invalid metadataPrefix %s", metadataPrefix), "GetRecord", identifier, metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
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
	ctx.XML(http.StatusOK, pmh)
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

func (s *Server) oaiHandlerListIdentifiers(ctx *gin.Context, partition *fair.Partition, context string, from, until time.Time, set, resumptionToken, metadataPrefix string) {
	var rData *resumptionData
	if resumptionToken != "" {
		data, err := s.resumptionTokenCache.Get(resumptionToken)
		if err != nil || data == nil {
			sendOAIError(ctx, oai.ErrorCodeBadResumptionToken, fmt.Sprintf("cannot load resumption data for %s: %v", resumptionToken, err), "ListenIdentifier", "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		var ok bool
		rData, ok = data.(*resumptionData)
		if !ok {
			s.log.Error().Msgf("invalid resumption data: %v", data)
			s.resumptionTokenCache.Remove(resumptionToken)
			NewResultMessage(ctx, http.StatusInternalServerError, errors.Wrapf(err, "invalid resumption data: %v", data))
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
		if count < partition.OAI.PageSize {
			var status oai.RecordHeaderStatusType = oai.RecordHeaderStatusOK
			if item.Status != fair.DataStatusActive {
				status = oai.RecordHeaderStatusDeleted
			}
			header := &oai.RecordHeader{
				Identifier: fmt.Sprintf("%s:%s:%s", partition.OAI.Scheme, partition.Domain, item.UUID),
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
				ExpirationDate:   time.Now().Add(partition.OAI.ResumptionTokenTimeout).Format("2006-01-02T15:04:05Z"),
				Value:            uuidStr,
				Cursor:           rData.cursor,
				CompleteListSize: rData.completeListSize,
			}
			rData.cursor++
			lastToken := rData.lastToken
			rData.lastToken = resumptionToken
			if err := s.resumptionTokenCache.SetWithExpire(uuidStr, rData, partition.OAI.ResumptionTokenTimeout); err != nil {
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
		if err := s.fair.GetItemsSeq(partition,
			rData.seq,
			until,
			[]fair.DataAccess{fair.DataAccessPublic, fair.DataAccessPublic},
			partition.OAI.PageSize+1,
			0,
			nil,
			itemFunc); err != nil {
			NewResultMessage(ctx, http.StatusInternalServerError, errors.Wrapf(err, "cannot read ItemsSeq content from database"))
			return
		}

	} else {
		if err := s.fair.GetItemsDatestamp(partition,
			rData.from,
			rData.until,
			set,
			[]fair.DataAccess{fair.DataAccessPublic, fair.DataAccessClosedData},
			partition.OAI.PageSize+1,
			0,
			&rData.completeListSize,
			itemFunc); err != nil {
			NewResultMessage(ctx, http.StatusInternalServerError, errors.Wrapf(err, "cannot read Datestamp content from database"))
			return
		}
	}

	if len(listIdentifiers.Header) == 0 {
		sendOAIError(ctx, oai.ErrorCodeNoRecordsMatch, "no records match", "ListenIdentifier", "", metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
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
	ctx.XML(http.StatusOK, pmh)

}

func (s *Server) oaiHandlerListRecords(ctx *gin.Context, partition *fair.Partition, context string, from, until time.Time, set, resumptionToken, metadataPrefix string) {
	if metadataPrefix == "" {
		metadataPrefix = "oai_dc"
	}
	switch metadataPrefix {
	case "oai_dc":
	case "oai_datacite":
	default:
		s.log.Info().Msgf("invalid metadataPrefix %s", metadataPrefix)
		sendOAIError(ctx, oai.ErrorCodeCannotDisseminateFormat, fmt.Sprintf("invalid metadataPrefix %s", metadataPrefix), "ListRecords", "", metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
		return
	}

	var rData *resumptionData
	if resumptionToken != "" {
		data, err := s.resumptionTokenCache.Get(resumptionToken)
		if err != nil || data == nil {
			sendOAIError(ctx, oai.ErrorCodeBadResumptionToken, fmt.Sprintf("cannot load resumption data for %s: %v", resumptionToken, err), "ListRecords", "", "", partition.AddrExt+"/"+oai.APIPATH)
			return
		}
		var ok bool
		rData, ok = data.(*resumptionData)
		if !ok {
			NewResultMessage(ctx, http.StatusInternalServerError, errors.Wrapf(err, "invalid resumption data: %v", data))
			s.log.Error().Msgf("invalid resumption data: %v", data)
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
		if count < partition.OAI.PageSize {
			var status oai.RecordHeaderStatusType = oai.RecordHeaderStatusOK
			if item.Status != fair.DataStatusActive {
				status = oai.RecordHeaderStatusDeleted
			}
			record := &oai.Record{
				Header: &oai.RecordHeader{
					Identifier: fmt.Sprintf("%s:%s:%s", partition.OAI.Scheme, partition.Domain, item.UUID),
					Datestamp:  item.Datestamp.Format("2006-01-02T15:04:05Z"),
					SetSpec:    item.Set,
					Status:     status,
				},
			}
			if item.Status == fair.DataStatusActive {
				metadata := &oai.Metadata{}
				switch metadataPrefix {
				case "oai_dc":
					dcmiData := &dcmi.DCMI{}
					dcmiData.InitNamespace()
					dcmiData.FromCore(item.Metadata)
					metadata.OAIDC = dcmiData
				case "oai_datacite":
					dataciteData := &dataciteModel.DataCite{}
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
				ExpirationDate:   time.Now().Add(partition.OAI.ResumptionTokenTimeout).Format("2006-01-02T15:04:05Z"),
				Value:            uuidStr,
				Cursor:           rData.cursor,
				CompleteListSize: rData.completeListSize,
			}
			rData.cursor++
			lastToken := rData.lastToken
			rData.lastToken = resumptionToken
			if err := s.resumptionTokenCache.SetWithExpire(uuidStr, rData, partition.OAI.ResumptionTokenTimeout); err != nil {
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
		if err := s.fair.GetItemsSeq(partition,
			rData.seq,
			until,
			[]fair.DataAccess{fair.DataAccessPublic, fair.DataAccessPublic},
			partition.OAI.PageSize+1,
			0,
			nil,
			itemFunc); err != nil {
			sendResult(s.log, ctx, http.StatusInternalServerError, errors.Wrapf(err, "cannot read content from database").Error(), nil)
			return
		}

	} else {
		if err := s.fair.GetItemsDatestamp(partition,
			rData.from,
			rData.until,
			set,
			[]fair.DataAccess{fair.DataAccessPublic, fair.DataAccessClosedData},
			partition.OAI.PageSize+1,
			0,
			&rData.completeListSize,
			itemFunc); err != nil {
			sendResult(s.log, ctx, http.StatusInternalServerError, errors.Wrapf(err, "cannot read content from database").Error(), nil)
			return
		}
	}

	if len(listRecords.Record) == 0 {
		sendOAIError(ctx, oai.ErrorCodeNoRecordsMatch, "no records match", "ListRecords", "", metadataPrefix, partition.AddrExt+"/"+oai.APIPATH)
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

	ctx.XML(http.StatusOK, pmh)
}

// todo: add paging with resumption token
func (s *Server) oaiHandlerListSets(ctx *gin.Context, partition *fair.Partition, context string) {
	sets, err := s.fair.GetSets(partition)
	if err != nil {
		NewResultMessage(ctx, http.StatusInternalServerError, errors.Wrapf(err, "cannot read sets from database"))
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
	ctx.XML(http.StatusOK, pmh)
}
