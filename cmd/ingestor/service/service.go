package service

import (
	"github.com/the4thamigo-uk/ingestor/pkg/log"
	"github.com/the4thamigo-uk/ingestor/pkg/server"
	"net"
)

// Service manages a service that allows clients to download CSV file datasources
type Service struct {
	addr string
	s    server.Server
	l    log.Logger
}

// New creates a new Service instance
func New(addr string, files []string, l log.Logger) (*Service, error) {
	s, err := server.New(server.WithLogger(l), server.WithSourceFiles(files))
	if err != nil {
		return nil, err
	}
	return &Service{
		addr: addr,
		l:    l,
		s:    s,
	}, nil
}

// Start begins listening on the given addr
func (s *Service) Start() error {
	l, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.s.Serve(l)
}

// Stop attempts to gracefully stop listening.
func (s *Service) Stop() {
	s.s.Shutdown()
}
