package httpclient

import "net/http"

type MockClient struct {
	client *http.Client
	url    string
}

func (c *MockClient) Get(url string) (*http.Response, error) {
	return c.client.Get(c.url)
}

func NewMock(client *http.Client, url string) HttpClient {
	return &MockClient{client, url}
}
