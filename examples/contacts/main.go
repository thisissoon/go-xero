package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	xero "github.com/thisissoon/go-xero"
	xauth "github.com/thisissoon/go-xero/authorizers/go-oauth"
)

func main() {
	pemfile := flag.String("pemfile", "", "Xero API PEM file location")
	token := flag.String("token", "", "Xero API Token")
	flag.Parse()
	f, err := os.Open(*pemfile)
	if err != nil {
		log.Fatal(err)
	}
	pk, err := xero.PrivateKey(f)
	if err != nil {
		log.Fatal(err)
	}
	client := xero.New(xauth.New(*token, pk))
	for i, contacts, err := client.Contacts(); err != io.EOF; i, contacts, err = i.Next() {
		if err != nil {
			log.Fatal(err)
		}
		for i := range contacts {
			fmt.Println(contacts[i].Name)
		}
	}
}
