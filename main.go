// Package main implements server for the iatafinder rpc services
package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"abelgoodwin1988/iata-finder/pkg/dataservice"
	"abelgoodwin1988/iata-finder/pkg/logger"
	"abelgoodwin1988/iata-finder/pkg/server"
)

var ctxLogger = logger.CtxLogger
var ds dataservice.Dataservice

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

	lis, s := server.Create(&ds, "configs/rpc.config.toml")

	// Start the rpc server and if it fails, log it and give up all hope
	if err := s.Serve(lis); err != nil {
		ctxLogger.Fatalf("Failed to start iatafinder server\n%v\n", err)
	}

}
