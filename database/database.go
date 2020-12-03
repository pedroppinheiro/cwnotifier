package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/pedroppinheiro/cwnotifier/config"
)

const (
	verifyQuerySQL string = "select 1"

	//Chamados prioritários (1,2) que foram encaminhados para a GERIN e que estão sem responsável. Ao se atribuir ao chamado a notificação deve parar
	getIncidentsWithoutOwnerNotificationQuery string = "select NumeroIncidente from Incidente where OwnedByTeam = :team and Prioridade in (1,2) and OwnerID = '' and Status in ('Encaminhado', 'Novo')"

	//Tarefas prioritárias (1 ou 2) para a GERIN que estão sem responsável ou atribuídas para mim. Ao iniciar a tarefa a notificação deve parar
	getTasksWithoutOwnerNotificationQuery string = `select t.ParentPublicID from Tarefas t
													where (t.EmailResponsavel = 'smaia@banparanet.com.br' or t.EmailResponsavel = '')
													and t.Status in ('Encaminhada', 'Nova')
													and t.OwnedByTeam = 'SUSIS - GERIN'
													and (select i.Prioridade from Incidente i where i.NumeroIncidente = t.ParentPublicID) in (1,2)`

	//Chamados prioritários (1 ou 2) para a GERIN que estão atribuídas para mim e que já podem ser concluídas. Ao concluir o chamado ou criar uma nova tarefa a notificação deve parar
	getIncidentsWithTasksQuery string = `select i.NumeroIncidente, i.Tarefas from Incidente i
															where i.OwnedByTeam = 'SUSIS - GERIN' 
															and i.Prioridade in (1,2) 
															and i.Status not in ('Resolvido', 'Fechado')
															and (i.OwnerID = '' or i.OwnedBy = 'Pedro Victor Pontes Pinheiro')`
	//TODO: pegar o texto e passar numa regex se da match de número de fechadas igual a tarefas
	//exemplo de resposta:
	//1240672                    1 Fechadas de 3 Tarefas

	//Chamados prioritários (1 ou 2) para a GERIN que estão atribuídas para mim e que já podem ser concluídas. Ao concluir o chamado ou criar uma nova tarefa a notificação deve parar
	getChangesThatNeedToBeValidated string = `select NumeroMudanca from Mudanca where Status = 'Resolvida' and CreatedBy like '%Silvana%'`
	//TODO: pegar o texto e passar numa regex se da match de número de fechadas igual a tarefas
	//exemplo de resposta:
	//1240672                    1 Fechadas de 3 Tarefas
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

// GetIncidentsWithoutOwner returns the number of priority incidents
func GetIncidentsWithoutOwner(teamName string) (int, string) {
	var (
		incidentNumber  string
		result          string
		numberOfResults int
	)

	rows, err := executeQuery(getIncidentsWithoutOwnerNotificationQuery, sql.Named("team", teamName))

	if err != nil {
		log.Panic("Error getting number of incidents without owner. ", err)
	}
	for rows.Next() {
		err := rows.Scan(&incidentNumber)
		if err != nil {
			log.Panic(err)
		}

		numberOfResults++

		if result == "" {
			result += incidentNumber
		} else {
			result += ", " + incidentNumber
		}
	}

	log.Printf("Found %v results: %v", numberOfResults, result)

	return numberOfResults, result
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
