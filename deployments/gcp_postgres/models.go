package main

// Airline is a struct representing openflights data for
//	IATA airline data
type Airline struct {
	ID       int
	Name     string
	Alias    string
	Iata     string
	Icao     string
	Callsign string
	Country  string
	Active   string
}

// Airlines is a set of Airline
type Airlines []Airline

// Airport is a struct representing openflight data for
//	IATA airport data
type Airport struct {
	ID                  int
	Name                string
	City                string
	Country             string
	Iata                string
	Icao                string
	Latitude            float64
	Longitude           float64
	Altitude            float64
	Timezone            string
	DaylightSavingsTime string
	Tz                  string
	TypeField           string
	Source              string
}

// Airports is a set of Airport
type Airports []Airport

// Values returns values of Airline
func (a Airline) Values() []interface{} {
	return []interface{}{a.ID, a.Name, a.Alias, a.Iata, a.Icao, a.Callsign, a.Country, a.Active}
}

// Values returns values of Airport
func (a Airport) Values() []interface{} {
	return []interface{}{a.ID, a.Name, a.City, a.Country, a.Iata, a.Icao, a.Latitude, a.Longitude, a.Altitude, a.Timezone, a.DaylightSavingsTime, a.Tz, a.TypeField, a.Source}
}

// Values returns values for Airports
func (a Airports) Values() [][]interface{} {
	rtrn := [][]interface{}{}
	for _, v := range a {
		vals := v.Values()
		rtrn = append(rtrn, vals)
	}
	return rtrn
}

// Values returns values for Airports
func (a Airlines) Values() [][]interface{} {
	rtrn := [][]interface{}{}
	for _, v := range a {
		vals := v.Values()
		rtrn = append(rtrn, vals)
	}
	return rtrn
}
