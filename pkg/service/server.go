package service

import (
	"context"
	"crypto/tls"
	"database/sql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	dcert "github.com/je4/utils/v2/pkg/cert"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type Server struct {
	host, port   string
	addrExt      *url.URL
	srv          *http.Server
	jwtKey       string
	jwtAlg       []string
	linkTokenExp time.Duration
	log          *logging.Logger
	accessLog    io.Writer
	db           *sql.DB
	dbschema     string
}

func NewServer(addr, addrExt string, log *logging.Logger, db *sql.DB, dbschema string, accessLog io.Writer, jwtKey string, jwtAlg []string, linkTokenExp time.Duration) (*Server, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot split address %s", addr)
	}
	extUrl, err := url.Parse(addrExt)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse external address %s", addrExt)
	}

	srv := &Server{
		host:         host,
		port:         port,
		addrExt:      extUrl,
		log:          log,
		db:           db,
		dbschema:     dbschema,
		accessLog:    accessLog,
		jwtKey:       jwtKey,
		jwtAlg:       jwtAlg,
		linkTokenExp: linkTokenExp,
	}
	return srv, nil
}

func (s *Server) ListenAndServe(cert, key string) (err error) {
	router := mux.NewRouter()

	router.Handle(
		"/item",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.createHandler) }()),
	).Methods("POST")

	loggedRouter := handlers.CombinedLoggingHandler(s.accessLog, handlers.ProxyHeaders(router))
	addr := net.JoinHostPort(s.host, s.port)
	s.srv = &http.Server{
		Handler: loggedRouter,
		Addr:    addr,
	}
	if s.addrExt == nil {
		s.addrExt, err = url.Parse(addr)
		if err != nil {
			return errors.Wrapf(err, "cannot parse addr %s", addr)
		}
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
