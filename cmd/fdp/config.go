package main

import (
	"github.com/BurntSushi/toml"
	"github.com/je4/utils/v2/pkg/config"
	"go.ub.unibas.ch/cloud/certloader/v2/pkg/loader"
	"log"
	"strings"
	"time"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type CfgDatabase struct {
	ServerType string
	DSN        config.EnvString
	ConnMax    int `toml:"connection_max"`
	Schema     string
}

type Endpoint struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type Forward struct {
	Local  Endpoint `toml:"local"`
	Remote Endpoint `toml:"remote"`
}

type SSHTunnel struct {
	User       string             `toml:"user"`
	PrivateKey string             `toml:"privatekey"`
	Endpoint   Endpoint           `toml:"endpoint"`
	Forward    map[string]Forward `toml:"forward"`
}

type OAI struct {
	Pagesize               int64    `toml:"pagesize"`
	ResumptionTokenTimeout duration `toml:"resumptiontokentimeout"`
	RepositoryName         string   `toml:"repositoryname"`
	AdminEmail             []string `toml:"adminemail"`
	SampleIdentifier       string   `toml:"sampleidentifier"`
	Delimiter              string   `toml:"delimiter"`
	Scheme                 string   `toml:"scheme"`
}

type Partition struct {
	Name         string   `toml:"name"`
	AddrExt      string   `toml:"addrext"`
	Description  string   `toml:"description"`
	Domain       string   `toml:"domain"`
	HandlePrefix string   `toml:"handleprefix"`
	OAI          OAI      `toml:"oai"`
	HandleID     string   `toml:"handleid"`
	JWTKey       string   `toml:"jwtkey"`
	JWTAlg       []string `toml:"jwtalg"`
}

type Handle struct {
	ServiceName    string `toml:"servicename"`
	Addr           string `toml:"addr"`
	JWTKey         string `toml:"jwtkey"`
	JWTAlg         string `toml:"jwtalg"`
	SkipCertVerify bool   `toml:"skipcertverify"`
}

type Datacite struct {
	Api      string `toml:"api"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	Prefix   string `toml:"prefix"`
}

type Config struct {
	ServiceName  string               `toml:"servicename"`
	Logfile      string               `toml:"logfile"`
	Loglevel     string               `toml:"loglevel"`
	Logformat    string               `toml:"logformat"`
	AccessLog    string               `toml:"accesslog"`
	Addr         string               `toml:"addr"`
	TLSConfig    *loader.Config       `toml:"tls"`
	JWTKey       string               `toml:"jwtkey"`
	JWTAlg       []string             `toml:"jwtalg"`
	LinkTokenExp duration             `toml:"linktokenexp"`
	DB           CfgDatabase          `toml:"database"`
	Tunnel       map[string]SSHTunnel `toml:"tunnel"`
	Partition    []Partition          `toml:"partition"`
	Handle       Handle               `toml:"handle"`
	Datacite     Datacite             `toml:"datacite"`
	UserName     string               `toml:"username"`
	Password     string               `toml:"password"`
}

func LoadConfig(filepath string) Config {
	var conf Config
	conf.ServiceName = "FairService"
	_, err := toml.DecodeFile(filepath, &conf)
	if err != nil {
		log.Fatalln("Error on loading config: ", err)
	}
	//	conf.AddrExt = strings.TrimRight(conf.AddrExt, "/")

	conf.Handle.Addr = strings.TrimRight(conf.Handle.Addr, "/")
	return conf
}
