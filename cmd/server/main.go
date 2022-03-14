package main

import (
	"flag"
	"fmt"

	"github.com/openinfradev/tks-common/pkg/argowf"
	"github.com/openinfradev/tks-common/pkg/grpc_client"
	"github.com/openinfradev/tks-common/pkg/grpc_server"
	"github.com/openinfradev/tks-common/pkg/log"

	"github.com/openinfradev/tks-contract/pkg/contract"
	pb "github.com/openinfradev/tks-proto/tks_pb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type server struct {
	pb.UnimplementedContractServiceServer
}

var (
	argowfClient     argowf.Client
	contractAccessor *contract.Accessor
	cspInfoClient    pb.CspInfoServiceClient
)

var (
	port               int
	tls                bool
	tlsClientCertPath  string
	tlsCertPath        string
	tlsKeyPath         string

	infoServiceAddress string
	infoServicePort    int
	argoAddress        string
	argoPort           int
	dbhost             string
	dbport             string
	dbuser             string
	dbpassword         string
)

func init() {
	flag.IntVar(&port, "port", 9110, "service port")
	flag.BoolVar(&tls, "tls", false, "enabled tls")
	flag.StringVar(&tlsClientCertPath, "tls-client-cert-path", "../../cert/tks-ca.crt", "path of ca cert file for tls")
	flag.StringVar(&tlsCertPath, "tls-cert-path", "../../cert/tks-server.crt", "path of cert file for tls")
	flag.StringVar(&tlsKeyPath, "tls-key-path", "../../cert/tks-server.key", "path of key file for tls")
	flag.StringVar(&infoServiceAddress, "info-address", "localhost", "service address for tks-info")
	flag.IntVar(&infoServicePort, "info-port", 9111, "service port for tks-info")
	flag.StringVar(&argoAddress, "argo-address", "localhost", "service address for argo-workflow")
	flag.IntVar(&argoPort, "argo-port", 2746, "service port for argo-workflow")
	flag.StringVar(&dbhost, "dbhost", "localhost", "host of postgreSQL")
	flag.StringVar(&dbport, "dbport", "5432", "port of postgreSQL")
	flag.StringVar(&dbuser, "dbuser", "postgres", "postgreSQL user")
	flag.StringVar(&dbpassword, "dbpassword", "password", "password for postgreSQL user")
}

func main() {
	flag.Parse()

	log.Info("*** Arguments *** ")
	log.Info("port : ", port)
	log.Info("tls : ", tls)
	log.Info("tlsClientCertPath : ", tlsClientCertPath)
	log.Info("tlsCertPath : ", tlsCertPath)
	log.Info("tlsKeyPath : ", tlsKeyPath)
	log.Info("infoServiceAddress : ", infoServiceAddress)
	log.Info("infoServicePort : ", infoServicePort)
	log.Info("argoAddress : ", argoAddress)
	log.Info("argoPort : ", argoPort)
	log.Info("dbhost : ", dbhost)
	log.Info("dbport : ", dbport)
	log.Info("dbuser : ", dbuser)
	log.Info("dbpassword : ", dbpassword)
	log.Info("****************** ")

	// initialize database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=tks port=%s sslmode=disable TimeZone=Asia/Seoul",
		dbhost, dbuser, dbpassword, dbport)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to open database ", err)
	}
	contractAccessor = contract.New(db)

	// initialize argo client
	_argowfClient, err := argowf.New(argoAddress, argoPort)
	if err != nil {
		log.Fatal("failed to create argowf client : ", err)
	}
	argowfClient = _argowfClient

	// initialize csp_info client
	cc, sc, err := grpc_client.CreateCspInfoClient(infoServiceAddress, infoServicePort, tls, tlsClientCertPath)
	if err != nil {
		log.Fatal("failed to create cspinfo client : ", err)
	}
	defer cc.Close()
	cspInfoClient = sc

	// start server
	s, conn, err := grpc_server.CreateServer(port, tls, tlsCertPath, tlsKeyPath)
	if err != nil {
		log.Fatal("failed to crate grpc_server : ", err)
	}

	// register & serve
	pb.RegisterContractServiceServer(s, &server{})
	if err := s.Serve(conn); err != nil {
		log.Fatal("failed to serve:", err)
	}

}
