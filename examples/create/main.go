package main

import (
	"flag"
	"fmt"
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
	// Contacts to create
	contacts := xero.Contacts{
		[]xero.Contact{
			{Name: "Foo"},
			{Name: "Bar"},
		},
	}
	// Contacts Response
	rsp := xero.ContactsResponse{}
	// Client
	client := xero.New(oauth.New(*token, pk))
	if err := client.CreateUpdate(xero.ContactsEndpoint, contacts, &rsp); err != nil {
		log.Fatal(err)
	}
	for _, c := range rsp.Contacts.Contacts {
		fmt.Println(c.ContactID, c.Name)
	}
}
