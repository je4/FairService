package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/service"
	lm "github.com/je4/utils/v2/pkg/logger"
	"github.com/je4/utils/v2/pkg/ssh"
	_ "github.com/lib/pq"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfgfile := flag.String("cfg", "./fdp.toml", "locations of config file")
	flag.Parse()
	config := LoadConfig(*cfgfile)

	// create logger instance
	logger, lf := lm.CreateLogger("FAIRService", config.Logfile, nil, config.Loglevel, config.Logformat)
	defer lf.Close()

	if config.SSHTunnel.User != "" && config.SSHTunnel.PrivateKey != "" {
		tunnels := map[string]*ssh.SourceDestination{}
		for _, t := range config.SSHTunnel.Tunnel {
			tunnels[t.Name] = &ssh.SourceDestination{
				Local: &ssh.Endpoint{
					Host: t.LocalEndpoint.Host,
					Port: t.LocalEndpoint.Port,
				},
				Remote: &ssh.Endpoint{
					Host: t.RemoteEndpoint.Host,
					Port: t.RemoteEndpoint.Port,
				},
			}
		}
		tunnel, err := ssh.NewSSHTunnel(
			config.SSHTunnel.User,
			config.SSHTunnel.PrivateKey,
			&ssh.Endpoint{
				Host: config.SSHTunnel.ServerEndpoint.Host,
				Port: config.SSHTunnel.ServerEndpoint.Port,
			},
			tunnels,
			logger,
		)
		if err != nil {
			logger.Errorf("cannot create sshtunnel %v@%v:%v - %v", config.SSHTunnel.User, config.SSHTunnel.ServerEndpoint.Host, &config.SSHTunnel.ServerEndpoint.Port, err)
			return
		}
		if err := tunnel.Start(); err != nil {
			logger.Errorf("cannot create sshtunnel %v - %v", tunnel.String(), err)
			return
		}
		defer tunnel.Close()
		time.Sleep(2 * time.Second)
	}

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

	var accesslog io.Writer
	if config.AccessLog == "" {
		accesslog = os.Stdout
	} else {
		f, err := os.OpenFile(config.AccessLog, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			logger.Panicf("cannot open file %s: %v", config.AccessLog, err)
			return
		}
		defer f.Close()
		accesslog = f
	}

	var partitions []*fair.Partition
	for _, pconf := range config.Partition {
		p, err := fair.NewPartition(pconf.Name, pconf.AddrExt, pconf.JWTKey, pconf.JWTAlg)
		if err != nil {
			logger.Panicf("cannot create partition %s: %v", pconf.Name, err)
			return
		}
		partitions = append(partitions, p)
	}

	fair, err := fair.NewFair(db, config.DB.Schema, logger)
	if err != nil {
		logger.Panicf("cannot initialize fair: %v", err)
	}
	for _, p := range partitions {
		fair.AddPartition(p)
	}

	srv, err := service.NewServer(config.Addr, logger, fair, accesslog, config.LinkTokenExp.Duration)
	if err != nil {
		logger.Panicf("cannot initialize server: %v", err)
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
		logger.Infof("shutdown requested")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.Shutdown(ctx)

		end <- true
	}()

	<-end
	logger.Info("server stopped")

}
