package contact

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"io"
	"os"
)

// Reader is an interface to extract a stream of Contact objects
type Reader interface {
	Read() (*Contact, error)
	ReadAll() ([]*Contact, error)
}

type reader struct {
	cr     *csv.Reader
	fields map[int]string
}

// NewReader creates a stream of Contact objects from an io.Reader providing CSV data.
// The first row of the CSV data must contain the field names and there must be at a
// minimum the following field names in the data :
//
// id,name,email,mobile_number
func NewReader(r io.Reader) (Reader, error) {
	cr := csv.NewReader(r)
	fields, err := readHeaderFields(cr)
	if err != nil {
		return nil, err
	}
	return &reader{
		cr:     cr,
		fields: fields,
	}, nil
}

// NewReaderFromFile is a helper function to avoid the caller having to create an io.Reader
// for the given filename.
func NewReaderFromFile(path string) (Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return NewReader(f)
}

// Read reads one Contact object from the CSV data.
func (r *reader) Read() (*Contact, error) {
	rec, err := r.cr.Read()
	if err != nil {
		return nil, err
	}
	return r.toContact(rec)
}

// Read reads all Contact objects from the CSV data.
func (r *reader) ReadAll() ([]*Contact, error) {
	recs, err := r.cr.ReadAll()
	if err != nil {
		return nil, err
	}
	var cs []*Contact
	for _, rec := range recs {
		c, err := r.toContact(rec)
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
	}
	return cs, nil
}

func (r *reader) toContact(rec []string) (*Contact, error) {
	var c Contact
	for i, value := range rec {
		field, ok := r.fields[i]
		if !ok {
			continue
		}
		err := c.setFieldValue(field, value)
		if err != nil {
			return nil, err
		}
	}
	err := c.validate()
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func readHeaderFields(cr *csv.Reader) (map[int]string, error) {
	rec, err := cr.Read()
	if err != nil {
		return nil, errors.Wrap(err, "No CSV header data found")
	}
	fields := map[int]string{}
	for i, field := range rec {
		fields[i] = field
	}
	return fields, nil
}
