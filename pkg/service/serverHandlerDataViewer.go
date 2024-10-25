package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/je4/utils/v2/pkg/datatable"
	"net/http"
)

func (s *Server) dataViewerHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("partition [%s] not found", pName), nil)
		return
	}

	tpl := s.templates["dataviewer"]
	if err := tpl.Execute(ctx.Writer, part); err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("error executing template %s in partition %s: %v", "partition", pName, err), nil)
		return
	}
}

func (s *Server) searchViewerHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("partition [%s] not found", pName), nil)
		return
	}

	dReq := &datatable.Request{}
	if err := dReq.FromRequest(ctx.Request); err != nil {
		sendResult(s.log, ctx, http.StatusBadRequest, fmt.Sprintf("cannot eval request parameter: %v", err), nil)
		return
	}

	result, num, total, err := s.fair.Search(part, dReq)
	if err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot search: %v", err), nil)
		return
	}

	rData := &DataTableResult{
		Draw:            dReq.Draw,
		RecordsTotal:    total,
		RecordsFiltered: num,
		Data:            result,
	}
	ctx.JSON(http.StatusOK, rData)

	return
}
