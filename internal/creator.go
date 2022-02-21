package internal

import (
	"amigo/pkg"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func CreateMigration(path string) {
	timestamp := time.Now().Unix()
	timeString := strconv.Itoa(int(timestamp))

	upFileName := path + timeString + "_up.sql"
	downFileName := path + timeString + "_down.sql"
	f1, err1 := os.Create(upFileName)
	f2, err2 := os.Create(downFileName)
	pkg.Ept(err1)
	pkg.Ept(err2)

	_ = f1.Close()
	_ = f2.Close()

	log.Info("Created up file: " + upFileName)
	log.Info("Created down file: " + downFileName)
}
