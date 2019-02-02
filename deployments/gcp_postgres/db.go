// Package main handles the initial set up and population of the database
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
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
	connCFG := GetConnectionConfiguration()
	fmt.Printf("%v", connCFG)
	clientCert := "../../certificates/client-cert.pem"
	clientKey := "../../certificates/client-key.pem"
	serverCA := "../../certificates/server-ca.pem"
	psqlConn := fmt.Sprintf("sslmode=verify-ca sslrootcert=%s sslcert=%s sslkey=%s host=%s port=%s user=%s dbname=%s password=%s",
		serverCA,
		clientCert,
		clientKey,
		connCFG.Host,
		connCFG.Port,
		connCFG.User,
		connCFG.Dbname,
		connCFG.Password,
	)
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Print("Connected Successfully")
}
