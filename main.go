package main

import (
	"github.com/fatih/color"
	"log"
	"os"
	"time"
)

//const configFilesPath = "/var/opx/properties/bac_sum"

const layout = "2006-01-02 15:04:05"

var config Config

func main() {

	// Profiling: https://flaviocopes.com/golang-profiling/
	//defer profile.Start().Stop()
	//defer profile.Start(profile.MemProfile).Stop()

	cpl := CommandLineParser{}
	err := cpl.Parse(os.Args[1:])

	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	config = Config{}
	config.Load()
	config.Environment = cpl.Environment
	config.Source = cpl.Source

	// validate environment
	_, ok := config.Databases[cpl.Environment]
	if !ok {
		color.Red("Environment is not available.")
		os.Exit(1)
	}

	start := time.Now()
	log.Printf("---------------------------------")
	log.Printf("Start time: " + time.Now().Format(layout))
	log.Printf("Environment: " + cpl.Environment)

	queue := make(chan string)

	go func() {
		//queue <- "http://hcmaslov.d-real.sci-nnov.ru/"
		queue <- "http://ukr.net/"
		//queue <- "https://alexandr-rogers.livejournal.com/1091602.html"
	}()

	for uri := range queue {
		Crawler{}.enqueue(uri, queue)
	}

	elapsed := time.Since(start)
	log.Printf("Took %s", elapsed)
	log.Printf("End time: " + time.Now().Format(layout))
	log.Printf("---------------------------------")

}
