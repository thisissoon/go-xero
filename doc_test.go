package xero_test

import (
	"fmt"
	"io"
	"log"
	"os"

	xero "github.com/thisissoon/go-xero"
	oauth "github.com/thisissoon/go-xero/authorizers/go-oauth"
)

func Example() {
	f, err := os.Open("/path/to/privatekey.pem")
	if err != nil {
		log.Fatal(err)
	}
	pk, err := xero.PrivateKey(f)
	if err != nil {
		log.Fatal(err)
	}
	client := xero.New(oauth.New("TOKEN", pk))
	rsp, err := client.Get("https://api.xero.com/api.xro/2.0/Invoices")
	if err != nil {
		log.Fatal(err)
	}
	defer rsp.Body.Close()
	log.Println(rsp.StatusCode)
}

func ExampleClient_Contacts() {
	f, err := os.Open("/path/to/privatekey.pem")
	if err != nil {
		log.Fatal(err)
	}
	pk, err := xero.PrivateKey(f)
	if err != nil {
		log.Fatal(err)
	}
	client := xero.New(oauth.New("TOKEN", pk))
	for i, contacts, err := client.Contacts(); err != io.EOF; i, contacts, err = i.Next() {
		if err != nil {
			log.Fatal(err)
		}
		for i := range contacts {
			fmt.Println(contacts[i].Name)
		}
	}
}
