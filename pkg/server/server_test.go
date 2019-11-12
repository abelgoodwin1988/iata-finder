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

var c iatafinder.IatafinderClient

func TestServer(t *testing.T) {
	/*
	* SETUP - use single client connection and data source for all tests
	 */
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

	c = iatafinder.NewIatafinderClient(cc)

	/*
	* SUBTESTS
	 */

	t.Run("GetAirportIATA", getAirportIATA)
	t.Run("GetAirportICAO", getAirportICAO)
	// t.Run("GetAirports", getAirports)
	// t.Run("GetAllAirports", getAllAirports)

	t.Run("GetAirlineIATA", getAirlineIATA)
	t.Run("GetAirlineICAO", getAirlineICAO)
	// t.Run("GetAirlines", getAirlines)
	// t.Run("GetAllAirlines", getAllAirlines)

	/*
	* TEARDOWN
	 */
	cc.Close()
	s.GracefulStop()
}

/*
** Individual testing function definitions
 */

func getAirportIATA(t *testing.T) {
	tests := []struct {
		iataIn string
		id     int32
		name   string
	}{
		{"ONT", 3734, "Ontario International Airport"},
		{"FRA", 340, "Frankfurt am Main Airport"},
		{"ICN", 3930, "Incheon International Airport"},
	}

	for _, test := range tests {
		req := &iatafinder.IATA{Iata: test.iataIn}
		res, err := c.GetAirportIATA(context.Background(), req)

		if err != nil {
			t.Errorf("error retrieving airport from IATA: %v\n", req)
			return
		}

		if res.Id != test.id {
			t.Errorf("GetAirportIATA(%v) - Expecting Id: %v / Got id = %v\n", req, test.id, res.Id)
		}

		if res.Name != test.name {
			t.Errorf("GetAirportIATA(%v) - Expecting Name: %v / Got Name = %v\n", req, test.name, res.Name)
		}
	}
}

func getAirportICAO(t *testing.T) {
	tests := []struct {
		icaoIn string
		id     int32
		name   string
	}{
		{"KONT", 3734, "Ontario International Airport"},
		{"EDDF", 340, "Frankfurt am Main Airport"},
		{"RKSI", 3930, "Incheon International Airport"},
	}

	for _, test := range tests {
		req := &iatafinder.ICAO{Icao: test.icaoIn}
		res, err := c.GetAirportICAO(context.Background(), req)

		if err != nil {
			t.Errorf("error retrieving airport from ICAO: %v\n", req)
			return
		}

		if res.Id != test.id {
			t.Errorf("GetAirportICAO(%v) - Expecting Id: %v / Got id = %v\n", req, test.id, res.Id)
		}

		if res.Name != test.name {
			t.Errorf("GetAirportICAO(%v) - Expecting Name: %v / Got Name = %v\n", req, test.name, res.Name)
		}
	}
}

func getAirports(t *testing.T) {
	tests := []struct {
		descriptor string
		length     int
	}{
		{"Budapest", 2},
		{"Penang", 1},
		{"Buenos Aires", 3},
	}

	for _, test := range tests {
		req := &iatafinder.AirportDescriptor{Descriptor_: test.descriptor}
		res, err := c.GetAirports(context.Background(), req)

		if err != nil {
			t.Errorf("error retrieving airports for descriptor: %v\n", req)
			return
		}

		if len(res.Airports) != test.length {
			t.Errorf("GetAirports(%v) - Expecting length: %v / Got length: %v", req, test.length, len(res.Airports))
		}
	}
}

func getAirlineIATA(t *testing.T) {
	tests := []struct {
		iataIn string
		id     int32
		name   string
	}{
		{"SQ", 4435, "Singapore Airlines"},
		{"EY", 2222, "Etihad Airways"},
		{"IB", 2822, "Iberia Airlines"},
	}

	for _, test := range tests {
		req := &iatafinder.IATA{Iata: test.iataIn}
		res, err := c.GetAirlineIATA(context.Background(), req)

		if err != nil {
			t.Errorf("error retrieving airline from IATA: %v\n", req)
			return
		}

		if res.ID != test.id {
			t.Errorf("GetAirlineIATA(%v) - Expecting Id: %v / Got id = %v\n", req, test.id, res.ID)
		}

		if res.Name != test.name {
			t.Errorf("GetAirlineIATA(%v) - Expecting Name: %v / Got Name = %v\n", req, test.name, res.Name)
		}
	}
}

func getAirlineICAO(t *testing.T) {
	tests := []struct {
		icaoIn string
		id     int32
		name   string
	}{
		{"SIA", 4435, "Singapore Airlines"},
		{"ETD", 2222, "Etihad Airways"},
		{"IBE", 2822, "Iberia Airlines"},
	}

	for _, test := range tests {
		req := &iatafinder.ICAO{Icao: test.icaoIn}
		res, err := c.GetAirlineICAO(context.Background(), req)

		if err != nil {
			t.Errorf("error retrieving airline from IATA: %v\n", req)
			return
		}

		if res.ID != test.id {
			t.Errorf("GetAirlineIATA(%v) - Expecting Id: %v / Got id = %v\n", req, test.id, res.ID)
		}

		if res.Name != test.name {
			t.Errorf("GetAirlineIATA(%v) - Expecting Name: %v / Got Name = %v\n", req, test.name, res.Name)
		}
	}
}
