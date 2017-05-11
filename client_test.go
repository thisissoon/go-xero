package xero

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testAuthorizer struct {
	err error
}

func (t *testAuthorizer) AuthorizeRequest(req *http.Request) error {
	return t.err
}

func TestClient_do(t *testing.T) {
	type testcase struct {
		tname          string
		method         string
		url            string
		authorizer     Authorizer
		client         *http.Client
		ts             *httptest.Server
		expectedStatus int
		expectedError  error
	}
	tt := []testcase{
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
			ts: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			})),
			expectedStatus: http.StatusOK,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			url := tc.url
			if tc.ts != nil {
				url = tc.ts.URL
				defer tc.ts.Close()
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
