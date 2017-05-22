package xero

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Accounts(t *testing.T) {
	type testcase struct {
		tname            string
		ts               func(t *testing.T) (*httptest.Server, *url.URL)
		expectedAccounts []Account
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
			expectedErr:      &xml.SyntaxError{Msg: "unexpected end element </uwotm8>", Line: 1},
			expectedAccounts: []Account{},
		},
		testcase{
			tname: "accounts returned",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response>
						<Accounts>
							<Account>
								<Name>Foo</Name>
							</Account>
							<Account>
								<Name>Bar</Name>
							</Account>
						</Accounts>
					</Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedAccounts: []Account{
				{Name: "Foo"},
				{Name: "Bar"},
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
			accounts, err := c.Accounts()
			assert.Equal(t, tc.expectedAccounts, accounts)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestClient_Account(t *testing.T) {
	type testcase struct {
		tname           string
		ts              func(t *testing.T) (*httptest.Server, *url.URL)
		expectedAccount Account
		expectedErr     error
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
			tname: "0 accounts",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response><Accounts></Accounts></Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedErr: fmt.Errorf("account %s not found", "foo"),
		},
		testcase{
			tname: "account returned",
			ts: func(t *testing.T) (*httptest.Server, *url.URL) {
				ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`<Response>
						<Accounts>
							<Account>
								<Name>Dwack</Name>
							</Account>
						</Accounts>
					</Response>`))
				}))
				u, err := url.Parse(ts.URL)
				assert.NoError(t, err)
				return ts, u
			},
			expectedAccount: Account{Name: "Dwack"},
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
			account, err := c.Account("foo")
			assert.Equal(t, tc.expectedAccount, account)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestAccountClass_MarshalXML(t *testing.T) {
	type testcase struct {
		tname        string
		accountClass AccountClass
		expectedXML  []byte
	}
	tt := []testcase{
		testcase{
			tname:        "ASSET",
			accountClass: AccountClassAsset,
			expectedXML:  []byte("<Response><Class>ASSET</Class></Response>"),
		},
		testcase{
			tname:        "EQUITY",
			accountClass: AccountClassEquity,
			expectedXML:  []byte("<Response><Class>EQUITY</Class></Response>"),
		},
		testcase{
			tname:        "EXPENSE",
			accountClass: AccountClassExp,
			expectedXML:  []byte("<Response><Class>EXPENSE</Class></Response>"),
		},
		testcase{
			tname:        "LIABILITY",
			accountClass: AccountClassLiab,
			expectedXML:  []byte("<Response><Class>LIABILITY</Class></Response>"),
		},
		testcase{
			tname:        "REVENUE",
			accountClass: AccountClassRevenue,
			expectedXML:  []byte("<Response><Class>REVENUE</Class></Response>"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name     `xml:"Response"`
				Class   AccountClass `xml:"Class"`
			}{
				Class: tc.accountClass,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestAccountClass_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname                string
		decoder              func(t *testing.T) elementDecoder
		expectedAccountClass AccountClass
		expectedErr          error
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
			expectedErr: fmt.Errorf("unsupported account class: %s", "foo"),
		},
		testcase{
			tname: "ASSET",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(accountClassAsset)
					return nil
				}}
			},
			expectedAccountClass: AccountClassAsset,
		},
		testcase{
			tname: "EQUITY",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(accountClassEquity)
					return nil
				}}
			},
			expectedAccountClass: AccountClassEquity,
		},
		testcase{
			tname: "EXPENSE",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(accountClassExp)
					return nil
				}}
			},
			expectedAccountClass: AccountClassExp,
		},
		testcase{
			tname: "LIABILITY",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(accountClassLiab)
					return nil
				}}
			},
			expectedAccountClass: AccountClassLiab,
		},
		testcase{
			tname: "REVENUE",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(accountClassRevenue)
					return nil
				}}
			},
			expectedAccountClass: AccountClassRevenue,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			a := AccountClass{}
			err := a.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedAccountClass, a)
		})
	}
}

func TestAccountClass_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname                string
		xml                  []byte
		expectedAccountClass AccountClass
	}
	tt := []testcase{
		testcase{
			tname:                "ASSET",
			xml:                  []byte("<Response><Class>ASSET</Class></Response>"),
			expectedAccountClass: AccountClassAsset,
		},
		testcase{
			tname:                "EQUITY",
			xml:                  []byte("<Response><Class>EQUITY</Class></Response>"),
			expectedAccountClass: AccountClassEquity,
		},
		testcase{
			tname:                "EXPENSE",
			xml:                  []byte("<Response><Class>EXPENSE</Class></Response>"),
			expectedAccountClass: AccountClassExp,
		},
		testcase{
			tname:                "LIABILITY",
			xml:                  []byte("<Response><Class>LIABILITY</Class></Response>"),
			expectedAccountClass: AccountClassLiab,
		},
		testcase{
			tname:                "REVENUE",
			xml:                  []byte("<Response><Class>REVENUE</Class></Response>"),
			expectedAccountClass: AccountClassRevenue,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name     `xml:"Response"`
				Class   AccountClass `xml:"Class"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedAccountClass, x.Class)
		})
	}
}

func TestAccountClass_String(t *testing.T) {
	type testcase struct {
		tname          string
		accountClass   AccountClass
		expectedString string
	}
	tt := []testcase{
		testcase{
			tname:          "ASSET",
			accountClass:   AccountClassAsset,
			expectedString: "ASSET",
		},
		testcase{
			tname:          "EQUITY",
			accountClass:   AccountClassEquity,
			expectedString: "EQUITY",
		},
		testcase{
			tname:          "EXPENSE",
			accountClass:   AccountClassExp,
			expectedString: "EXPENSE",
		},
		testcase{
			tname:          "LIABILITY",
			accountClass:   AccountClassLiab,
			expectedString: "LIABILITY",
		},
		testcase{
			tname:          "REVENUE",
			accountClass:   AccountClassRevenue,
			expectedString: "REVENUE",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.accountClass.String())
		})
	}
}

func TestAccountType_MarshalXML(t *testing.T) {
	type testcase struct {
		tname       string
		accountType AccountType
		expectedXML []byte
	}
	tt := []testcase{
		testcase{
			tname:       "BANK",
			accountType: AccountTypeBank,
			expectedXML: []byte("<Response><Type>BANK</Type></Response>"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name    `xml:"Response"`
				Type    AccountType `xml:"Type"`
			}{
				Type: tc.accountType,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestAccountType_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname               string
		decoder             func(t *testing.T) elementDecoder
		expectedAccountType AccountType
		expectedErr         error
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
			tname: "invalid account type",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString("foo")
					return nil
				}}
			},
			expectedErr: fmt.Errorf("unsupported account type: %s", "foo"),
		},
		testcase{
			tname: "BANK",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(accountTypeBank)
					return nil
				}}
			},
			expectedAccountType: AccountTypeBank,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			a := AccountType{}
			err := a.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedAccountType, a)
		})
	}
}

func TestAccountType_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname               string
		xml                 []byte
		expectedAccountType AccountType
	}
	tt := []testcase{
		testcase{
			tname:               "BANK",
			xml:                 []byte("<Response><Type>BANK</Type></Response>"),
			expectedAccountType: AccountTypeBank,
		},
		testcase{
			tname:               "CURRENT",
			xml:                 []byte("<Response><Type>CURRENT</Type></Response>"),
			expectedAccountType: AccountTypeCurrent,
		},
		testcase{
			tname:               "CURRLIAB",
			xml:                 []byte("<Response><Type>CURRLIAB</Type></Response>"),
			expectedAccountType: AccountTypeCurrLiab,
		},
		testcase{
			tname:               "DEPRECIATN",
			xml:                 []byte("<Response><Type>DEPRECIATN</Type></Response>"),
			expectedAccountType: AccountTypeDepreciatn,
		},
		testcase{
			tname:               "DIRECTCOSTS",
			xml:                 []byte("<Response><Type>DIRECTCOSTS</Type></Response>"),
			expectedAccountType: AccountTypeDC,
		},
		testcase{
			tname:               "EQUITY",
			xml:                 []byte("<Response><Type>EQUITY</Type></Response>"),
			expectedAccountType: AccountTypeEquity,
		},
		testcase{
			tname:               "EXPENSE",
			xml:                 []byte("<Response><Type>EXPENSE</Type></Response>"),
			expectedAccountType: AccountTypeExp,
		},
		testcase{
			tname:               "FIXED",
			xml:                 []byte("<Response><Type>FIXED</Type></Response>"),
			expectedAccountType: AccountTypeFixed,
		},
		testcase{
			tname:               "INVENTORY",
			xml:                 []byte("<Response><Type>INVENTORY</Type></Response>"),
			expectedAccountType: AccountTypeInventory,
		},
		testcase{
			tname:               "LIABILITY",
			xml:                 []byte("<Response><Type>LIABILITY</Type></Response>"),
			expectedAccountType: AccountTypeLiab,
		},
		testcase{
			tname:               "NONCURRENT",
			xml:                 []byte("<Response><Type>NONCURRENT</Type></Response>"),
			expectedAccountType: AccountTypeNonCurrent,
		},
		testcase{
			tname:               "OTHERINCOME",
			xml:                 []byte("<Response><Type>OTHERINCOME</Type></Response>"),
			expectedAccountType: AccountTypeOther,
		},
		testcase{
			tname:               "OVERHEADS",
			xml:                 []byte("<Response><Type>OVERHEADS</Type></Response>"),
			expectedAccountType: AccountTypeOverhead,
		},
		testcase{
			tname:               "PREPAYMENT",
			xml:                 []byte("<Response><Type>PREPAYMENT</Type></Response>"),
			expectedAccountType: AccountTypePrepay,
		},
		testcase{
			tname:               "REVENUE",
			xml:                 []byte("<Response><Type>REVENUE</Type></Response>"),
			expectedAccountType: AccountTypeRevenue,
		},
		testcase{
			tname:               "SALES",
			xml:                 []byte("<Response><Type>SALES</Type></Response>"),
			expectedAccountType: AccountTypeSale,
		},
		testcase{
			tname:               "TERMLIAB",
			xml:                 []byte("<Response><Type>TERMLIAB</Type></Response>"),
			expectedAccountType: AccountTypeTermLiab,
		},
		testcase{
			tname:               "PAYGLIABILITY",
			xml:                 []byte("<Response><Type>PAYGLIABILITY</Type></Response>"),
			expectedAccountType: AccountTypePAYG,
		},
		testcase{
			tname:               "SUPERANNUATIONEXPENSE",
			xml:                 []byte("<Response><Type>SUPERANNUATIONEXPENSE</Type></Response>"),
			expectedAccountType: AccountTypeSAExp,
		},
		testcase{
			tname:               "SUPERANNUATIONLIABILITY",
			xml:                 []byte("<Response><Type>SUPERANNUATIONLIABILITY</Type></Response>"),
			expectedAccountType: AccountTypeSALiab,
		},
		testcase{
			tname:               "WAGESEXPENSE",
			xml:                 []byte("<Response><Type>WAGESEXPENSE</Type></Response>"),
			expectedAccountType: AccountTypeWageExp,
		},
		testcase{
			tname:               "WAGESPAYABLELIABILITY",
			xml:                 []byte("<Response><Type>WAGESPAYABLELIABILITY</Type></Response>"),
			expectedAccountType: AccountTypeWageLiab,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name    `xml:"Response"`
				Type    AccountType `xml:"Type"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedAccountType, x.Type)
		})
	}
}

func TestAccountType_String(t *testing.T) {
	type testcase struct {
		tname          string
		accountType    AccountType
		expectedString string
	}
	tt := []testcase{
		testcase{
			tname:          "BANK",
			accountType:    AccountTypeBank,
			expectedString: "BANK",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.accountType.String())
		})
	}
}

func TestAccountStatus_MarshalXML(t *testing.T) {
	type testcase struct {
		tname         string
		accountStatus AccountStatus
		expectedXML   []byte
	}
	tt := []testcase{
		testcase{
			tname:         "BANK",
			accountStatus: AccountStatusActive,
			expectedXML:   []byte("<Response><Type>ACTIVE</Type></Response>"),
		},
		testcase{
			tname:         "BANK",
			accountStatus: AccountStatusArchive,
			expectedXML:   []byte("<Response><Type>ARCHIVED</Type></Response>"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name      `xml:"Response"`
				Type    AccountStatus `xml:"Type"`
			}{
				Type: tc.accountStatus,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestAccountStatus_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname                 string
		decoder               func(t *testing.T) elementDecoder
		expectedAccountStatus AccountStatus
		expectedErr           error
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
			tname: "invalid account status",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString("foo")
					return nil
				}}
			},
			expectedErr: fmt.Errorf("unsupported account status: %s", "foo"),
		},
		testcase{
			tname: "ACTIVE",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(accountStatusActive)
					return nil
				}}
			},
			expectedAccountStatus: AccountStatusActive,
		},
		testcase{
			tname: "ARCHIVED",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(accountStatusArchive)
					return nil
				}}
			},
			expectedAccountStatus: AccountStatusArchive,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			a := AccountStatus{}
			err := a.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedAccountStatus, a)
		})
	}
}

func TestAccountStatus_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname                 string
		xml                   []byte
		expectedAccountStatus AccountStatus
	}
	tt := []testcase{
		testcase{
			tname: "ACTIVE",
			xml:   []byte("<Response><Status>ACTIVE</Status></Response>"),
			expectedAccountStatus: AccountStatusActive,
		},
		testcase{
			tname: "ARCHIVED",
			xml:   []byte("<Response><Status>ARCHIVED</Status></Response>"),
			expectedAccountStatus: AccountStatusArchive,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name      `xml:"Response"`
				Status  AccountStatus `xml:"Status"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedAccountStatus, x.Status)
		})
	}
}

func TestAccountStatus_String(t *testing.T) {
	type testcase struct {
		tname          string
		accountStatus  AccountStatus
		expectedString string
	}
	tt := []testcase{
		testcase{
			tname:          "ACTIVE",
			accountStatus:  AccountStatusActive,
			expectedString: "ACTIVE",
		},
		testcase{
			tname:          "ARCHIVED",
			accountStatus:  AccountStatusArchive,
			expectedString: "ARCHIVED",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.accountStatus.String())
		})
	}
}

//
func TestBankAccountType_MarshalXML(t *testing.T) {
	type testcase struct {
		tname           string
		bankAccountType BankAccountType
		expectedXML     []byte
	}
	tt := []testcase{
		testcase{
			tname:           "BANK",
			bankAccountType: BankAccountTypeBank,
			expectedXML:     []byte("<Response><BankAccountType>BANK</BankAccountType></Response>"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName         xml.Name        `xml:"Response"`
				BankAccountType BankAccountType `xml:"BankAccountType"`
			}{
				BankAccountType: tc.bankAccountType,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestBankAccountType_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname                   string
		decoder                 func(t *testing.T) elementDecoder
		expectedBankAccountType BankAccountType
		expectedErr             error
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
			tname: "invalid bank account type",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString("foo")
					return nil
				}}
			},
			expectedErr: fmt.Errorf("unsupported bank account type: %s", "foo"),
		},
		testcase{
			tname: "BANK",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(bankAccountTypeBank)
					return nil
				}}
			},
			expectedBankAccountType: BankAccountTypeBank,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			a := BankAccountType{}
			err := a.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedBankAccountType, a)
		})
	}
}

func TestBankAccountType_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname                   string
		xml                     []byte
		expectedBankAccountType BankAccountType
	}
	tt := []testcase{
		testcase{
			tname: "BANK",
			xml:   []byte("<Response><BankAccountType>BANK</BankAccountType></Response>"),
			expectedBankAccountType: BankAccountTypeBank,
		},
		testcase{
			tname: "CREDITCARD",
			xml:   []byte("<Response><BankAccountType>CREDITCARD</BankAccountType></Response>"),
			expectedBankAccountType: BankAccountTypeCC,
		},
		testcase{
			tname: "PAYPAL",
			xml:   []byte("<Response><BankAccountType>PAYPAL</BankAccountType></Response>"),
			expectedBankAccountType: BankAccountTypePaypal,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName         xml.Name        `xml:"Response"`
				BankAccountType BankAccountType `xml:"BankAccountType"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedBankAccountType, x.BankAccountType)
		})
	}
}

func TestBankAccountType_String(t *testing.T) {
	type testcase struct {
		tname           string
		bankAccountType BankAccountType
		expectedString  string
	}
	tt := []testcase{
		testcase{
			tname:           "BANK",
			bankAccountType: BankAccountTypeBank,
			expectedString:  "BANK",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.bankAccountType.String())
		})
	}
}
