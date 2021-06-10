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

type Config struct {
	Logfile      string   `toml:"logfile"`
	Loglevel     string   `toml:"loglevel"`
	Logformat    string   `toml:"logformat"`
	AccessLog    string   `toml:"accesslog"`
	Addr         string   `toml:"addr"`
	AddrExt      string   `toml:"addrext"`
	CertPEM      string   `toml:"certpem"`
	KeyPEM       string   `toml:"keypem"`
	JWTKey       string   `toml:"jwtkey"`
	JWTAlg       []string `toml:"jwtalg"`
	LinkTokenExp duration `toml:"linktokenexp"`
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
