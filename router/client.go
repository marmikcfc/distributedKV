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

func NewClient(addr string, timeout time.Duration) (*Client, error) {
	connection, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return nil, err
	}
	var r *Router
	r =New()
	var ok bool
	r.AddStore(address, &ok)
	return &Client{connection: rpc.NewClient(connection), Route:r}, nil
}

func (c *Client) Get(key string) (*StoreItem, error) {
	var item *StoreItem
	err := c.connection.Call("Router.Get", key, &item)
	return item, err
}

func (c *Client) Put(item *StoreItem) (bool, error) {
	var added bool
	println("inside Rpuer Client Put ")
	//var ok bool
	//added, err := rt.Put(item, &ok)
	err := c.connection.Call("Router.Put", item, &added)
	return added, err
}

func (c *Client) Delete(key string) (bool, error) {
	var deleted bool
	err := c.connection.Call("Router.Delete", key, &deleted)
	return deleted, err
}

func (c *Client) Clear() (bool, error) {
	var cleared bool
	err := c.connection.Call("Router.Clear", true, &cleared)
	return cleared, err
}

func (c *Client) AddStore(addr string) (bool, error) {
	var ok bool
	err := c.connection.Call("Router.AddStore", addr, &ok)
	return ok, err
}
