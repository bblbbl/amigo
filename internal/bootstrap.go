package internal

import (
	"amigo/pkg"
	"database/sql"
)

func CreateMigrationTable() {
	query := queryMap[GetDBProvider()]

	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		pkg.Ept(err)
	}(connection)

	_, err := connection.Exec(query)
	pkg.Ept(err)
}

var queryMap = map[string]string{
	"mariadb":  mysqlQuery,
	"mysql":    mysqlQuery,
	"postgres": postgresQuery,
}

var mysqlQuery = `CREATE TABLE IF NOT EXISTS migration (
					id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    				version VARCHAR(30) NOT NULL 
					)`

var postgresQuery = `CREATE TABLE IF NOT EXISTS migration (
    					id SERIAL PRIMARY KEY,
    					version VARCHAR(30) NOT NULL 
    			  	)`
