package xero

import "encoding/xml"

type decoder interface {
	DecodeElement(v interface{}, start *xml.StartElement) error
}
