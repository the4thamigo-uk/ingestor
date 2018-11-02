package identity

import (
	"github.com/satori/go.uuid"
)

// New creates a new global identifier.
func New() (string, error) {
	// TODO: next release of satori/go.uuid will return an error
	//u, err := uuid.NewV4()
	//if err != nil {
	//	return "", err
	//}

	u := uuid.NewV4()
	return u.String(), nil
}
