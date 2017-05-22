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
		get(string, interface{}) error
	}
	// The poster interface is implemented by the Client type and used internally
	poster interface {
		post(string, io.Reader, interface{}) error
	}
	// The putter interface is implemented by the Client type and used internally
	putter interface {
		put(string, io.Reader, interface{}) error
	}
)

// The Endpoint type defines official Xero API endpoints
type Endpoint string

// String implements the Stringer interface
func (e Endpoint) String() string {
	return string(e)
}

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
func (c *Client) url(endpoint Endpoint, extra ...string) *url.URL {
	parts := append([]string{endpoint.String()}, extra...)
	return &url.URL{
		Scheme: c.scheme,
		Host:   c.host,
		Path:   path.Join(parts...),
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

// doDecode performs a HTTP request to the Xero API and automatically decodes
// the response into a destination interface
func (c *Client) doDecode(method, urlStr string, body io.Reader, dst interface{}) error {
	rsp, err := c.do(method, urlStr, body)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	decoder := xml.NewDecoder(rsp.Body)
	if err := decoder.Decode(dst); err != nil {
		return err
	}
	return nil
}

// get performs a HTTP GET request to the Xero API and decodes the response
// into a destination interface
func (c *Client) get(urlStr string, dst interface{}) error {
	return c.doDecode(http.MethodGet, urlStr, nil, dst)
}

// Get sends a HTTP GET request for the given URL, no body is sent
func (c *Client) Get(urlStr string) (*http.Response, error) {
	return c.do(http.MethodGet, urlStr, nil)
}

// post performs a HTTP POST request to the Xero API and decodes the response
// into a destination interface
func (c *Client) post(urlStr string, body io.Reader, dst interface{}) error {
	return c.doDecode(http.MethodPost, urlStr, body, dst)
}

// Post sends a HTTP POST request for the given url and request body
func (c *Client) Post(urlStr string, body io.Reader) (*http.Response, error) {
	return c.do(http.MethodPost, urlStr, body)
}

// put performs a HTTP PUT request to the Xero API and decodes the response
// into a destination interface
func (c *Client) put(urlStr string, body io.Reader, dst interface{}) error {
	return c.doDecode(http.MethodPut, urlStr, body, dst)
}

// Post sends a HTTP PUT request for the given url and request body
func (c *Client) Put(urlStr string, body io.Reader) (*http.Response, error) {
	return c.do(http.MethodPut, urlStr, body)
}
