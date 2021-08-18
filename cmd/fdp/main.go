package main

import (
	"context"
	"crypto/sha512"
	"database/sql"
	"flag"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/service"
	"github.com/je4/utils/v2/pkg/JWTInterceptor"
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
	cfgFile := flag.String("cfg", "/etc/fdp.toml", "locations of config file")
	flag.Parse()
	config := LoadConfig(*cfgFile)

	// create logger instance
	logger, lf := lm.CreateLogger("FAIRService", config.Logfile, nil, config.Loglevel, config.Logformat)
	defer lf.Close()

	var tunnels []*ssh.SSHtunnel
	for name, tunnel := range config.Tunnel {
		logger.Infof("starting tunnel %s", name)

		forwards := make(map[string]*ssh.SourceDestination)
		for fwName, fw := range tunnel.Forward {
			forwards[fwName] = &ssh.SourceDestination{
				Local: &ssh.Endpoint{
					Host: fw.Local.Host,
					Port: fw.Local.Port,
				},
				Remote: &ssh.Endpoint{
					Host: fw.Remote.Host,
					Port: fw.Remote.Port,
				},
			}
		}

		t, err := ssh.NewSSHTunnel(
			tunnel.User,
			tunnel.PrivateKey,
			&ssh.Endpoint{
				Host: tunnel.Endpoint.Host,
				Port: tunnel.Endpoint.Port,
			},
			forwards,
			logger,
		)
		if err != nil {
			logger.Errorf("cannot create tunnel %v@%v:%v - %v", tunnel.User, tunnel.Endpoint.Host, tunnel.Endpoint.Port, err)
			return
		}
		if err := t.Start(); err != nil {
			logger.Errorf("cannot create configfile %v - %v", t.String(), err)
			return
		}
		tunnels = append(tunnels, t)
	}
	defer func() {
		for _, t := range tunnels {
			t.Close()
		}
	}()
	// if tunnels are made, wait until connection is established
	if len(config.Tunnel) > 0 {
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

	var accessLog io.Writer
	var f *os.File
	if config.AccessLog == "" {
		accessLog = os.Stdout
	} else {
		f, err = os.OpenFile(config.AccessLog, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			logger.Panicf("cannot open file %s: %v", config.AccessLog, err)
			return
		}
		defer f.Close()
		accessLog = f
	}

	var handle *fair.HandleServiceClient
	if config.Handle.Addr != "" {
		tr, err := JWTInterceptor.NewJWTTransport(nil,
			sha512.New(),
			"",
			config.Handle.JWTKey,
			config.Handle.JWTAlg,
			30*time.Second)
		if err != nil {
			logger.Panicf("cannot create JWTInterceptor.JWTTransport: %v", err)
			return
		}

		handle, err = fair.NewHandleServiceClient(config.Handle.Addr, tr)
		if err != nil {
			logger.Panicf("cannot create handle service: %v", err)
			return
		}
	}
	var partitions []*fair.Partition
	for _, pconf := range config.Partition {
		p, err := fair.NewPartition(
			pconf.Name,
			pconf.AddrExt,
			pconf.Domain,
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
			logger.Panicf("cannot create partition %s: %v", pconf.Name, err)
			return
		}
		partitions = append(partitions, p)
	}

	fair, err := fair.NewFair(db, config.DB.Schema, handle, logger)
	if err != nil {
		logger.Panicf("cannot initialize fair: %v", err)
	}
	for _, p := range partitions {
		fair.AddPartition(p)
	}

	srv, err := service.NewServer(config.Addr, logger, fair, accessLog, config.JWTKey, config.JWTAlg, config.LinkTokenExp.Duration)
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
