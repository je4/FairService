package service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/je4/FairService/v2/pkg/model/myfair"
	"github.com/lib/pq"
	"net/http"
	"reflect"
	"sort"
)

type CreateResultStatus struct {
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
	UUID    string `json:"uuid,omitempty"`
}

type CreateData struct {
	Source    string      `json:"source"`
	Signature string      `json:"signature"`
	Metadata  myfair.Core `json:"metadata"`
	Set       []string    `json:"set"`
	Catalog   []string    `json:"catalog"`
	Public    bool        `json:"public"`
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
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
	sort.Strings(data.Catalog)
	sort.Strings(data.Set)

	src, err := s.GetSourceByName(data.Source)
	if err != nil {
		sendResult("error", fmt.Sprintf("cannot get source: %v", err), "")
		return
	}

	sqlstr := fmt.Sprintf("SELECT uuidStr, metadata, setspec, catalog, public FROM %s.oai WHERE source=$1 AND signature=$2", s.dbschema)
	params := []interface{}{src.ID, data.Signature}
	row := s.db.QueryRow(sqlstr, params)

	var metaStr string
	var uuidStr string
	var set, catalog []string
	var public bool
	if err := row.Scan(&uuidStr, &metaStr, pq.Array(&set), pq.Array(&catalog), &public); err != nil {
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
		coreBytes, err := json.Marshal(data.Metadata)
		if err != nil {
			sendResult("error", fmt.Sprintf("cannot marshal core data [%v]: %v", data.Metadata, err), "")
			return
		}
		sqlstr := fmt.Sprintf("INSERT INTO %s.oai "+
			"(uuid, datestamp, setspec, metadata, dirty, signature, source, public, catalog, seq) "+
			"VALUES($1, NOW(), $2, $3, false, $4, $5, $6, $7, NEXTVAL('lastchange'))", s.dbschema)
		params := []interface{}{
			uuidStr, // uuid
			// datestamp
			pq.Array(data.Set), // setspec
			string(coreBytes),  // metadata
			// dirty
			data.Signature,         // signature
			src.ID,                 // source
			data.Public,            // public
			pq.Array(data.Catalog), // catalog
			// seq
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
		sort.Strings(catalog)
		sort.Strings(set)
		// do update here
		var oldMeta myfair.Core
		if err := json.Unmarshal([]byte(metaStr), &oldMeta); err != nil {
			sendResult("error", fmt.Sprintf("cannot unmarshal core [%s]: %v", metaStr, err), uuidStr)
			return
		}
		if reflect.DeepEqual(oldMeta, data.Metadata) &&
			equalStrings(set, data.Set) &&
			equalStrings(catalog, data.Catalog) &&
			public == data.Public {
			sendResult("ok", "no update", uuidStr)
			return
		}

		dataMetaBytes, err := json.Marshal(data.Metadata)
		if err != nil {
			sendResult("error", fmt.Sprintf("cannot unmarshal data core: %v", err), uuidStr)
			return
		}

		sqlstr = fmt.Sprintf("UPDATE %s.data "+
			"SET datestamp=NOW(), setspec=$1, metadata=$2, dirty=FALSE, public=$3, catalog=$4, seq=NEXTVAL('lastchange') "+
			"WHERE uuidStr=$5", s.dbschema)
		params := []interface{}{
			pq.Array(data.Set),
			string(dataMetaBytes),
			data.Public,
			pq.Array(data.Catalog),
			uuidStr}
		if _, err := s.db.Exec(sqlstr, params); err != nil {
			sendResult("error", fmt.Sprintf("cannot update [%s] - [%v]: %v", sqlstr, params, err), uuidStr)
			return
		}
		sendResult("ok", "update done", uuidStr)
		return
	}
}
