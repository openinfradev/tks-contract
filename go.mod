module tks-contract

go 1.16

require (
	github.com/google/uuid v1.2.0
	github.com/openinfradev/tks-contract v0.0.0
	github.com/openinfradev/tks-proto v0.0.3
	google.golang.org/grpc v1.36.1
	google.golang.org/protobuf v1.26.0
)

replace github.com/openinfradev/tks-contract => ./
