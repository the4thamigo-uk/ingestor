package service

import (
	"github.com/the4thamigo-uk/ingestor/pkg/client"
	"github.com/the4thamigo-uk/ingestor/pkg/log"
	"github.com/the4thamigo-uk/ingestor/pkg/store"
	"sort"
	"time"
)

// TODO: make poll period configurable
const pollSeconds = 5

// Service manages a persistent connection to download data sources from a remote party at the address 'addr'
// The service continuously polls for the existence of data sources that it has not previously downloaded.
type Service struct {
	addr  string
	st    store.Store
	c     *client.Client
	l     log.Logger
	close chan bool
}

// New creates a new instance of the service
func New(addr string, st store.Store, l log.Logger) *Service {
	return &Service{
		addr:  addr,
		st:    st,
		close: make(chan bool),
		l:     l,
	}
}

// Start begins polling the remote party for new data sources.
func (s *Service) Start() error {
	// TODO: might be better to listen to a stream rather than polling!
	var err error
	s.c, err = client.New(s.addr)
	if err != nil {
		return err
	}
	defer s.c.Close()
	for {
		s.l.Println("Polling")
		err := s.readAll(s.addr)
		if err != nil {
			s.l.Printf("Poll error '%s'", err)
		}

		if !s.wait(time.Second * pollSeconds) {
			break
		}
	}
	return nil
}

// Stop attempts to shutdown a running invocation of Start().
// It is the caller's responsibility to wait for Start() to complete.
func (s *Service) Stop() {
	if s.c != nil {
		s.c.Close()
		s.close <- true
	}
}

func (s *Service) wait(d time.Duration) bool {
	select {
	case <-s.close:
		return false
	case <-time.After(d):
		return true
	}
}

func (s *Service) readAll(addr string) error {
	var err error
	s.c, err = client.New(s.addr)
	if err != nil {
		s.l.Printf("Connect error '%s'", err)
	}
	defer s.c.Close()

	availIDs, err := s.c.ListSources()
	if err != nil {
		return err
	}
	doneIDs, err := s.st.Sources()
	if err != nil {
		return err
	}
	sort.Strings(doneIDs)

	for _, id := range availIDs {
		i := sort.SearchStrings(doneIDs, id)
		if i < len(doneIDs) {
			continue
		}

		// TODO: should consider batching to real database i.e. unit-of-work etc
		err = s.c.ReadSource(id,
			func(c client.Contact) error {
				s.l.Printf("Importing '%s,%s,%s'", c.Name(), c.Email(), c.Mobile())
				return s.st.PutContact(&store.Contact{
					Name:   c.Name(),
					Mobile: c.Mobile(),
					Email:  c.Email(),
				})
			})
		if err != nil {
			return err
		}

		err = s.st.PutSource(id)
		if err != nil {
			return err
		}
	}
	return nil
}
