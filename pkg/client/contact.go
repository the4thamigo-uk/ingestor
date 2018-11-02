package client

import (
	pb "github.com/the4thamigo-uk/ingestor/pkg/protocol"
)

// Contact wraps the protocol response to avoid caller's to import grpc (which is an implementation detail).
type Contact struct {
	rsp *pb.ReadResponse
}

// Name returns the name of the contact
func (c *Contact) Name() string {
	return c.rsp.Name
}

// Email returns the email of the contact
func (c *Contact) Email() string {
	return c.rsp.Email
}

// Mobile returns the mobile of the contact
func (c *Contact) Mobile() string {
	return c.rsp.Mobile
}
