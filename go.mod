module github.com/je4/FairService/v2

go 1.16

require (
	emperror.dev/errors v0.8.0
	github.com/BurntSushi/toml v0.3.1
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/bluele/gcache v0.0.2
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/google/uuid v1.1.2
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/je4/utils/v2 v2.0.0-20210702125424-8c1cdd3f1ccc
	github.com/lib/pq v1.10.2
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/pkg/errors v0.9.1
)

replace github.com/je4/FairService/v2 => ./
