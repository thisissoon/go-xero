package xero

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContactIterator_url(t *testing.T) {
	type testcase struct {
		tname       string
		page        int
		expectedURL string
	}
	tt := []testcase{
		testcase{
			tname:       "page 1",
			page:        1,
			expectedURL: "https://api.xero.com/api.xro/2.0/Contacts?page=1",
		},
		testcase{
			tname:       "page 2",
			page:        2,
			expectedURL: "https://api.xero.com/api.xro/2.0/Contacts?page=2",
		},
		testcase{
			tname:       "page 3",
			page:        3,
			expectedURL: "https://api.xero.com/api.xro/2.0/Contacts?page=3",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			i := ContactIterator{tc.page, &Client{}, &url.URL{
				Scheme: "https",
				Host:   "api.xero.com",
				Path:   "/api.xro/2.0/Contacts",
			}}
			assert.Equal(t, tc.expectedURL, i.url())
		})
	}
}

func TestContactIterator_Next(t *testing.T) {
	type testcase struct {
		tname            string
		ts               func(t *testing.T) (*httptest.Server, *url.URL)
		expectedContacts []Contact
		expectedErr      error
	}
	tt := []testcase{
		testcase{
			tname: "bad xml",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`</uwot>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedErr: &xml.SyntaxError{
				Msg:  "unexpected end element </uwot>",
				Line: 1,
			},
		},
		testcase{
			tname: "no contacts EOF",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response>
						<Contacts></Contacts>
					</Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedErr: io.EOF,
		},
		testcase{
			tname: "returns contacts",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response>
						<Contacts>
							<Contact>
								<Name>Foo</Name>
							</Contact>
						</Contacts>
					</Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedContacts: []Contact{{Name: "Foo"}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			ts, tsUrl := tc.ts(t)
			defer ts.Close()
			c := &Client{
				authorizer: new(testAuthorizer),
				scheme:     tsUrl.Scheme,
				host:       tsUrl.Host,
				root:       tsUrl.Path,
			}
			i := ContactIterator{1, c, c.url("/")}
			_, contacts, err := i.Next()
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedContacts, contacts)
		})
	}
}
