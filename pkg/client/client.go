package client

import (
	"context"
	pb "github.com/the4thamigo-uk/ingestor/pkg/protocol"
	"google.golang.org/grpc"
	"io"
)

// Client manages a single connection to a remote server.
type Client struct {
	cn *grpc.ClientConn
}

// ReadHandler is a callback function used to iterate over the data on the remote server.
type ReadHandler func(c Contact) error

// New creates a new Client instance and connects to the remote address.
func New(addr string) (*Client, error) {
	cn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{
		cn: cn,
	}, nil
}

// Close initiates a shutdown of the connection.
// Following this, it is the caller's responsibility to wait for any of the other Client methods to complete.
func (c *Client) Close() error {
	if c != nil {
		return c.cn.Close()
	}
	return nil
}

// ListSources returns the set of ids of the data sources available on the remote server.
func (c *Client) ListSources() ([]string, error) {
	cl := pb.NewIngestorClient(c.cn)
	rsp, err := cl.ListSources(context.Background(), &pb.ListRequest{})
	if err != nil {
		return nil, err
	}
	return toIds(rsp.Sources), nil
}

// ReadSource reads the data from the remote data source given by the provided id.
// The client receives an invocation of h for each data item found.
func (c *Client) ReadSource(id string, h ReadHandler) error {
	cl := pb.NewIngestorClient(c.cn)
	req := pb.ReadRequest{Id: id}
	st, err := cl.ReadSource(context.Background(), &req)
	if err != nil {
		return err
	}
	for {
		rsp, err := st.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		err = h(Contact{rsp})
		if err != nil {
			return err
		}
	}
	return nil
}

func toIds(ss []*pb.Source) []string {
	var ids []string
	for _, s := range ss {
		ids = append(ids, s.Id)
	}
	return ids
}
