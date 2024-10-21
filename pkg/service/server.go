package service

import (
	"context"
	"crypto/sha512"
	"crypto/tls"
	"fmt"
	"github.com/bluele/gcache"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	fair "github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/utils/v2/pkg/JWTInterceptor"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/pkg/errors"
	"html/template"
	"io"
	"io/fs"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Server struct {
	service              string
	host, port           string
	name, password       string
	srv                  *http.Server
	linkTokenExp         time.Duration
	jwtKey               string
	jwtAlg               []string
	log                  zLogger.ZLogger
	accessLog            io.Writer
	fair                 *fair.Fair
	templates            map[string]*template.Template
	resumptionTokenCache gcache.Cache
}

func NewServer(service, addr, name, password string, log zLogger.ZLogger, fair *fair.Fair, accessLog io.Writer, jwtKey string, jwtAlg []string, linkTokenExp time.Duration) (*Server, error) {
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
		service:              service,
		host:                 host,
		port:                 port,
		name:                 name,
		password:             password,
		log:                  log,
		accessLog:            accessLog,
		linkTokenExp:         linkTokenExp,
		fair:                 fair,
		jwtKey:               jwtKey,
		jwtAlg:               jwtAlg,
		templates:            map[string]*template.Template{},
		resumptionTokenCache: gcache.New(500).ARC().Build(),
	}

	return srv, srv.InitTemplates()
}

func (s *Server) InitTemplates() error {
	for key, val := range templateFiles {
		tpl, err := template.ParseFS(templateFS, val)
		if err != nil {
			return errors.New(fmt.Sprintf("cannot parse template %s: %s", key, val))
		}
		s.templates[key] = tpl
	}
	return nil
}

func (s *Server) ListenAndServe(tlsConfig *tls.Config) (err error) {
	if tlsConfig == nil {
		return errors.New("TLS config is nil")
	}
	router := mux.NewRouter()

	fsys, err := fs.Sub(staticFS, "static")
	if err != nil {
		return errors.Wrap(err, "cannot get subtree of embedded static")
	}
	httpStaticServer := http.FileServer(http.FS(fsys))
	router.PathPrefix("/{partition}/static").Handler(
		handlers.CompressHandler(func(prefix string, h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				vars := mux.Vars(r)
				partition, ok := vars["partition"]
				if !ok {
					http.NotFound(w, r)
					return
				}
				fullPrefix := fmt.Sprintf("/%s%s", partition, prefix)
				p := strings.TrimPrefix(r.URL.Path, fullPrefix)
				rp := strings.TrimPrefix(r.URL.RawPath, fullPrefix)
				if len(p) < len(r.URL.Path) && (r.URL.RawPath == "" || len(rp) < len(r.URL.RawPath)) {
					r2 := new(http.Request)
					*r2 = *r
					r2.URL = new(url.URL)
					*r2.URL = *r.URL
					r2.URL.Path = p
					r2.URL.RawPath = rp
					w.Header().Set("Cache-Control", "max-age=3600")
					h.ServeHTTP(w, r2)
				} else {
					http.NotFound(w, r)
				}
			})

		}("/static", httpStaticServer),
		// http.StripPrefix("/static", httpStaticServer)
		),
	).Methods("GET")

	router.Handle(
		"/{partition}/oai/{context}",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.oaiHandler) }()),
	).Methods("GET")

	router.HandleFunc(
		"/{partition}",
		func(w http.ResponseWriter, req *http.Request) {
			vars := mux.Vars(req)
			pName := vars["partition"]

			part, err := s.fair.GetPartition(pName)
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				w.Header().Set("Content-type", "text/plain")
				w.Write([]byte(fmt.Sprintf("partition [%s] not found", pName)))
				return
			}
			http.Redirect(w, req, part.AddrExt+"/", http.StatusPermanentRedirect)
		},
	).Methods("GET")
	router.Handle(
		"/{partition}/ping",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.pingHandler) }()),
	).Methods("GET")
	router.Handle(
		"/{partition}/",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.partitionHandler) }()),
	).Methods("GET")
	router.Handle(
		"/{partition}/viewer",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.dataViewerHandler) }()),
	).Methods("GET")
	router.Handle(
		"/{partition}/viewer/search",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.searchViewerHandler) }()),
	).Methods("GET")
	router.Handle(
		"/{partition}/viewer/item/{uuid}",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.detailHandler) }()),
	).Methods("GET")
	router.Handle(
		"/{partition}/oai/",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.partitionOAIHandler) }()),
	).Methods("GET")
	router.Handle(
		"/{partition}/item",
		handlers.CompressHandler(
			JWTInterceptor.JWTInterceptor(
				s.service,
				"CreateItem",
				JWTInterceptor.Secure,
				func() http.Handler { return http.HandlerFunc(s.createHandler) }(),
				s.jwtKey,
				s.jwtAlg,
				sha512.New(),
				s.log,
			))).
		Methods("POST")
	router.Handle(
		"/{partition}/source",
		handlers.CompressHandler(
			JWTInterceptor.JWTInterceptor(
				s.service,
				"setSource",
				JWTInterceptor.Secure,
				func() http.Handler { return http.HandlerFunc(s.setSourceHandler) }(),
				s.jwtKey,
				s.jwtAlg,
				sha512.New(),
				s.log,
			))).
		Methods("POST")
	router.Handle(
		"/{partition}/item/{uuid}/originaldata",
		handlers.CompressHandler(
			JWTInterceptor.JWTInterceptor(
				s.service,
				"originalDataWrite",
				JWTInterceptor.Secure,
				func() http.Handler { return http.HandlerFunc(s.originalDataWriteHandler) }(),
				s.jwtKey,
				s.jwtAlg,
				sha512.New(),
				s.log,
			))).
		Methods("POST")
	router.Handle(
		"/{partition}/item/{uuid}/originaldata",
		handlers.CompressHandler(
			JWTInterceptor.JWTInterceptor(
				s.service,
				"originalDataRead",
				JWTInterceptor.Secure,
				func() http.Handler { return http.HandlerFunc(s.originalDataReadHandler) }(),
				s.jwtKey,
				s.jwtAlg,
				sha512.New(),
				s.log,
			))).
		Methods("GET")
	router.Handle(
		"/{partition}/startupdate",
		handlers.CompressHandler(
			JWTInterceptor.JWTInterceptor(
				s.service,
				"StartUpdate",
				JWTInterceptor.Secure,
				func() http.Handler { return http.HandlerFunc(s.startUpdateHandler) }(),
				s.jwtKey,
				s.jwtAlg,
				sha512.New(),
				s.log,
			))).
		Methods("POST")
	router.Handle(
		"/{partition}/endupdate",
		handlers.CompressHandler(
			JWTInterceptor.JWTInterceptor(
				s.service,
				"EndUpdate",
				JWTInterceptor.Secure,
				func() http.Handler { return http.HandlerFunc(s.endUpdateHandler) }(),
				s.jwtKey,
				s.jwtAlg,
				sha512.New(),
				s.log,
			))).
		Methods("POST")
	router.Handle(
		"/{partition}/archive",
		handlers.CompressHandler(
			JWTInterceptor.JWTInterceptor(
				s.service,
				"AddArchive",
				JWTInterceptor.Secure,
				func() http.Handler { return http.HandlerFunc(s.createArchiveHandler) }(),
				s.jwtKey,
				s.jwtAlg,
				sha512.New(),
				s.log,
			))).
		Methods("POST")
	router.Handle(
		"/{partition}/archive/{archive}",
		handlers.CompressHandler(
			JWTInterceptor.JWTInterceptor(
				s.service,
				"AddArchiveItem",
				JWTInterceptor.Secure,
				func() http.Handler { return http.HandlerFunc(s.addArchiveItemHandler) }(),
				s.jwtKey,
				s.jwtAlg,
				sha512.New(),
				s.log,
			))).
		Methods("POST")
	router.Handle(
		"/{partition}/archive/{archive}",
		handlers.CompressHandler(
			JWTInterceptor.JWTInterceptor(
				s.service,
				"GetArchiveItem",
				JWTInterceptor.Secure,
				func() http.Handler { return http.HandlerFunc(s.getArchiveItemHandler) }(),
				s.jwtKey,
				s.jwtAlg,
				sha512.New(),
				s.log,
			))).
		Methods("GET")
	router.Handle(
		"/{partition}/abortupdate",
		handlers.CompressHandler(
			JWTInterceptor.JWTInterceptor(
				s.service,
				"AbortUpdate",
				JWTInterceptor.Secure,
				func() http.Handler { return http.HandlerFunc(s.abortUpdateHandler) }(),
				s.jwtKey,
				s.jwtAlg,
				sha512.New(),
				s.log,
			))).
		Methods("POST")
	router.Handle(
		"/{partition}/item/{uuid}/{outputType}",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.itemHandler) }()),
	).Methods("GET")
	router.Handle(
		"/{partition}/item/{uuid}",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.itemHandler) }()),
	).Methods("GET")
	router.Handle(
		"/{partition}/createdoi/{uuid}",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.createDOIHandler) }()),
	).Methods("GET")

	router.Handle(
		"/{partition}/redir/{uuid}",
		handlers.CompressHandler(func() http.Handler { return http.HandlerFunc(s.redirectHandler) }()),
	).Methods("GET")

	loggedRouter := handlers.CombinedLoggingHandler(s.accessLog, handlers.ProxyHeaders(router))
	addr := net.JoinHostPort(s.host, s.port)
	s.srv = &http.Server{
		Handler:   loggedRouter,
		Addr:      addr,
		TLSConfig: tlsConfig,
	}

	for _, part := range s.fair.GetPartitions() {
		s.log.Info().Msgf("starting FAIR Data Point at %v - https://%s/%s", part.AddrExt, addr, part.Name)
	}
	return s.srv.ListenAndServeTLS("", "")
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
