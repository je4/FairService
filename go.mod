module github.com/je4/FairService/v2

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/BurntSushi/toml v0.4.1
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/bluele/gcache v0.0.2
	github.com/google/uuid v1.3.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/je4/HandleCreator/v2 v2.0.2
	github.com/je4/utils/v2 v2.0.2
	github.com/je4/zsearch/v2 v2.0.0-20211030141626-5b0d9fcfe7e2
	github.com/lib/pq v1.10.3
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/pkg/errors v0.9.1
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20211102192858-4dd72447c267 // indirect
)

replace github.com/je4/FairService/v2 => ./
