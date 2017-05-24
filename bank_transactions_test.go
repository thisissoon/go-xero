package xero

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBankTransactionIterator_url(t *testing.T) {
	type testcase struct {
		tname       string
		page        int
		expectedURL string
	}
	tt := []testcase{
		testcase{
			tname:       "page 1",
			page:        1,
			expectedURL: "https://api.xero.com/api.xro/2.0/BankTransactions?page=1",
		},
		testcase{
			tname:       "page 2",
			page:        2,
			expectedURL: "https://api.xero.com/api.xro/2.0/BankTransactions?page=2",
		},
		testcase{
			tname:       "page 3",
			page:        3,
			expectedURL: "https://api.xero.com/api.xro/2.0/BankTransactions?page=3",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			i := BankTransactionIterator{tc.page, &Client{}, &url.URL{
				Scheme: "https",
				Host:   "api.xero.com",
				Path:   "/api.xro/2.0/BankTransactions",
			}}
			assert.Equal(t, tc.expectedURL, i.url())
		})
	}
}

func TestBankTransactionIterator_Next(t *testing.T) {
	type testcase struct {
		tname                string
		getter               testGetter
		ts                   func(t *testing.T) (*httptest.Server, *url.URL)
		expectedTransactions []BankTransaction
		expectedErr          error
	}
	tt := []testcase{
		testcase{
			tname: "request error",
			getter: testGetter(func(string, interface{}) error {
				return errors.New("request error")
			}),
			expectedErr: errors.New("request error"),
		},
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
						<BankTransactions></BankTransactions>
					</Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedErr: io.EOF,
		},
		testcase{
			tname: "returns transactions",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response>
						<BankTransactions>
							<BankTransaction>
								<Reference>Foo</Reference>
							</BankTransaction>
						</BankTransactions>
					</Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedTransactions: []BankTransaction{{Reference: "Foo"}},
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			var u = new(url.URL)
			if tc.ts != nil {
				ts, tsUrl := tc.ts(t)
				u = tsUrl
				defer ts.Close()
			}
			var c getter
			if tc.getter != nil {
				c = tc.getter
			} else {
				c = &Client{authorizer: new(testAuthorizer)}
			}
			i := BankTransactionIterator{1, c, u}
			_, items, err := i.Next()
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedTransactions, items)
		})
	}
}

func TestClient_BankTransactions(t *testing.T) {
	reqCount := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqCount += 1
		w.WriteHeader(http.StatusOK)
		switch reqCount {
		case 1:
			w.Write([]byte(`<Response>
				<BankTransactions>
					<BankTransaction>
						<Reference>Foo</Reference>
					</BankTransaction>
				</BankTransactions>
			</Response>`))

		case 2:
			w.Write([]byte(`<Response>
				<BankTransactions>
					<BankTransaction>
						<Reference>Bar</Reference>
					</BankTransaction>
				</BankTransactions>
			</Response>`))
		default:
			w.Write([]byte(`<Response><BankTransactions></BankTransactions></Response>`))
		}
	}))
	u, err := url.Parse(ts.URL)
	assert.NoError(t, err)
	c := &Client{
		authorizer: new(testAuthorizer),
		scheme:     u.Scheme,
		host:       u.Host,
		root:       u.Path,
	}
	x := 1
	receivedTrans := make(map[int][]BankTransaction)
	for i, items, err := c.BankTransactions(); err != io.EOF; i, items, err = i.Next() {
		receivedTrans[x] = items
		x++
	}
	assert.Equal(t, 3, x)
	assert.Equal(t, map[int][]BankTransaction{
		1: []BankTransaction{{Reference: "Foo"}},
		2: []BankTransaction{{Reference: "Bar"}},
	}, receivedTrans)
}

func TestClient_BankTransaction(t *testing.T) {
	type testcase struct {
		tname         string
		ts            func(t *testing.T) (*httptest.Server, *url.URL)
		expectedTrans BankTransaction
		expectedErr   error
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
			tname: "0 transactions",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response><BankTransactions></BankTransactions></Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedErr: fmt.Errorf("transaction %s not found", "foo"),
		},
		testcase{
			tname: "transaction returned",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response>
						<BankTransactions>
							<BankTransaction>
								<Reference>Dwack</Reference>
							</BankTransaction>
						</BankTransactions>
					</Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedTrans: BankTransaction{Reference: "Dwack"},
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
			trans, err := c.BankTransaction("foo")
			assert.Equal(t, tc.expectedTrans, trans)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestBankTransactionType_MarshalXML(t *testing.T) {
	type testcase struct {
		tname               string
		bankTransactionType BankTransactionType
		expectedXML         []byte
	}
	tt := []testcase{
		testcase{
			tname:               bankTransTypeReceive,
			bankTransactionType: BankTransTypeReceive,
			expectedXML:         []byte("<Response><Type>RECEIVE</Type></Response>"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name            `xml:"Response"`
				Type    BankTransactionType `xml:"Type"`
			}{
				Type: tc.bankTransactionType,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestBankTransactionType_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname        string
		decoder      func(t *testing.T) elementDecoder
		expectedType BankTransactionType
		expectedErr  error
	}
	tt := []testcase{
		testcase{
			tname: "decoder error",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					return errors.New("decoder error")
				}}
			},
			expectedErr: errors.New("decoder error"),
		},
		testcase{
			tname: "invalid account class",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString("foo")
					return nil
				}}
			},
			expectedErr: fmt.Errorf("unsupported bank transaction type: %s", "foo"),
		},
		testcase{
			tname: "RECEIVE",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(bankTransTypeReceive)
					return nil
				}}
			},
			expectedType: BankTransTypeReceive,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			a := BankTransactionType{}
			err := a.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedType, a)
		})
	}
}

func TestBankTransactionType_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname        string
		xml          []byte
		expectedType BankTransactionType
	}
	tt := []testcase{
		testcase{
			tname:        "RECEIVE",
			xml:          []byte("<Response><Type>RECEIVE</Type></Response>"),
			expectedType: BankTransTypeReceive,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name            `xml:"Response"`
				Type    BankTransactionType `xml:"Type"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedType, x.Type)
		})
	}
}

func TestBankTransactionType_String(t *testing.T) {
	type testcase struct {
		tname               string
		bankTransactionType BankTransactionType
		expectedString      string
	}
	tt := []testcase{
		testcase{
			tname:               "RECEIVE",
			bankTransactionType: BankTransTypeReceive,
			expectedString:      "RECEIVE",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.bankTransactionType.String())
		})
	}
}

func TestBankTransactionStatus_MarshalXML(t *testing.T) {
	type testcase struct {
		tname       string
		status      BankTransactionStatus
		expectedXML []byte
	}
	tt := []testcase{
		testcase{
			tname:       "AUTHORISED",
			status:      BankTransStatusAuth,
			expectedXML: []byte("<Response><Status>AUTHORISED</Status></Response>"),
		},
		testcase{
			tname:       "DELETED",
			status:      BankTransStatusDel,
			expectedXML: []byte("<Response><Status>DELETED</Status></Response>"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name              `xml:"Response"`
				Status  BankTransactionStatus `xml:"Status"`
			}{
				Status: tc.status,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestBankTransactionStatus_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname          string
		decoder        func(t *testing.T) elementDecoder
		expectedStatus BankTransactionStatus
		expectedErr    error
	}
	tt := []testcase{
		testcase{
			tname: "decoder error",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					return errors.New("decoder error")
				}}
			},
			expectedErr: errors.New("decoder error"),
		},
		testcase{
			tname: "invalid transaction status",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString("foo")
					return nil
				}}
			},
			expectedErr: fmt.Errorf("unsupported bank transaction status: %s", "foo"),
		},
		testcase{
			tname: "AUTHORISED",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(bankTransStatusAuth)
					return nil
				}}
			},
			expectedStatus: BankTransStatusAuth,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			a := BankTransactionStatus{}
			err := a.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedStatus, a)
		})
	}
}

func TestBankTransactionStatus_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname          string
		xml            []byte
		expectedStatus BankTransactionStatus
	}
	tt := []testcase{
		testcase{
			tname:          "AUTHORISED",
			xml:            []byte("<Response><Status>AUTHORISED</Status></Response>"),
			expectedStatus: BankTransStatusAuth,
		},
		testcase{
			tname:          "DELETED",
			xml:            []byte("<Response><Status>DELETED</Status></Response>"),
			expectedStatus: BankTransStatusDel,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name              `xml:"Response"`
				Status  BankTransactionStatus `xml:"Status"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, x.Status)
		})
	}
}

func TestBankTransactionStatus_String(t *testing.T) {
	type testcase struct {
		tname          string
		status         BankTransactionStatus
		expectedString string
	}
	tt := []testcase{
		testcase{
			tname:          "AUTHORISED",
			status:         BankTransStatusAuth,
			expectedString: "AUTHORISED",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.status.String())
		})
	}
}

//
func TestLineAmountType_MarshalXML(t *testing.T) {
	type testcase struct {
		tname          string
		lineAmountType LineAmountType
		expectedXML    []byte
	}
	tt := []testcase{
		testcase{
			tname:          "Exclusive",
			lineAmountType: LineAmountTypeExc,
			expectedXML:    []byte("<Response><Type>Exclusive</Type></Response>"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name       `xml:"Response"`
				Type    LineAmountType `xml:"Type"`
			}{
				Type: tc.lineAmountType,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestLineAmountType_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname        string
		decoder      func(t *testing.T) elementDecoder
		expectedType LineAmountType
		expectedErr  error
	}
	tt := []testcase{
		testcase{
			tname: "decoder error",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					return errors.New("decoder error")
				}}
			},
			expectedErr: errors.New("decoder error"),
		},
		testcase{
			tname: "invalid line amount type",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString("foo")
					return nil
				}}
			},
			expectedErr: fmt.Errorf("unsupported line amount type: %s", "foo"),
		},
		testcase{
			tname: "Exclusive",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(lineAmountTypeExc)
					return nil
				}}
			},
			expectedType: LineAmountTypeExc,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			a := LineAmountType{}
			err := a.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedType, a)
		})
	}
}

func TestLineAmountType_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname        string
		xml          []byte
		expectedType LineAmountType
	}
	tt := []testcase{
		testcase{
			tname:        "Exclusive",
			xml:          []byte("<Response><Type>Exclusive</Type></Response>"),
			expectedType: LineAmountTypeExc,
		},
		testcase{
			tname:        "Inclusive",
			xml:          []byte("<Response><Type>Inclusive</Type></Response>"),
			expectedType: LineAmountTypeInc,
		},
		testcase{
			tname:        "NoTax",
			xml:          []byte("<Response><Type>NoTax</Type></Response>"),
			expectedType: LineAmountTypeNoTax,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name       `xml:"Response"`
				Type    LineAmountType `xml:"Type"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedType, x.Type)
		})
	}
}

func TestLineAmountType_String(t *testing.T) {
	type testcase struct {
		tname          string
		lineAmountType LineAmountType
		expectedString string
	}
	tt := []testcase{
		testcase{
			tname:          "Exclusive",
			lineAmountType: LineAmountTypeExc,
			expectedString: "Exclusive",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.lineAmountType.String())
		})
	}
}
