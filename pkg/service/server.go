package service

import (
	"context"
	"crypto/sha512"
	"crypto/tls"
	"fmt"
	"github.com/bluele/gcache"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
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
	extUrl               *url.URL
}

func NewServer(service, addr, addrExt, name, password string, log zLogger.ZLogger, fair *fair.Fair, accessLog io.Writer, jwtKey string, jwtAlg []string, linkTokenExp time.Duration) (*Server, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot split address %s", addr)
	}
	extUrl, err := url.Parse(strings.TrimSuffix(addrExt, "/"))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot parse external address %s", addrExt)
	}

	srv := &Server{
		extUrl:               extUrl,
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
	//router := mux.NewRouter()

	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.GET("/resolver/*pid", s.resolverHandler)
	router.GET("/redir/:uuid", s.redirectHandler)

	partition := router.Group("/:partition", cors.Default())
	partitionAuth := partition.Group("/", gin.BasicAuth(gin.Accounts{
		s.name: s.password,
	}))

	fsys, err := fs.Sub(staticFS, "static")
	if err != nil {
		return errors.Wrap(err, "cannot get subtree of embedded static")
	}
	router.StaticFS("/static", http.FS(fsys))
	for _, part := range s.fair.GetPartitions() {
		router.StaticFS(fmt.Sprintf("/%s/static", part.Name), http.FS(fsys))
	}
	//partition.StaticFS("/static", http.FS(fsys))

	/*
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
	*/

	partition.Group("/oai", func() gin.HandlerFunc {
		return func(c *gin.Context) {
			if c.NegotiateFormat(gin.MIMEXML) == gin.MIMEXML {
				c.Writer.Write([]byte(xmlHeader))
			}
			c.Next()
		}
	}()).GET("/:context", s.oaiHandler)

	partition.GET("/ping", s.pingHandler)
	partition.GET("/", s.partitionHandler)
	partitionAuth.GET("/viewer", s.dataViewerHandler)
	partitionAuth.GET("/viewer/search", s.searchViewerHandler)
	partitionAuth.GET("/viewer/item/:uuid", s.detailHandler)
	partition.GET("/oai/", s.partitionOAIHandler)
	partition.POST("/item", JWTInterceptor.JWTInterceptorGIN(
		s.service,
		"CreateItem",
		JWTInterceptor.Secure,
		s.jwtKey,
		s.jwtAlg,
		sha512.New(),
		s.log,
	), s.createHandler)
	partition.POST("/source", JWTInterceptor.JWTInterceptorGIN(
		s.service,
		"SetSource",
		JWTInterceptor.Secure,
		s.jwtKey,
		s.jwtAlg,
		sha512.New(),
		s.log,
	), s.setSourceHandler)
	partition.POST("/item/:uuid/originaldata", JWTInterceptor.JWTInterceptorGIN(
		s.service,
		"OriginalDataWrite",
		JWTInterceptor.Secure,
		s.jwtKey,
		s.jwtAlg,
		sha512.New(),
		s.log,
	), s.originalDataWriteHandler)
	partition.GET("/item/:uuid/originaldata", JWTInterceptor.JWTInterceptorGIN(
		s.service,
		"OriginalDataRead",
		JWTInterceptor.Secure,
		s.jwtKey,
		s.jwtAlg,
		sha512.New(),
		s.log,
	), s.originalDataReadHandler)
	partition.POST("/startupdate", JWTInterceptor.JWTInterceptorGIN(
		s.service,
		"StartUpdate",
		JWTInterceptor.Secure,
		s.jwtKey,
		s.jwtAlg,
		sha512.New(),
		s.log,
	), s.startUpdateHandler)
	partition.POST("/endupdate", JWTInterceptor.JWTInterceptorGIN(
		s.service,
		"EndUpdate",
		JWTInterceptor.Secure,
		s.jwtKey,
		s.jwtAlg,
		sha512.New(),
		s.log,
	), s.endUpdateHandler)
	partition.POST("/abortupdate", JWTInterceptor.JWTInterceptorGIN(
		s.service,
		"AbortUpdate",
		JWTInterceptor.Secure,
		s.jwtKey,
		s.jwtAlg,
		sha512.New(),
		s.log,
	), s.abortUpdateHandler)
	partition.POST("/archive", JWTInterceptor.JWTInterceptorGIN(
		s.service,
		"AddArchive",
		JWTInterceptor.Secure,
		s.jwtKey,
		s.jwtAlg,
		sha512.New(),
		s.log,
	), s.createArchiveHandler)
	partition.POST("/archive/:archive", JWTInterceptor.JWTInterceptorGIN(
		s.service,
		"AddArchiveItem",
		JWTInterceptor.Secure,
		s.jwtKey,
		s.jwtAlg,
		sha512.New(),
		s.log,
	), s.addArchiveItemHandler)
	partition.GET("/archive/:archive", JWTInterceptor.JWTInterceptorGIN(
		s.service,
		"GetArchiveItem",
		JWTInterceptor.Secure,
		s.jwtKey,
		s.jwtAlg,
		sha512.New(),
		s.log,
	), s.getArchiveItemHandler)
	partition.GET("/item/:uuid/:outputType", s.itemHandler)
	partition.GET("/item/:uuid", s.itemHandler)
	partitionAuth.GET("/createdoi/:uuid", s.createDOIHandler)
	partition.GET("/redir/:uuid", s.redirectHandler)

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
