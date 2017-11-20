package router

import (
	"net"
	"net/rpc"
	"time"
)

type (
	Client struct {
		connection *rpc.Client
		Route *Router
	}
)

const(
	address     = "localhost:50051"
)

func NewClient(addr string, timeout time.Duration, number_of_nodes int) (*Client, error) {
	connection, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	var r *Router
	r =New()
	for i := 0; i < number_of_nodes; i++ {
	var ok bool
	r.AddStore(address, &ok)		
	}

	return &Client{connection: rpc.NewClient(connection), Route:r}, nil
}
