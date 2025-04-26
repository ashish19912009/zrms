module github.com/ashish19912009/zrms/services/authZ

go 1.23.6

toolchain go1.23.8

require (
	github.com/rs/zerolog v1.34.0
	google.golang.org/grpc v1.71.0
	google.golang.org/protobuf v1.36.5
	honnef.co/go/tools v0.6.1
)

require (
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	go.opentelemetry.io/otel v1.35.0 // indirect
	go.opentelemetry.io/otel/sdk v1.35.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	golang.org/x/tools v0.30.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250218202821-56aae31c358a // indirect
)

replace github.com/ashish19912009/services/authZ => ./ // 👈 tell Go "use local path"
