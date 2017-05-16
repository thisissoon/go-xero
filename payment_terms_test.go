package xero

import (
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaymentTerm_MarshalXML(t *testing.T) {
	type testcase struct {
		tname       string
		paymentTerm PaymentTerm
		expectedXML []byte
	}
	tt := []testcase{
		testcase{
			tname:       "DAYSAFTERBILLDATE",
			paymentTerm: PaymentTermDaysAfterBillDate,
			expectedXML: []byte("<Response><PaymentTerm>DAYSAFTERBILLDATE</PaymentTerm></Response>"),
		},
		testcase{
			tname:       "DAYSAFTERBILLMONTH",
			paymentTerm: PaymentTermSaysAfterBillMonth,
			expectedXML: []byte("<Response><PaymentTerm>DAYSAFTERBILLMONTH</PaymentTerm></Response>"),
		},
		testcase{
			tname:       "OFCURRENTMONTH",
			paymentTerm: PaymentTermOfCurrentMonth,
			expectedXML: []byte("<Response><PaymentTerm>OFCURRENTMONTH</PaymentTerm></Response>"),
		},
		testcase{
			tname:       "OFFOLLOWINGMONTH",
			paymentTerm: PaymentTermOfFollowingMonth,
			expectedXML: []byte("<Response><PaymentTerm>OFFOLLOWINGMONTH</PaymentTerm></Response>"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName     xml.Name    `xml:"Response"`
				PaymentTerm PaymentTerm `xml:"PaymentTerm"`
			}{
				PaymentTerm: tc.paymentTerm,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestPaymentTerm_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname               string
		decoder             func(t *testing.T) elementDecoder
		expectedPaymentTerm PaymentTerm
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
			tname: "invalid payment term",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString("foo")
					return nil
				}}
			},
			expectedErr: fmt.Errorf("unsupported payment term: %s", "foo"),
		},
		testcase{
			tname: "DAYSAFTERBILLDATE",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(paymentTermDaysAfterBillDate)
					return nil
				}}
			},
			expectedPaymentTerm: PaymentTermDaysAfterBillDate,
		},
		testcase{
			tname: "DAYSAFTERBILLMONTH",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(paymentTermSaysAfterBillMonth)
					return nil
				}}
			},
			expectedPaymentTerm: PaymentTermSaysAfterBillMonth,
		},
		testcase{
			tname: "OFCURRENTMONTH",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(paymentTermOfCurrentMonth)
					return nil
				}}
			},
			expectedPaymentTerm: PaymentTermOfCurrentMonth,
		},
		testcase{
			tname: "OFFOLLOWINGMONTH",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(paymentTermOfFollowingMonth)
					return nil
				}}
			},
			expectedPaymentTerm: PaymentTermOfFollowingMonth,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			pt := PaymentTerm{}
			err := pt.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedPaymentTerm, pt)
		})
	}
}

func TestPaymentTerm_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname               string
		xml                 []byte
		expectedPaymentTerm PaymentTerm
	}
	tt := []testcase{
		testcase{
			tname:               "DAYSAFTERBILLDATE",
			xml:                 []byte("<Response><PaymentTerm>DAYSAFTERBILLDATE</PaymentTerm></Response>"),
			expectedPaymentTerm: PaymentTermDaysAfterBillDate,
		},
		testcase{
			tname:               "DAYSAFTERBILLMONTH",
			xml:                 []byte("<Response><PaymentTerm>DAYSAFTERBILLMONTH</PaymentTerm></Response>"),
			expectedPaymentTerm: PaymentTermSaysAfterBillMonth,
		},
		testcase{
			tname:               "OFCURRENTMONTH",
			xml:                 []byte("<Response><PaymentTerm>OFCURRENTMONTH</PaymentTerm></Response>"),
			expectedPaymentTerm: PaymentTermOfCurrentMonth,
		},
		testcase{
			tname:               "OFFOLLOWINGMONTH",
			xml:                 []byte("<Response><PaymentTerm>OFFOLLOWINGMONTH</PaymentTerm></Response>"),
			expectedPaymentTerm: PaymentTermOfFollowingMonth,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName     xml.Name    `xml:"Response"`
				PaymentTerm PaymentTerm `xml:"PaymentTerm"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPaymentTerm, x.PaymentTerm)
		})
	}
}

func TestPaymentTerm_String(t *testing.T) {
	type testcase struct {
		tname          string
		paymentTerm    PaymentTerm
		expectedString string
	}
	tt := []testcase{
		testcase{
			tname:          "DAYSAFTERBILLDATE",
			paymentTerm:    PaymentTermDaysAfterBillDate,
			expectedString: "DAYSAFTERBILLDATE",
		},
		testcase{
			tname:          "DAYSAFTERBILLMONTH",
			paymentTerm:    PaymentTermSaysAfterBillMonth,
			expectedString: "DAYSAFTERBILLMONTH",
		},
		testcase{
			tname:          "OFCURRENTMONTH",
			paymentTerm:    PaymentTermOfCurrentMonth,
			expectedString: "OFCURRENTMONTH",
		},
		testcase{
			tname:          "OFFOLLOWINGMONTH",
			paymentTerm:    PaymentTermOfFollowingMonth,
			expectedString: "OFFOLLOWINGMONTH",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.paymentTerm.String())
		})
	}
}
