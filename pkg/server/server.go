package server

import (
	"context"
	"github.com/pkg/errors"
	"github.com/the4thamigo-uk/ingestor/pkg/bimap"
	"github.com/the4thamigo-uk/ingestor/pkg/contact"
	"github.com/the4thamigo-uk/ingestor/pkg/identity"
	"github.com/the4thamigo-uk/ingestor/pkg/log"
	pb "github.com/the4thamigo-uk/ingestor/pkg/protocol"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"net"
	"os"
)

// Server is the public interface used to hide the implementation details of server.
type Server interface {
	Serve(l net.Listener) error
	Shutdown()
}

// server manages a server that allows clients to access data in CSV file datasources.
type server struct {
	s *grpc.Server
	m bimap.BiMap
	l log.Logger
}

// New creates a new instance of a server using the given options
func New(opts ...Option) (Server, error) {
	s := &server{
		s: grpc.NewServer(),
		m: bimap.New(),
		l: log.Disabled(),
	}
	pb.RegisterIngestorServer(s.s, s)
	reflection.Register(s.s)
	for _, opt := range opts {
		err := opt(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

// Serve begins listening for client connections, on the given listener
func (s *server) Serve(l net.Listener) error {
	return s.s.Serve(l)
}

// Shutdown attempts a graceful shutdown of any executing server() operation.
// Following this call, it is the caller's responsibility to wait for any executing Serve() function complete.
func (s *server) Shutdown() {
	s.s.GracefulStop()
}

// AddSource is the implementation of the callback from grpc when it receives a request from the client.
func (s *server) AddSource(ctx context.Context, req *pb.AddRequest) (*pb.AddResponse, error) {
	s.l.Printf("AddSource request received")
	id, err := s.addSourceFile(req.Filename)
	if err != nil {
		s.l.Printf("AddSource error '%s'", err)
		return nil, err
	}
	return &pb.AddResponse{
		Source: &pb.Source{
			Id:       id,
			Filename: req.Filename,
		},
	}, nil
}

// ListSources is the implementation of the callback from grpc when it receives a request from the client.
func (s *server) ListSources(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	s.l.Printf("ListSources request received")
	var sources []*pb.Source
	s.m.Iterate(
		func(id string, file string) {
			sources = append(sources,
				&pb.Source{
					Id:       id,
					Filename: file,
				})
		})
	return &pb.ListResponse{
		Sources: sources,
	}, nil

}

// ReadSource is the implementation of the callback from grpc when it receives a request from the client.
func (s *server) ReadSource(req *pb.ReadRequest, st pb.Ingestor_ReadSourceServer) error {
	s.l.Printf("ReadSource request received")
	err := s.readSource(req.Id, st)
	if err != nil {
		s.l.Printf("ReadSource error '%s'", err)
		return err
	}
	return nil
}

func (s *server) addSourceFiles(files []string) ([]string, error) {
	var ids []string
	for _, file := range files {
		id, err := s.addSourceFile(file)
		if err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (s *server) addSourceFile(file string) (string, error) {
	var err error
	id, ok := s.m.GetByVal(file)
	if !ok {
		if _, err := os.Stat(file); os.IsNotExist(err) {
			return "", err
		}
		id, err = identity.New()
		if err != nil {
			return "", err
		}
		s.m.Put(id, file)
		s.l.Printf("Source Added for '%s' with id '%s", file, id)
	}
	return id, nil
}

func (s *server) readSource(id string, st pb.Ingestor_ReadSourceServer) error {
	r, err := s.getReader(id)
	if err != nil {
		return err
	}
	return writeToStream(r, st)
}

func (s *server) getReader(id string) (contact.Reader, error) {
	file, ok := s.m.Get(id)
	if !ok {
		return nil, errors.Errorf("The source with id '%s' is not available", id)
	}
	return contact.NewReaderFromFile(file)
}

func writeToStream(r contact.Reader, st pb.Ingestor_ReadSourceServer) error {
	for {
		c, err := r.Read()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		rsp := &pb.ReadResponse{
			Mobile: c.Mobile,
			Name:   c.Name,
			Email:  c.Email,
		}
		err = st.Send(rsp)
		if err != nil {
			return err
		}
	}
}
