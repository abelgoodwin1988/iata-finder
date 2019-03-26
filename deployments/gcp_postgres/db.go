// Package main handles the initial set up and population of the database
package main

import (
	"crypto/tls"
	"crypto/x509"
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
	db, err := connect(connCFG)
	if err != nil {
		log.Fatalf("Failed to create database Connection: %v", err)
	}
	defer db.Close()

	// Create the tables by loading the creation script into
	//	a string var, and executing it against the db.
	createTables, _ := ioutil.ReadFile("./scripts/0001_01_create_tables.sql")
	createTablesS := string(createTables)
	if _, err := db.Exec(createTablesS); err != nil {
		err = errors.Wrapf(err, "Table creation query failed (%s)", createTablesS)
	}

	// Use the postgres copy protocol for a quicker bulk insert of data.
	//	First copy in airports, and then copy in airlines.
	if copyCount, err := db.CopyFrom(
		pgx.Identifier{"airport"},
		[]string{"id", "name", "city", "country", "iata", "icao", "latitude", "longitude", "altitude", "timezone", "daylight_savings_time", "tz", "type", "source"},
		pgx.CopyFromRows(airports.Values()),
	); err != nil {
		fmt.Errorf("Failed to use postgres copy protocol to copy data to airports: %v", err)
	} else {
		fmt.Printf("Inserted %v records into airport\n", copyCount)
	}
	if copyCount, err := db.CopyFrom(
		pgx.Identifier{"airline"},
		[]string{"id", "name", "alias", "iata", "icao", "callsign", "country", "active"},
		pgx.CopyFromRows(airlines.Values()),
	); err != nil {
		fmt.Errorf("Failed to use postgres copy protocol to copy data to airlines: %v", err)
	} else {
		fmt.Printf("Inserted %v records into airline\n", copyCount)
	}
	fmt.Println("Success!")
}

func connect(connCFG ConnectionConfiguration) (*pgx.Conn, error) {
	// Enumerate the certs required with relative paths.
	clientCert := "../../certificates/client-cert.pem"
	clientKey := "../../certificates/client-key.pem"
	serverCA := "../../certificates/server-ca.pem"
	// Create keypair for client key/key
	client, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		fmt.Errorf("Failed to create X509 Key Pair: %v", err)
		return nil, err
	}

	// Create CA cert pool for root ca
	roots := x509.NewCertPool()
	certCA, _ := ioutil.ReadFile(serverCA)
	roots.AppendCertsFromPEM(certCA)
	conn := pgx.ConnConfig{
		Host:     connCFG.Host,
		Port:     connCFG.Port,
		Database: connCFG.Dbname,
		User:     connCFG.User,
		Password: connCFG.Password,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         connCFG.Host,
			Certificates:       []tls.Certificate{client},
			RootCAs:            roots,
		},
	}
	db, err := pgx.Connect(conn)
	if err != nil {
		fmt.Errorf("Failed to connect: %v", err)
		return nil, err
	}

	fmt.Printf("Connected Successfully to %s\n", connCFG.Host)
	return db, nil
}
