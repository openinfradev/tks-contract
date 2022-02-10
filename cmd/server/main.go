package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/openinfradev/tks-common/pkg/grpc_client"
	"github.com/openinfradev/tks-common/pkg/log"
	"github.com/openinfradev/tks-contract/pkg/contract"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	port               int    = 9110
	infoServiceAddress string = ""
	infoServicePort    int    = 9111
	argoAddress        string = ""
	argoPort           int    = 2746
	dbhost             string = "localhost"
	dbport             string = "5432"
	dbuser             string = "postgres"
	dbpassword         string = "password"
)

type server struct {
	pb.UnimplementedContractServiceServer
}

func init() {
	setFlags()
}

func setFlags() {
	flag.IntVar(&port, "port", 9110, "service port")
	flag.StringVar(&infoServiceAddress, "info-address", "", "service address for tks-info")
	flag.IntVar(&infoServicePort, "info-port", 9111, "service port for tks-info")
	flag.StringVar(&argoAddress, "argo-address", "", "service address for argo-workflow")
	flag.IntVar(&argoPort, "argo-port", 2746, "service port for argo-workflow")
	flag.StringVar(&dbhost, "dbhost", "localhost", "host of postgreSQL")
	flag.StringVar(&dbport, "dbport", "5432", "port of postgreSQL")
	flag.StringVar(&dbuser, "dbuser", "postgres", "postgreSQL user")
	flag.StringVar(&dbpassword, "dbpassword", "password", "password for postgreSQL user")
}

func main() {
	flag.Parse()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=tks port=%s sslmode=disable TimeZone=Asia/Seoul",
		dbhost, dbuser, dbpassword, dbport)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("an error occurred ", err)
	}
	contractAccessor = contract.New(db)

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	log.Info("Starting to listen port ", port)
	if err != nil {
		log.Fatal("failed to listen:", err)
	}
	cc, sc, err := grpc_client.CreateCspInfoClient(infoServiceAddress, infoServicePort, "tks-contract")
	if err != nil {
		log.Error()
	}
	defer cc.Close()
	cspInfoClient = sc

	s := grpc.NewServer()

	log.Info("Started to listen port ", port)
	log.Info("****************************")

	InitHandlers(argoAddress, argoPort)

	pb.RegisterContractServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatal("failed to serve:", err)
	}
}
