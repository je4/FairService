package main

import (
	"context"
	"crypto/tls"
	"flag"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/service"
	"github.com/je4/FairService/v2/pkg/service/datacite"
	hcClient "github.com/je4/HandleCreator/v2/pkg/client"
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
	//	pgxConf.TLSConfig = &tls.Config{InsecureSkipVerify: true, ServerName: "dd-pdb3.ub.unibas.ch"}
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
	/*
		// get database connection handle
		db, err := sql.Open(config.DB.ServerType, string(config.DB.DSN))
		if err != nil {
			log.Fatalf("error opening database: %v", err)
		}
		// close on shutdown
		defer db.Close()

		// Open doesn't open a connection. Validate DSN data:
		err = db.Ping()
		if err != nil {
			log.Fatalf("error pinging database: %v", err)
		}

	*/

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

	var dataciteClient *datacite.Client
	if config.Datacite.Api != "" {
		dataciteClient, err = datacite.NewClient(
			config.Datacite.Api,
			config.Datacite.User,
			config.Datacite.Password,
			config.Datacite.Prefix)
		if err != nil {
			logger.Fatal().Msgf("cannot create datacite client: %v", err)
			return
		}

		if err := dataciteClient.Heartbeat(); err != nil {
			logger.Fatal().Msgf("cannot check datacite heartbeat: %v", err)
			return
		}

		/*
			r, err := dataciteClient.RetrieveDOI("10.5438/0012")
			if err != nil {
				logger.Fatal().Msgf("cannot get doi: %v", err)
				return
			}
			logger.Info().Msgf("doi: %v", r)
		*/
	}
	var handle *hcClient.HandleCreatorClient
	if config.Handle.Addr != "" {
		handle, err = hcClient.NewHandleCreatorClient(config.Handle.ServiceName, config.Handle.Addr, config.Handle.JWTKey, config.Handle.JWTAlg, config.Handle.SkipCertVerify, logger)
		if err != nil {
			logger.Fatal().Msgf("cannot create handle service: %v", err)
			return
		}
	} else {
		logger.Info().Msg("no handle creator configured")
	}
	var partitions []*fair.Partition
	for _, pconf := range config.Partition {
		p, err := fair.NewPartition(
			pconf.Name,
			pconf.AddrExt,
			pconf.Domain,
			pconf.HandlePrefix,
			pconf.OAI.RepositoryName,
			pconf.OAI.AdminEmail,
			pconf.OAI.SampleIdentifier,
			pconf.OAI.Delimiter,
			pconf.OAI.Scheme,
			pconf.HandleID,
			pconf.Description,
			pconf.OAI.Pagesize,
			pconf.OAI.ResumptionTokenTimeout.Duration,
			pconf.JWTKey,
			pconf.JWTAlg)
		if err != nil {
			logger.Fatal().Msgf("cannot create partition %s: %v", pconf.Name, err)
			return
		}
		partitions = append(partitions, p)
	}

	fair, err := fair.NewFair(db, config.DB.Schema, handle, dataciteClient, logger)
	if err != nil {
		logger.Fatal().Msgf("cannot initialize fair: %v", err)
	}
	for _, p := range partitions {
		fair.AddPartition(p)
	}

	// create TLS Certificate.
	// the certificate MUST contain <package>.<service> as DNS name
	serverTLSConfig, serverLoader, err := loader.CreateServerLoader(true, config.TLSConfig, nil, logger)
	if err != nil {
		logger.Fatal().Err(err).Msg("cannot create server loader")
	}
	defer serverLoader.Close()

	srv, err := service.NewServer(config.ServiceName, config.Addr, config.UserName, config.Password, logger, fair, accessLog, config.JWTKey, config.JWTAlg, config.LinkTokenExp.Duration)
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
