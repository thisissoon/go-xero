# `go-xero`

A `go` package for interfacing with the Xero accounting software HTTP API.

``` go
package main

import (
	"flag"
	"log"
	"os"

	xero "github.com/thisissoon/go-xero"
	oauth "github.com/thisissoon/go-xero/authorizers/go-oauth"
)

func main() {
	// CLI Flags
	pemfile := flag.String("pemfile", "", "Xero API PEM file location")
	token := flag.String("token", "", "Xero API Token")
	flag.Parse()
	// Read Private Key
	f, err := os.Open(*pemfile)
	if err != nil {
		log.Fatal(err)
	}
	pk, err := xero.PrivateKey(f)
	if err != nil {
		log.Fatal(err)
	}
	// Client
	client := xero.New(oauth.New(*token, pk))
	rsp, err := client.Get("https://api.xero.com/api.xro/2.0/Invoices")
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()
	log.Println(rsp.StatusCode)
}
```
