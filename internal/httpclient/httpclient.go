package httpclient

import (
	"net/http"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type Client struct {
	client *http.Client
}

func (c *Client) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

func New(client *http.Client) HttpClient {
	return &Client{client}
}
