package endpoint

import (

	"net/http"
)

type (
	Client struct {
		connection *http.Client
		url        string
		Response   *http.Response
	}
	Key struct {
		Namespace string
		Group     string
		Id        string
	}
)

func NewClient(addr string) (*Client, error) {
	c := &http.Client{}
	url := addr + "/"
	return &Client{connection: c, url: url}, nil
}
