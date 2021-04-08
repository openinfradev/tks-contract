package main

import (
	"net"

	"github.com/openinfradev/tks-contract/pkg/contract"
	"github.com/openinfradev/tks-contract/pkg/log"
	pb "github.com/openinfradev/tks-proto/pbgo"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedContractServiceServer
}

func init() {
	log.Initialize("tks-contract")
	contractAccessor = contract.NewContractAccessor()
	InsertMockupContracts(contractAccessor)
}

func main() {
	lis, err := net.Listen("tcp", port)
	log.Println("Starting to listen port ", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterContractServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
