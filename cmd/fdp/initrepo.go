package main

import (
	"github.com/je4/FairService/v2/pkg/re3data"
	"time"
)

func initRepo() re3data.Repository {
	repo := re3data.Repository{
		Identifiers: re3data.RepositoryIdentifiers{
			Re3Data: "",
			DOI:     "",
		},
		Name: re3data.StringLang{
			Value: "Integrated Catalog Mediathek HGK FHNW",
		},
		URL: "https://mediathek.hgk.fhnw.ch/amp",
		Type: []re3data.RepositoryType{
			"interdisciplinary",
		},
		Updated:  time.Now().Format("2006/01/02"),
		Language: []string{"eng"},
		Subject: []re3data.RepositorySubject{
			re3data.RepositorySubject{
				Scheme: "",
				Id:     "",
				Name:   "",
			},
		},
		ProviderType:          nil,
		Institution:           nil,
		DataAccess:            re3data.RepositoryDatabaseAccess{},
		DataLicense:           re3data.RepositoryDataLicense{},
		DataUpload:            nil,
		Versioning:            "",
		EnhancedPublication:   "",
		QualityManagement:     "",
		EntryDate:             "",
		LastUpdate:            "",
		AdditionalName:        re3data.StringLang{},
		RepositoryIdentifiers: nil,
		Description:           re3data.StringLang{},
		Contact:               nil,
		Keyword:               nil,
	}
	repo.InitNamespace()
	return repo
}
