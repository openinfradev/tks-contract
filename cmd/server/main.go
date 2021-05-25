package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"

	"github.com/sktelecom/tks-contract/pkg/contract"
	gcInfo "github.com/sktelecom/tks-contract/pkg/grpc-client"
	"github.com/sktelecom/tks-contract/pkg/log"
	pb "github.com/sktelecom/tks-proto/pbgo"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	port               int    = 9110
	infoServiceAddress string = ""
	infoServicePort    int    = 9111
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
	cc, sc, err := gcInfo.CreateClientsObject(infoServiceAddress, infoServicePort, false, "")
	infoClient = gcInfo.New(cc, sc)
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
