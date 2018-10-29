package main

import (
	pb "github.com/caquillo07/grpc-demo-shipping-containers/consignment-service/proto/consignment"

	"fmt"
	vesselProto "github.com/caquillo07/grpc-demo-shipping-containers/vessel-service/proto/vessel"
	"github.com/micro/go-micro"
	"log"
	"os"
)

const (
	defaultHost = "localhost:27017"
)

func main() {
	// Database host from the environment variables
	host := os.Getenv("DB_HOST")

	if host == "" {
		host = defaultHost
	}

	session, err := CreateSession(host)

	// Mgo creates a 'master' session, we need to end that session
	// before the main function closes.
	defer session.Close()

	if err != nil {
		log.Panicf("could not connect to the datastore with hose %s - %v\n", host, err)
	}

	// create a new service. Optionally include the options here.
	srv := micro.NewService(
		micro.Name("go.micro.srv.consignment"),
		micro.Version("latest"),
	)

	vesselClient := vesselProto.NewVesselServiceClient("go.micro.srv.vessel", srv.Client())

	// init will parse the command line flags.
	srv.Init()

	// register handler
	pb.RegisterShippingServiceHandler(srv.Server(), &handler{session, vesselClient})

	// run server!
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
