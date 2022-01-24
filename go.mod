module github.com/openinfradev/tks-contract

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.3.3
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/lib/pq v1.10.3
	github.com/openinfradev/tks-cluster-lcm v0.0.0-20211102040100-00b4e15d3e4c
	github.com/openinfradev/tks-proto v0.0.6-0.20211015003551-ed8f9541f40d
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0 // indirect
	google.golang.org/grpc v1.41.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/postgres v1.1.2
	gorm.io/gorm v1.21.16
)

replace github.com/openinfradev/tks-contract => ./

//replace github.com/openinfradev/tks-proto => ../tks-proto
replace github.com/openinfradev/tks-cluster-lcm => ../tks-cluster-lcm
