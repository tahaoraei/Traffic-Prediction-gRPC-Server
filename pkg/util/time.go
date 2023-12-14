package util

import (
	"time"
	_ "time/tzdata"
	"timeMachine/pkg/logger"
)

var log = logger.Get()

func Now(tz string) string {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatal().Msgf("can't find location timezone ", err)
	}
	return time.Now().In(loc).Format("2006-01-02 15:04:05")
}

func Weekday(tz string) int {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatal().Msgf("can't find location timezone ", err)
	}
	return int(time.Now().In(loc).Weekday())
}
