package oauth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Simply calls the xero oauth endpoints with a fake private key and token
// whilst the Authorization header generated will not be valid we know that the
// flow is working
func TestAuthorizer(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		t.Fatal(err)
	}
	pk, err := x509.ParsePKCS1PrivateKey(x509.MarshalPKCS1PrivateKey(key))
	req, err := http.NewRequest(http.MethodGet, "https://api.xero.com/api.xro/2.0/Invoices", nil)
	if err != nil {
		t.Fatal(err)
	}
	a := New("token", pk)
	assert.NoError(t, a.AuthorizeRequest(req))
	assert.NotEqual(t, "", req.Header.Get("Authorization"))
}
