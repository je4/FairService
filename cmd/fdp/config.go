package main

import (
	"github.com/BurntSushi/toml"
	"github.com/je4/utils/v2/pkg/config"
	"github.com/je4/utils/v2/pkg/stashconfig"
	"go.ub.unibas.ch/cloud/certloader/v2/pkg/loader"
	"log"
	"strings"
)

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
	Pagesize               int64           `toml:"pagesize"`
	ResumptionTokenTimeout config.Duration `toml:"resumptiontokentimeout"`
	RepositoryName         string          `toml:"repositoryname"`
	AdminEmail             []string        `toml:"adminemail"`
	SampleIdentifier       string          `toml:"sampleidentifier"`
	Delimiter              string          `toml:"delimiter"`
	Scheme                 string          `toml:"scheme"`
}

type Partition struct {
	Name        string           `toml:"name"`
	AddrExt     string           `toml:"addrext"`
	Description string           `toml:"description"`
	Domain      string           `toml:"domain"`
	OAI         OAI              `toml:"oai"`
	ARK         PartitionARK     `toml:"ark"`
	Handle      Handle           `toml:"handle"`
	Datacite    Datacite         `toml:"datacite"`
	JWTKey      config.EnvString `toml:"jwtkey"`
	JWTAlg      []string         `toml:"jwtalg"`
}

type PartitionARK struct {
	NAAN     string `toml:"naan"`
	Shoulder string `toml:"shoulder"`
	Prefix   string `toml:"prefix"`
}

type Handle struct {
	ServiceName    string           `toml:"servicename"`
	Addr           string           `toml:"addr"`
	JWTKey         config.EnvString `toml:"jwtkey"`
	JWTAlg         string           `toml:"jwtalg"`
	SkipCertVerify bool             `toml:"skipcertverify"`
	ID             string           `toml:"id"`
	Prefix         string           `toml:"prefix"`
}

type Datacite struct {
	Api      string           `toml:"api"`
	User     config.EnvString `toml:"user"`
	Password config.EnvString `toml:"password"`
	Prefix   string           `toml:"prefix"`
}

type Config struct {
	ServiceName  string             `toml:"servicename"`
	Log          stashconfig.Config `toml:"log"`
	AccessLog    string             `toml:"accesslog"`
	Addr         string             `toml:"addr"`
	AddrExt      string             `toml:"addrext"`
	TLSConfig    *loader.Config     `toml:"tls"`
	JWTKey       config.EnvString   `toml:"jwtkey"`
	JWTAlg       []string           `toml:"jwtalg"`
	AdminBearer  config.EnvString   `toml:"adminbearer"`
	LinkTokenExp config.Duration    `toml:"linktokenexp"`
	DB           CfgDatabase        `toml:"database"`
	//Tunnel       map[string]SSHTunnel `toml:"tunnel"`
	Partition []Partition      `toml:"partition"`
	Handle    Handle           `toml:"handle"`
	UserName  string           `toml:"username"`
	Password  config.EnvString `toml:"password"`
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
