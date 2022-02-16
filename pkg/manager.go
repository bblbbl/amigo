package pkg

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"os"
	"time"
)

var (
	DbProvider = flag.String("dbProvider", "mysql", "Database")
	dbName     = flag.String("dbName", "", "Database name")
	dbUser     = flag.String("dbUser", "", "Database user")
	dbPassword = flag.String("dbPassword", "", "Database password")
	dbPort     = flag.String("dbPort", "3306", "Database port")
	dbHost     = flag.String("dbHost", "", "Database port")
)

func GetConnection() *sql.DB {
	var (
		user     = os.Getenv("DB_USER")
		password = os.Getenv("DB_PASSWORD")
		name     = os.Getenv("DB_NAME")
		port     = os.Getenv("DB_PORT")
		host     = os.Getenv("DB_HOST")
		provider = os.Getenv("DB_PROVIDER")
	)

	if user == "" || password == "" || name == "" || host == "" || provider == "" {
		user = *dbUser
		password = *dbPassword
		name = *dbName
		port = *dbPort
		host = *dbHost
		provider = *DbProvider
	}

	var connection *sql.DB

	switch provider {
	case "mariadb":
		connection = getMysql(user, password, name)
	case "mysql":
		connection = getMysql(user, password, name)
	case "postgres":
		connection = getPostgres(user, password, name, host, port)
	}

	return connection
}

func getPostgres(user string, password string, name string, host string, port string) *sql.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		name,
	)

	db, err := sql.Open("postgres", psqlInfo)
	Ept(err)

	err = db.Ping()
	Ept(err)

	return db
}

func getMysql(user string, password string, name string) *sql.DB {
	connectionString := fmt.Sprintf(
		"%s:%s@/%s",
		user,
		password,
		name,
	)

	db, err := sql.Open("mysql", connectionString)
	Ept(err)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
