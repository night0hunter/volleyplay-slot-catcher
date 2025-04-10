package utils

import (
	"log"
	"time"
)

const (
	TimeFormat         = "02.01.2006"
	MoscowLocationName = "Europe/Moscow"
)

var MoscowLocation *time.Location

func init() {
	var err error
	MoscowLocation, err = time.LoadLocation(MoscowLocationName)
	if err != nil {
		log.Fatal(err)
	}
}
