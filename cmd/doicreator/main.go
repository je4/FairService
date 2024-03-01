package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/service"
	"github.com/je4/FairService/v2/pkg/service/datacite"
	hcClient "github.com/je4/HandleCreator/v2/pkg/client"
	"github.com/je4/utils/v2/pkg/zLogger"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfgFile := flag.String("cfg", "/etc/fdp.toml", "locations of config file")
	flag.Parse()
	config := LoadConfig(*cfgFile)

	var out io.Writer = os.Stdout
	if config.Logfile != "" {
		fp, err := os.OpenFile(config.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalf("cannot open logfile %s: %v", config.Logfile, err)
		}
		defer fp.Close()
		out = fp
	}

	output := zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}
	_logger := zerolog.New(output).With().Timestamp().Logger()
	_logger.Level(zLogger.LogLevel(config.Loglevel))
	var logger zLogger.ZLogger = &_logger

	// get database connection handle
	db, err := sql.Open(config.DB.ServerType, config.DB.DSN)
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

	var accessLog io.Writer
	var f *os.File
	if config.AccessLog == "" {
		accessLog = os.Stdout
	} else {
		f, err = os.OpenFile(config.AccessLog, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			logger.Panic().Msgf("cannot open file %s: %v", config.AccessLog, err)
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
			logger.Panic().Msgf("cannot create datacite client: %v", err)
			return
		}

		if err := dataciteClient.Heartbeat(); err != nil {
			logger.Panic().Msgf("cannot check datacite heartbeat: %v", err)
			return
		}

		/*
			r, err := dataciteClient.RetrieveDOI("10.5438/0012")
			if err != nil {
				logger.Panic().Msgf("cannot get doi: %v", err)
				return
			}
			logger.Infof("doi: %v", r)
		*/
	}
	var handle *hcClient.HandleCreatorClient
	if config.Handle.Addr != "" {
		handle, err = hcClient.NewHandleCreatorClient(config.Handle.ServiceName, config.Handle.Addr, config.Handle.JWTKey, config.Handle.JWTAlg, config.Handle.SkipCertVerify, logger)
		if err != nil {
			logger.Panic().Msgf("cannot create handle service: %v", err)
			return
		}
		if err := handle.Ping(); err != nil {
			logger.Fatal().Msgf("cannot ping handle server on %s: %v", config.Handle.Addr, err)
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
			logger.Panic().Msgf("cannot create partition %s: %v", pconf.Name, err)
			return
		}
		partitions = append(partitions, p)
	}

	fair, err := fair.NewFair(db, config.DB.Schema, handle, dataciteClient, logger)
	if err != nil {
		logger.Panic().Msgf("cannot initialize fair: %v", err)
	}
	for _, p := range partitions {
		fair.AddPartition(p)
	}

	srv, err := service.NewServer(config.ServiceName, config.Addr, config.UserName, config.Password, logger, fair, accessLog, config.JWTKey, config.JWTAlg, config.LinkTokenExp.Duration)
	if err != nil {
		logger.Panic().Msgf("cannot initialize server: %v", err)
	}
	go func() {
		if err := srv.ListenAndServe(config.CertPEM, config.KeyPEM); err != nil {
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
		logger.Info().Msg("shutdown requested")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.Shutdown(ctx)

		end <- true
	}()

	<-end
	logger.Info().Msg("server stopped")

}
