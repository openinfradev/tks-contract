module github.com/sktelecom/tks-contract

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/golang/mock v1.6.0
	github.com/google/uuid v1.2.0
	github.com/jackc/pgproto3/v2 v2.0.7 // indirect
	github.com/lib/pq v1.10.2
	github.com/sirupsen/logrus v1.8.1
	github.com/sktelecom/tks-proto v0.0.6-0.20210622012523-ded9f951101f
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a // indirect
	golang.org/x/net v0.0.0-20210525063256-abc453219eb5 // indirect
	golang.org/x/sys v0.0.0-20210521203332-0cec03c779c1 // indirect
	google.golang.org/genproto v0.0.0-20210524171403-669157292da3 // indirect
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.10
)

replace github.com/sktelecom/tks-contract => ./
