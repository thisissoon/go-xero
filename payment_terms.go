package xero

import (
	"encoding/xml"
	"fmt"
)

// Xero payment term string values:
// https://developer.xero.com/documentation/api/types#PaymentTerms
const (
	paymentTermDaysAfterBillDate  = "DAYSAFTERBILLDATE"
	paymentTermSaysAfterBillMonth = "DAYSAFTERBILLMONTH"
	paymentTermOfCurrentMonth     = "OFCURRENTMONTH"
	paymentTermOfFollowingMonth   = "OFFOLLOWINGMONTH"
)

// Xero Payment Term types
var (
	PaymentTermDaysAfterBillDate  = PaymentTerm{paymentTermDaysAfterBillDate}
	PaymentTermSaysAfterBillMonth = PaymentTerm{paymentTermSaysAfterBillMonth}
	PaymentTermOfCurrentMonth     = PaymentTerm{paymentTermOfCurrentMonth}
	PaymentTermOfFollowingMonth   = PaymentTerm{paymentTermOfFollowingMonth}
)

// Xero PaymentTerm type
type PaymentTerm struct {
	value string
}

// String implements the Stringer interface returning the string representation
// of the PaymentTerm
func (pt *PaymentTerm) String() string {
	return pt.value
}

// MarshalXML marshals a PaymentTerm into valid XML for Xero
func (pt *PaymentTerm) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(pt.value, start)
}

// unmarshalXML handles converting raw Xero Payment Term XML data into valid Payment Term
func (pt *PaymentTerm) unmarshalXML(decoder elementDecoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}
	switch value {
	case paymentTermDaysAfterBillDate:
		*pt = PaymentTermDaysAfterBillDate
	case paymentTermSaysAfterBillMonth:
		*pt = PaymentTermSaysAfterBillMonth
	case paymentTermOfCurrentMonth:
		*pt = PaymentTermOfCurrentMonth
	case paymentTermOfFollowingMonth:
		*pt = PaymentTermOfFollowingMonth
	default:
		return fmt.Errorf("unsupported payment term: %s", value)
	}
	return nil
}

// UnmarshalXML handles converting raw Xero Payment Term XML data into valid Payment Term
func (pt *PaymentTerm) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return pt.unmarshalXML(decoder, start)
}
