package xero

import (
	"encoding/xml"
	"io"
)

// elementDecoder defines an interface implementd by xml.Decoder
type elementDecoder interface {
	DecodeElement(v interface{}, start *xml.StartElement) error
}

// Interfaces encoders and decoders must implement
type Encoder interface {
	Encode(w io.Writer) error
}

// Convenience method for encoding xml into a io writer
func encode(w io.Writer, v interface{}) error {
	enc := xml.NewEncoder(w)
	return enc.Encode(v)
}
