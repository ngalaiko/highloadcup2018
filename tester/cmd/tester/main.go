package main

import (
	"flag"
	"log"

	"github.com/ngalayko/highloadcup/tester"
)

var (
	appEndpoint         = flag.String("app_endpoint", "", "endpoint of the app")
	dataPath            = flag.String("data_path", "", "path to a test file")
	logPath             = flag.String("log_path", "/tmp/tester/result.log", "path to write log")
	healthCheckEndpoint = flag.String("healthcheck_endpoint", "", "application healthcheck endpoint")
)

func main() {
	flag.Parse()

	tester, err := tester.New(*appEndpoint, *dataPath)
	if err != nil {
		log.Panic(err)
	}

	if err := tester.Run(*healthCheckEndpoint, *logPath); err != nil {
		log.Panic(err)
	}
}
