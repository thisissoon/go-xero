package xero

import (
	"encoding/xml"
	"fmt"
)

// Xero phone types as strings:
// https://developer.xero.com/documentation/api/types#PhoneTypes
const (
	phoneTypeDefault = "DEFAULT"
	phoneTypeDDI     = "DDI"
	phoneTypeMobile  = "MOBILE"
	phoneTypeFax     = "FAX"
)

// Xero Phone types
var (
	PhoneTypeDefault = PhoneType{phoneTypeDefault}
	PhoneTypeDDI     = PhoneType{phoneTypeDDI}
	PhoneTypeMobile  = PhoneType{phoneTypeMobile}
	PhoneTypeFax     = PhoneType{phoneTypeFax}
)

// The PhoneType used for storing a Phone records type as defined in Xero.
// - Default
// - DDI
// - MOBILE
// - FAX
type PhoneType struct {
	value string
}

// String implements the Stringer interface returning the string representation
// of the PhoneType
func (pt PhoneType) String() string {
	return pt.value
}

// MarshalXML marshals a PhoneType into valid XML for Xero
func (pt *PhoneType) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(pt.value, start)
}

// unmarshalXML handles converting raw Xero Payment Term XML data into valid Payment Term
func (pt *PhoneType) unmarshalXML(decoder decoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}
	switch value {
	case phoneTypeDefault:
		*pt = PhoneTypeDefault
	case phoneTypeDDI:
		*pt = PhoneTypeDDI
	case phoneTypeMobile:
		*pt = PhoneTypeMobile
	case phoneTypeFax:
		*pt = PhoneTypeFax
	default:
		return fmt.Errorf("unsupported phone type: %s", value)
	}
	return nil
}

// UnmarshalXML handles converting raw Xero Payment Term XML data into valid Payment Term
func (pt *PhoneType) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return pt.unmarshalXML(decoder, start)
}

// The Phone type holds data about phone numbers from Xero.
type Phone struct {
	PhoneType        PhoneType `xml:"PhoneType,omitempty"`
	PhoneNumber      string    `xml:"PhoneNumber,omitempty"`
	PhoneAreaCode    string    `xml:"PhoneAreaCode,omitempty"`
	PhoneCountryCode string    `xml:"PhoneCountryCode,omitempty"`
}
