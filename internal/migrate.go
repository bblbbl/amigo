package internal

import (
	"amigo/pkg"
	"database/sql"
	log "github.com/sirupsen/logrus"
)

type migration struct {
	Id      uint64
	Version string
}

func Migrate(step int, migrationPath string) {
	existVersionList := getExistVersionList()
	migrationFiles := GetMigrationFiles(migrationPath)

	counter := 0
	appliedMigration := 0
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
			log.Warning("Migrating: " + version)
			ExecuteMigration(path)
			insertVersion(version)
			counter++
			appliedMigration++
			log.Info("Successfully migrated " + version)
		}
	}

	if appliedMigration == 0 {
		log.Warning("Nothing to migrate")
	} else {
		log.Infof("Applied %d migrations", appliedMigration)
	}
}

func getExistVersionList() []string {
	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		pkg.Ept(err)
	}(connection)

	resultSet, err := connection.Query("SELECT id, version FROM migration ORDER BY version")
	pkg.Ept(err)

	defer resultSet.Close()

	var result []string
	for resultSet.Next() {
		var m migration
		err = resultSet.Scan(&m.Id, &m.Version)
		pkg.Ept(err)

		result = append(result, m.Version)
	}

	return result
}

func insertVersion(version string) {
	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		pkg.Ept(err)
	}(connection)

	query := PrepareQuery("INSERT INTO migration (version) VALUES (?)")
	_, err := connection.Exec(query, version)
	if err != nil {
		pkg.Ept(err)
	}
}
