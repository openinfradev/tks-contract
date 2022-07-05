module github.com/openinfradev/tks-contract

go 1.16

require (
	github.com/golang/mock v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/lib/pq v1.10.4
	github.com/openinfradev/tks-common v0.0.0-20220713082320-955751826e06 // indirect
	github.com/openinfradev/tks-proto v0.0.6-0.20220406043255-9fffe49c4625
	github.com/stretchr/testify v1.7.0
	google.golang.org/genproto v0.0.0-20211013025323-ce878158c4d4 // indirect
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/postgres v1.1.2
	gorm.io/gorm v1.21.16
)

replace github.com/openinfradev/tks-contract => ./

//replace github.com/openinfradev/tks-proto => ../tks-proto
//replace github.com/openinfradev/tks-common => ./tks-common
