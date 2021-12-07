package service

import (
	"crypto/subtle"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/FairService/v2/pkg/model/dcmi"
	"github.com/op/go-logging"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type CreateResultStatus struct {
	Status  string           `json:"status"`
	Message string           `json:"message,omitempty"`
	Item    *fair.ItemData   `json:"uuid,omitempty"`
	Items   []*fair.ItemData `json:"uuids,omitempty"`
}

func sendCreateResult(log *logging.Logger, w http.ResponseWriter, t string, message string, item *fair.ItemData) {
	if item != nil {
		if t == "ok" {
			log.Infof(fmt.Sprintf("%s: %s", message, item.UUID))
		} else {
			log.Error(fmt.Sprintf("%s: %s", message, item.UUID))
		}
	} else {
		if t == "ok" {
			log.Infof(fmt.Sprintf("%s", message))
		} else {
			log.Error(fmt.Sprintf("%s", message))
		}
	}
	w.Header().Set("Content-type", "text/json")
	data, _ := json.MarshalIndent(CreateResultStatus{Status: t, Message: message, Item: item}, "", "  ")
	w.Write(data)
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
	data, err := s.fair.GetItem(part, uuidStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error loading item %s/%s: %v", pName, uuidStr, err)))
		return
	}
	if data.Status != fair.DataStatusActive {
		tpl := s.templates["deleted"]
		if err := tpl.Execute(w, struct {
			Part *fair.Partition
			Data *fair.ItemData
		}{Part: part, Data: data}); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-type", "text/plain")
			w.Write([]byte(fmt.Sprintf("error executing template %s in partition %s: %v", "partition", pName, err)))
			return
		}

	} else {
		source, err := s.fair.GetSourceByName(part, data.Source)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Set("Content-type", "text/plain")
			w.Write([]byte(fmt.Sprintf("error loading source %s for item %s/%s: %v", data.Source, pName, uuidStr, err)))
			return
		}

		targetURL := strings.Replace(source.DetailURL, "{signature}", data.Signature, -1)

		http.Redirect(w, req, targetURL, 301)
	}
}

func (s *Server) detailHandler(w http.ResponseWriter, req *http.Request) {
	if !BasicAuth(w, req, s.name, s.password, "FAIR Service") {
		return
	}

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

	doiError := ""
	message := ""
	if _, ok := req.URL.Query()["createdoi"]; ok {
		targetUrl := fmt.Sprintf("%s/redir/%s", part.AddrExt, uuidStr)
		_, err := s.fair.CreateDOI(part, uuidStr, targetUrl)
		if err != nil {
			doiError = err.Error()
		} else {
			message = "Draft DOI successfully created"
		}
	}

	data, err := s.fair.GetItem(part, uuidStr)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error loading item %s/%s: %v", pName, uuidStr, err)))
		return
	}
	tpl := s.templates["detail"]
	if err := tpl.Execute(w, struct {
		Error   string
		Message string
		Part    *fair.Partition
		Data    *fair.ItemData
	}{Error: doiError, Message: message, Part: part, Data: data}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("error executing template %s in partition %s: %v", "partition", pName, err)))
		return
	}
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

	vars := mux.Vars(req)
	pName := vars["partition"]
	uuidStr := vars["uuid"]
	outputType := vars["outputType"]
	if outputType == "" {
		outputType = "json"
	}

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("partition [%s] not found", pName)))
		return
	}

	data, err := s.fair.GetItem(part, uuidStr)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("error loading item %v: %v", uuidStr, err), nil)
		return
	}
	if data == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("item [%s] not found in partition %s", uuidStr, pName)))
		return
	}
	if data.Access != fair.DataAccessOpenAccess && data.Access != fair.DataAccessPublic {
		w.WriteHeader(http.StatusForbidden)
		sendCreateResult(s.log, w, "error", fmt.Sprintf("no public access for %v: %v", uuidStr, data.Access), nil)
		return
	}
	if data.Status != fair.DataStatusActive {
		w.WriteHeader(http.StatusForbidden)
		sendCreateResult(s.log, w, "error", fmt.Sprintf("status of %v not active: %v", uuidStr, data.Status), nil)
		return
	}
	switch outputType {
	case "json":
		w.Header().Set("Content-type", "text/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot marshal data for %v", uuidStr), nil)
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
		if err := enc.Encode(dcmiData); err != nil {
			sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot marshal data for %v", uuidStr), nil)
			return
		}
		return
	case "datacite":
		dataciteData := &dataciteModel.DataCite{}
		dataciteData.InitNamespace()
		dataciteData.FromCore(data.Metadata)
		w.Header().Set("Content-type", "text/xml")
		enc := xml.NewEncoder(w)
		enc.Indent("", "  ")
		if err := enc.Encode(dataciteData); err != nil {
			sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot marshal data for %v", uuidStr), nil)
			return
		}
		return
	default:
		sendCreateResult(s.log, w, "error", fmt.Sprintf("invalid output type %s for %v", outputType, uuidStr), nil)
		return
	}
}

func (s *Server) createDOIHandler(w http.ResponseWriter, req *http.Request) {
	if !BasicAuth(w, req, s.name, s.password, "FAIR Service") {
		return
	}

	vars := mux.Vars(req)
	pName := vars["partition"]
	uuidStr := vars["uuid"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot get partition %s for %v", pName, uuidStr), nil)
		return
	}

	targetUrl := fmt.Sprintf("%s/redir/%s", part.AddrExt, uuidStr)
	doiResult, err := s.fair.CreateDOI(part, uuidStr, targetUrl)
	if err != nil {
		sendCreateResult(s.log, w, "error", err.Error(), nil)
		return
	}

	w.Header().Set("Content-type", "text/json")
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(doiResult); err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot marshal data for %v", uuidStr), nil)
		return
	}
	return
}

func (s *Server) startUpdateHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pName := vars["partition"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot get partition %s", pName), nil)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var data fair.SourceData
	if err := decoder.Decode(&data); err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot parse request body: %v", err), nil)
		return
	}

	if err := s.fair.StartUpdate(part, data.Source); err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot start update for %s on %s: %v", data.Source, pName, err), nil)
		return
	}
	sendCreateResult(s.log, w, "ok", fmt.Sprintf("starting update for %s on %s", data.Source, pName), nil)
}

func (s *Server) pingHandler(w http.ResponseWriter, req *http.Request) {
	sendCreateResult(s.log, w, "ok", "pong", nil)
}

func (s *Server) endUpdateHandler(w http.ResponseWriter, req *http.Request) {

	vars := mux.Vars(req)
	pName := vars["partition"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot get partition %s", pName), nil)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var data fair.SourceData
	if err := decoder.Decode(&data); err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot parse request body: %v", err), nil)
		return
	}

	if err := s.fair.EndUpdate(part, data.Source); err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot end update for %s on %s: %v", data.Source, pName, err), nil)
		return
	}
	sendCreateResult(s.log, w, "ok", fmt.Sprintf("end update for %s on %s", data.Source, pName), nil)
}

func (s *Server) abortUpdateHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot get partition %s", pName), nil)
		return
	}

	decoder := json.NewDecoder(req.Body)
	var data fair.SourceData
	if err := decoder.Decode(&data); err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot parse request body: %v", err), nil)
		return
	}

	if err := s.fair.AbortUpdate(part, data.Source); err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot end update for %s on %s: %v", data.Source, pName, err), nil)
		return
	}
	sendCreateResult(s.log, w, "ok", fmt.Sprintf("abort update for %s on %s", data.Source, pName), nil)
}

func (s *Server) originalDataReadHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]
	uuidStr := vars["uuid"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot get partition %s", pName), nil)
		return
	}

	data, t, err := s.fair.GetOriginalData(part, uuidStr)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot read original data: %v", err), nil)
		return
	}
	switch t {
	case "XML":
		w.Header().Set("Content-Type", "text/xml")
	case "JSON":
		w.Header().Set("Content-Type", "application/json")
	default:
		w.Header().Set("Content-Type", "text/plain")
	}
	w.Write(data)
}

func (s *Server) originalDataWriteHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]
	uuidStr := vars["uuid"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot get partition %s", pName), nil)
		return
	}

	bdata, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.log.Errorf("cannot read request body: %v", err)
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot read request body: %v", err), nil)
		return
	}

	var t string
	var data interface{}
	if err := json.Unmarshal(bdata, data); err == nil {
		t = "JSON"
	} else {
		if err := xml.Unmarshal(bdata, data); err == nil {
			t = "XML"
		} else {
			t = "Other"
		}
	}

	if err := s.fair.SetOriginalData(part, uuidStr, t, bdata); err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot set original data for %s: %v", uuidStr, err), nil)
		return
	}
	sendCreateResult(s.log, w, "ok", fmt.Sprintf("original metadata for %s stored", uuidStr), nil)
}

func (s *Server) createHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("partition [%s] not found", pName)))
		return
	}

	var data = &fair.ItemData{}

	/*
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&data)
	*/
	bdata, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.log.Errorf("cannot read request body: %v", err)
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot read request body: %v", err), nil)
		return
	}

	if err := json.Unmarshal(bdata, data); err != nil {
		s.log.Errorf("cannot unmarshal request body [%s]: %v", string(bdata), err)
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot unmarshal request body [%s]: %v", string(bdata), err), nil)
		return
	}

	item, err := s.fair.CreateItem(part, data)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot create item: %v", err), nil)
		return
	}
	sendCreateResult(s.log, w, "ok", "update done", item)
	return
}

func (s *Server) setSourceHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]

	var data = &fair.Source{}
	bdata, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.log.Errorf("cannot read request body: %v", err)
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot read request body: %v", err), nil)
		return
	}

	if err := json.Unmarshal(bdata, data); err != nil {
		s.log.Errorf("cannot unmarshal request body [%s]: %v", string(bdata), err)
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot unmarshal request body [%s]: %v", string(bdata), err), nil)
		return
	}
	if data.Partition != pName {
		s.log.Errorf("source and partition do not match %s != %s", data.Partition, pName)
		sendCreateResult(s.log, w, "error", fmt.Sprintf("source and partition do not match %s != %s", data.Partition, pName), nil)
		return
	}

	if err := s.fair.SetSource(data); err != nil {
		s.log.Errorf("cannot create source %v: %v", data, err)
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot create source %v: %v", data, err), nil)
		return
	}
	sendCreateResult(s.log, w, "ok", fmt.Sprintf("cannot create source %v: %v", data, err), nil)

}

type DataTableResult struct {
	Draw            int64               `json:"draw"`
	RecordsTotal    int64               `json:"recordsTotal"`
	RecordsFiltered int64               `json:"recordsFiltered"`
	Data            []map[string]string `json:"data"`
	Sql             string              `json:"sql"`
}

var columnsParam = regexp.MustCompile(`columns\[([0-9]+)\]\[([a-z]+)\]`)
