package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/the4thamigo-uk/ingestor/cmd/reader/service"
	"github.com/the4thamigo-uk/ingestor/pkg/store"
	"github.com/the4thamigo-uk/ingestor/pkg/store/cassandra"
	"github.com/the4thamigo-uk/ingestor/pkg/store/memory"
	"github.com/the4thamigo-uk/interrupter"
	"os"
)

func main() {
	addr := pflag.StringP("address", "c", ":8080", "Address that the server should connect to")
	sAddr := pflag.StringP("store", "s", "", "Address of Cassandra host")
	pflag.Parse()

	l := log.New()

	var err error
	var st store.Store

	if *sAddr == "" {
		st = memory.New()
	} else {
		c, err := cassandra.New("reader", *sAddr)
		if err != nil {
			handleError(err)
			return
		}
		defer c.Close()
		st = c
	}

	s := service.New(*addr, st, l)

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
