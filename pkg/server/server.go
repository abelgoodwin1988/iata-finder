package server

import (
	configmodels "abelgoodwin1988/iata-finder/configs/models"
	"abelgoodwin1988/iata-finder/pkg/dataservice"
	"abelgoodwin1988/iata-finder/pkg/logger"
	iatafinder "abelgoodwin1988/iata-finder/rpc"
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var ctxLogger = logger.CtxLogger
var ds *dataservice.Dataservice

type server struct{}

// Create initializes the rpc server
func Create(dataservice *dataservice.Dataservice, configPath string) (net.Listener, *grpc.Server) {
	ctxLogger.Info("Starting iata-finder service")

	if dataservice == nil {
		// Will warn for now
		ctxLogger.Warningln("No Dataservice provided")
	}

	ds = dataservice

	// load config values for rpc
	rpcConfig := configmodels.RPCConfig{}
	loadRPCConfig(&rpcConfig, configPath)

	// Configure and create rpc server
	lis := rpcListen(&rpcConfig)

	s := grpc.NewServer()
	iatafinder.RegisterIatafinderServer(s, &server{})

	ctxLogger.WithFields(logrus.Fields{"at": fmt.Sprintf("%s:%s", rpcConfig.IP, rpcConfig.Port)}).Info("iata-finder service configured to be served")

	// return listener and registered server. We will serve outside of server
	// since we may need to run server inside a go routine in test case
	return lis, s
}

func loadRPCConfig(rpcConfig *configmodels.RPCConfig, configPath string) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		ctxLogger.WithError(err).Error("Failed to read rpc.config.go")
	}
	if _, err := toml.Decode(string(data), &rpcConfig); err != nil {
		ctxLogger.WithError(err).Error("Failed to decode rpc.config.go")
	}
}

func rpcListen(rpcConfig *configmodels.RPCConfig) net.Listener {
	network := rpcConfig.Network
	ip := rpcConfig.IP
	port := rpcConfig.Port
	address := fmt.Sprintf("%s:%s", ip, port)
	lis, err := net.Listen(network, address)
	if err != nil {
		ctxLogger.Fatalf("Listener Failed: %v", err)
	}

	return lis
}

func (*server) GetAirports(ctx context.Context, in *iatafinder.AirportDescriptor) (*iatafinder.Airports, error) {
	descriptor := in.GetDescriptor_()
	airports := []*iatafinder.Airport{}
	for _, airport := range ds.Data.Airports.Airports {
		// name, city, or country partial matches
		if strings.Contains(airport.GetCity(), descriptor) ||
			strings.Contains(airport.GetCountry(), descriptor) ||
			strings.Contains(airport.GetName(), descriptor) {
			airports = append(airports, airport)
		}
	}
	ctxLogger.WithFields(logrus.Fields{
		"Method": "GetAirports",
		"Found":  len(airports) > 0,
	}).Debug()
	airportss := &iatafinder.Airports{Airports: airports}
	return airportss, nil
}

func (*server) GetAirportIATA(ctx context.Context, in *iatafinder.IATA) (*iatafinder.Airport, error) {
	iata := in.GetIata()
	for _, airport := range ds.Data.Airports.Airports {
		if airport.Iata == iata {
			ctxLogger.WithFields(logrus.Fields{
				"Method": "GetAirportIATA",
				"Found":  true,
				"IATA":   iata,
			}).Debug()
			return airport, nil
		}
	}
	ctxLogger.Errorf("Failed to find %s in dataset for IATA's", iata)
	return nil, fmt.Errorf("Failed to find %s in dataset for IATA's", iata)
}

func (*server) GetAirportICAO(ct context.Context, in *iatafinder.ICAO) (*iatafinder.Airport, error) {
	icao := in.GetIcao()
	for _, airport := range ds.Data.Airports.Airports {
		if airport.Icao == icao {
			ctxLogger.WithFields(logrus.Fields{
				"Method": "GetAirportIATA",
				"Found":  true,
				"ICAO":   icao,
			}).Debug()
			return airport, nil
		}
	}
	ctxLogger.Errorf("Failed to find %s in dataset for ICAO's", icao)
	return nil, fmt.Errorf("Failed to find %s in dataset for ICAO's", icao)
}

func (*server) GetAllAirports(ctx context.Context, in *iatafinder.EmptyRequest) (*iatafinder.Airports, error) {
	ctxLogger.WithFields(logrus.Fields{
		"Method": "GetAllAirports",
		"Found":  true,
	}).Debug()
	return &ds.Data.Airports, nil
}

func (*server) GetAirlineIATA(ctx context.Context, in *iatafinder.IATA) (*iatafinder.Airline, error) {
	iata := in.GetIata()
	for _, airline := range ds.Data.Airlines.Airlines {
		if airline.Iata == iata {
			ctxLogger.WithFields(logrus.Fields{
				"Method": "GetAirlineIATA",
				"Found":  true,
				"IATA":   iata,
			}).Debug()
			return airline, nil
		}
	}
	ctxLogger.Errorf("Failed to find %s in dataset for Airport IATA's", iata)
	return nil, fmt.Errorf("Failed to find %s in dataset for Airport IATA's", iata)
}

func (*server) GetAirlineICAO(ctx context.Context, in *iatafinder.ICAO) (*iatafinder.Airline, error) {
	icao := in.GetIcao()
	for _, airline := range ds.Data.Airlines.Airlines {
		if airline.Icao == icao {
			ctxLogger.WithFields(logrus.Fields{
				"Method": "GetAirlineICAO",
				"Found":  true,
				"ICAO":   icao,
			}).Debug()
			return airline, nil
		}
	}
	ctxLogger.Errorf("Failed to find %s in dataset for Airport ICAO's", icao)
	return nil, fmt.Errorf("Failed to find %s in dataset for Airport ICAO's", icao)
}

func (*server) GetAllAirlines(ctx context.Context, in *iatafinder.EmptyRequest) (*iatafinder.Airlines, error) {
	ctxLogger.WithFields(logrus.Fields{
		"Method": "GetAllAirlines",
		"Found":  true,
	}).Debug()
	return &ds.Data.Airlines, nil
}

func (*server) GetAirlines(ctx context.Context, in *iatafinder.AirlineDescriptor) (*iatafinder.Airlines, error) {
	descriptor := in.GetDescriptor_()
	airlines := []*iatafinder.Airline{}
	for _, airline := range ds.Data.Airlines.Airlines {
		if strings.Contains(airline.GetAlias(), descriptor) ||
			strings.Contains(airline.GetCountry(), descriptor) ||
			strings.Contains(airline.GetCallsign(), descriptor) ||
			strings.Contains(airline.GetName(), descriptor) {
			airlines = append(airlines, airline)
		}
	}
	ctxLogger.WithFields(logrus.Fields{
		"Method": "GetAirlines",
		"Found":  len(airlines) > 0,
	}).Debug()

	return &iatafinder.Airlines{Airlines: airlines}, nil
}
