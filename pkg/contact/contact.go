package contact

import (
	"github.com/pkg/errors"
	"strings"
)

// Contact contains fields for a contact
type Contact struct {
	ID     string
	Name   string
	Email  string
	Mobile string
}

var fields = map[string]func(c *Contact, val string) error{
	"id":            setID,
	"name":          setName,
	"email":         setEmail,
	"mobile_number": setMobile,
}

func (c *Contact) setFieldValue(name string, val string) error {
	f, ok := fields[name]
	if !ok {
		// ignore extra fields that we dont know about
		return nil
	}
	return f(c, val)
}

func setID(c *Contact, val string) error {
	c.ID = val
	return nil
}

func setName(c *Contact, val string) error {
	c.Name = val
	return nil
}

func setEmail(c *Contact, val string) error {
	c.Email = val
	return nil
}

func setMobile(c *Contact, val string) error {
	// put phone number into standard international format
	val = strings.Replace(val, " ", "", -1)
	val = strings.Replace(val, "(", "", -1)
	val = strings.Replace(val, ")", "", -1)
	if val[0] == '0' {
		val = strings.TrimLeft(val, "0")
		val = "+44" + val
	}
	c.Mobile = val
	return nil
}

func (c *Contact) validate() error {
	if c.ID == "" {
		return errors.New("id is not set")
	}
	if c.Name == "" {
		return errors.New("name is not set")
	}
	if c.Email == "" {
		return errors.New("email is not set")
	}
	if c.Mobile == "" {
		return errors.New("mobile is not set")
	}
	return nil
}
