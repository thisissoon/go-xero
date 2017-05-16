package xero

import (
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressType_MarshalXML(t *testing.T) {
	type testcase struct {
		tname       string
		addressType AddressType
		expectedXML []byte
	}
	tt := []testcase{
		testcase{
			tname:       "POBOX",
			addressType: AddressTypePOBox,
			expectedXML: []byte("<Response><AddressType>POBOX</AddressType></Response>"),
		},
		testcase{
			tname:       "STREET",
			addressType: AddressTypeStreet,
			expectedXML: []byte("<Response><AddressType>STREET</AddressType></Response>"),
		},
		testcase{
			tname:       "DELIVERY",
			addressType: AddressTypeDelivery,
			expectedXML: []byte("<Response><AddressType>DELIVERY</AddressType></Response>"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName     xml.Name    `xml:"Response"`
				AddressType AddressType `xml:"AddressType"`
			}{
				AddressType: tc.addressType,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestAddressType_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname               string
		decoder             func(t *testing.T) elementDecoder
		expectedAddressType AddressType
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
			tname: "invalid address type",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString("foo")
					return nil
				}}
			},
			expectedErr: fmt.Errorf("unsupported address type: %s", "foo"),
		},
		testcase{
			tname: "POBOX",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(addressTypePOBox)
					return nil
				}}
			},
			expectedAddressType: AddressTypePOBox,
		},
		testcase{
			tname: "STREET",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(addressTypeStreet)
					return nil
				}}
			},
			expectedAddressType: AddressTypeStreet,
		},
		testcase{
			tname: "DELIVERY",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString(addressTypeDelivery)
					return nil
				}}
			},
			expectedAddressType: AddressTypeDelivery,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			a := AddressType{}
			err := a.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedAddressType, a)
		})
	}
}

func TestAddressType_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname               string
		xml                 []byte
		expectedAddressType AddressType
	}
	tt := []testcase{
		testcase{
			tname:               "POBOX",
			xml:                 []byte("<Response><AddressType>POBOX</AddressType></Response>"),
			expectedAddressType: AddressTypePOBox,
		},
		testcase{
			tname:               "STREET",
			xml:                 []byte("<Response><AddressType>STREET</AddressType></Response>"),
			expectedAddressType: AddressTypeStreet,
		},
		testcase{
			tname:               "DELIVERY",
			xml:                 []byte("<Response><AddressType>DELIVERY</AddressType></Response>"),
			expectedAddressType: AddressTypeDelivery,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName     xml.Name    `xml:"Response"`
				AddressType AddressType `xml:"AddressType"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedAddressType, x.AddressType)
		})
	}
}

func TestAddressType_String(t *testing.T) {
	type testcase struct {
		tname          string
		addressType    AddressType
		expectedString string
	}
	tt := []testcase{
		testcase{
			tname:          "POBOX",
			addressType:    AddressTypePOBox,
			expectedString: "POBOX",
		},
		testcase{
			tname:          "STREET",
			addressType:    AddressTypeStreet,
			expectedString: "STREET",
		},
		testcase{
			tname:          "DELIVERY",
			addressType:    AddressTypeDelivery,
			expectedString: "DELIVERY",
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedString, tc.addressType.String())
		})
	}
}
