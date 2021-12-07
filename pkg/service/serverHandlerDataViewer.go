package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/je4/FairService/v2/pkg/datatable"
	"net/http"
)

func (s *Server) dataViewerHandler(w http.ResponseWriter, req *http.Request) {
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

func (s *Server) searchViewerHandler(w http.ResponseWriter, req *http.Request) {
	if !BasicAuth(w, req, s.name, s.password, "FAIR Service") {
		return
	}

	vars := mux.Vars(req)
	pName := vars["partition"]

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("invalid partition %s", pName), nil)
		return
	}

	dReq := &datatable.Request{}
	if err := dReq.FromRequest(req); err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot eval request parameter: %v", err), nil)
		return
	}

	result, num, total, err := s.fair.Search(part, dReq)
	if err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot search: %v", err), nil)
		return
	}

	rData := &DataTableResult{
		Draw:            dReq.Draw,
		RecordsTotal:    total,
		RecordsFiltered: num,
		Data:            result,
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(rData); err != nil {
		sendCreateResult(s.log, w, "error", fmt.Sprintf("cannot search: %v", err), nil)
		return
	}
	return
}
