package main

import (
	"flag"

	"github.com/robfig/cron"
	log "github.com/sirupsen/logrus"
)

const (
	defaultSchedule = "0 * * * * *"
	defaultLogLevel = "info"
)

var (
	config Config
)

type Config struct {
	schedule cron.Schedule
}

func init() {
	level := flag.String("log-level", defaultLogLevel,
		"Log level: 'debug' | 'info' | 'warn' | 'error' | 'fatal' | 'panic'")
	schedule := flag.String("schedule", defaultSchedule,
		"Cron schedule: '<second> <minute> <hour> <dom> <month>'")
	flag.Parse()

	// set log level
	logLevel, err := log.ParseLevel(*level)
	if err != nil {
		log.Fatalf("Log level parsing error: %v", err)
	}
	log.SetLevel(logLevel)

	// register cron job
	sched, err := cron.Parse(*schedule)
	if err != nil {
		log.Fatalf("Cron schedule parsing error: %v", err)
	}
	config.schedule = sched
}
