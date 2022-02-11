package internal

import (
	"amigo/pkg"
	"database/sql"
	"log"
)

func Rollback(step int, migrationDirectory string) {
	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {

		}
	}(connection)

	if step == -1 {
		step = 1
	}

	resultSet, err := connection.Query("SELECT id, version FROM migration ORDER BY version DESC LIMIT ?", step)
	var versionList []string
	for resultSet.Next() {
		var m migration
		err = resultSet.Scan(&m.Id, &m.Version)
		if err != nil {
			log.Fatal(err)
		}
		versionList = append(versionList, m.Version)
	}

	migrationFiles := GetMigrationFiles(migrationDirectory)

	for _, version := range versionList {
		for k, path := range migrationFiles {
			if k == 0 {
				continue
			}

			v, isUp := GetMigrationVersion(path)

			if isUp {
				continue
			}

			if version == v {
				ExecuteMigration(path)
				deleteMigration(version)
			}
		}
	}
}

func deleteMigration(version string) {
	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(connection)

	_, err := connection.Exec("DELETE FROM migration WHERE version = ?", version)
	if err != nil {
		return
	}
}
