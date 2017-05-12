package xero

import (
	"encoding/xml"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationStatus_MarshalXMLAttr(t *testing.T) {
	type testcase struct {
		tname         string
		status        ValidationStatus
		expectedXML   []byte
		expectedError error
	}
	tt := []testcase{
		testcase{
			tname:         "invalid value",
			status:        ValidationStatus{"foo"},
			expectedError: errors.New("invalid validation type: foo"),
		},
		testcase{
			tname:       "ok value",
			status:      ValidationStatusOK,
			expectedXML: []byte(`<Foo status="OK"></Foo>`),
		},
		testcase{
			tname:       "error value",
			status:      ValidationStatusError,
			expectedXML: []byte(`<Foo status="ERROR"></Foo>`),
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			body := struct {
				XMLName xml.Name         `xml:"Foo"`
				Status  ValidationStatus `xml:"status,attr"`
			}{
				Status: tc.status,
			}
			b, err := xml.Marshal(&body)
			assert.Equal(t, tc.expectedError, err)
			assert.Equal(t, tc.expectedXML, b)
		})
	}
}

func TestValidationStatus_UnmarshalXMLAttr(t *testing.T) {
	type testcase struct {
		tname          string
		xml            []byte
		expectedError  error
		expectedStatus ValidationStatus
	}
	tt := []testcase{
		testcase{
			tname:         "invalid value",
			xml:           []byte(`<foo status="BAR"></foo>`),
			expectedError: fmt.Errorf("unknown validation status: BAR"),
		},
		testcase{
			tname:          "ok value",
			xml:            []byte(`<foo status="OK"></foo>`),
			expectedStatus: ValidationStatusOK,
		},
		testcase{
			tname:          "error value",
			xml:            []byte(`<foo status="ERROR"></foo>`),
			expectedStatus: ValidationStatusError,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			body := struct {
				Status ValidationStatus `xml:"status,attr"`
			}{}
			err := xml.Unmarshal(tc.xml, &body)
			assert.Equal(t, tc.expectedStatus, body.Status)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

func TestValidation(t *testing.T) {
	raw := []byte(`<Response>
		<Foos>
			<Foo status="OK">
				<Name>Foo</Name>
			</Foo>
			<Foo status="ERROR">
				<Name>Bar</Name>
				<ValidationErrors>
					<ValidationError>
						<Message>you forgot something</Message>
					</ValidationError>
				</ValidationErrors>
			</Foo>
		</Foos>
	</Response>`)
	body := struct {
		XMLName xml.Name `xml:"Response"`
		Foos    []struct {
			ValidationErrors
			Name string `xml:"Name"`
		} `xml:"Foos>Foo"`
	}{}
	expected := struct {
		XMLName xml.Name `xml:"Response"`
		Foos    []struct {
			ValidationErrors
			Name string `xml:"Name"`
		} `xml:"Foos>Foo"`
	}{
		XMLName: xml.Name{Local: "Response"},
		Foos: []struct {
			ValidationErrors
			Name string `xml:"Name"`
		}{
			{
				ValidationErrors: ValidationErrors{
					Status: ValidationStatusOK,
				},
				Name: "Foo",
			},
			{
				ValidationErrors: ValidationErrors{
					Status: ValidationStatusError,
					Errors: []ValidationError{
						ValidationError{"you forgot something"},
					},
				},
				Name: "Bar",
			},
		},
	}
	assert.NoError(t, xml.Unmarshal(raw, &body))
	assert.Equal(t, expected, body)
}
