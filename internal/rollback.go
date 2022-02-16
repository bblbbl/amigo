package internal

import (
	"amigo/pkg"
	"database/sql"
)

func Rollback(step int, migrationDirectory string) {
	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		pkg.Ept(err)
	}(connection)

	if step == -1 {
		step = 1
	}

	resultSet, err := connection.Query("SELECT id, version FROM migration ORDER BY version DESC LIMIT ?", step)
	defer resultSet.Close()

	var versionList []string
	for resultSet.Next() {
		var m migration
		err = resultSet.Scan(&m.Id, &m.Version)
		pkg.Ept(err)

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
		pkg.Ept(err)
	}(connection)

	_, err := connection.Exec("DELETE FROM migration WHERE version = ?", version)
	pkg.Ept(err)
}
