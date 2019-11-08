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

	// instantiate data source
	ds.DataCollector()
	ds.ParseHandler()

	// create server which can use custom config
	lis, s := Create(&ds, "../../configs/rpc.config.test.toml")

	// Start the rpc server in goroutine so tests can run
	go func() {
		if err := s.Serve(lis); err != nil {
			ctxLogger.Fatalf("Failed to start iatafinder server\n%v\n", err)
		}
	}()
	// create a client connection for all tests
	cc, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		ctxLogger.Errorf("could not connect: %v\n", err)
	}

	c := iatafinder.NewIatafinderClient(cc)

	// SUBTESTS
	t.Run("GetAirportIATA", func(t *testing.T) {
		req := &iatafinder.IATA{Iata: "ONT"}
		res, err := c.GetAirportIATA(context.Background(), req)
		expected := "Ontario"

		if err != nil {
			t.Errorf("error retrieving airport from IATA: %v", req)
			return
		}

		if res.City != expected {
			t.Errorf("GetAirportIATA(%v)=%v. Expecting: %v", req, res.City, expected)
		}
	})

	// TEARDOWN
	cc.Close()
}
