package xero

import (
	"encoding/xml"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUTCDate_MarshalXML(t *testing.T) {
	now := time.Now().UTC()
	type testcase struct {
		tname       string
		utcDate     UTCDate
		expectedXML []byte
	}
	tt := []testcase{
		testcase{
			tname:       "marshal xml",
			utcDate:     UTCDate{now},
			expectedXML: []byte(fmt.Sprintf("<Response><Date>%s</Date></Response>", now.Format(utcDateLayout))),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name `xml:"Response"`
				Date    UTCDate  `xml:"Date"`
			}{
				Date: tc.utcDate,
			}
			b, err := xml.Marshal(&x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestUTCDate_unmarshalXML(t *testing.T) {
	type testcase struct {
		tname           string
		decoder         func(t *testing.T) elementDecoder
		expectedUTCDate UTCDate
		expectedErr     error
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
			tname: "time parse error",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					v = "foo" // Invalid format
					return nil
				}}
			},
			expectedErr: &time.ParseError{Layout: utcDateLayout, Value: "", LayoutElem: "2006", ValueElem: "", Message: ""},
		},
		testcase{
			tname: "ok",
			decoder: func(t *testing.T) elementDecoder {
				return &testDecoder{t: t, fn: func(t *testing.T, v interface{}, s *xml.StartElement) error {
					val := reflect.ValueOf(v).Elem()
					val.SetString("2009-05-14T01:44:26.747")
					return nil
				}}
			},
			expectedUTCDate: UTCDate{time.Date(2009, 5, 14, 01, 44, 26, 747000000, time.UTC)},
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			d := UTCDate{}
			err := d.unmarshalXML(tc.decoder(t), xml.StartElement{})
			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expectedUTCDate, d)
		})
	}
}

func TestUTCDate_UnmarshalXML(t *testing.T) {
	type testcase struct {
		tname           string
		xml             []byte
		expectedUTCDate UTCDate
	}
	tt := []testcase{
		testcase{
			tname:           "unmarshal xml",
			xml:             []byte("<Response><Date>2009-05-14T01:44:26.747</Date></Response>"),
			expectedUTCDate: UTCDate{time.Date(2009, 5, 14, 01, 44, 26, 747000000, time.UTC)},
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			x := struct {
				XMLName xml.Name `xml:"Response"`
				Date    UTCDate  `xml:"Date"`
			}{}
			err := xml.Unmarshal(tc.xml, &x)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedUTCDate, x.Date)
		})
	}
}

func TestUTCDate_Time(t *testing.T) {
	dt := time.Date(2009, 5, 14, 01, 44, 26, 747000000, time.UTC)
	type testcase struct {
		tname        string
		utcDate      UTCDate
		expectedTime time.Time
	}
	tt := []testcase{
		testcase{
			tname:        "returns time",
			utcDate:      UTCDate{dt},
			expectedTime: dt,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			assert.Equal(t, tc.expectedTime, tc.utcDate.Time())
		})
	}
}

func TestNewUTCDate(t *testing.T) {
	now := time.Now().UTC()
	assert.Equal(t, UTCDate{now}, NewUTCDate(now))
}
