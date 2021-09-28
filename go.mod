module github.com/openinfradev/tks-contract

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/jackc/pgx/v4 v4.13.0 // indirect
	github.com/lib/pq v1.10.3
	github.com/openinfradev/tks-cluster-lcm v0.0.0-20210908061731-769e64f93271
	github.com/openinfradev/tks-proto v0.0.6-0.20210901093202-5e0db3fa3d4f
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210915214749-c084706c2272 // indirect
	golang.org/x/net v0.0.0-20210916014120-12bc252f5db8 // indirect
	golang.org/x/sys v0.0.0-20210915083310-ed5796bab164 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20210916144049-3192f974c780 // indirect
	google.golang.org/grpc v1.40.0
	google.golang.org/protobuf v1.27.1
	gorm.io/driver/postgres v1.1.1
	gorm.io/gorm v1.21.15
)

replace github.com/openinfradev/tks-contract => ./

replace github.com/openinfradev/tks-cluster-lcm => ../tks-cluster-lcm
