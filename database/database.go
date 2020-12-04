package database

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"github.com/pedroppinheiro/cwnotifier/config"
)

const (
	verifyQuerySQL string = "select 1"

	//Chamados prioritários (1,2) que foram encaminhados para a GERIN e que estão sem responsável. Ao se atribuir ao chamado a notificação deve parar
	getIncidentsWithoutOwnerQuery string = "select NumeroIncidente from Incidente where OwnedByTeam = :team and Prioridade in (1,2) and OwnerID = '' and Status in ('Encaminhado', 'Novo')"

	//Tarefas prioritárias (1 ou 2) para a GERIN que estão sem responsável ou atribuídas para mim. Ao iniciar a tarefa a notificação deve parar
	getTasksWithoutOwnerQuery string = `select t.ParentPublicID from Tarefas t
where t.OwnedByTeam = :team
and t.Status in ('Encaminhada', 'Nova')
and (t.EmailResponsavel = :email or t.EmailResponsavel = '')
and (select i.Prioridade from Incidente i where i.NumeroIncidente = t.ParentPublicID) in (1,2)`

	//Chamados prioritários (1 ou 2) para a GERIN que estão atribuídas para mim e que já podem ser concluídas. Ao concluir o chamado ou criar uma nova tarefa a notificação deve parar
	getIncidentsWithTasksQuery string = `select i.NumeroIncidente, i.Tarefas, count(*) from Incidente i
join Tarefas t on i.NumeroIncidente = t.ParentPublicID
where i.OwnedByTeam = :team
and i.Prioridade in (1,2) 
and i.Status not in ('Resolvido', 'Fechado')
and (i.OwnerID = '' or i.OwnedBy = :userName)
and t.Status = 'Fechada'
group by i.NumeroIncidente, i.Tarefas, t.Status`

	//Chamados prioritários (1 ou 2) para a GERIN que estão atribuídas para mim e que já podem ser concluídas. Ao concluir o chamado ou criar uma nova tarefa a notificação deve parar
	getChangesThatNeedToBeValidatedQuery string = `select NumeroMudanca from Mudanca where Status = 'Resolvida' and CreatedBy = :userName`
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

func executeQuery(query string, args ...interface{}) (*sql.Rows, error) {
	log.Printf("Executing query \"%v\".", query)
	return connection.Query(query, args...)
}

// GetIncidentsWithoutOwner returns the incidents without owner
func GetIncidentsWithoutOwner(teamName string) []string {
	var (
		incidentNumber string
		results        []string
	)

	rows, err := executeQuery(getIncidentsWithoutOwnerQuery, sql.Named("team", teamName))

	if err != nil {
		log.Panic("Error getting incidents without owner. ", err)
	}
	for rows.Next() {
		err := rows.Scan(&incidentNumber)
		if err != nil {
			log.Panic(err)
		}

		results = append(results, incidentNumber)

	}

	log.Printf("GetIncidentsWithoutOwner: Found %v results: %v", len(results), results)
	return results
}

// GetTasksWithoutOwner returns the tasks without owner
func GetTasksWithoutOwner(teamName string, email string) []string {
	var (
		taskNumber string
		results    []string
	)

	rows, err := executeQuery(getTasksWithoutOwnerQuery, sql.Named("team", teamName), sql.Named("email", email))

	if err != nil {
		log.Panic("Error getting tasks without owner. ", err)
	}
	for rows.Next() {
		err := rows.Scan(&taskNumber)
		if err != nil {
			log.Panic(err)
		}

		results = append(results, taskNumber)

	}

	log.Printf("GetTasksWithoutOwner: Found %v results: %v", len(results), results)
	return results
}

// GetIncidentsWithClosedTasks returns incidents with tasks
func GetIncidentsWithClosedTasks(teamName string, userName string) []string {
	var (
		incidentNumber      string
		taskDescription     string
		numberOfClosedTasks string
		results             []string
	)

	rows, err := executeQuery(getIncidentsWithTasksQuery, sql.Named("team", teamName), sql.Named("userName", userName))

	if err != nil {
		log.Panic("Error getting incidents with tasks. ", err)
	}
	for rows.Next() {
		err := rows.Scan(&incidentNumber, &taskDescription, &numberOfClosedTasks)
		if err != nil {
			log.Panic(err)
		}

		if getTotalTasksFromTaskDescription(taskDescription) == numberOfClosedTasks {
			results = append(results, incidentNumber)
		}
	}

	log.Printf("GetIncidentsWithClosedTasks: Found %v results: %v", len(results), results)
	return results
}

var taskDescriptionRegex = regexp.MustCompile(`(?mi)\d+ Fechadas de (?P<totalTasks>\d+) Tarefas`)

func getTotalTasksFromTaskDescription(taskDescription string) string {
	match := taskDescriptionRegex.FindStringSubmatch(taskDescription)

	if match == nil || len(match) == 0 {
		log.Panicf("invalid task description found, was expecting \"\\d Fechadas de \\d Tarefas\", but found \"%v\"", taskDescription)
	}

	return match[1]
}

// GetChangesThatNeedToBeValidated returns changes that need to be validated
func GetChangesThatNeedToBeValidated(userName string) []string {
	var (
		changeNumber string
		results      []string
	)

	rows, err := executeQuery(getChangesThatNeedToBeValidatedQuery, sql.Named("userName", userName))

	if err != nil {
		log.Panic("Error getting changes that need to be validated. ", err)
	}
	for rows.Next() {
		err := rows.Scan(&changeNumber)
		if err != nil {
			log.Panic(err)
		}

		results = append(results, changeNumber)
	}

	log.Printf("GetChangesThatNeedToBeValidated: Found %v results: %v", len(results), results)
	return results
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
