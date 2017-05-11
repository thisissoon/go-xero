package oauth

import (
	"crypto/rsa"
	"net/http"
	"net/url"
	"path"

	"github.com/garyburd/go-oauth/oauth"
)

// The Authorizer type provides a wrapper to the github.com/garyburd/go-oauth/oauth
// that fullfills the xero.Authorizer interface for authorizing xero API requests
type Authorizer struct {
	token           string                // Xero API token
	pk              *rsa.PrivateKey       // Private Key
	scheme          string                // https
	host            string                // api.xero.com
	signatureMethod oauth.SignatureMethod // oauth.RSASHA1
}

// url constructs a valid xero api url
func (a Authorizer) url(parts ...string) string {
	u := url.URL{
		Scheme: a.scheme,
		Host:   a.host,
		Path:   path.Join(parts...),
	}
	return u.String()
}

// AuthorizeRequest fullfills the xero.Authorizer interface, authorizing a HTTP request
func (a *Authorizer) AuthorizeRequest(req *http.Request) error {
	client := &oauth.Client{
		Credentials: oauth.Credentials{
			Token: a.token,
		},
		TemporaryCredentialRequestURI: a.url("oauth", "RequestToken"),
		ResourceOwnerAuthorizationURI: a.url("oauth", "Authorize"),
		TokenRequestURI:               a.url("oauth", "AccessToken"),
		SignatureMethod:               a.signatureMethod,
		PrivateKey:                    a.pk,
	}
	return client.SetAuthorizationHeader(
		req.Header,
		&client.Credentials,
		req.Method,
		req.URL,
		req.URL.Query())
}

// New constructs a new Authorizer
func New(token string, pk *rsa.PrivateKey) *Authorizer {
	a := &Authorizer{
		token:           token,
		pk:              pk,
		scheme:          "https",
		host:            "api.xero.com",
		signatureMethod: oauth.RSASHA1,
	}
	return a
}
