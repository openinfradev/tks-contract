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
	port         int  = 5000
	enableMockup bool = false
)

type server struct {
	pb.UnimplementedContractServiceServer
}

func init() {
	setFlags()

	contractAccessor = contract.NewContractAccessor()
}

func setFlags() {
	flag.IntVar(&port, "port", 50051, "service port")
	flag.BoolVar(&enableMockup, "enable-mockup", false, "enable mockup contracts")
}

func main() {
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	log.Info("Starting to listen port ", port)
	flag.Parse()
	if err != nil {
		log.Fatal("failed to listen:", err)
	}
	if enableMockup {
		if err := InsertMockupContracts(contractAccessor); err != nil {
			log.Warn("failed to create mockup data:", err)
		}
	}
	s := grpc.NewServer()
	pb.RegisterContractServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
