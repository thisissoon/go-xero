package xero

import (
	"encoding/xml"
	"time"
)

// Xero date time layouts
const (
	utcDateLayout = "2006-01-02T15:04:05.000"
)

// The UTCDate type is used for storing Xero UTC date field values
type UTCDate struct {
	time time.Time
}

// MarshalXML is handles converting UTCDate time to Xero XML format
func (d UTCDate) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	format := d.time.Format(utcDateLayout)
	return encoder.EncodeElement(format, start)
}

// unmarshalXML handles converting raw Xero UTC Date XML data into valid time
func (d *UTCDate) unmarshalXML(decoder decoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}
	t, err := time.Parse(utcDateLayout, value)
	if err != nil {
		return err
	}
	*d = UTCDate{t.UTC()}
	return nil
}

// UnmarshalXML handles converting raw Xero UTC Date XML data into valid time
func (d *UTCDate) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return d.unmarshalXML(decoder, start)
}

// Time returns the UTC Date as a time.Time
func (d UTCDate) Time() time.Time {
	return d.time
}

// NewUTCDate constructs a new UTCDate from a time.Time
func NewUTCDate(t time.Time) UTCDate {
	return UTCDate{t}
}
