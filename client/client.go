package client

import (
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	url   string
	token string
}

// New creates a new Client
func New(url, token string) *Client {
	return &Client{url, token}
}

func (c *Client) request(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.url+url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	return http.DefaultClient.Do(req)
}
