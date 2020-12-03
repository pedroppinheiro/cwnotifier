package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/pedroppinheiro/cwnotifier/config"
)

const (
	verifyQuerySQL                    string = "select 1"
	getNumberOfPriorityIncidentsQuery string = "select count(*) from Incidente where OwnedByTeam = 'SUSIS - GERIN' and Prioridade in (1,2) and OwnerID = '' and Status = 'Encaminhado'"
)

var connection *sql.DB

// Connect connects to the database
func Connect(databaseConfig config.Database) {
	var errConnect, errVerifyConnection error

	connection, errConnect = sql.Open(
		"mssql",
		fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v;", databaseConfig.Server, databaseConfig.Port, databaseConfig.User, databaseConfig.Password, databaseConfig.DatabaseName),
	)

	if errConnect != nil {
		log.Panic("Error creating connection object. ", errConnect)
	}

	errVerifyConnection = verifyConnection()

	if errVerifyConnection != nil {
		log.Panic("Error connecting to database. ", errVerifyConnection)
	}

	log.Println("Connected successfully to database.")
}

func verifyConnection() error {
	_, err := executeQuery(verifyQuerySQL)
	return err
}

func executeQuery(query string) (*sql.Rows, error) {
	log.Printf("Executing query \"%v\".", query)
	return connection.Query(query)
}

// GetNumberOfPriorityIncidents returns the number of priority incidents
func GetNumberOfPriorityIncidents() int {
	var queryResult string

	rows, err := executeQuery(getNumberOfPriorityIncidentsQuery)

	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		err := rows.Scan(&queryResult)
		if err != nil {
			log.Panic(err)
		}
		log.Printf("Query returned: %v\n", queryResult)
	}

	numberOfImportantIncidents, err := strconv.Atoi(queryResult)

	if err != nil {
		log.Panic(err)
	}

	return numberOfImportantIncidents
}

// CloseConnection closes the connection
func CloseConnection() {
	log.Println("Closing the connection with database.")
	err := connection.Close()

	if err != nil {
		log.Println("Error during closing connection with db. ", err)
	} else {
		log.Println("Connection with database was closed.")
	}
}
