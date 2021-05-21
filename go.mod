module github.com/sktelecom/tks-contract

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/google/uuid v1.2.0
	github.com/lib/pq v1.10.2
	github.com/sirupsen/logrus v1.8.1
	github.com/sktelecom/tks-proto v0.0.4-0.20210419072147-cbafa000deab
	github.com/stretchr/testify v1.7.0 // indirect
	google.golang.org/grpc v1.36.1
	google.golang.org/protobuf v1.26.0
)

replace github.com/sktelecom/tks-contract => ./
