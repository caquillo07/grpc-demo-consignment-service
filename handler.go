package main

import (
	pb "github.com/caquillo07/grpc-demo-shipping-containers/consignment-service/proto/consignment"
	vesselProto "github.com/caquillo07/grpc-demo-shipping-containers/vessel-service/proto/vessel"
	"golang.org/x/net/context"
	"gopkg.in/mgo.v2"
	"log"
)

type handler struct {
	session      *mgo.Session
	vesselClient vesselProto.VesselServiceClient
}

func (s *handler) GetRepo() Repository {
	return &ConsignmentRepository{s.session.Clone()}
}

func (s *handler) CreateConsignment(ctx context.Context, req *pb.Consignment, res *pb.Response) error {
	repo := s.GetRepo()
	defer repo.Close()

	// Here we call a client instance of our vessel handler with our consignment weight,
	// and the amount of containers as the capacity value
	vesselResponse, err := s.vesselClient.FindAvailable(context.Background(), &vesselProto.Specification{
		MaxWeight: req.Weight,
		Capacity:  int32(len(req.Containers)),
	})
	if err != nil {
		return err
	}

	log.Printf("Found vessel: %s \n", vesselResponse.Vessel.Name)

	// We set the VesselId as the vessel we got back from our
	// vessel handler
	req.VesselId = vesselResponse.Vessel.Id

	// Save our consignment
	err = repo.Create(req)
	if err != nil {
		return err
	}

	// Return matching the `Response` message we created in our
	// protobuf definition.
	res.Created = true
	res.Consignment = req
	return nil
}

func (s *handler) GetConsignments(ctx context.Context, req *pb.GetRequest, res *pb.Response) error {
	repo := s.GetRepo()
	defer repo.Close()
	consignments, err := repo.GetAll()

	if err != nil {
		return err
	}

	res.Consignments = consignments
	return nil
}
