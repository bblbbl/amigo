package internal

import (
	"amigo/pkg"
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type migration struct {
	Id      uint64
	Version string
}

func CreateMigrationTable() {
	query := `CREATE TABLE IF NOT EXISTS migration (
    			id INT NOT NULL PRIMARY KEY AUTO_INCREMENT,
    			version VARCHAR(30) NOT NULL 
    			)`

	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			panic(err)
		}
	}(connection)

	_, err := connection.Exec(query)
	if err != nil {
		panic(err)
	}
}

func CreateMigration(path string) {
	timestamp := time.Now().Unix()
	timeString := strconv.Itoa(int(timestamp))

	f1, err1 := os.Create(path + timeString + "_up.sql")
	f2, err2 := os.Create(path + timeString + "_down.sql")

	if err1 != nil || err2 != nil {
		log.Fatal(err1, err2)
	}

	_ = f1.Close()
	_ = f2.Close()
}

func Migrate(step int, migrationPath string) {
	existVersionList := getExistVersionList()
	migrationFiles := getMigrationFiles(migrationPath)

	counter := 0
	for k, path := range migrationFiles {
		if counter == step && counter != -1 {
			break
		}

		if k == 0 {
			continue
		}

		version, isUp := getMigrationVersion(path)

		if !isUp {
			continue
		}

		if !inArray(version, existVersionList) {
			executeMigration(path)
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

func getMigrationFiles(migrationDirectory string) []string {
	var files []string

	err := filepath.Walk(migrationDirectory, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files
}

func getMigrationVersion(path string) (string, bool) {
	pathPartList := strings.Split(path, "/")
	requiredPart := pathPartList[len(pathPartList)-1]
	resultPartList := strings.Split(requiredPart, "_")

	return resultPartList[0], resultPartList[1] == "up.sql"
}

func inArray(val string, array []string) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}

	return false
}

func executeMigration(path string) {
	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(connection)

	rawQuery, err := ioutil.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	queryList := strings.Split(string(rawQuery), ";")

	for _, query := range queryList {
		_, err = connection.Exec(query)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func insertVersion(version string) {
	connection := pkg.GetConnection()

	defer func(connection *sql.DB) {
		err := connection.Close()
		if err != nil {

		}
	}(connection)

	_, err := connection.Exec("INSERT INTO migration (version) VALUES (?)", version)
	if err != nil {
		log.Fatal(err)
	}
}

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

	migrationFiles := getMigrationFiles(migrationDirectory)

	for _, version := range versionList {
		for k, path := range migrationFiles {
			if k == 0 {
				continue
			}

			v, isUp := getMigrationVersion(path)

			if isUp {
				continue
			}

			if version == v {
				executeMigration(path)
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
