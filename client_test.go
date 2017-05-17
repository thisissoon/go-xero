package xero

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testGetter func(string) (*http.Response, error)

func (fn testGetter) Get(urlStr string) (*http.Response, error) {
	return fn(urlStr)
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
