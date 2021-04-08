package main

import (
	"flag"
	"net"
	"strconv"

	"github.com/openinfradev/tks-contract/pkg/contract"
	"github.com/openinfradev/tks-contract/pkg/log"
	pb "github.com/openinfradev/tks-proto/pbgo"
	"google.golang.org/grpc"
)

var (
	port         int
	enableMockup bool
)

type server struct {
	pb.UnimplementedContractServiceServer
}

func init() {
	getFlags()

	contractAccessor = contract.NewContractAccessor()
	if enableMockup {
		InsertMockupContracts(contractAccessor)
	}
}

func getFlags() {
	flag.IntVar(&port, "port", 50051, "service port")
	flag.BoolVar(&enableMockup, "enable-mockup", false, "enable mockup contracts")
	flag.Parse()
}

func main() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	log.Info("Starting to listen port ", port)
	if err != nil {
		log.Fatal("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterContractServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve: %v", err)
	}
}
