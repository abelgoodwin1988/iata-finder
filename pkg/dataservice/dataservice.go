package dataservice

import (
	"abelgoodwin1988/iata-finder/pkg/logger"
	iatafinder "abelgoodwin1988/iata-finder/rpc"
	"encoding/csv"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var ctxLogger = logger.CtxLogger.WithField("package", "dataservice")

type csvType struct {
	Reader io.Reader
	Name   string
}

type data struct {
	// wg       sync.WaitGroup
	Updated  time.Time
	Airports iatafinder.Airports
	Airlines iatafinder.Airlines
	CSVs     []csvType
}

// Dataservice exposes the methods and data gathered by
//  the dataservice package
type Dataservice struct {
	URLTargets      []string
	DataDestination string
	FileType        string
	Data            data
	Interval        time.Duration
}

// Init Initalizes the dataservice collector
func (d *Dataservice) Init(wg *sync.WaitGroup) {
	go func() {
		for {
			// Collect the data
			d.dataCollector()
			// Parse the data
			d.parseHandler()
			wg.Done()
			// Let's be kind to our friends open flight data and consume infrequently
			time.Sleep(d.Interval)
		}
	}()
}

// dataCollector requests the csv's for airports and airlines
func (d *Dataservice) dataCollector() {
	ctxLogger.WithFields(logrus.Fields{"Dataservice:": d}).Debugf("Starting datacollector")
	for _, urlTarget := range d.URLTargets {
		ctxLogger.Debugf("Fetching %s", urlTarget)
		// Get the file
		resp, err := http.Get(urlTarget)
		if err != nil {
			ctxLogger.WithError(err).Errorf("Failed to Get asset from URL: %s", urlTarget)
		}
		d.Data.CSVs = append(d.Data.CSVs, csvType{resp.Body, urlTarget})
	}
}

// ParseHandler routes and manages the parsing process from
//	csv's to go data structs and returns them
func (d *Dataservice) parseHandler() {
	ctxLogger.Debug("Starting Parse Handler")
	for _, csvReader := range d.Data.CSVs {
		reader := csv.NewReader(csvReader.Reader)
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				ctxLogger.WithError(err).Error("Error reading reader %s", csvReader.Name)
			}

			if strings.Contains(csvReader.Name, "airport") {
				d.Data.Airports.Airports = append(d.Data.Airports.Airports, &iatafinder.Airport{
					Id:                  mustAtoi(line[0]),
					Name:                line[1],
					City:                line[2],
					Country:             line[3],
					Iata:                line[4],
					Icao:                line[5],
					Latitude:            mustFloat(line[6]),
					Longitude:           mustFloat(line[7]),
					Altitude:            mustFloat(line[8]),
					Timezone:            line[9],
					DaylightSavingsTime: line[10],
					Tz:                  line[11],
					Type:                line[12],
					Source:              line[13],
				})
			} else if strings.Contains(csvReader.Name, "airline") {
				d.Data.Airlines.Airlines = append(d.Data.Airlines.Airlines, &iatafinder.Airline{
					ID:       mustAtoi(line[0]),
					Name:     line[1],
					Alias:    line[2],
					Iata:     line[3],
					Icao:     line[4],
					Callsign: line[5],
					Country:  line[6],
					Active:   line[7],
				})
			}
		}

	}

	ctxLogger.WithFields(logrus.Fields{"Airports": len(d.Data.Airports.Airports), "Airlines": len(d.Data.Airlines.Airlines)}).Debug("Values Read In")
	ctxLogger.Info("Finished Parse Handling")
}

// GetAirlines returns the dataservice current airlines
func (d *Dataservice) GetAirlines() iatafinder.Airlines {
	return d.Data.Airlines
}

// GetAirports returns the dataservice current airports
func (d *Dataservice) GetAirports() iatafinder.Airports {
	return d.Data.Airports
}

// GetUpdate returns the dataservice current airlines
func (d *Dataservice) GetUpdate() time.Time {
	return d.Data.Updated
}

func mustAtoi(s string) int32 {
	i, err := strconv.Atoi(s)
	if err != nil {
		ctxLogger.WithError(err).Errorf("Failed to convert string %s to int", s)
		var zero int32
		return zero
	}
	return int32(i)
}

func mustFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		ctxLogger.WithError(err).Errorf("Failed to convert string %s to float64", s)
		var zero float64
		return zero
	}
	return f
}
