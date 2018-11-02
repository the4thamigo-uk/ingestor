package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/the4thamigo-uk/ingestor/cmd/reader/service"
	"github.com/the4thamigo-uk/ingestor/pkg/store/memory"
	"github.com/the4thamigo-uk/interrupter"
	"os"
)

func main() {
	addr := pflag.StringP("address", "c", ":8080", "Address that the server should connect to")

	l := log.New()
	ms := memory.New()
	s := service.New(*addr, ms, l)

	irpt := interrupter.New(s.Stop)
	defer irpt.Close()

	err := s.Start()
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
