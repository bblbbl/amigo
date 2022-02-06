package main

import (
	"amigo/internal"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var (
	dir      = flag.String("dir", "migration", "Directory for store/read project migration")
	create   = flag.Bool("create", false, "Create migration")
	migrate  = flag.Bool("migrate", false, "Run migration")
	rollback = flag.Bool("rollback", false, "Rollback migrations")
	step     = flag.Int("step", -1, "Count migration to migrate/rollback")
)

func main() {
	flag.Parse()
	_ = godotenv.Load()

	internal.CreateMigrationTable()

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	migrationDirectory := currentDir + "/" + *dir + "/"

	if *create {
		internal.CreateMigration(migrationDirectory)
	} else if *migrate {
		internal.Migrate(*step, migrationDirectory)
	} else if *rollback {
		internal.Rollback(*step, migrationDirectory)
	}
}
