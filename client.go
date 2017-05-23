package xero

import (
	"bytes"
	"encoding/xml"
	"fmt"
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

// An Exceotion is returned by the API for none 200 responses
// This may contain extra information about the error such as in a 400
// bad request, others may not contain any other information
// This type stores the origional HTTP request and response status code
// as well as the API exception returned by Xero
type Error struct {
	Status       int
	Request      *http.Request
	APIException APIException
}

// Error returns the string representation of the Error
func (e Error) Error() string {
	return fmt.Sprintf(
		"Xero API Exception: [%d] %d: %s",
		e.Status,
		e.APIException.ErrorNumber,
		e.APIException.Message)
}

// An APIException is returned when the API errors
//   <ApiException>
//     <ErrorNumber>10</ErrorNumber>
//     <Type>ValidationException</Type>
//     <Message>A validation exception occurred</Message>
//     <Elements>
//       <DataContractBase xsi:type="Invoice">
//         <ValidationErrors>
//           <ValidationError>
//             <Message>Email address must be valid.</Message>
//           </ValidationError>
//         </ValidationErrors>
//      </DataContractBase>
//     </Elements>
//   </ApiException>
type APIException struct {
	ErrorNumber      int              `xml:"ErrorNumber"`
	Type             string           `xml:"Type"`
	Message          string           `xml:"Message"`
	DataContractBase DataContractBase `xml:"DataContractBase"`
}

// DataContactBase holds the type the API exception was for
// and any ValidationError's that occured that need to be corrected
type DataContractBase struct {
	Type             string            `xml:"xsi:type,attr"`
	ValidationErrors []ValidationError `xml:"ValidationErrors>ValidationError,omitempty"`
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

// doEncode encodes the encoder into a http body and makes the request to the API
// the response body is not processed and is automatically closed
func (c *Client) doEncode(method, urlStr string, enc Encoder) error {
	var body = new(bytes.Buffer)
	if err := enc.Encode(body); err != nil {
		return err
	}
	rsp, err := c.do(method, urlStr, body)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	return nil
}

// doEncodeDecode encodes the encoder into a http body and makes the request to the API and
// returns the responses of doDecode
func (c *Client) doEncodeDecode(method, urlStr string, enc Encoder, dst interface{}) error {
	var body = new(bytes.Buffer)
	if err := enc.Encode(body); err != nil {
		return err
	}
	return c.doDecode(method, urlStr, body, dst)
}

// get performs a HTTP GET request to the Xero API and decodes the response
// into a destination interface
func (c *Client) get(urlStr string, dst interface{}) error {
	return c.doDecode(http.MethodGet, urlStr, nil, dst)
}

// post performs a HTTP POST request to the Xero API and decodes the response
// into a destination interface
func (c *Client) post(urlStr string, enc Encoder, dst interface{}) error {
	return c.doEncodeDecode(http.MethodPost, urlStr, enc, dst)
}

// put performs a HTTP PUT request to the Xero API and decodes the response
// into a destination interface
func (c *Client) put(urlStr string, enc Encoder, dst interface{}) error {
	return c.doEncodeDecode(http.MethodPut, urlStr, enc, dst)
}

// Get sends a HTTP GET request for the given URL, no body is sent
func (c *Client) Get(urlStr string) (*http.Response, error) {
	return c.do(http.MethodGet, urlStr, nil)
}

// Post sends a HTTP POST request for the given url and request body and
// returns the raw HTTP response
func (c *Client) Post(urlStr string, body io.Reader) (*http.Response, error) {
	return c.do(http.MethodPost, urlStr, body)
}

// Post sends a HTTP PUT request for the given url and request body returning
// the raw HTTP response
func (c *Client) Put(urlStr string, body io.Reader) (*http.Response, error) {
	return c.do(http.MethodPut, urlStr, body)
}

// Use Create to send PUT requests to the xero API, encoding the request data
// into XML and decoding the response XML into the destination interface
func (c *Client) Create(ep Endpoint, enc Encoder, dst interface{}) error {
	return c.put(c.url(ep).String(), enc, dst)
}

// Use CreateUpdate to send POST requests to the xero API, encoding the request data
// into XML and decoding the response XML into the destination interface
func (c *Client) CreateUpdate(ep Endpoint, enc Encoder, dst interface{}) error {
	return c.post(c.url(ep).String(), enc, dst)
}
