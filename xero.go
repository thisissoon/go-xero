package xero

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

// New constructs a new Xero API Client
func New(authorizer Authorizer) *Client {
	return &Client{
		authorizer: authorizer,
	}
}

// ReadPrivateKey reads a pem file, decodes it and returns a rsa.PrivateKey
// for Xero HTTP request Authorization
func ReadPrivateKey(path string) (*rsa.PrivateKey, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, extra := pem.Decode(b)
	if block == nil || len(extra) > 0 {
		return nil, errors.New("failed to decode PEM")
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
