package xero

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type errReader struct {
	err error
}

func (r errReader) Read([]byte) (int, error) {
	return 0, r.err
}

func TestPrivateKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatal(err)
	}
	var buff bytes.Buffer
	wrtr := bufio.NewWriter(&buff)
	if err := pem.Encode(wrtr, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}); err != nil {
		t.Fatal(err)
	}
	wrtr.Flush()
	parsedPk, err := x509.ParsePKCS1PrivateKey(x509.MarshalPKCS1PrivateKey(key))
	if err != nil {
		t.Fatal(err)
	}
	type testcase struct {
		tname              string
		r                  io.Reader
		expectedPrivateKey *rsa.PrivateKey
		expectedError      error
	}
	tt := []testcase{
		testcase{
			"read failure",
			errReader{errors.New("read error")},
			nil,
			errors.New("read error"),
		},
		testcase{
			"invalid pem",
			bytes.NewReader([]byte("not a valid pk")),
			nil,
			errors.New("failed to decode PEM"),
		},
		testcase{
			"valid pem",
			bytes.NewReader(buff.Bytes()),
			parsedPk,
			nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			pk, err := PrivateKey(tc.r)
			assert.Equal(t, tc.expectedPrivateKey, pk)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}

type fakeAuthorizer struct {
	err error
}

func (a fakeAuthorizer) AuthorizeRequest(*http.Request) error {
	return a.err
}

func TestNew(t *testing.T) {
	type testcase struct {
		tname          string
		authorizer     Authorizer
		expectedClient *Client
	}
	tt := []testcase{
		testcase{
			"constructs new client",
			fakeAuthorizer{},
			&Client{
				authorizer: fakeAuthorizer{},
				scheme:     "https",
				host:       "api.xero.com",
				root:       "/api.xro/2.0",
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.tname, func(t *testing.T) {
			client := New(tc.authorizer)
			assert.Equal(t, tc.expectedClient, client)
		})
	}
}
