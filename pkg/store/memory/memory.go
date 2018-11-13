package memory

import (
	"github.com/the4thamigo-uk/ingestor/pkg/store"
)

type memorystore struct {
	contacts map[string]*store.Contact
	sources  map[string]bool
}

// New creates a new in-memory store
func New() store.Store {
	return &memorystore{
		contacts: map[string]*store.Contact{},
		sources:  map[string]bool{},
	}
}

func (ms *memorystore) PutContact(c *store.Contact) error {
	return ms.PutContacts([]*store.Contact{c})
}

func (ms *memorystore) PutSource(id string) error {
	return ms.PutSources([]string{id})
}

func (ms *memorystore) PutContacts(cs []*store.Contact) error {
	// store contacts keyed on the mobile id
	for _, c := range cs {
		ms.contacts[c.Mobile] = c
	}
	return nil
}

func (ms *memorystore) PutSources(ss []string) error {
	// store contacts keyed on the mobile id
	for _, s := range ss {
		ms.sources[s] = true
	}
	return nil
}

func (ms *memorystore) Contacts() ([]*store.Contact, error) {
	var cs []*store.Contact
	for _, c := range ms.contacts {
		cc := *c
		cs = append(cs, &cc)
	}
	return cs, nil
}

func (ms *memorystore) Sources() ([]string, error) {
	var ss []string
	for s := range ms.sources {
		ss = append(ss, s)
	}
	return ss, nil
}

func (ms *memorystore) Close() error {
	return nil
}
