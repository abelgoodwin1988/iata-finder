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
}

func (*server) GetAirport(ctx context.Context, in *iatafinder.AirportDescriptor) (*iatafinder.Airports, error) {
	return &iatafinder.Airports{}, nil
}

func (*server) GetAirportIATA(ctx context.Context, in *iatafinder.IATA) (*iatafinder.Airport, error) {
	return &iatafinder.Airport{}, nil
}
