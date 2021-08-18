package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/dcmi"
	"io/ioutil"
	"net/http"
	"strings"
)

type CreateResultStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	UUID    string `json:"uuid,omitempty"`
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

	var data fair.ItemData

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

	if err := json.Unmarshal(bdata, &data); err != nil {
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
