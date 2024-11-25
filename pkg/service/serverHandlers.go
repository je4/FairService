package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/FairService/v2/pkg/model/dcmi"
	"github.com/je4/utils/v2/pkg/zLogger"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type FairResultMessage struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

type FairResultStatus struct {
	Status       string              `json:"status"`
	Message      string              `json:"message,omitempty"`
	Item         *fair.ItemData      `json:"uuid,omitempty"`
	ArchiveItems []*fair.ArchiveItem `json:"archiveitems,omitempty"`
}

func sendResult(log zLogger.ZLogger, ctx *gin.Context, status int, message string, item *fair.ItemData) {
	if item != nil {
		if status == http.StatusOK {
			log.Info().Msgf("%s: %s", message, item.UUID)
		} else {
			log.Error().Msgf("%s: %s", message, item.UUID)
		}
	} else {
		if status == http.StatusOK {
			log.Info().Msgf("%s", message)
		} else {
			log.Error().Msgf("%s", message)
		}
	}
	message = fmt.Sprintf("%s: %s", ctx.HandlerName(), message)
	if item == nil {
		ctx.JSON(status, FairResultMessage{Status: http.StatusText(status), Message: message})
	} else {
		ctx.JSON(status, FairResultStatus{Status: http.StatusText(status), Message: message, Item: item})
	}
}

/*
func BasicAuth(ctx *gin.Context, username, password, realm string) bool {

	user, pass, ok := ctx.Request.BasicAuth()

	if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
		w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
		w.WriteHeader(401)
		w.Write([]byte("Unauthorised.\n"))
		return false
	}

	return true
}

*/

func (s *Server) resolverHandler(ctx *gin.Context) {
	partition := ctx.Param("partition")
	pid := strings.Trim(ctx.Param("pid"), "/")
	if ctx.Request.URL.RawQuery != "" {
		pid += "?" + ctx.Request.URL.RawQuery
	}
	var data string
	var _type fair.ResolveResultType
	var err error
	if partition == "" {
		data, _type, err = s.fair.Resolve(pid)
		if err != nil {
			sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("cannot resolve %s: %v", pid, err), nil)
			return
		}
	} else {
		part, err := s.fair.GetPartition(partition)
		if err != nil {
			sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("cannot get partition %s: %v", partition, err), nil)
			return
		}
		data, _type, err = part.Resolve(pid)
		if err != nil {
			sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("cannot resolve %s: %v", pid, err), nil)
			return
		}
	}
	switch _type {
	case fair.ResolveResultTypeRedirect:
		ctx.Redirect(http.StatusMovedPermanently, data)
	case fair.ResolveResultTypeTextPlain:
		ctx.String(http.StatusOK, data)
	case fair.ResolveResultTypeApplicationXML:
		ctx.Header("Content-Type", "application/xml")
		ctx.String(http.StatusOK, data)
	case fair.ResolveResultTypeApplicationJSON:
		ctx.Header("Content-Type", "application/json")
		ctx.String(http.StatusOK, data)
	case fair.ResolveResultTypeApplicationYAML:
		ctx.Header("Content-Type", "text/yaml")
		ctx.String(http.StatusOK, data)
	case fair.ResolveResultTypeUnknown:
		ctx.String(http.StatusOK, data)
	default:
		ctx.String(http.StatusOK, data)
	}
}

func (s *Server) redirectHandler(ctx *gin.Context) {
	pName := ctx.Param("partition") + ctx.GetString("partition")
	uuidStr := ctx.Param("uuid") + ctx.GetString("uuid")
	suffix := ctx.Param("suffix") + ctx.GetString("suffix")

	var part *fair.Partition
	var err error
	if pName != "" {
		part, err = s.fair.GetPartition(pName)
		if err != nil {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("%s: partition [%s] not found", ctx.HandlerName(), pName))
			return
		}
	}
	data, err := s.fair.GetItem(part, uuidStr)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("error loading item %s/%s: %v", pName, uuidStr, err))
		return
	}
	if part == nil {
		part, err = s.fair.GetPartition(data.Partition)
		if err != nil {
			ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("%s: partition [%s] not found", ctx.HandlerName(), data.Partition))
			return
		}
	}
	if data.Status == fair.DataStatusDeletedMeta {
		tpl := s.templates["detail"]
		if err := tpl.Execute(ctx.Writer, struct {
			BaseURL string
			Part    *fair.Partition
			Data    *fair.ItemData
		}{BaseURL: part.AddrExt + "/", Part: part, Data: data}); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error executing template %s in partition %s: %v", "partition", pName, err))
			return
		}
		return
	}
	if data.Status != fair.DataStatusActive {
		tpl := s.templates["deleted"]
		if err := tpl.Execute(ctx.Writer, struct {
			BaseURL string
			Part    *fair.Partition
			Data    *fair.ItemData
		}{BaseURL: part.AddrExt + "/", Part: part, Data: data}); err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error executing template %s in partition %s: %v", "partition", pName, err))
			return
		}
		return
	}
	source, err := s.fair.GetSourceByName(data.Partition, data.Source)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error loading source %s for item %s/%s: %v", data.Source, pName, uuidStr, err))
		return
	}

	var targetURL string
	if data.URL != "" {
		targetURL = data.URL
	} else {
		targetURL = strings.Replace(source.DetailURL, "{signature}", data.Signature, -1)
	}
	targetURL += suffix
	ctx.Redirect(http.StatusMovedPermanently, targetURL)
}

func (s *Server) detailHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")
	uuidStr := ctx.Param("uuid")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("%s: partition [%s] not found", ctx.HandlerName(), pName))
		return
	}

	doiError := ""
	message := ""

	if _, ok := ctx.GetQuery("createdoi"); ok {
		//		targetUrl := fmt.Sprintf("%s/redir/%s", part.AddrExt, uuidStr)
		pid, err := part.CreatePID(uuidStr, dataciteModel.RelatedIdentifierTypeDOI)
		if err != nil {
			doiError = err.Error()
		} else {
			message = fmt.Sprintf("DOI %s successfully created", pid)
		}
	}

	data, err := s.fair.GetItem(part, uuidStr)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("error loading item %s/%s: %v", pName, uuidStr, err))
		return
	}
	tpl := s.templates["detail"]
	if err := tpl.Execute(ctx.Writer, struct {
		BaseURL string
		Error   string
		Message string
		Part    *fair.Partition
		Data    *fair.ItemData
	}{BaseURL: part.AddrExt + "/", Error: doiError, Message: message, Part: part, Data: data}); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error executing template %s in partition %s: %v", "partition", pName, err))
		return
	}
}

func (s *Server) partitionHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("%s: partition [%s] not found", ctx.HandlerName(), pName))
		return
	}

	tpl := s.templates["partition"]
	if err := tpl.Execute(ctx.Writer, struct {
		BaseURL string
		Part    *fair.Partition
	}{
		BaseURL: part.AddrExt + "/",
		Part:    part,
	}); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error executing template %s in partition %s: %v", "partition", pName, err))
		return
	}
}

func (s *Server) partitionOAIHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("%s: partition [%s] not found", ctx.HandlerName(), pName))
		return
	}

	tpl := s.templates["oai"]
	if err := tpl.Execute(ctx.Writer, struct {
		BaseURL string
		Part    *fair.Partition
	}{BaseURL: part.AddrExt + "/", Part: part}); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("error executing template %s in partition %s: %v", "partition", pName, err))
		return
	}
}

func (s *Server) itemHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")
	uuidStr := ctx.Param("uuid")
	outputType := ctx.Param("outputType")
	if outputType == "" {
		outputType = "json"
	}

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("%s: partition [%s] not found", ctx.HandlerName(), pName))
		return
	}

	data, err := s.fair.GetItem(part, uuidStr)
	if err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("error loading item %v: %v", uuidStr, err), nil)
		return
	}
	if data == nil {
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("item [%s] not found in partition %s", uuidStr, pName))
		return
	}
	if data.Access != fair.DataAccessOpenAccess && data.Access != fair.DataAccessPublic {
		sendResult(s.log, ctx, http.StatusForbidden, fmt.Sprintf("no public access for %v: %v", uuidStr, data.Access), nil)
		return
	}
	if data.Status != fair.DataStatusActive {
		ctx.AbortWithError(http.StatusForbidden, fmt.Errorf("status of %v not active: %v", uuidStr, data.Status))
		return
	}
	switch outputType {
	case "json":
		ctx.JSON(http.StatusOK, data)
		return
	case "dcmi":
		dcmiData := &dcmi.DCMI{}
		dcmiData.InitNamespace()
		dcmiData.FromCore(data.Metadata)
		ctx.XML(http.StatusOK, dcmiData)
		return
	case "datacite":
		dataciteData := &dataciteModel.DataCite{}
		dataciteData.InitNamespace()
		dataciteData.FromCore(data.Metadata)
		ctx.XML(http.StatusOK, dataciteData)
		return
	default:
		sendResult(s.log, ctx, http.StatusBadRequest, fmt.Sprintf("invalid output type %s for %v", outputType, uuidStr), nil)
		return
	}
}

func (s *Server) createDOIHandler(ctx *gin.Context) {

	pName := ctx.Param("partition")
	uuidStr := ctx.Param("uuid")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("cannot get partition %s for %v", pName, uuidStr), nil)
		return
	}

	doiResult, err := part.CreatePID(uuidStr, dataciteModel.RelatedIdentifierTypeDOI)
	if err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	ctx.JSON(http.StatusOK, doiResult)
	return
}

// startUpdate godoc
// @Summary      starts update transaction
// @ID			 post-start-update
// @Description  starts update transaction for a source
// @Tags         fairservice
// @Security 	 BearerAuth
// @Produce      json
// @Param 		 partition path string true "Partition"
// @Param 		 source       body fair.SourceData true "source to start update"
// @Success      200  {object}  FairResultMessage
// @Failure      400  {object}  FairResultMessage
// @Failure      401  {object}  FairResultMessage
// @Failure      404  {object}  FairResultMessage
// @Failure      500  {object}  FairResultMessage
// @Router       /{partition}/startupdate [post]
func (s *Server) startUpdateHandler(ctx *gin.Context) {

	pName := ctx.Param("partition")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("cannot get partition %s", pName), nil)
		return
	}

	var data fair.SourceData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		sendResult(s.log, ctx, http.StatusBadRequest, fmt.Sprintf("cannot parse request body: %v", err), nil)
		return
	}

	if err := s.fair.StartUpdate(part, data.Source); err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot start update for %s on %s: %v", data.Source, pName, err), nil)
		return
	}
	sendResult(s.log, ctx, http.StatusOK, fmt.Sprintf("starting update for %s on %s", data.Source, pName), nil)
}

// ping godoc
// @Summary      does pong
// @ID			 get-ping
// @Description  for testing if server is running
// @Tags         mediaserver
// @Param 		 domain path string true "Domain"
// @Produce      plain
// @Success      200  {string}  string
// @Router       /{domain}/ping [get]
func (s *Server) pingHandler(ctx *gin.Context) {
	sendResult(s.log, ctx, http.StatusOK, "pong", nil)
}

// endUpdate godoc
// @Summary      ends update transaction
// @ID			 post-end-update
// @Description  ends update transaction for a source with commit
// @Tags         fairservice
// @Security 	 BearerAuth
// @Produce      json
// @Param 		 partition path string true "Partition"
// @Param 		 source       body fair.SourceData true "source to end update"
// @Success      200  {object}  FairResultMessage
// @Failure      400  {object}  FairResultMessage
// @Failure      401  {object}  FairResultMessage
// @Failure      404  {object}  FairResultMessage
// @Failure      500  {object}  FairResultMessage
// @Router       /{partition}/endupdate [post]
func (s *Server) endUpdateHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("cannot get partition %s", pName), nil)
		return
	}

	var data fair.SourceData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		sendResult(s.log, ctx, http.StatusBadRequest, fmt.Sprintf("cannot parse request body: %v", err), nil)
		return
	}

	if err := s.fair.EndUpdate(part, data.Source); err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot end update for %s on %s: %v", data.Source, pName, err), nil)
		return
	}
	sendResult(s.log, ctx, http.StatusOK, fmt.Sprintf("end update for %s on %s", data.Source, pName), nil)
}

// abortUpdate godoc
// @Summary      aborts update transaction
// @ID			 post-abort-update
// @Description  ends aborts transaction for a source without removal of missing items
// @Tags         fairservice
// @Security 	 BearerAuth
// @Produce      json
// @Param 		 partition path string true "Partition"
// @Param 		 source       body fair.SourceData true "source to abort update"
// @Success      200  {object}  FairResultMessage
// @Failure      400  {object}  FairResultMessage
// @Failure      401  {object}  FairResultMessage
// @Failure      404  {object}  FairResultMessage
// @Failure      500  {object}  FairResultMessage
// @Router       /{partition}/abortupdate [post]
func (s *Server) abortUpdateHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("cannot get partition %s", pName), nil)
		return
	}

	var data fair.SourceData
	if err := ctx.ShouldBindJSON(&data); err != nil {
		sendResult(s.log, ctx, http.StatusBadRequest, fmt.Sprintf("cannot parse request body: %v", err), nil)
		return
	}

	if err := s.fair.AbortUpdate(part, data.Source); err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot end update for %s on %s: %v", data.Source, pName, err), nil)
		return
	}
	sendResult(s.log, ctx, http.StatusOK, fmt.Sprintf("abort update for %s on %s", data.Source, pName), nil)
}

func (s *Server) originalDataReadHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")
	uuidStr := ctx.Param("uuid")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("cannot get partition %s", pName), nil)
		return
	}

	data, t, err := s.fair.GetOriginalData(part, uuidStr)
	if err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot read original data: %v", err), nil)
		return
	}
	switch t {
	case "XML":
		ctx.Header("Content-Type", "text/xml")
	case "JSON":
		ctx.Header("Content-Type", "application/json")
	default:
		ctx.Header("Content-Type", "text/plain")
	}
	ctx.Writer.Write(data)
}

func (s *Server) originalDataWriteHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")
	uuidStr := ctx.Param("uuid")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("cannot get partition %s", pName), nil)
		return
	}

	bdata, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot read request body: %v", err), nil)
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
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot set original data for %s: %v", uuidStr, err), nil)
		return
	}
	sendResult(s.log, ctx, http.StatusOK, fmt.Sprintf("original data for %s stored", uuidStr), nil)
}

// createItem godoc
// @Summary      creates a new item
// @ID			 post-create-item
// @Description  creates a new item within a transaction
// @Tags         fairservice
// @Security 	 BearerAuth
// @Produce      json
// @Param 		 partition path string true "Partition"
// @Param 		 source       body fair.ItemData true "source to abort update"
// @Success      200  {object}  FairResultMessage
// @Failure      400  {object}  FairResultMessage
// @Failure      401  {object}  FairResultMessage
// @Failure      404  {object}  FairResultMessage
// @Failure      500  {object}  FairResultMessage
// @Router       /{partition}/item [post]
func (s *Server) createHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")

	part, err := s.fair.GetPartition(pName)
	if err != nil {
		sendResult(s.log, ctx, http.StatusNotFound, fmt.Sprintf("partition [%s] not found", pName), nil)
		return
	}

	var data = &fair.ItemData{}

	if err := ctx.ShouldBindJSON(data); err != nil {
		sendResult(s.log, ctx, http.StatusBadRequest, fmt.Sprintf("cannot parse request body: %v", err), nil)
		return
	}

	item, err := s.fair.CreateItem(part, data)
	if err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot create item: %v", err), nil)
		return
	}
	sendResult(s.log, ctx, http.StatusOK, "update done", item)
	return
}

// setSource godoc
// @Summary      update or create source
// @ID			 post-set-source
// @Description  updates or creates source for a partition
// @Tags         fairservice
// @Security 	 BearerAuth
// @Produce      json
// @Param 		 partition path string true "Partition"
// @Param 		 source       body fair.Source true "source to set"
// @Success      200  {object}  FairResultMessage
// @Failure      400  {object}  FairResultMessage
// @Failure      401  {object}  FairResultMessage
// @Failure      404  {object}  FairResultMessage
// @Failure      500  {object}  FairResultMessage
// @Router       /{partition}/source [post]
func (s *Server) setSourceHandler(ctx *gin.Context) {
	pName := ctx.Param("partition")

	var data = &fair.Source{}
	if err := ctx.ShouldBindJSON(data); err != nil {
		sendResult(s.log, ctx, http.StatusBadRequest, fmt.Sprintf("cannot parse request body: %v", err), nil)
		return
	}
	if data.Partition != pName {
		sendResult(s.log, ctx, http.StatusBadRequest, fmt.Sprintf("source and partition do not match %s != %s", data.Partition, pName), nil)
		return
	}

	if err := s.fair.SetSource(data); err != nil {
		sendResult(s.log, ctx, http.StatusInternalServerError, fmt.Sprintf("cannot set source %v: %v", data, err), nil)
		return
	}
	sendResult(s.log, ctx, http.StatusOK, fmt.Sprintf("source %v set", data), nil)
}

type DataTableResult struct {
	Draw            int64               `json:"draw"`
	RecordsTotal    int64               `json:"recordsTotal"`
	RecordsFiltered int64               `json:"recordsFiltered"`
	Data            []map[string]string `json:"data"`
	Sql             string              `json:"sql"`
}

var columnsParam = regexp.MustCompile(`columns\[([0-9]+)\]\[([a-z]+)\]`)
