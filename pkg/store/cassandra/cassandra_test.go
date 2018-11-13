package cassandra

import (
	_ "github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/require"
	_ "github.com/the4thamigo-uk/ingestor/pkg/store"
	_ "testing"
)

/*
func testContact() *store.Contact {
	return &store.Contact{
		Mobile: "123",
		Name:   "name",
		Email:  "person@domain.com",
	}
}

func TestCassandra_WriteContact(t *testing.T) {
	s, err := New("test", "127.0.0.1")
	require.Nil(t, err)
	err = s.CreateTables()
	require.Nil(t, err)
	c := testContact()
	err = s.PutContact(c)
	require.Nil(t, err)
	cs, err := s.Contacts()
	require.Len(t, cs, 1)
	assert.Equal(t, c, cs[0])
}

func TestCassandra_WriteSource(t *testing.T) {
	s, err := New("test", "127.0.0.1")
	require.Nil(t, err)
	err = s.CreateTables()
	require.Nil(t, err)
	err = s.PutSource("123")
	require.Nil(t, err)
	ids, err := s.Sources()
	require.Len(t, ids, 1)
	assert.Equal(t, "123", ids[0])
}
*/
