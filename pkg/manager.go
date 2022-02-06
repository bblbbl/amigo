package pkg

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var (
	dbName     = flag.String("dbName", "", "Database name")
	dbUser     = flag.String("dbUser", "", "Database user")
	dbPassword = flag.String("dbPassword", "", "Database password")
)

func GetConnection() *sql.DB {
	var (
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		name     = os.Getenv("DB_NAME")
	)

	if user == "" || password == "" || name == "" {
		user = *dbUser
		password = *dbPassword
		name = *dbName
	}

	connectionString := fmt.Sprintf(
		"%s:%s@/%s",
		user,
		password,
		name,
	)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
