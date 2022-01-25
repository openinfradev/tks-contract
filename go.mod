module github.com/openinfradev/tks-contract

go 1.16

require (
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/lib/pq v1.10.4
	github.com/openinfradev/tks-common v0.0.0-00010101000000-000000000000
	github.com/openinfradev/tks-proto v0.0.6-0.20211015003551-ed8f9541f40d
	github.com/stretchr/testify v1.7.0
	golang.org/x/sys v0.0.0-20211013075003-97ac67df715c // indirect
	google.golang.org/genproto v0.0.0-20211013025323-ce878158c4d4 // indirect
	google.golang.org/grpc v1.43.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/postgres v1.1.2
	gorm.io/gorm v1.21.16
)

replace github.com/openinfradev/tks-contract => ./

//replace github.com/openinfradev/tks-proto => ../tks-proto
//replace github.com/openinfradev/tks-cluster-lcm => ../tks-cluster-lcm
//replace github.com/openinfradev/tks-common => ../../tks-common
