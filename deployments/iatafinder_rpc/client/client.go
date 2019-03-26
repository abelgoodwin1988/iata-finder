// Package main implements iatafinder rpc client
package main

import (
	"context"
	"fmt"
	"log"

	iatafinder "github.com/abelgoodwin1988/IATA-FINDER/deployments/iatafinder_rpc"

	"google.golang.org/grpc"
)

func main() {
	fmt.Printf("Entered iatafinder client\n")

	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Error creating client connection to rpc server:\n%v\n", err)
	}
	defer cc.Close()

	c := iatafinder.NewIatafinderClient(cc)

	request := &iatafinder.IATA{
		Iata: "JFK",
	}

	response, err := c.GetAirportIATA(context.Background(), request)
	if err != nil {
		log.Fatalf("Error making request for airport by IATA:\n%v\n", err)
	}

	fmt.Printf("Success: %v\n", response.GetName())
}
