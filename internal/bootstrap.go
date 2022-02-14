package internal

import (
	"amigo/pkg"
	"database/sql"
	"log"
)

func CreateMigrationTable() {
	query := queryMap[*pkg.DbProvider]

	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			panic(err)
		}
	}(connection)

	_, err := connection.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
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
