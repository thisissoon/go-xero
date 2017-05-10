package oauth

import (
	"crypto/rsa"
	"net/http"
	"net/url"
	"path"

	"github.com/garyburd/go-oauth/oauth"
)

// The Option type provides a function signature for configuring the Authorizer
type Option func(a *Authorizer)

// WithScheme configures the Xero API request protocol (e.g https)
func WithScheme(scheme string) Option {
	return func(a *Authorizer) {
		a.scheme = scheme
	}
}

// WithHost configures the Xero API request host (e.g api.xero.com)
func WithHost(host string) Option {
	return func(a *Authorizer) {
		a.host = host
	}
}

// WithSignatureMethod configures the Xero API oAuth signature method (e.g RSASHA1)
func WithSignatureMethod(method oauth.SignatureMethod) Option {
	return func(a *Authorizer) {
		a.signatureMethod = method
	}
}

// The Authorizer type provides a wrapper to the github.com/garyburd/go-oauth/oauth
// that fullfills the xero.Authorizer interface for authorizing xero API requests
type Authorizer struct {
	// Required
	token string
	pk    *rsa.PrivateKey
	// Configurable
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
func New(token string, pk *rsa.PrivateKey, opts ...Option) *Authorizer {
	a := &Authorizer{
		// Required
		token: token,
		pk:    pk,
		// Configurable
		scheme:          "https",
		host:            "api.xero.com",
		signatureMethod: oauth.RSASHA1,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}
