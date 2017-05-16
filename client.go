package xero

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

// Internal interface tyes implemented by the Client type
type (
	// The getter interface is implemented by the Client type and used internally
	getter interface {
		Get(string) (*http.Response, error)
	}
	// The poster interface is implemented by the Client type and used internally
	poster interface {
		Post(string, io.Reader) (*http.Response, error)
	}
	// The putter interface is implemented by the Client type and used internally
	putter interface {
		Put(string, io.Reader) (*http.Response, error)
	}
)

// The Authorizer interface defines  common interface for authorising Xero HTTP
// requests using oAuth. The AuthorizeRequest takes the HTTP request that requires
// authorization.
type Authorizer interface {
	AuthorizeRequest(request *http.Request) error
}

// The Response type defines the XML response body wrapper
//   <Response>
//       <Id>...</Id>
//       <Status>...</Status>
//       <ProviderName>...</ProviderName
//       <DateTimeUTC>...</DateTimeUTC>
//       ...
//   </Response>
type Response struct {
	XMLName      xml.Name  `xml:"Response"`
	Id           string    `xml:"Id"`
	Status       string    `xml:"Status"`
	ProviderName string    `xml:"ProviderName"`
	DateTimeUTC  time.Time `xml:"DateTimeUTC"`
}

// A Client is a Xero API client. It provides methods for calling Xero API endpoints.
// This type should be constructed with the New() method.
type Client struct {
	authorizer Authorizer
	client     *http.Client

	scheme string // Xero API Protocol Scheme (https)
	host   string // Xero API Host (api.xero.com)
	root   string // Xero API Root (/api.xro/2.0)
}

// url constructs a valid Xero API url. The scheme, host and api root are
// automatically appended to the url path
func (c *Client) url(endpoint string) *url.URL {
	return &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path.Join(c.root, endpoint),
	}
}

// do calls the Xero API
func (c *Client) do(method, urlStr string, body io.Reader) (*http.Response, error) {
	switch method {
	case http.MethodPost, http.MethodPut:
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		v := u.Query()
		v.Set("SummarizeErrors", "false")
		u.RawQuery = v.Encode()
		urlStr = u.String()
	}
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
