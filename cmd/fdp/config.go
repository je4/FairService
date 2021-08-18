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
	Name        string   `toml:"name"`
	AddrExt     string   `toml:"addrext"`
	Description string   `toml:"description"`
	Domain      string   `toml:"domain"`
	OAI         OAI      `toml:"oai"`
	HandleID    string   `toml:"handleid"`
	JWTKey      string   `toml:"jwtkey"`
	JWTAlg      []string `toml:"jwtalg"`
}

type Handle struct {
	Addr   string `toml:"addr"`
	JWTKey string `toml:"jwtkey"`
	JWTAlg string `toml:"jwtalg"`
}

type Config struct {
	Logfile      string               `toml:"logfile"`
	Loglevel     string               `toml:"loglevel"`
	Logformat    string               `toml:"logformat"`
	AccessLog    string               `toml:"accesslog"`
	Addr         string               `toml:"addr"`
	CertPEM      string               `toml:"certpem"`
	KeyPEM       string               `toml:"keypem"`
	LinkTokenExp duration             `toml:"linktokenexp"`
	DB           CfgDatabase          `toml:"database"`
	Tunnel       map[string]SSHTunnel `toml:"tunnel"`
	Partition    []Partition          `toml:"partition"`
	Handle       Handle               `toml:"handle"`
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
