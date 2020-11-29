package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
)

const (
	server   string = "localhost"
	port     string = "1433"
	user     string = "sa"
	password string = "redeinf123"
	database string = "cherwell"

	verifyQuerySQL                string = "select 1"
	getNumberOfPriorityTasksQuery string = "select count(*) from cherwell where prioridade = 1"
)

var connection *sql.DB

// Connect connects to the database
func Connect() {
	var errConnect, errVerifyConnection error

	connection, errConnect = sql.Open(
		"mssql",
		fmt.Sprintf("server=%v;port=%v;user id=%v;password=%v;database=%v;", server, port, user, password, database),
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

// GetNumberOfPriorityTasks returns the number of priority tasks
func GetNumberOfPriorityTasks() int {
	var queryResult string

	rows, err := executeQuery(getNumberOfPriorityTasksQuery)

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

	numberOfImportantTasks, err := strconv.Atoi(queryResult)

	if err != nil {
		log.Panic(err)
	}

	return numberOfImportantTasks
}

func executeQuery(query string) (*sql.Rows, error) {
	log.Printf("Executing query \"%v\".", query)
	return connection.Query(query)
}

func verifyConnection() error {
	_, err := executeQuery(verifyQuerySQL)
	return err
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