package main

import (
	"context"
	"crypto/tls"
	"flag"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/model/dataciteModel"
	"github.com/je4/FairService/v2/pkg/service"
	"github.com/je4/utils/v2/pkg/zLogger"
	ublogger "gitlab.switch.ch/ub-unibas/go-ublogger/v2"
	"go.ub.unibas.ch/cloud/certloader/v2/pkg/loader"
	"io"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
)

type queryTracer struct {
	log zLogger.ZLogger
}

func (tracer *queryTracer) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	tracer.log.Debug().Msgf("postgreSQL command start '%s' - %v", data.SQL, data.Args)
	return ctx
}

func (tracer *queryTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		tracer.log.Error().Err(data.Err).Msgf("postgreSQL command error")
		return
	}
	tracer.log.Debug().Msgf("postgreSQL command end: %s (%d)", data.CommandTag.String(), data.CommandTag.RowsAffected())
}

func main() {
	cfgFile := flag.String("cfg", "/etc/fdp.toml", "locations of config file")
	flag.Parse()
	config := LoadConfig(*cfgFile)

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("cannot get hostname: %v", err)
	}

	var loggerTLSConfig *tls.Config
	var loggerLoader io.Closer
	if config.Log.Stash.TLS != nil {
		loggerTLSConfig, loggerLoader, err = loader.CreateClientLoader(config.Log.Stash.TLS, nil)
		if err != nil {
			log.Fatalf("cannot create client loader: %v", err)
		}
		defer loggerLoader.Close()
	}

	_logger, _logstash, _logfile, err := ublogger.CreateUbMultiLoggerTLS(config.Log.Level, config.Log.File,
		ublogger.SetDataset(config.Log.Stash.Dataset),
		ublogger.SetLogStash(config.Log.Stash.LogstashHost, config.Log.Stash.LogstashPort, config.Log.Stash.Namespace, config.Log.Stash.LogstashTraceLevel),
		ublogger.SetTLS(config.Log.Stash.TLS != nil),
		ublogger.SetTLSConfig(loggerTLSConfig),
	)
	if err != nil {
		log.Fatalf("cannot create logger: %v", err)
	}
	if _logstash != nil {
		defer _logstash.Close()
	}
	if _logfile != nil {
		defer _logfile.Close()
	}

	l2 := _logger.With().Timestamp().Str("host", hostname).Logger() //.Output(output)
	var logger zLogger.ZLogger = &l2

	pgxConf, err := pgxpool.ParseConfig(string(config.DB.DSN))
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot parse db connection string")
	}
	//	pgxConf.TLSConfig = &tls.ARKConfig{InsecureSkipVerify: true, ServerName: "dd-pdb3.ub.unibas.ch"}
	// create prepared queries on each connection
	/*
		pgxConf.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
			return service.AfterConnectFunc(ctx, conn, logger)
		}
	*/
	pgxConf.BeforeConnect = func(ctx context.Context, cfg *pgx.ConnConfig) error {
		cfg.Tracer = &queryTracer{log: logger}
		return nil
	}
	var conn *pgxpool.Pool
	var dbstrRegexp = regexp.MustCompile(`^postgres://postgres:([^@]+)@.+$`)
	pws := dbstrRegexp.FindStringSubmatch(string(config.DB.DSN))
	if len(pws) == 2 {
		logger.Info().Msgf("connecting to database: %s", strings.Replace(string(config.DB.DSN), pws[1], "xxxxxxxx", -1))
	} else {
		logger.Info().Msgf("connecting to database")
	}
	db, err := pgxpool.NewWithConfig(context.Background(), pgxConf)
	//conn, err = pgx.ConnectConfig(context.Background(), pgxConf)
	if err != nil {
		logger.Fatal().Err(err).Msgf("cannot connect to database: %s", config.DB.DSN)
	}
	defer conn.Close()

	if err := db.Ping(context.Background()); err != nil {
		logger.Error().Err(err).Msg("cannot ping database")
	} else {
		logger.Info().Msg("database connection established")
	}

	var accessLog io.Writer
	var f *os.File
	if config.AccessLog == "" {
		accessLog = os.Stdout
	} else {
		f, err = os.OpenFile(config.AccessLog, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			logger.Fatal().Msgf("cannot open file %s: %v", config.AccessLog, err)
			return
		}
		defer f.Close()
		accessLog = f
	}

	fairService, err := fair.NewFair(db, config.DB.Schema, logger)
	if err != nil {
		logger.Fatal().Msgf("cannot initialize fair: %v", err)
	}

	//	var partitions []*fair.Partition
	for _, pconf := range config.Partition {
		partition, err := fair.NewPartition(
			fairService,
			pconf.Name,
			pconf.AddrExt,
			pconf.Domain,
			&fair.OAIConfig{
				RepositoryName:         pconf.OAI.RepositoryName,
				AdminEmail:             pconf.OAI.AdminEmail,
				SampleIdentifier:       pconf.OAI.SampleIdentifier,
				Delimiter:              pconf.OAI.Delimiter,
				Scheme:                 pconf.OAI.Scheme,
				PageSize:               pconf.OAI.Pagesize,
				ResumptionTokenTimeout: time.Duration(pconf.OAI.ResumptionTokenTimeout),
			},
			pconf.Description,
			string(pconf.JWTKey),
			pconf.JWTAlg,
			logger)
		if err != nil {
			logger.Fatal().Msgf("cannot create partition %s: %v", pconf.Name, err)
			return
		}
		mr, err := fair.NewResolver(partition, logger)
		if err != nil {
			logger.Fatal().Err(err).Msg("cannot create resolver")
			return
		}
		if _, err := fair.NewARKService(mr, &fair.ARKConfig{
			Shoulder: pconf.ARK.Shoulder,
			Prefix:   pconf.ARK.Prefix,
			NAAN:     pconf.ARK.NAAN,
		}, logger); err != nil {
			logger.Fatal().Msgf("cannot create ark service: %v", err)
			return
		}
		if _, err := fair.NewHandleService(mr, &fair.HandleConfig{
			ServiceName:    pconf.Handle.ServiceName,
			Addr:           pconf.Handle.Addr,
			JWTKey:         pconf.Handle.JWTKey.String(),
			JWTAlg:         pconf.Handle.JWTAlg,
			SkipCertVerify: pconf.Handle.SkipCertVerify,
			ID:             pconf.Handle.ID,
			Prefix:         pconf.Handle.Prefix,
		}, logger); err != nil {
			logger.Fatal().Msgf("cannot create handle service: %v", err)
		}
		mr.CreateAll(partition, dataciteModel.RelatedIdentifierTypeARK)
		if _, err := fair.NewDataciteService(mr, fair.DataciteConfig{
			Api:      pconf.Datacite.Api,
			User:     pconf.Datacite.User.String(),
			Password: pconf.Datacite.Password.String(),
			Prefix:   pconf.Datacite.Prefix,
		}, logger); err != nil {
			logger.Fatal().Msgf("cannot create datacite service: %v", err)
		}
		partition.AddResolver(mr)
		fairService.AddPartition(partition)
	}

	// create TLS Certificate.
	// the certificate MUST contain <package>.<service> as DNS name
	serverTLSConfig, serverLoader, err := loader.CreateServerLoader(false, config.TLSConfig, nil, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server loader")
	}
	defer serverLoader.Close()

	srv, err := service.NewServer(config.ServiceName,
		config.Addr,
		config.AddrExt,
		config.UserName,
		config.Password.String(),
		config.AdminBearer.String(),
		logger,
		fairService,
		accessLog,
		config.JWTKey.String(),
		config.JWTAlg,
		time.Duration(config.LinkTokenExp))
	if err != nil {
		logger.Fatal().Msgf("cannot initialize server: %v", err)
	}

	go func() {
		if err := srv.ListenAndServe(serverTLSConfig); err != nil {
			log.Fatalf("server died: %v", err)
		}
	}()

	end := make(chan bool, 1)

	// process waiting for interrupt signal (TERM or KILL)
	go func() {
		sigint := make(chan os.Signal, 1)

		// interrupt signal sent from terminal
		signal.Notify(sigint, os.Interrupt)

		signal.Notify(sigint, syscall.SIGTERM)
		signal.Notify(sigint, syscall.SIGKILL)

		<-sigint

		// We received an interrupt signal, shut down.
		logger.Info().Msgf("shutdown requested")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.Shutdown(ctx)

		end <- true
	}()

	<-end
	logger.Info().Msg("server stopped")

}
