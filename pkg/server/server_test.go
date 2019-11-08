package server

import (
	"abelgoodwin1988/iata-finder/pkg/dataservice"
	iatafinder "abelgoodwin1988/iata-finder/rpc"
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func TestServer(t *testing.T) {
	// SETUP
	pwd, _ := os.Getwd()

	ds := dataservice.Dataservice{
		URLTargets: []string{
			"https://raw.githubusercontent.com/jpatokal/openflights/master/data/airports.dat",
			"https://raw.githubusercontent.com/jpatokal/openflights/master/data/airlines.dat",
		},
		DataDestination: fmt.Sprintf("%s/assets", pwd),
		FileType:        ".csv",
		Interval:        time.Hour * 24,
	}

	s, err := Create(&ds, "../../configs/rpc.config.test.toml")
	if err != nil {
		t.Errorf("failed to create gRPC server with provided DataService: %v", err)
	}

	// create a client connection for all tests
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		t.Errorf("could not connect: %v", err)
	}
	cc.Close()

	c := iatafinder.NewIatafinderClient(cc)
	// SUBTESTS
	t.Run("GetAirportIATA", func(t *testing.T) {
		req := &iatafinder.IATA{Iata: "ONT"}
		res, err := c.GetAirportIATA(context.Background(), req)
		expected := "Ontario"

		if err != nil {
			t.Errorf("error get airport from IATA: %v", req)
		}

		if res.City != expected {
			t.Errorf("GetAirportIATA(%v)=%v. Expecting: %v", req, res.City, expected)
		}
	})

	// TEARDOWN
	s.Stop()
}
