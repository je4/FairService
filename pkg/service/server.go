package service

import (
	"context"
	"crypto/tls"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	fair "github.com/je4/FairService/v2/pkg/fair"
	dcert "github.com/je4/utils/v2/pkg/cert"
	"github.com/op/go-logging"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"time"
)

type Server struct {
	host, port   string
	srv          *http.Server
	linkTokenExp time.Duration
	log          *logging.Logger
	accessLog    io.Writer
	fair         *fair.Fair
}

func NewServer(addr string, log *logging.Logger, fair *fair.Fair, accessLog io.Writer, linkTokenExp time.Duration) (*Server, error) {
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
		accessLog:    accessLog,
		linkTokenExp: linkTokenExp,
		fair:         fair,
	}

	return srv, nil
}

func (s *Server) ListenAndServe(cert, key string) (err error) {
	router := mux.NewRouter()

	router.Handle(
		"/{partition}/oai2",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.oaiHandler) }()),
	).Methods("GET")
	router.Handle(
		"/{partition}/item",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.createHandler) }()),
	).Methods("POST")
	router.Handle(
		"/{partition}/startupdate",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.startUpdateHandler) }()),
	).Methods("POST")
	router.Handle(
		"/{partition}/endupdate",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.endUpdateHandler) }()),
	).Methods("POST")
	router.Handle(
		"/{partition}/abortupdate",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.abortUpdateHandler) }()),
	).Methods("POST")
	router.Handle(
		"/{partition}/item/{uuid}/{outputType}",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.itemHandler) }()),
	).Methods("GET")
	router.Handle(
		"/{partition}/item/{uuid}",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.itemHandler) }()),
	).Methods("GET")

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
