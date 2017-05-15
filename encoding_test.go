package xero

import (
	"encoding/xml"
	"testing"
)

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
