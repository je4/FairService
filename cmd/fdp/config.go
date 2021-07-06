package main

import (
	"github.com/BurntSushi/toml"
	"log"
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
	DSN        string
	ConnMax    int `toml:"connection_max"`
	Schema     string
}

type Endpoint struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type Tunnel struct {
	Name           string   `toml:"name"`
	LocalEndpoint  Endpoint `toml:"localendpoint"`
	RemoteEndpoint Endpoint `toml:"remoteendpoint"`
}

type SSHTunnel struct {
	User           string   `toml:"user"`
	PrivateKey     string   `toml:"privatekey"`
	ServerEndpoint Endpoint `toml:"serverendpoint"`
	Tunnel         []Tunnel `toml:"tunnel"`
}

type Partition struct {
	Name    string   `toml:"name"`
	AddrExt string   `toml:"addrext"`
	JWTKey  string   `toml:"jwtkey"`
	JWTAlg  []string `toml:"jwtalg"`
}

type Config struct {
	Logfile      string      `toml:"logfile"`
	Loglevel     string      `toml:"loglevel"`
	Logformat    string      `toml:"logformat"`
	AccessLog    string      `toml:"accesslog"`
	Addr         string      `toml:"addr"`
	CertPEM      string      `toml:"certpem"`
	KeyPEM       string      `toml:"keypem"`
	LinkTokenExp duration    `toml:"linktokenexp"`
	DB           CfgDatabase `toml:"database"`
	SSHTunnel    SSHTunnel   `toml:"sshtunnel"`
	Partition    []Partition `toml:"partition"`
}

func LoadConfig(filepath string) Config {
	var conf Config
	_, err := toml.DecodeFile(filepath, &conf)
	if err != nil {
		log.Fatalln("Error on loading config: ", err)
	}
	//	conf.AddrExt = strings.TrimRight(conf.AddrExt, "/")

	return conf
}
