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
type (
	Encoder interface {
		Encode(v interface{}) error
	}
	Decoder interface {
		Decode(v interface{}) error
	}
)

// Func types for consutructing new Encoders and Decoders
type (
	NewDecoderFn func(r io.Reader) Decoder
	NewEncoderFn func(w io.Writer) Encoder
)

// Default Decoder is JSON
var defaultNewDecoder = func(r io.Reader) Decoder {
	return xml.NewDecoder(r)
}

// Set the global NewDecoder to use default new decoder
var NewDecoder NewDecoderFn = defaultNewDecoder

// Default Encoder is JSON
var defaultNewEncoder = func(w io.Writer) Encoder {
	return xml.NewEncoder(w)
}

// Set the global NewEncoder to use default new encoder
var NewEncoder NewEncoderFn = defaultNewEncoder
