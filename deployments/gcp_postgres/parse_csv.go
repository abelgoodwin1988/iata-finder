package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

// ParseHandler routes and manages the parsing process from
//	csv's to go data structs and returns them
func ParseHandler() (Airports, Airlines) {
	var airport Airports
	var airline Airlines
	csvF, _ := os.Open("../../assets/airports.csv")
	reader := csv.NewReader(bufio.NewReader(csvF))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		// line = setNullValues(line, "\N")
		airport = append(airport, Airport{
			ID:                  mustAtoi(line[0]),
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
			TypeField:           line[12],
			Source:              line[13],
		})
	}
	csvF, _ = os.Open("../../assets/airlines.csv")
	reader = csv.NewReader(bufio.NewReader(csvF))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		// line = setNullValues(line, "\N")
		airline = append(airline, Airline{
			ID:       mustAtoi(line[0]),
			Iata:     line[1],
			Name:     line[2],
			Icao:     line[3],
			Callsign: line[4],
			Country:  line[5],
			Active:   line[6],
		})
	}
	fmt.Printf("Airports read in: %v \nAirlines read in: %v\n", len(airport), len(airline))
	return airport, airline
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Panic(err)
	}
	return i
}

func mustFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Panic(err)
	}
	return f
}
