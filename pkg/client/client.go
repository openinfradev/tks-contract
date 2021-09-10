package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	//"google.golang.org/grpc/credentials"

	"github.com/openinfradev/tks-contract/pkg/log"
	pb "github.com/openinfradev/tks-proto/pbgo"
)

var (
	conn *grpc.ClientConn
	contractClient pb.ContractServiceClient
)

func RequestLogging() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		end := time.Now()

		log.Info(fmt.Sprintf("[GRPC:%s][START:%s][END:%s][ERR:%v]", method, start.Format(time.RFC3339), end.Format(time.RFC3339), err))
		log.Debug(fmt.Sprintf("[GRPC:%s][REQUEST %s][REPLY %s]", method, req, reply))
		
		return err
	}
}

func GetConnection(host string) (*grpc.ClientConn, error) {
	if conn == nil {
		_conn, err := grpc.Dial(
			host,
			grpc.WithInsecure(),
			grpc.WithUnaryInterceptor(
				grpc_middleware.ChainUnaryClient(
					RequestLogging(),
				),
			),
		)
		if err != nil {
			return nil, err
		}
		conn = _conn
	} 
	return conn, nil
}


func GetContractClient(address string, port int, caller string) (pb.ContractServiceClient, error) {
	conn, err := GetConnection( fmt.Sprintf("%s:%d", address, port) )
	if err != nil {
		return nil, err
	} 

	contractClient = pb.NewContractServiceClient(conn)
	return contractClient, nil
}

