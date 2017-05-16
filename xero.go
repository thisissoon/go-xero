package xero

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
)

// New constructs a new Xero API Client
func New(authorizer Authorizer) *Client {
	return &Client{
		authorizer: authorizer,
		scheme:     "https",
		host:       "api.xero.com",
		root:       "/api.xro/2.0",
	}
}

// PrivateKey decodes a private key pem and returns a rsa.PrivateKey
// for Xero HTTP request Authorization
func PrivateKey(r io.Reader) (*rsa.PrivateKey, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	block, extra := pem.Decode(b)
	if block == nil || len(extra) > 0 {
		return nil, errors.New("failed to decode PEM")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
