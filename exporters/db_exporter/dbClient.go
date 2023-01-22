package db_exporter

import (
	"database/sql"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var dbConnection *sql.DB

func SetupDBConnection(logName string) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		"URL", 30000, "USER", "PASSWORD", "TABLENAME")

	connection, err := sql.Open("pgx", psqlInfo)

	if err != nil {
		fmt.Println(err)
	}
	dbConnection = connection
	_, err = dbConnection.Exec(`CREATE TABLE IF NOT EXISTS ` + logName + `(
		id serial PRIMARY KEY, 
		nodeId varchar(255),
		nodeName varchar(255),
		value varchar(255),
		timestamp timestamp,
		logName varchar(255),
		server varchar(255)
	)`)

	if err != nil {
		fmt.Println(err)
	}
}

func InsertValues(namespace string, nodeId string, nodeName string, value interface{}, timestamp time.Time, logName string, server string) {

	sqlStatement := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	insert, val, _ := sqlStatement.Insert(namespace).Columns("nodeId", "nodeName", "value", "timestamp", "logName", "server").Values(nodeId, nodeName, fmt.Sprint(value), timestamp, logName, server).ToSql()

	_, err := dbConnection.Exec(insert, val...)

	if err != nil {
		fmt.Println(err)
	}

}
