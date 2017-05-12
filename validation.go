package xero

import (
	"encoding/xml"
	"fmt"
)

// Standard XML validation status response values from Xero
const (
	validationStatusOK    = "OK"
	validationStatusError = "ERROR"
)

// Validation status values returned by Xero
var (
	ValidationStatusOK    = ValidationStatus{validationStatusOK}
	ValidationStatusError = ValidationStatus{validationStatusError}
)

// The elementDecoder interface is used when handling decoding custom types
// from xml strings to their actual types
type elementDecoder interface {
	DecodeElement(interface{}, *xml.StartElement) error
}

// The ValidationStatus type holds the validation status xml attribute, e.g:
//   <Response>
//       <Invoices>
//           <Invoice status="OK">
//               ...
//           </Invoice>
//           <Invoice status="ERROR">
//               ...
//           </Invoice>
//       </Invoices>
//   </Response>
type ValidationStatus struct {
	status string
}

// String implements the fmt.Stringer interface
func (v ValidationStatus) String() string {
	return v.status
}

// MarshalXMLAttr handles marshaling the validation status into an xml attribute
func (v ValidationStatus) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	switch v {
	case ValidationStatusOK, ValidationStatusError:
		return xml.Attr{name, v.String()}, nil
	default:
		return xml.Attr{}, fmt.Errorf("invalid validation type: %s", v.String())
	}
}

// UnmarshalXMLAttr handles unmarshaling the raw "status" xml attribute value into
// a ValidatuoinStatus type
func (v *ValidationStatus) UnmarshalXMLAttr(attr xml.Attr) error {
	switch attr.Value {
	case validationStatusOK:
		*v = ValidationStatusOK
	case validationStatusError:
		*v = ValidationStatusError
	default:
		return fmt.Errorf("unknown validation status: %s", attr.Value)
	}
	return nil
}

// The Validation type is used for validating PUT/POST requests
// to the Xero API. Each type, such as Invoice embeds this common
// validation type which can be checked against in the response
// from Xero for each type posted/put to the API
// See the validation.go example
type Validation struct {
	Status ValidationStatus  `xml:"status,attr,omitempty"`
	Errors []ValidationError `xml:"ValidationErrors>ValidationError,omitempty"`
}

// The ValidationError type holds a individual validation error
// message from the Xero API
type ValidationError struct {
	Message string `xml:"Message"`
}
