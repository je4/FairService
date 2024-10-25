package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/je4/FairService/v2/pkg/fair"
	"io/ioutil"
	"net/http"
)

type Archive struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) createArchiveHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("partition [%s] not found", pName)))
		return
	}

	var data = &Archive{}

	/*
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&data)
	*/
	bdata, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.log.Error().Msgf("cannot read request body: %v", err)
		sendResult(s.log, w, "error", fmt.Sprintf("cannot read request body: %v", err), nil)
		return
	}

	if err := json.Unmarshal(bdata, data); err != nil {
		s.log.Error().Msgf("cannot unmarshal request body [%s]: %v", string(bdata), err)
		sendResult(s.log, w, "error", fmt.Sprintf("cannot unmarshal request body [%s]: %v", string(bdata), err), nil)
		return
	}

	name := fmt.Sprintf("%s.%s", part.Name, data.Name)
	if err := s.fair.AddArchive(part, fmt.Sprintf("%s.%s", part.Name, name), data.Description); err != nil {
		sendResult(s.log, w, "error", fmt.Sprintf("cannot create item: %v", err), nil)
		return
	}
	sendResult(s.log, w, "ok", fmt.Sprintf("archive %s created", name), nil)
	return
}

func (s *Server) getArchiveItemHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]
	archive := vars["archive"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		s.log.Error().Msgf("partition [%s] not found", pName)
		sendResult(s.log, w, "error", fmt.Sprintf("partition [%s] not found", pName), nil)
		return
	}
	var items = []*fair.ArchiveItem{}
	if err := s.fair.GetArchiveItems(part, fmt.Sprintf("%s.%s", part.Name, archive), false, func(item *fair.ArchiveItem) error {
		items = append(items, item)
		return nil
	}); err != nil {
		s.log.Error().Msgf("cannot get archive items: %v", err)
		sendResult(s.log, w, "error", fmt.Sprintf("cannot get archive items: %v", err), nil)
		return
	}

	w.Header().Set("Content-type", "text/json")
	data, _ := json.MarshalIndent(FairResultStatus{Status: "ok", Message: fmt.Sprintf("%v items found", len(items)), ArchiveItems: items}, "", "  ")
	w.Write(data)
}

func (s *Server) addArchiveItemHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	pName := vars["partition"]
	archive := vars["archive"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-type", "text/plain")
		w.Write([]byte(fmt.Sprintf("partition [%s] not found", pName)))
		return
	}

	var uuid string

	bdata, err := ioutil.ReadAll(req.Body)
	if err != nil {
		s.log.Error().Msgf("cannot read request body: %v", err)
		sendResult(s.log, w, "error", fmt.Sprintf("cannot read request body: %v", err), nil)
		return
	}

	if err := json.Unmarshal(bdata, &uuid); err != nil {
		s.log.Error().Msgf("cannot unmarshal request body [%s]: %v", string(bdata), err)
		sendResult(s.log, w, "error", fmt.Sprintf("cannot unmarshal request body [%s]: %v", string(bdata), err), nil)
		return
	}

	item, err := s.fair.GetItem(part, uuid)
	if err != nil {
		s.log.Error().Msgf("cannot get item %s/%s: %v", part.Name, uuid, err)
		sendResult(s.log, w, "error", fmt.Sprintf("cannot get item %s/%s: %v", part.Name, uuid, err), nil)
		return
	}

	/*
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&data)
	*/
	if err := s.fair.AddArchiveItem(part, fmt.Sprintf("%s.%s", part.Name, archive), item); err != nil {
		sendResult(s.log, w, "error", fmt.Sprintf("cannot add item %s to %s: %v", item.UUID, archive, err), nil)
		return
	}
	sendResult(s.log, w, "ok", "archive %s created", nil)
	return
}
