module github.com/je4/FairService/v2

go 1.21

replace github.com/je4/FairService/v2 => ./

//replace github.com/je4/utils/v2 => ../utils/

require (
	github.com/BurntSushi/toml v1.3.2
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de
	github.com/bluele/gcache v0.0.2
	github.com/google/uuid v1.6.0
	github.com/gorilla/handlers v1.5.2
	github.com/gorilla/mux v1.8.1
	github.com/je4/HandleCreator/v2 v2.0.5
	github.com/je4/utils/v2 v2.0.23
	github.com/lib/pq v1.10.9
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.32.0
)

require (
	emperror.dev/errors v0.8.1 // indirect
	github.com/blend/go-sdk v1.20220411.3 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/machinebox/progress v0.2.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/pkg/sftp v1.13.6 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.20.0 // indirect
	golang.org/x/sys v0.17.0 // indirect
)
