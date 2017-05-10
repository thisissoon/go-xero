package oauth_test

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"net/http"

	oauth "github.com/thisissoon/go-xero/authorizers/go-oauth"
)

func Example() {
	f, err := ioutil.ReadFile("/path/to/pem")
	if err != nil {
		log.Fatal(err)
	}
	block, extra := pem.Decode(f)
	if block == nil || len(extra) > 0 {
		log.Fatal("failed to decode PEM")
	}
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}
	authorizer := oauth.New("TOKEN", pk)
	req, err := http.NewRequest(
		http.MethodGet,
		"https://api.xero.com/api.xro/2.0/Invoices",
		nil)
	if err != nil {
		log.Fatal(err)
	}
	if err := authorizer.AuthorizeRequest(req); err != nil {
		log.Fatal(err)
	}
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()
	log.Println(rsp.Status)
}
