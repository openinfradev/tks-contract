package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/sktelecom/tks-contract/pkg/log"
	pb "github.com/sktelecom/tks-proto/pbgo"
)

type InfoClient struct {
	cc *grpc.ClientConn
	sc pb.InfoServiceClient
}

func NewInfoClient(address string, port int, tls bool, caFile string) (*InfoClient, error) {
	opts := grpc.WithInsecure()
	if tls {
		if caFile == "" {
			log.Fatal("CA file is null. CA file must be exist when tls is on.")
			return &InfoClient{}, fmt.Errorf("CA file not found.")
		} else {
			creds, err := credentials.NewServerTLSFromFile(caFile, "")
			if err != nil {
				log.Fatal("Error while loading CA trust certificate: ", err.Error())
				return &InfoClient{}, err
			}
			opts = grpc.WithTransportCredentials(creds)
		}
	}
	address = fmt.Sprintf("%s:%d", address, port)
	log.Info("gRPC server address is ", address)
	cc, err := grpc.Dial(address, opts)
	if err != nil {
		log.Fatal("Could not connect to gRPC server", err)
		return &InfoClient{}, err
	}
	sc := pb.NewInfoServiceClient(cc)
	return &InfoClient{
		cc: cc,
		sc: sc,
	}, nil
}

func (c *InfoClient) CreateCSPInfo(ctx context.Context, contractId string,
	cspName string, auth string) (*pb.IDResponse, error) {
	return c.sc.CreateCSPInfo(ctx, &pb.CreateCSPInfoRequest{
		ContractId: contractId,
		CspName:    cspName,
		Auth:       auth,
	})
}

func (c *InfoClient) Close() error {
	return c.cc.Close()
}
