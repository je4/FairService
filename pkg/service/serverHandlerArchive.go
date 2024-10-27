package service

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/je4/FairService/v2/pkg/fair"
	"net/http"
)

type Archive struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *Server) createArchiveHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("partition [%s] not found", pName), nil)
		return
	}

	var data = &Archive{}

	/*
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&data)
	*/
	bdata, err := ctx.GetRawData()
	if err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot read request body: %v", err), nil)
		return
	}

	if err := json.Unmarshal(bdata, data); err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot unmarshal request body [%s]: %v", string(bdata), err), nil)
		return
	}

	name := fmt.Sprintf("%s.%s", part.Name, data.Name)
	if err := s.fair.AddArchive(part, fmt.Sprintf("%s.%s", part.Name, name), data.Description); err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot create archive item: %v", err), nil)
		return
	}
	sendResult(s.log, ctx, http.StatusOK, fmt.Sprintf("archive %s created", name), nil)
	return
}

func (s *Server) getArchiveItemHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")
	archive := ctx.Param("archive")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("partition [%s] not found", pName), nil)
		return
	}
	var items = []*fair.ArchiveItem{}
	if err := s.fair.GetArchiveItems(part, fmt.Sprintf("%s.%s", part.Name, archive), false, func(item *fair.ArchiveItem) error {
		items = append(items, item)
		return nil
	}); err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot get archive items: %v", err), nil)
		return
	}

	ctx.JSON(http.StatusOK, FairResultStatus{Status: "ok", Message: fmt.Sprintf("%v items found", len(items)), ArchiveItems: items})
}

func (s *Server) addArchiveItemHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")
	archive := ctx.Param("archive")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("partition [%s] not found", pName), nil)
		return
	}

	var uuid string

	if err := ctx.ShouldBindJSON(&uuid); err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot unmarshal request body: %v", err), nil)
		return
	}

	item, err := s.fair.GetItem(part, uuid)
	if err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot get item %s/%s: %v", part.Name, uuid, err), nil)
		return
	}

	if err := s.fair.AddArchiveItem(part, fmt.Sprintf("%s.%s", part.Name, archive), item); err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot add item %s to archive %s: %v", item.UUID, archive, err), nil)
		return
	}
	sendResult(s.log, ctx, http.StatusOK, fmt.Sprintf("item %s added to archive %s", item.UUID, archive), nil)
	return
}
