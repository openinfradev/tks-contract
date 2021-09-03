package client

import (
	"sync"
	"fmt"

	"google.golang.org/grpc"
	//"google.golang.org/grpc/credentials"

	//"github.com/sktelecom/tks-contract/pkg/log"
	pb "github.com/openinfradev/tks-proto/pbgo"
)

var (
	once sync.Once
	contractClient pb.ContractServiceClient
)

func GetContractClient(address string, port int, caller string) pb.ContractServiceClient {
	host := fmt.Sprintf("%s:%d", address, port)
	once.Do(func() {
		conn, _ := grpc.Dial(
			host,
			grpc.WithInsecure(),
		)

		contractClient = pb.NewContractServiceClient(conn)
	})
	return contractClient
}

