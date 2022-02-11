package internal

import (
	"amigo/pkg"
	"database/sql"
	"log"
)

func CreateMigrationTable() {
	var query string
	switch *pkg.DbProvider {
	case "mysql":
		query = mysqlQuery
		break
	case "postgres":
		query = postgresQuery
		break
	}

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

var mysqlQuery = `CREATE TABLE IF NOT EXISTS migration (
					id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    				version VARCHAR(30) NOT NULL 
					)`

var postgresQuery = `CREATE TABLE IF NOT EXISTS migration (
    					id SERIAL PRIMARY KEY,
    					version VARCHAR(30) NOT NULL 
    			  	)`
