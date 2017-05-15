package xero

import (
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhoneType_MarshalXML(t *testing.T) {
	type testcase struct {
		tname       string
		phoneType   PhoneType
		expectedXML []byte
	}
	tt := []testcase{
		testcase{
			tname:       "DEFAILT",
			phoneType:   PhoneTypeDefault,
			expectedXML: []byte("<Response><PhoneType>DEFAULT</PhoneType></Response>"),
		},
		testcase{
			tname:       "DDI",
			phoneType:   PhoneTypeDDI,
			expectedXML: []byte("<Response><PhoneType>DDI</PhoneType></Response>"),
		},
		testcase{
			tname:       "MOBILE",
			phoneType:   PhoneTypeMobile,
			expectedXML: []byte("<Response><PhoneType>MOBILE</PhoneType></Response>"),
		},
		testcase{
			tname:       "FAX",
			phoneType:   PhoneTypeFax,
			expectedXML: []byte("<Response><PhoneType>FAX</PhoneType></Response>"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName   xml.Name  `xml:"Response"`
				PhoneType PhoneType `xml:"PhoneType"`
			}{
				PhoneType: tc.phoneType,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestPhoneType_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname             string
		decoder           func(t *testing.T) decoder
		expectedPhoneType PhoneType
		expectedErr       error
	}
	tt := []testcase{
		testcase{
			tname: "decoder error",
			decoder: func(t *testing.T) decoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					return errors.New("decoder error")
				}}
			},
			expectedErr: errors.New("decoder error"),
		},
		testcase{
			tname: "invalid phone type",
			decoder: func(t *testing.T) decoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString("foo")
					return nil
				}}
			},
			expectedErr: fmt.Errorf("unsupported phone type: %s", "foo"),
		},
		testcase{
			tname: "DEFAULT",
			decoder: func(t *testing.T) decoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(phoneTypeDefault)
					return nil
				}}
			},
			expectedPhoneType: PhoneTypeDefault,
		},
		testcase{
			tname: "DDI",
			decoder: func(t *testing.T) decoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(phoneTypeDDI)
					return nil
				}}
			},
			expectedPhoneType: PhoneTypeDDI,
		},
		testcase{
			tname: "MOBILE",
			decoder: func(t *testing.T) decoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(phoneTypeMobile)
					return nil
				}}
			},
			expectedPhoneType: PhoneTypeMobile,
		},
		testcase{
			tname: "FAX",
			decoder: func(t *testing.T) decoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(phoneTypeFax)
					return nil
				}}
			},
			expectedPhoneType: PhoneTypeFax,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			pt := PhoneType{}
			err := pt.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedPhoneType, pt)
		})
	}
}

func TestPhoneType_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname             string
		xml               []byte
		expectedPhoneType PhoneType
	}
	tt := []testcase{
		testcase{
			tname:             "DEFAULT",
			xml:               []byte("<Response><PhoneType>DEFAULT</PhoneType></Response>"),
			expectedPhoneType: PhoneTypeDefault,
		},
		testcase{
			tname:             "DDI",
			xml:               []byte("<Response><PhoneType>DDI</PhoneType></Response>"),
			expectedPhoneType: PhoneTypeDDI,
		},
		testcase{
			tname:             "MOBILE",
			xml:               []byte("<Response><PhoneType>MOBILE</PhoneType></Response>"),
			expectedPhoneType: PhoneTypeMobile,
		},
		testcase{
			tname:             "FAX",
			xml:               []byte("<Response><PhoneType>FAX</PhoneType></Response>"),
			expectedPhoneType: PhoneTypeFax,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName   xml.Name  `xml:"Response"`
				PhoneType PhoneType `xml:"PhoneType"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedPhoneType, x.PhoneType)
		})
	}
}

func TestPhoneType_String(t *testing.T) {
	type testcase struct {
		tname          string
		phoneType      PhoneType
		expectedString string
	}
	tt := []testcase{
		testcase{
			tname:          "DEFAULT",
			phoneType:      PhoneTypeDefault,
			expectedString: "DEFAULT",
		},
		testcase{
			tname:          "DDI",
			phoneType:      PhoneTypeDDI,
			expectedString: "DDI",
		},
		testcase{
			tname:          "MOBILE",
			phoneType:      PhoneTypeMobile,
			expectedString: "MOBILE",
		},
		testcase{
			tname:          "FAX",
			phoneType:      PhoneTypeFax,
			expectedString: "FAX",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.phoneType.String())
		})
	}
}
