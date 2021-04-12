module github.com/sktelecom/tks-contract

go 1.16

require (
	github.com/google/uuid v1.2.0
	github.com/sirupsen/logrus v1.8.1
	github.com/sktelecom/tks-proto v0.0.2
	google.golang.org/grpc v1.36.1
	google.golang.org/protobuf v1.26.0
)

replace github.com/sktelecom/tks-contract => ./
