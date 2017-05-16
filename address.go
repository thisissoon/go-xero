package xero

import (
	"encoding/xml"
	"fmt"
)

// Predefined address types from Xero
// https://developer.xero.com/documentation/api/types#AddressTypes
const (
	addressTypePOBox    = "POBOX"
	addressTypeStreet   = "STREET"
	addressTypeDelivery = "DELIVERY"
)

// Xero Address types
var (
	AddressTypePOBox    = AddressType{addressTypePOBox}
	AddressTypeStreet   = AddressType{addressTypeStreet}
	AddressTypeDelivery = AddressType{addressTypeDelivery}
)

// The AddressType type defines the speicific address types within Xero:
// - POBOX
// - STREET
// - DELIVERY
type AddressType struct {
	value string
}

// String implements the Stringer interface returning the string representation
// of the AddressType
func (a AddressType) String() string {
	return a.value
}

// MarshalXML marshals a AddressType into valid XML for Xero
func (a *AddressType) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(a.value, start)
}

// unmarshalXML handles converting raw Xero Payment Term XML data into valid Payment Term
func (a *AddressType) unmarshalXML(decoder elementDecoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}
	switch value {
	case addressTypePOBox:
		*a = AddressTypePOBox
	case addressTypeStreet:
		*a = AddressTypeStreet
	case addressTypeDelivery:
		*a = AddressTypeDelivery
	default:
		return fmt.Errorf("unsupported address type: %s", value)
	}
	return nil
}

// UnmarshalXML handles converting raw Xero Payment Term XML data into valid Payment Term
func (a *AddressType) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return a.unmarshalXML(decoder, start)
}

// A Address type holds data for individual Xero addresses
type Address struct {
	AddressType  AddressType `xml:"AddressType,omitempty"`
	AddressLine1 string      `xml:"AddressLine1,omitempty"`
	AddressLine2 string      `xml:"AddressLine2,omitempty"`
	AddressLine3 string      `xml:"AddressLine3,omitempty"`
	AddressLine4 string      `xml:"AddressLine4,omitempty"`
	City         string      `xml:"City,omitempty"`
	Region       string      `xml:"Region,omitempty"`
	PostalCode   string      `xml:"PostalCode,omitempty"`
	Country      string      `xml:"Country,omitempty"`
	AttentionTo  string      `xml:"AttentionTo,omitempty"`
}
