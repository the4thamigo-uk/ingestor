package cassandra

import (
	"github.com/gocql/gocql"
	"github.com/the4thamigo-uk/ingestor/pkg/store"
	"time"
)

type Cassandra interface {
	store.Store
}

type cassandra struct {
	s  *gocql.Session
	ks string
}

func New(keySpace string, hosts ...string) (store.Store, error) {
	c := gocql.NewCluster(hosts...)
	c.Timeout = time.Second * 10
	c.Keyspace = keySpace

	s, err := c.CreateSession()
	if err != nil {
		return nil, err
	}
	return &cassandra{
		s:  s,
		ks: keySpace,
	}, nil
}

func (s *cassandra) PutContact(c *store.Contact) error {
	return s.PutContacts([]*store.Contact{c})
}

func (s *cassandra) PutSource(id string) error {
	return s.PutSources([]string{id})

}
func (s *cassandra) PutContacts(cs []*store.Contact) error {
	b := s.s.NewBatch(gocql.LoggedBatch)
	for _, c := range cs {
		// TODO cql injection issues?
		b.Query(`INSERT INTO contacts (mobile, name, email) VALUES (?,?,?)`, c.Mobile, c.Name, c.Email)
	}
	return s.s.ExecuteBatch(b)

}
func (s *cassandra) PutSources(ids []string) error {
	b := s.s.NewBatch(gocql.LoggedBatch)
	for _, id := range ids {
		b.Query(`INSERT INTO sources (source) VALUES (?)`, id)
	}
	return s.s.ExecuteBatch(b)
}
func (s *cassandra) Contacts() ([]*store.Contact, error) {
	q := s.s.Query(`SELECT mobile, name, email FROM contacts`)
	it := q.Iter()

	var cs []*store.Contact
	for {
		m := map[string]interface{}{}
		if !it.MapScan(m) {
			break
		}
		cs = append(cs, &store.Contact{
			Mobile: m["mobile"].(string),
			Name:   m["name"].(string),
			Email:  m["email"].(string),
		})
	}
	return cs, nil
}
func (s *cassandra) Sources() ([]string, error) {
	q := s.s.Query(`SELECT source FROM sources`)
	it := q.Iter()

	var ids []string
	for {
		m := map[string]interface{}{}
		if !it.MapScan(m) {
			break
		}
		ids = append(ids, m["source"].(string))
	}
	return ids, nil
}

func (s *cassandra) Close() error {
	s.s.Close()
	return nil
}
