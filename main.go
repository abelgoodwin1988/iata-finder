// Package main implements server for the iatafinder rpc services
package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"time"

	configmodels "abelgoodwin1988/iata-finder/configs/models"
	"abelgoodwin1988/iata-finder/pkg/dataservice"
	"abelgoodwin1988/iata-finder/pkg/logger"
	iatafinder "abelgoodwin1988/iata-finder/rpc"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var ctxLogger = logger.CtxLogger
var ds dataservice.Dataservice

type server struct{}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	ctxLogger.Info("Starting iata-finder data service")
	pwd, _ := os.Getwd()
	ds = dataservice.Dataservice{
		URLTargets: []string{
			"https://raw.githubusercontent.com/jpatokal/openflights/master/data/airports.dat",
			"https://raw.githubusercontent.com/jpatokal/openflights/master/data/airlines.dat",
		},
		DataDestination: fmt.Sprintf("%s/assets", pwd),
		FileType:        ".csv",
		Interval:        time.Hour * 24,
	}
	ds.Init(&wg)
	wg.Wait()

	ctxLogger.Info("Starting iata-finder service")
	// load config values for rpc
	var rpcConfig configmodels.RPCConfig
	loadRPCConfig(&rpcConfig)

	// Configure and create rpc server
	lis, s := rpcListenAndServe(&rpcConfig)

	iatafinder.RegisterIatafinderServer(s, &server{})

	ctxLogger.WithFields(logrus.Fields{"at": fmt.Sprintf("%s:%s", rpcConfig.IP, rpcConfig.Port)}).Info("Serving iata-finder service")
	// Start the rpc server and if it fails, log it and give up all hope
	if err := s.Serve(lis); err != nil {
		ctxLogger.Fatalf("Failed to start iatafinder server\n%v\n", err)
	}
}

func loadRPCConfig(rpcConfig *configmodels.RPCConfig) {
	data, err := ioutil.ReadFile("configs/rpc.config.toml")
	if err != nil {
		ctxLogger.WithError(err).Error("Failed to read rpc.config.go")
	}
	if _, err := toml.Decode(string(data), &rpcConfig); err != nil {
		ctxLogger.WithError(err).Error("Failed to decode rpc.config.go")
	}
}

func rpcListenAndServe(rpcConfig *configmodels.RPCConfig) (net.Listener, *grpc.Server) {
	network := rpcConfig.Network
	ip := rpcConfig.IP
	port := rpcConfig.Port
	address := fmt.Sprintf("%s:%s", ip, port)
	lis, err := net.Listen(network, address)
	if err != nil {
		ctxLogger.Fatalf("Listener Failed: %v", err)
	}

	return lis, grpc.NewServer()
}

func (*server) GetAirport(ctx context.Context, in *iatafinder.AirportDescriptor) (*iatafinder.Airports, error) {
	return &iatafinder.Airports{}, nil
}

func (*server) GetAirportIATA(ctx context.Context, in *iatafinder.IATA) (*iatafinder.Airport, error) {
	iata := in.GetIata()
	for _, airport := range ds.Data.Airports {
		if airport.Iata == iata {
			ctxLogger.WithFields(logrus.Fields{
				"Method": "GetAirportIata",
				"Found":  true,
			}).Debugf("\nFound Airport by IATA code:\n %s\n", airport)
			return &iatafinder.Airport{
				Id:                  int32(airport.ID),
				Name:                airport.Name,
				City:                airport.City,
				Country:             airport.Country,
				Iata:                airport.Iata,
				Icao:                airport.Icao,
				Latitude:            airport.Latitude,
				Longitude:           airport.Longitude,
				Altitude:            airport.Altitude,
				Timezone:            airport.Tz,
				DaylightSavingsTime: airport.DaylightSavingsTime,
				Type:                airport.TypeField,
				Source:              airport.Source,
			}, nil
		}
	}
	ctxLogger.Errorf("Failed to find %s in dataset for IATA's", iata)
	return nil, fmt.Errorf("Failed to find %s in dataset for IATA's", iata)
}

func (*server) GetAirportICAO(ct context.Context, in *iatafinder.ICAO) (*iatafinder.Airport, error) {
	return &iatafinder.Airport{}, nil
}

func (*server) GetAirports(ctx context.Context, in *iatafinder.EmptyRequest) (*iatafinder.Airports, error) {
	return &iatafinder.Airports{}, nil
}
