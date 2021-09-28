module github.com/openinfradev/tks-contract

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.1-0.20190118093823-f849b5445de4
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/lib/pq v1.10.3
	github.com/openinfradev/tks-proto v0.0.6-0.20210924020717-178698d59e9d
	github.com/sirupsen/logrus v1.8.1
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/net v0.0.0-20210908191846-a5e095526f91 // indirect
	golang.org/x/sys v0.0.0-20210909193231-528a39cd75f3 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210909211513-a8c4777a87af // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.15
)

replace github.com/openinfradev/tks-contract => ./
