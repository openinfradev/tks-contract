package main

import (
	"flag"
	"net"
	"strconv"

	"github.com/sktelecom/tks-contract/pkg/contract"
	gcInfo "github.com/sktelecom/tks-contract/pkg/grpc-client"
	"github.com/sktelecom/tks-contract/pkg/log"
	pb "github.com/sktelecom/tks-proto/pbgo"
	"google.golang.org/grpc"
)

var (
	port               int    = 9110
	enableMockup       bool   = false
	infoServiceAddress string = ""
	infoServicePort    int    = 9111
)

type server struct {
	pb.UnimplementedContractServiceServer
}

func init() {
	setFlags()

	contractAccessor = contract.NewContractAccessor()
}

func setFlags() {
	flag.IntVar(&port, "port", 9110, "service port")
	flag.BoolVar(&enableMockup, "enable-mockup", false, "enable mockup contracts")
	flag.StringVar(&infoServiceAddress, "info-address", "", "service address for tks-info")
	flag.IntVar(&infoServicePort, "info-port", 9111, "service port for tks-info")
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

	infoClient, err = gcInfo.NewInfoClient(infoServiceAddress, infoServicePort, false, "")
	if err != nil {
		log.Error()
	}
	defer infoClient.Close()

	s := grpc.NewServer()
	pb.RegisterContractServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
