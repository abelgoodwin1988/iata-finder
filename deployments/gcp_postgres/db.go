// Package main handles the initial set up and population of the database
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"
)

/**
 * db is the main handler and entry point for the program to:
 *	connect to the database
 *	delete any schema/tables that exist in the place where
	we'll be creating our datastore
 *	populate the tables with the files found in the assets
	directory
*/
func main() {
	// get airports and airlines
	airports, airlines := ParseHandler()
	// fetch the connection information
	connCFG := GetConnectionConfiguration()
	// connect to and defer close of db
	db := connect(connCFG)
	defer db.Close()
	// Create the tables by loading the creation script into
	//	a string var, and executing it against the db.
	createTables, _ := ioutil.ReadFile("./scripts/0001_01_create_tables.sql")
	createTablesS := string(createTables)
	if _, err := db.Exec(createTablesS); err != nil {
		err = errors.Wrapf(err, "Table creation query failed (%s)", createTablesS)
	}

	if copyCount, err := db.CopyFrom(
		pgx.Identifier{"airport"},
		[]string{
			"id",
			"name",
			"city",
			"country",
			"iata",
			"icao",
			"latitude",
			"longitude",
			"altitude",
			"timezone",
			"daylight_savings_time",
			"tz",
			"type",
			"source",
		},
		pgx.CopyFromRows(airports.Values()),
	); err != nil {
		log.Panic(err)
	} else {
		fmt.Printf("Inserted %v records into airport\n", copyCount)
	}

	if copyCount, err := db.CopyFrom(
		pgx.Identifier{"airline"},
		[]string{
			"id",
			"name",
			"iata",
			"icao",
			"callsign",
			"country",
			"active",
		},
		pgx.CopyFromRows(airlines.Values()),
	); err != nil {
		log.Panic(err)
	} else {
		fmt.Printf("Inserted %v records into airline\n", copyCount)
	}
}

func connect(connCFG ConnectionConfiguration) *pgx.Conn {
	// clientCert := "../../certificates/client-cert.pem"
	// clientKey := "../../certificates/client-key.pem"
	// serverCA := "../../certificates/server-ca.pem"
	// sslmode=verify-ca sslrootcert=%s sslcert=%s sslkey=%s
	psqlConn := fmt.Sprintf("host=%s port=%v user=%s dbname=%s password=%s",
		// serverCA,
		// clientCert,
		// clientKey,
		connCFG.Host,
		connCFG.Port,
		connCFG.User,
		connCFG.Dbname,
		connCFG.Password,
	)
	conn, _ := pgx.ParseConnectionString(psqlConn)
	db, err := pgx.Connect(conn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Connected Successfully to %s\n", connCFG.Host)
	return db
}
