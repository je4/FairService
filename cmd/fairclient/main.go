/*
Copyright 2020 Center for Digital Matter HGK FHNW, Basel.
Copyright 2020 info-age GmbH, Basel.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS-IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"flag"
	"github.com/je4/FairService/v2/pkg/fair"
	"github.com/je4/FairService/v2/pkg/fairclient"
	"github.com/je4/FairService/v2/pkg/model/myfair"
	"github.com/je4/utils/v2/pkg/zLogger"
	"github.com/rs/zerolog"
	"io"
	"os"
	"time"
)

func main() {
	cfgfile := flag.String("cfg", "./fairclient.toml", "locations of config file")
	flag.Parse()
	config := LoadConfig(*cfgfile)

	var out io.Writer = os.Stdout

	output := zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339}
	_logger := zerolog.New(output).With().Timestamp().Logger()
	_logger.Level(zLogger.LogLevel(config.Loglevel))
	var logger zLogger.ZLogger = &_logger

	fservice, err := fairclient.NewFairService(
		config.ServiceName,
		config.Address,
		config.CertSkipVerify,
		config.JwtKey,
		config.JwtAlg,
		30*time.Second,
	)
	if err != nil {
		logger.Fatal().Msgf("cannot instantiate fair service: %v", err)
	}
	if err := fservice.Ping(); err != nil {
		logger.Fatal().Msgf("cannot ping fair service: %v", err)
	}

	var item = &fair.ItemData{
		Source:    "github",
		Partition: "mediathek",
		Signature: "je4/FairService",
		Metadata: myfair.Core{
			Identifier: []myfair.Identifier{},
			Person:     []myfair.Person{},
			Title: []myfair.Title{
				{
					Data: "FairService",
					Type: myfair.TitleTypeMain,
				},
			},
			Publisher:       "self",
			PublicationYear: "2021",
			ResourceType:    "code",
			Rights:          "",
			License:         "Apache-2.0",
			Media:           []*myfair.Media{},
			Poster:          nil,
		},
		Set:        []string{},
		Catalog:    []string{},
		Identifier: []string{},
		Access:     "public",
		Status:     "",
		Seq:        0,
		UUID:       "",
		Datestamp:  time.Now(),
		URL:        "",
	}

	newItem, err := fservice.Create(item)
	if err != nil {
		logger.Fatal().Msgf("cannot create item: %v", err)
	}
	logger.Info().Msgf("new item: %v", newItem)

}
