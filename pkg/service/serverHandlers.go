package service

import (
	"crypto/subtle"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/dcmi"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type CreateResultStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	UUID    string `json:"uuid,omitempty"`
}

func BasicAuth(w http.ResponseWriter, r *http.Request, username, password, realm string) bool {

	user, pass, ok := r.BasicAuth()

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return false
	}

	return true
}

func (s *Server) redirectHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]
	uuidStr := vars["uuid"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("partition [%s] not found", pName)))
		return
	}
	data, err := s.fair.GetItem(pName, uuidStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error loading item %s/%s: %v", pName, uuidStr, err)))
		return
	}
	source, err := s.fair.GetSourceByName(data.Source, part.Name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error loading source %s for item %s/%s: %v", data.Source, pName, uuidStr, err)))
		return
	}

	targetURL := strings.Replace(source.DetailURL, "{signature}", data.Signature, -1)

	http.Redirect(w, req, targetURL, 301)
}

func (s *Server) partitionHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("partition [%s] not found", pName)))
		return
	}

	tpl := s.templates["partition"]
	if err := tpl.Execute(w, part); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error executing template %s in partition %s: %v", "partition", pName, err)))
		return
	}
}

func (s *Server) dataviewerHandler(w http.ResponseWriter, req *http.Request) {
	if !BasicAuth(w, req, s.name, s.password, "FAIR Service") {
		return
	}

	vars := mux.Vars(req)
	pName := vars["partition"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("partition [%s] not found", pName)))
		return
	}

	tpl := s.templates["dataviewer"]
	if err := tpl.Execute(w, part); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error executing template %s in partition %s: %v", "partition", pName, err)))
		return
	}
}

func (s *Server) partitionOAIHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("partition [%s] not found", pName)))
		return
	}

	tpl := s.templates["oai"]
	if err := tpl.Execute(w, part); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error executing template %s in partition %s: %v", "partition", pName, err)))
		return
	}
}

func (s *Server) itemHandler(w http.ResponseWriter, req *http.Request) {
	sendResult := func(t string, message string, uuid string) {
		if t == "ok" {
			s.log.Infof(fmt.Sprintf("%s: %s", message, uuid))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error(fmt.Sprintf("%s: %s", message, uuid))
		}
		w.Header().Set("Content-type", "text/json")
		data, _ := json.MarshalIndent(CreateResultStatus{Status: t, Message: message, UUID: uuid}, "", "  ")
		w.Write(data)
	}
	vars := mux.Vars(req)
	pName := vars["partition"]
	uuidStr := vars["uuid"]
	outputType := vars["outputType"]
	if outputType == "" {
		outputType = "json"
	}

	data, err := s.fair.GetItem(pName, uuidStr)
	if err != nil {
		sendResult("error", fmt.Sprintf("error loading item: %v", err), uuidStr)
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("item [%s] not found in partition %s", uuidStr, pName)))
		return
	}

	switch outputType {
	case "json":
		w.Header().Set("Content-type", "text/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			sendResult("error", fmt.Sprintf("cannot marshal data"), uuidStr)
			return
		}
		return
	case "dcmi":
		dcmiData := &dcmi.DCMI{}
		dcmiData.InitNamespace()
		dcmiData.FromCore(data.Metadata)
		w.Header().Set("Content-type", "text/xml")
		enc := xml.NewEncoder(w)
		enc.Indent("", "  ")
		w.Header().Set("Content-type", "text/json")
		if err := enc.Encode(dcmiData); err != nil {
			sendResult("error", fmt.Sprintf("cannot marshal data"), uuidStr)
			return
		}
		return
	default:
		sendResult("error", fmt.Sprintf("invalid output type %s", outputType), uuidStr)
		return
	}
}

func (s *Server) createDOIHandler(w http.ResponseWriter, req *http.Request) {
	if !BasicAuth(w, req, s.name, s.password, "FAIR Service") {
		return
	}

	sendResult := func(t string, message string, uuid string) {
		if t == "ok" {
			s.log.Infof(fmt.Sprintf("%s: %s", message, uuid))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			s.log.Error(fmt.Sprintf("%s: %s", message, uuid))
		}
		w.Header().Set("Content-type", "text/json")
		data, _ := json.MarshalIndent(CreateResultStatus{Status: t, Message: message, UUID: uuid}, "", "  ")
		w.Write(data)
	}

	vars := mux.Vars(req)
	pName := vars["partition"]
	uuidStr := vars["uuid"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult("error", fmt.Sprintf("cannot get partition %s", pName), uuidStr)
		return
	}

	targetUrl := fmt.Sprintf("%s/redir/%s", part.AddrExt, uuidStr)
	doiResult, err := s.fair.CreateDOI(pName, uuidStr, targetUrl)
	if err != nil {
		sendResult("error", err.Error(), uuidStr)
		return
	}

	w.Header().Set("Content-type", "text/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doiResult); err != nil {
		sendResult("error", fmt.Sprintf("cannot marshal data"), uuidStr)
		return
	}
	return
}

func (s *Server) startUpdateHandler(w http.ResponseWriter, req *http.Request) {
	sendResult := func(t string, message string) {
		if t == "ok" {
			s.log.Infof(fmt.Sprintf("%s", message))
		} else {
			s.log.Error(fmt.Sprintf("%s", message))
		}
		w.Header().Set("Content-type", "text/json")
		data, _ := json.MarshalIndent(CreateResultStatus{Status: t, Message: message, UUID: ""}, "", "  ")
		w.Write(data)
	}

	vars := mux.Vars(req)
	pName := vars["partition"]

	decoder := json.NewDecoder(req.Body)
	var data fair.SourceData
	err := decoder.Decode(&data)
	if err != nil {
		sendResult("error", fmt.Sprintf("cannot parse request body: %v", err))
		return
	}

	if err := s.fair.StartUpdate(pName, data.Source); err != nil {
		sendResult("error", fmt.Sprintf("cannot start update for %s on %s: %v", data.Source, pName, err))
		return
	}
	sendResult("ok", fmt.Sprintf("starting update for %s on %s", data.Source, pName))
}

func (s *Server) endUpdateHandler(w http.ResponseWriter, req *http.Request) {
	sendResult := func(t string, message string) {
		if t == "ok" {
			s.log.Infof(fmt.Sprintf("%s", message))
		} else {
			s.log.Error(fmt.Sprintf("%s", message))
		}
		w.Header().Set("Content-type", "text/json")
		data, _ := json.MarshalIndent(CreateResultStatus{Status: t, Message: message, UUID: ""}, "", "  ")
		w.Write(data)
	}

	vars := mux.Vars(req)
	pName := vars["partition"]

	decoder := json.NewDecoder(req.Body)
	var data fair.SourceData
	err := decoder.Decode(&data)
	if err != nil {
		sendResult("error", fmt.Sprintf("cannot parse request body: %v", err))
		return
	}

	if err := s.fair.EndUpdate(pName, data.Source); err != nil {
		sendResult("error", fmt.Sprintf("cannot end update for %s on %s: %v", data.Source, pName, err))
		return
	}
	sendResult("ok", fmt.Sprintf("end update for %s on %s", data.Source, pName))
}

func (s *Server) abortUpdateHandler(w http.ResponseWriter, req *http.Request) {
	sendResult := func(t string, message string) {
		if t == "ok" {
			s.log.Infof(fmt.Sprintf("%s", message))
		} else {
			s.log.Error(fmt.Sprintf("%s", message))
		}
		w.Header().Set("Content-type", "text/json")
		data, _ := json.MarshalIndent(CreateResultStatus{Status: t, Message: message, UUID: ""}, "", "  ")
		w.Write(data)
	}

	vars := mux.Vars(req)
	pName := vars["partition"]

	decoder := json.NewDecoder(req.Body)
	var data fair.SourceData
	err := decoder.Decode(&data)
	if err != nil {
		sendResult("error", fmt.Sprintf("cannot parse request body: %v", err))
		return
	}

	if err := s.fair.AbortUpdate(pName, data.Source); err != nil {
		sendResult("error", fmt.Sprintf("cannot end update for %s on %s: %v", data.Source, pName, err))
		return
	}
	sendResult("ok", fmt.Sprintf("abort update for %s on %s", data.Source, pName))
}

func (s *Server) createHandler(w http.ResponseWriter, req *http.Request) {
	sendResult := func(t string, message string, uuid string) {
		if t == "ok" {
			s.log.Infof(fmt.Sprintf("%s: %s", message, uuid))
		} else {
			s.log.Error(fmt.Sprintf("%s: %s", message, uuid))
		}
		w.Header().Set("Content-type", "text/json")
		data, _ := json.MarshalIndent(CreateResultStatus{Status: t, Message: message, UUID: uuid}, "", "  ")
		w.Write(data)
	}

	vars := mux.Vars(req)
	pName := vars["partition"]

	var data = &fair.ItemData{}

	/*
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&data)
	*/
	bdata, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.log.Errorf("cannot read request body: %v", err)
		sendResult("error", fmt.Sprintf("cannot read request body: %v", err), "")
		return
	}

	if err := json.Unmarshal(bdata, data); err != nil {
		s.log.Errorf("cannot unmarshal request body [%s]: %v", string(bdata), err)
		sendResult("error", fmt.Sprintf("cannot unmarshal request body [%s]: %v", string(bdata), err), "")
		return
	}

	uuidStr, err := s.fair.CreateItem(pName, data)
	if err != nil {
		sendResult("error", fmt.Sprintf("cannot create item: %v", err), "")
		return
	}
	sendResult("ok", "update done", uuidStr)
	return
}

type SearchQueryColumnSearch struct {
	Value string `json:"value"`
	Regex bool   `json:"regex"`
}

type SearchQueryColumn struct {
	Data       string                  `json:"data"`
	Name       string                  `json:"name"`
	Searchable bool                    `json:"searchable"`
	Orderable  bool                    `json:"orderable"`
	Search     SearchQueryColumnSearch `json:"search"`
}

type SearchQueryOrder struct {
	Column int64  `json:"column"`
	Dir    string `json:"dir"`
}

type SearchQuery struct {
	Columns []SearchQueryColumn     `json:"columns"`
	Order   []SearchQueryOrder      `json:"order"`
	Start   int64                   `json:"start"`
	Length  int64                   `json:"length"`
	Search  SearchQueryColumnSearch `json:"search"`
}

type DataTableResult struct {
	Draw            int64               `json:"draw"`
	RecordsTotal    int64               `json:"recordsTotal"`
	RecordsFiltered int64               `json:"recordsFiltered"`
	Data            []map[string]string `json:"data"`
}

var columnsParam = regexp.MustCompile(`columns\[([0-9]+)\]\[([a-z]+)\]`)

func (s *Server) searchDatatableHandler(w http.ResponseWriter, req *http.Request) {
	if !BasicAuth(w, req, s.name, s.password, "FAIR Service") {
		return
	}

	sendResult := func(t string, message string, uuid string) {
		if t == "ok" {
			s.log.Infof(fmt.Sprintf("%s: %s", message, uuid))
		} else {
			s.log.Error(fmt.Sprintf("%s: %s", message, uuid))
		}
		w.Header().Set("Content-type", "text/json")
		data, _ := json.MarshalIndent(CreateResultStatus{Status: t, Message: message, UUID: uuid}, "", "  ")
		w.Write(data)
	}

	vars := mux.Vars(req)
	pName := vars["partition"]

	_, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult("error", fmt.Sprintf("invalid partition %s", pName), "")
		return
	}

	values := req.URL.Query()
	/*
		sq := SearchQuery{
			Columns: []SearchQueryColumn{},
			Order:   []SearchQueryOrder{},
			Search:  SearchQueryColumnSearch{},
		}
	*/
	var start, length, draw int64
	var search string
	for key, vals := range values {
		if len(vals) == 0 {
			continue
		}
		val := vals[0]
		/*
			columns := columnsParam.FindAllString(key, -1)
			if columns != nil {

			} else {

		*/
		switch key {
		case "start":
			if start, err = strconv.ParseInt(val, 10, 64); err != nil {
				sendResult("error", fmt.Sprintf("invalid parameter %s: %s", key, val), "")
				return
			}
		case "length":
			if length, err = strconv.ParseInt(val, 10, 64); err != nil {
				sendResult("error", fmt.Sprintf("invalid parameter %s: %s", key, val), "")
				return
			}
		case "draw":
			if draw, err = strconv.ParseInt(val, 10, 64); err != nil {
				sendResult("error", fmt.Sprintf("invalid parameter %s: %s", key, val), "")
				return
			}
		case "search[value]":
			search = val
		}
	}

	result, num, err := s.fair.Search(pName, search, start, length)
	if err != nil {
		sendResult("error", fmt.Sprintf("cannot search: %v", err), "")
		return
	}

	rData := &DataTableResult{
		Draw:            draw,
		RecordsTotal:    num,
		RecordsFiltered: num - int64(len(result)),
		Data:            result,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(rData); err != nil {
		sendResult("error", fmt.Sprintf("cannot search: %v", err), "")
		return
	}
	return
}
