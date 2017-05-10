package xero

import (
	"io"
	"net/http"
)

// The Authorizer interface defines  common interface for authorising Xero HTTP
// requests using oAuth. The AuthorizeRequest takes the HTTP request that requires
// authorization.
type Authorizer interface {
	AuthorizeRequest(request *http.Request) error
}

// A Client is a Xero API client. It provides methods for calling Xero API endpoints.
// This type should be constructed with the New() method.
type Client struct {
	authorizer Authorizer
	client     *http.Client
}

// do calls the Xero API
func (c *Client) do(method, urlStr string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/xml")
	if err := c.authorizer.AuthorizeRequest(req); err != nil {
		return nil, err
	}
	client := c.client
	if client == nil {
		client = http.DefaultClient
	}
	return client.Do(req)
}

// Get sends a HTTP GET request for the given URL, no body is sent
func (c *Client) Get(urlStr string) (*http.Response, error) {
	return c.do(http.MethodGet, urlStr, nil)
}

// Post sends a HTTP POST request for the given url and request body
func (c *Client) Post(urlStr string, body io.Reader) (*http.Response, error) {
	return c.do(http.MethodPost, urlStr, body)
}

// Post sends a HTTP PUT request for the given url and request body
func (c *Client) Put(urlStr string, body io.Reader) (*http.Response, error) {
	return c.do(http.MethodPut, urlStr, body)
}
