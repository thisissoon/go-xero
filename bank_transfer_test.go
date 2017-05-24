package xero

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_BankTransfers(t *testing.T) {
	type testcase struct {
		tname             string
		ts                func(t *testing.T) (*httptest.Server, *url.URL)
		expectedTransfers []BankTransfer
		expectedErr       error
	}
	tt := []testcase{
		testcase{
			tname: "bad xml",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`</uwotm8>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedErr:       &xml.SyntaxError{Msg: "unexpected end element </uwotm8>", Line: 1},
			expectedTransfers: []BankTransfer{},
		},
		testcase{
			tname: "transfers returned",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response>
						<BankTransfers>
							<BankTransfer>
								<Amount>20.00</Amount>
							</BankTransfer>
							<BankTransfer>
								<Amount>20.00</Amount>
							</BankTransfer>
						</BankTransfers>
					</Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedTransfers: []BankTransfer{
				{Amount: 20},
				{Amount: 20},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			ts, u := tc.ts(t)
			defer ts.Close()
			c := &Client{
				authorizer: new(testAuthorizer),
				scheme:     u.Scheme,
				host:       u.Host,
				root:       u.Path,
			}
			transfers, err := c.BankTransfers()
			assert.Equal(t, tc.expectedTransfers, transfers)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestClient_BankTransfer(t *testing.T) {
	type testcase struct {
		tname            string
		ts               func(t *testing.T) (*httptest.Server, *url.URL)
		expectedTransfer BankTransfer
		expectedErr      error
	}
	tt := []testcase{
		testcase{
			tname: "bad xml",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`</uwotm8>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedErr: &xml.SyntaxError{Msg: "unexpected end element </uwotm8>", Line: 1},
		},
		testcase{
			tname: "0 transfers",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response><BankTransfers></BankTransfers></Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedErr: fmt.Errorf("transfer %s not found", "foo"),
		},
		testcase{
			tname: "transfer returned",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response>
						<BankTransfers>
							<BankTransfer>
								<Amount>20.00</Amount>
							</BankTransfer>
						</BankTransfers>
					</Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedTransfer: BankTransfer{Amount: 20},
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			ts, u := tc.ts(t)
			defer ts.Close()
			c := &Client{
				authorizer: new(testAuthorizer),
				scheme:     u.Scheme,
				host:       u.Host,
				root:       u.Path,
			}
			transfer, err := c.BankTransfer("foo")
			assert.Equal(t, tc.expectedTransfer, transfer)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
