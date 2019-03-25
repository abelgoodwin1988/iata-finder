// Package main implements server for the iatafinder rpc services
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	iatafinder "github.com/abelgoodwin1988/IATA-FINDER/deployments/iatafinder_rpc"
	"google.golang.org/grpc"
)

type server struct{}

func main() {
	fmt.Printf("iatafinder service started.\n")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Listener Failed: %v", err)
	}

	s := grpc.NewServer()

	iatafinder.RegisterIatafinderServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start iatafinder server\n%v\n", err)
	}
}

func (*server) GetAirport(ctx context.Context, in *iatafinder.AirportDescriptor) (*iatafinder.Airports, error) {
	return &iatafinder.Airports{}, nil
}

func (*server) GetAirportIATA(ctx context.Context, in *iatafinder.IATA) (*iatafinder.Airport, error) {
	return &iatafinder.Airport{}, nil
}

func (*server) GetAirportICAO(ct context.Context, in *iatafinder.ICAO) (*iatafinder.Airport, error) {
	return &iatafinder.Airport{}, nil
}

func (*server) GetAirports(ctx context.Context, in *iatafinder.EmptyRequest) (*iatafinder.Airports, error) {
	return &iatafinder.Airports{}, nil
}
