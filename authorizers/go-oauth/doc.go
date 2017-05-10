/*
The package implements the xero.Authorizer interface using the https://github.com/garyburd/go-oauth library.

The example below shows the raw usage of this authorizer outside of the xero package.

	package main

	import (
		"crypto/rsa"
		"crypto/x509"
		"encoding/pem"
		"errors"
		"flag"
		"io/ioutil"
		"log"
		"net/http"

		oauth "github.com/thisissoon/go-xero/authorizers/go-oauth"
	)

	// Decode private key pem file
	func pk(loc string) (*rsa.PrivateKey, error) {
		f, err := ioutil.ReadFile(loc)
		if err != nil {
			return nil, err
		}
		block, extra := pem.Decode(f)
		if block == nil || len(extra) > 0 {
			return nil, errors.New("failed to decode PEM")
		}
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	}

	func main() {
		// CLI Flags
		pem := flag.String("pem", "", "Xero API PEM file location")
		token := flag.String("token", "", "Xero API Token")
		flag.Parse()
		// Construct Authorizer
		key, err := pk(*pem)
		if err != nil {
			log.Fatal(err)
		}
		authorizer := oauth.New(*token, key)
		// Build a Xero API request
		req, err := http.NewRequest(
			http.MethodGet,
			"https://api.xero.com/api.xro/2.0/Invoices",
			nil)
		if err != nil {
			log.Fatal(err)
		}
		// Authorize Request
		if err := authorizer.AuthorizeRequest(req); err != nil {
			log.Fatal(err)
		}
		// Make Request
		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer rsp.Body.Close()
		log.Println(rsp.Status)
	}
*/
package oauth
