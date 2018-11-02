package server

import (
	"github.com/the4thamigo-uk/ingestor/pkg/log"
)

// Option defines the interface to use to set options on the Server instance
type Option func(*server) error

// WithLogger sets which logger instance the server should use
func WithLogger(l log.Logger) Option {
	return func(s *server) error {
		s.l = l
		return nil
	}
}

// WithSourceFiles sets up the server with some initial data source files.
// TODO: each time the server restarts it creates a new ID for the same
// datafile. This means clients will download the same files multiple times.
// We could instead use file hashes as IDs, but this might be inefficient for
// large data files. Alternatively, using filenames as IDs exposes implementation
// details to the clients of our server.
func WithSourceFiles(files []string) Option {
	return func(s *server) error {
		_, err := s.addSourceFiles(files)
		return err
	}
}
