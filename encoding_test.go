package xero

import (
	"encoding/xml"
	"io"
	"testing"
)

type testEncoder struct {
	t  *testing.T
	fn func(t *testing.T, w io.Writer) error
}

func (t *testEncoder) Encode(w io.Writer) error {
	if t.fn != nil {
		return t.fn(t.t, w)
	}
	return nil
}

type testDecoder struct {
	t  *testing.T
	fn func(t *testing.T, v interface{}, s *xml.StartElement) error
}

func (t *testDecoder) DecodeElement(v interface{}, s *xml.StartElement) error {
	if t.fn != nil {
		return t.fn(t.t, v, s)
	}
	return nil
}
