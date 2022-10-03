package myutil

import (
	"log"
	"time"
)

func TrackTime(msg string) (string, time.Time) {
	return msg, time.Now()
}

func LogElapsed(msg string, start time.Time) {
	log.Printf("%v: took: %v\n", msg, time.Since(start))
}