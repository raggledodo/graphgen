package main

import (
	"flag"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	defaultFrequency = "30s"
	defaultLogLevel  = "info"
	defaultOutDir    = "tmp"
)

var (
	config Config
)

type Config struct {
	freq   time.Duration
	outdir string
}

func init() {
	level := flag.String("log-level", defaultLogLevel,
		"Log level: 'debug' | 'info' | 'warn' | 'error' | 'fatal' | 'panic'")
	freq := flag.String("freq", defaultFrequency,
		"Script generation frequency")
	flag.StringVar(&config.outdir, "outdir", defaultOutDir,
		"Temporary directory for holding proto and data generating script")
	flag.Parse()

	// set log level
	logLevel, err := log.ParseLevel(*level)
	if err != nil {
		log.Fatalf("Log level parsing error: %v", err)
	}
	log.SetLevel(logLevel)

	dur, err := time.ParseDuration(*freq)
	if err != nil {
		log.Fatalf("Failed to parse frequency: %v", err)
	}
	config.freq = dur
}
