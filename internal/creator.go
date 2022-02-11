package internal

import (
	"log"
	"os"
	"strconv"
	"time"
)

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
