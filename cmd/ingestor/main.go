package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/the4thamigo-uk/ingestor/cmd/ingestor/service"
	"github.com/the4thamigo-uk/interrupter"
	"os"
)

func main() {
	addr := pflag.StringP("address", "l", ":8080", "Address that the server should listen on")
	files := pflag.StringSliceP("files", "f", nil, "Path of a source data file to ingest")
	pflag.Parse()

	l := log.New()

	s, err := service.New(*addr, *files, l)

	if err != nil {
		handleError(err)
		return
	}

	irpt := interrupter.New(s.Stop)
	defer irpt.Close()

	err = s.Start()
	if err != nil {
		handleError(err)
		return
	}
}

func handleError(err error) {
	fmt.Println(err)
	pflag.PrintDefaults()
	os.Exit(1)
}
