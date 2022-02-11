package internal

import (
	"amigo/pkg"
	"database/sql"
	"log"
)

type migration struct {
	Id      uint64
	Version string
}

func Migrate(step int, migrationPath string) {
	existVersionList := getExistVersionList()
	migrationFiles := GetMigrationFiles(migrationPath)

	counter := 0
	for k, path := range migrationFiles {
		if counter == step && counter != -1 {
			break
		}

		if k == 0 {
			continue
		}

		version, isUp := GetMigrationVersion(path)

		if !isUp {
			continue
		}

		if !InArray(version, existVersionList) {
			ExecuteMigration(path)
			insertVersion(version)
			counter++
		}
	}
}

func getExistVersionList() []string {
	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(connection)

	resultSet, err := connection.Query("SELECT id, version FROM migration ORDER BY version")
	if err != nil {
		log.Fatal(err)
	}

	var result []string
	for resultSet.Next() {
		var m migration
		err = resultSet.Scan(&m.Id, &m.Version)
		if err != nil {
			log.Fatal(err)
		}

		result = append(result, m.Version)
	}

	return result
}

func insertVersion(version string) {
	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {

		}
	}(connection)

	query := PrepareQuery("INSERT INTO migration (version) VALUES (?)")
	_, err := connection.Exec(query, version)
	if err != nil {
		log.Println("Insert migration err")
		log.Fatal(err)
	}
}
