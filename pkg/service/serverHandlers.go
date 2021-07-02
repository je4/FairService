package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/je4/FairService/v2/pkg/model/myfair"
	"net/http"
	"reflect"
)

type CreateResultStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	UUID    string `json:"uuid,omitempty"`
}

type CreateData struct {
	Source           string      `json:"source"`
	SourceIdentifier string      `json:"sourceID"`
	Core             myfair.Core `json:"core"`
	DOI              string      `json:"doi,omitempty"`
}

func (s *Server) createHandler(w http.ResponseWriter, req *http.Request) {
	sendResult := func(t string, message string, uuid string) {
		s.log.Error(message)
		w.Header().Set("Content-type", "text/json")
		data, _ := json.MarshalIndent(CreateResultStatus{Status: t, Message: message}, "", "  ")
		w.Write(data)

	}

	vars := mux.Vars(req)
	s.log.Infof("vars: %v", vars)

	decoder := json.NewDecoder(req.Body)
	var data CreateData
	err := decoder.Decode(&data)
	if err != nil {
		sendResult("error", fmt.Sprintf("cannot parse request body: %v", err), "")
		return
	}

	sqlstr := fmt.Sprintf("SELECT uuidStr, doi, core, seq FROM %s.data WHERE source=$1 AND sourceid=$2", s.dbschema)
	params := []string{data.Source, data.SourceIdentifier}
	row := s.db.QueryRow(sqlstr, params)

	var coreStr string
	var uuidStr string
	var doi sql.NullString
	var seq int64
	if err := row.Scan(&uuidStr, &doi, &coreStr, &seq); err != nil {
		if err != sql.ErrNoRows {
			sendResult("error", fmt.Sprintf("cannot execute query [%s] - [%v]: %v", sqlstr, params, err), "")
			return
		}
		// do insert here

		uuidVal, err := uuid.NewUUID()
		if err != nil {
			sendResult("error", fmt.Sprintf("cannot generate uuid: %v", err), "")
			return
		}
		uuidStr := uuidVal.String()
		coreBytes, err := json.Marshal(data.Core)
		if err != nil {
			sendResult("error", fmt.Sprintf("cannot marshal core data [%v]: %v", data.Core, err), "")
			return
		}
		var sqlDOI sql.NullString
		if data.DOI == "" {
			sqlDOI.Valid = false
		} else {
			sqlDOI.String = data.DOI
			sqlDOI.Valid = true
		}
		sqlstr := fmt.Sprintf("INSERT INTO %s.data (uuidStr, source, sourceid, doi, core, seq) VALUES($1, $2, $3, $4, $5, NEXTVAL('lastchange'))", s.dbschema)
		params := []interface{}{
			uuidStr, data.Source, data.SourceIdentifier, sqlDOI, string(coreBytes),
		}
		ret, err := s.db.Exec(sqlstr, params)
		if err != nil {
			sendResult("error", fmt.Sprintf("cannot execute query [%s] - [%v]: %v", sqlstr, params, err), "")
			return
		}
		num, err := ret.RowsAffected()
		if err != nil {
			sendResult("error", fmt.Sprintf("could not get affected rows: %v", err), uuidStr)
			return
		}
		if num == 0 {
			sendResult("error", fmt.Sprintf("no affected rows"), uuidStr)
			return
		}
	} else {
		// do update here
		var oldCore myfair.Core
		if err := json.Unmarshal([]byte(coreStr), &oldCore); err != nil {
			sendResult("error", fmt.Sprintf("cannot unmarshal core [%s]: %v", coreStr, err), uuidStr)
			return
		}
		if reflect.DeepEqual(oldCore, data.Core) && doi.String == data.DOI {
			sendResult("ok", "no update", uuidStr)
			return
		}

		dataCoreByte, err := json.Marshal(data.Core)
		if err != nil {
			sendResult("error", fmt.Sprintf("cannot unmarshal data core: %v", err), uuidStr)
			return
		}

		sqlstr = fmt.Sprintf("UPDATE %s.data SET doi=$1, core=$2, seq=NEXTVAL('lastchange') WHERE uuidStr=$3", s.dbschema)
		params := []string{data.DOI, string(dataCoreByte), uuidStr}
		if _, err := s.db.Exec(sqlstr, params); err != nil {
			sendResult("error", fmt.Sprintf("cannot update [%s] - [%v]: %v", sqlstr, params, err), uuidStr)
			return
		}
		sendResult("ok", "update done", uuidStr)
		return
	}
}
