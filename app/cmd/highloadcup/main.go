package main

import (
	"flag"
	"log"

	"github.com/ngalayko/highloadcup/app"
)

var (
	dataPath    = flag.String("data_path", "", "path to initial data")
	listenAddr  = flag.String("addr", ":80", "addr to listen")
	profileAddr = flag.String("profile_addr", "", "enable profile")
)

func main() {
	flag.Parse()

	a, err := app.New(*dataPath)
	if err != nil {
		log.Panic(err.Error())
	}

	go func() {
		if *profileAddr == "" {
			return
		}

		if err := a.ListenAndServeProfile(*profileAddr); err != nil {
			log.Printf("profile server stopped with error: %s", err)
		}
	}()

	if err := a.ListenAndServe(*listenAddr); err != nil {
		log.Printf("web server stopped with error: %s", err)
	}
}
