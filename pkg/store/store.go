package store

// Contact is a contact object to use with the store
type Contact struct {
	Name   string
	Email  string
	Mobile string
}

// Store is an interface used to persist Contact objects and Source identifiers
type Store interface {
	PutContact(cs *Contact) error    // store contact object in the store allowing for batching if required
	PutSource(id string) error       // store id for a source that has been processed
	PutContacts(cs []*Contact) error // store contact objects in the store allowing for batching if required
	PutSources(ids []string) error   // store ids for the sources that have been processed
	Contacts() ([]*Contact, error)   // get all contacts
	Sources() ([]string, error)      // get all sources
}
