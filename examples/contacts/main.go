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
	// Iteration
	fmt.Println("Contact Iteration")
	fmt.Println("-----------------")
	for i, contacts, err := client.Contacts(); err != io.EOF; i, contacts, err = i.Next() {
		if err != nil {
			log.Fatal(err)
		}
		for i := range contacts {
			fmt.Println(contacts[i].ContactID, contacts[i].Name)
		}
	}
	// Get Singular Contact
	fmt.Println("Single Contact")
	fmt.Println("--------------")
	contact, err := client.Contact("847d9131-ac1e-4819-af61-b3f1cf854009")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(contact.ContactID, contact.Name)
}
