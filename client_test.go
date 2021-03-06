package xero

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testRoundTrip func(*http.Request) (*http.Response, error)

func (fn testRoundTrip) RoundTrip(r *http.Request) (*http.Response, error) {
	return fn(r)
}

type testGetter func(string, interface{}) error

func (fn testGetter) get(urlStr string, dst interface{}) error {
	return fn(urlStr, dst)
}

type testAuthorizer struct {
	err error
}

func (t *testAuthorizer) AuthorizeRequest(req *http.Request) error {
	return t.err
}

type tHTTPHandler struct {
	t       *testing.T
	handler func(*testing.T, http.ResponseWriter, *http.Request)
}

func (t *tHTTPHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.handler(t.t, w, r)
}

func TestClient_do(t *testing.T) {
	type testcase struct {
		tname          string
		method         string
		url            string
		authorizer     Authorizer
		client         *http.Client
		ts             func(t *testing.T) *httptest.Server
		expectedStatus int
		expectedError  error
	}
	tt := []testcase{
		testcase{
			tname:      "POST SummarizeErrors = false",
			method:     http.MethodPost,
			authorizer: &testAuthorizer{},
			ts: func(t *testing.T) *httptest.Server {
				handler := &tHTTPHandler{
					t: t,
					handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
						assert.Equal(t, http.MethodPost, r.Method)
						assert.Equal(t, "false", r.URL.Query().Get("SummarizeErrors"))
						w.WriteHeader(http.StatusOK)
					},
				}
				return httptest.NewServer(handler)
			},
			expectedStatus: http.StatusOK,
		},
		testcase{
			tname:      "PUT SummarizeErrors = false",
			method:     http.MethodPut,
			authorizer: &testAuthorizer{},
			ts: func(t *testing.T) *httptest.Server {
				handler := &tHTTPHandler{
					t: t,
					handler: func(t *testing.T, w http.ResponseWriter, r *http.Request) {
						assert.Equal(t, http.MethodPut, r.Method)
						assert.Equal(t, "false", r.URL.Query().Get("SummarizeErrors"))
						w.WriteHeader(http.StatusOK)
					},
				}
				return httptest.NewServer(handler)
			},
			expectedStatus: http.StatusOK,
		},
		testcase{
			tname:  "SummarizeErrors bad url",
			method: http.MethodPost,
			url:    "://invalid",
			expectedError: &url.Error{
				Op:  "parse",
				URL: "://invalid",
				Err: errors.New("missing protocol scheme"),
			},
		},
		testcase{
			tname:         "new http request error",
			method:        "bad method",
			expectedError: errors.New("net/http: invalid method \"bad method\""),
		},
		testcase{
			tname:         "authorizer error",
			authorizer:    &testAuthorizer{err: errors.New("authorizer error")},
			expectedError: errors.New("authorizer error"),
		},
		testcase{
			tname:      "ok",
			authorizer: &testAuthorizer{},
			ts: func(*testing.T) *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
			},
			expectedStatus: http.StatusOK,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			url := tc.url
			if tc.ts != nil {
				ts := tc.ts(t)
				url = ts.URL
				defer ts.Close()
			}
			client := &Client{
				authorizer: tc.authorizer,
				client:     tc.client,
			}
			rsp, err := client.do(tc.method, url, nil)
			assert.Equal(t, tc.expectedError, err)
			if rsp != nil {
				assert.Equal(t, tc.expectedStatus, rsp.StatusCode)
			}
		})
	}
}

func TestClient_doDecode(t *testing.T) {
	type testcase struct {
		tname         string
		client        *http.Client
		method        string
		urlStr        string
		body          io.Reader
		dst           Response
		expectedError error
		expectedDst   interface{}
	}
	tt := []testcase{
		testcase{
			tname: "request error",
			client: &http.Client{
				Transport: testRoundTrip(func(*http.Request) (*http.Response, error) {
					return nil, errors.New("request error")
				}),
			},
			expectedError: &url.Error{Op: "Get", URL: "", Err: errors.New("request error")},
			expectedDst:   Response{},
		},
		testcase{
			tname: "invalid xml",
			dst:   Response{},
			client: &http.Client{
				Transport: testRoundTrip(func(*http.Request) (*http.Response, error) {
					r := bytes.NewBuffer([]byte("</uwotm8>"))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(r),
					}, nil
				}),
			},
			expectedError: &xml.SyntaxError{Msg: "unexpected end element </uwotm8>", Line: 1},
			expectedDst:   Response{},
		},
		testcase{
			tname: "ok",
			dst:   Response{},
			client: &http.Client{
				Transport: testRoundTrip(func(*http.Request) (*http.Response, error) {
					r := bytes.NewBuffer([]byte(`<Response><ProviderName>Foo</ProviderName></Response>`))
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(r),
					}, nil
				}),
			},
			expectedError: nil,
			expectedDst:   Response{XMLName: xml.Name{Local: "Response"}, ProviderName: "Foo"},
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			client := &Client{authorizer: new(testAuthorizer), client: tc.client}
			err := client.doDecode(tc.method, tc.urlStr, tc.body, &tc.dst)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedDst, tc.dst)
		})
	}
}

func TestClient_doEncode(t *testing.T) {
	type testcase struct {
		name          string
		enc           func(t *testing.T) Encoder
		client        *http.Client
		expectedError error
	}
	tt := []testcase{
		testcase{
			name: "encode error",
			enc: func(t *testing.T) Encoder {
				return &testEncoder{t, func(t *testing.T, w io.Writer) error {
					return errors.New("encoding error")
				}}
			},
			expectedError: errors.New("encoding error"),
		},
		testcase{
			name: "encode error",
			enc: func(t *testing.T) Encoder {
				return &testEncoder{t, func(t *testing.T, w io.Writer) error {
					return nil
				}}
			},
			client: &http.Client{
				Transport: testRoundTrip(func(*http.Request) (*http.Response, error) {
					return nil, errors.New("request error")
				}),
			},
			expectedError: &url.Error{
				Op:  "Post",
				URL: "/?SummarizeErrors=false",
				Err: errors.New("request error"),
			},
		},
		testcase{
			name: "ok",
			enc: func(t *testing.T) Encoder {
				return &testEncoder{t, func(t *testing.T, w io.Writer) error {
					return nil
				}}
			},
			client: &http.Client{
				Transport: testRoundTrip(func(*http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte(""))),
					}, nil
				}),
			},
			expectedError: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			client := &Client{authorizer: new(testAuthorizer), client: tc.client}
			err := client.doEncode(http.MethodPost, "/", tc.enc(t))
			assert.Equal(t, tc.expectedError, err, "%s", err)
		})
	}
}

func TestCheckResponse(t *testing.T) {
	type testcase struct {
		tname            string
		rsp              *http.Response
		expectedError    error
		expectedResponse *http.Response
	}
	tt := []testcase{
		{
			tname: "200 OK",
			rsp: &http.Response{
				StatusCode: http.StatusOK,
			},
			expectedError: nil,
			expectedResponse: &http.Response{
				StatusCode: http.StatusOK,
			},
		},
		{
			tname: "400 Bad Request",
			rsp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body: ioutil.NopCloser(bytes.NewReader([]byte(`
					<ApiException>
						<ErrorNumber>10</ErrorNumber>
						<Type>ValidationException</Type>
						<Message>A validation exception occurred</Message>
						<Elements>
							<DataContractBase xsi:type="Invoice">
								<ValidationErrors>
								<ValidationError>
									<Message>Email address must be valid.</Message>
								</ValidationError>
							  </ValidationErrors>
						   </DataContractBase>
						</Elements>
					</ApiException>`))),
			},
			expectedError: APIException{
				ErrorNumber: 10,
				Type:        "ValidationException",
				Message:     "A validation exception occurred",
				Elements: []DataContractBase{
					{
						Type: "Invoice",
						ValidationErrors: []ValidationError{
							{
								Message: "Email address must be valid.",
							},
						},
					},
				},
			},
			expectedResponse: nil,
		},
		{
			tname: "503",
			rsp: &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Body:       ioutil.NopCloser(bytes.NewReader([]byte("The Organisation is offline"))),
				Request: &http.Request{
					Method: http.MethodGet,
					URL:    &url.URL{Path: "/foo"},
				},
			},
			expectedError: fmt.Errorf(
				"%d: %s (%s: %s) %s",
				http.StatusServiceUnavailable,
				http.StatusText(http.StatusServiceUnavailable),
				http.MethodGet,
				"/foo",
				[]byte("The Organisation is offline")),
			expectedResponse: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			rsp, err := checkResponse(tc.rsp)
			assert.Equal(t, tc.expectedResponse, rsp)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
