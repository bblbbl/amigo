package internal

import (
	"amigo/pkg"
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GetMigrationFiles(migrationDirectory string) []string {
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

func GetMigrationVersion(path string) (string, bool) {
	pathPartList := strings.Split(path, "/")
	requiredPart := pathPartList[len(pathPartList)-1]
	resultPartList := strings.Split(requiredPart, "_")

	return resultPartList[0], resultPartList[1] == "up.sql"
}

func InArray(val string, array []string) bool {
	for _, v := range array {
		if val == v {
			return true
		}
	}

	return false
}

func ExecuteMigration(path string) {
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
