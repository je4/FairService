package main

import (
	"github.com/BurntSushi/toml"
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
	DSN        string
	ConnMax    int `toml:"connection_max"`
	Schema     string
}

type Endpoint struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type SSHTunnel struct {
	User           string   `toml:"user"`
	PrivateKey     string   `toml:"privatekey"`
	LocalEndpoint  Endpoint `toml:"localendpoint"`
	ServerEndpoint Endpoint `toml:"serverendpoint"`
	RemoteEndpoint Endpoint `toml:"remoteendpoint"`
}

type Config struct {
	Logfile      string      `toml:"logfile"`
	Loglevel     string      `toml:"loglevel"`
	Logformat    string      `toml:"logformat"`
	AccessLog    string      `toml:"accesslog"`
	Addr         string      `toml:"addr"`
	AddrExt      string      `toml:"addrext"`
	CertPEM      string      `toml:"certpem"`
	KeyPEM       string      `toml:"keypem"`
	JWTKey       string      `toml:"jwtkey"`
	JWTAlg       []string    `toml:"jwtalg"`
	LinkTokenExp duration    `toml:"linktokenexp"`
	DB           CfgDatabase `toml:"database"`
	SSHTunnel    SSHTunnel   `toml:"sshtunnel"`
}

func LoadConfig(filepath string) Config {
	var conf Config
	_, err := toml.DecodeFile(filepath, &conf)
	if err != nil {
		log.Fatalln("Error on loading config: ", err)
	}
	conf.AddrExt = strings.TrimRight(conf.AddrExt, "/")

	return conf
}
