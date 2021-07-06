package service

import (
	"context"
	"crypto/tls"
	"database/sql"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	dcert "github.com/je4/utils/v2/pkg/cert"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

type Source struct {
	ID          int64
	Name        string
	DetailURL   string
	Description string
}

type Server struct {
	host, port   string
	srv          *http.Server
	jwtKey       string
	jwtAlg       []string
	linkTokenExp time.Duration
	log          *logging.Logger
	accessLog    io.Writer
	db           *sql.DB
	dbschema     string
	sourcesMutex sync.RWMutex
	sources      map[int64]*Source
	Partitions   map[string]*Partition
}

func NewServer(addr string, log *logging.Logger, db *sql.DB, dbschema string, accessLog io.Writer, linkTokenExp time.Duration) (*Server, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot split address %s", addr)
	}
	/*
		extUrl, err := url.Parse(addrExt)
		if err != nil {
			return nil, errors.Wrapf(err, "cannot parse external address %s", addrExt)
		}
	*/

	srv := &Server{
		host:         host,
		port:         port,
		log:          log,
		db:           db,
		dbschema:     dbschema,
		accessLog:    accessLog,
		Partitions:   make(map[string]*Partition),
		linkTokenExp: linkTokenExp,
		sourcesMutex: sync.RWMutex{},
	}

	if err := srv.LoadSources(); err != nil {
		return nil, errors.Wrap(err, "cannot load sources")
	}

	return srv, nil
}

func (s *Server) AddPartition(p *Partition) {
	p.s = s
	s.Partitions[p.Name] = p
}

func (s *Server) LoadSources() error {
	sqlstr := fmt.Sprintf("SELECT sourceid, name, detailurl, description FROM %s.source", s.dbschema)
	rows, err := s.db.Query(sqlstr)
	if err != nil {
		return errors.Wrapf(err, "cannot execute %s", sqlstr)
	}
	defer rows.Close()
	s.sourcesMutex.Lock()
	defer s.sourcesMutex.Unlock()
	s.sources = make(map[int64]*Source)
	for rows.Next() {
		src := &Source{}
		if err := rows.Scan(&src.ID, &src.Name, &src.DetailURL, &src.Description); err != nil {
			return errors.Wrap(err, "cannot scan values")
		}
		s.sources[src.ID] = src
	}
	return nil
}

func (s *Server) GetSourceById(id int64) (*Source, error) {
	s.sourcesMutex.RLock()
	defer s.sourcesMutex.RUnlock()
	if s, ok := s.sources[id]; ok {
		return s, nil
	} else {
		return nil, errors.New(fmt.Sprintf("source #%v not found", id))
	}
}

func (s *Server) GetSourceByName(name string) (*Source, error) {
	s.sourcesMutex.RLock()
	defer s.sourcesMutex.RUnlock()
	for _, src := range s.sources {
		if src.Name == name {
			return src, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("source %s not found", name))
}

func (s *Server) ListenAndServe(cert, key string) (err error) {
	router := mux.NewRouter()

	router.Handle(
		"/{partition}/item",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.createHandler) }()),
	).Methods("POST")

	loggedRouter := handlers.CombinedLoggingHandler(s.accessLog, handlers.ProxyHeaders(router))
	addr := net.JoinHostPort(s.host, s.port)
	s.srv = &http.Server{
		Handler: loggedRouter,
		Addr:    addr,
	}

	if cert == "auto" || key == "auto" {
		s.log.Info("generating new certificate")
		cert, err := dcert.DefaultCertificate()
		if err != nil {
			return errors.Wrap(err, "cannot generate default certificate")
		}
		s.srv.TLSConfig = &tls.Config{Certificates: []tls.Certificate{*cert}}
		s.log.Infof("starting FAIR Data Point at https://%v/", addr)
		return s.srv.ListenAndServeTLS("", "")
	} else if cert != "" && key != "" {
		s.log.Infof("starting FAIR Data Point at https://%v", addr)
		return s.srv.ListenAndServeTLS(cert, key)
	} else {
		s.log.Infof("starting FAIR Data Point at http://%v", addr)
		return s.srv.ListenAndServe()
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
