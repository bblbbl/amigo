package pkg

import log "github.com/sirupsen/logrus"

func Ept(err error) {
	if err != nil {
		log.Error("Error: " + err.Error())
		panic(err)
	}
}
