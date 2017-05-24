# `go-xero`

[![CircleCI](https://img.shields.io/circleci/project/github/thisissoon/go-xero.svg)](https://circleci.com/gh/thisissoon/go-xero)
[![Coverage](https://img.shields.io/codecov/c/github/thisissoon/go-xero.svg)](https://codecov.io/gh/thisissoon/go-xero)
[![GoDoc](https://godoc.org/github.com/thisissoon/go-xero?status.svg)](https://godoc.org/github.com/thisissoon/go-xero)

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

## Task List

- [x] Base API Request Client
- [x] `go-oauth` Authorizer
- [x] Simple `GET|POST|PUT` support
- [x] Base Test Suite / CI
- [x] PUT/POST Error Handling
- [ ] Attchments
  - [ ] `GET`
- [ ] Accounts (@jamesjwarren)
  - [x] `GET`
  - [ ] `DELETE`
- [x] Bank Transactions (@jamesjwarren)
  - [x] `GET`
- [x] Bank Transfer (@jamesjwarren)
  - [x] `GET`
- [ ] Branding Themes
  - [ ] `GET`
- [ ] Contacts (@krak3n)
  - [x] `GET`
- [ ] Contact Groups
  - [ ] `GET`
  - [ ] `DELETE`
- [ ] Credit Notes
  - [ ] `GET`
- [ ] Currencies
  - [ ] `GET`
- [ ] Employees
  - [ ] `GET`
- [ ] Expense Claims
  - [ ] `GET`
- [ ] Invoices
  - [ ] `GET`
- [ ] Invoice Reminders
  - [ ] `GET`
- [ ] Items
  - [ ] `GET`
  - [ ] `DELETE`
- [ ] Journals
  - [ ] `GET`
- [ ] Linked Transactions
  - [ ] `GET`
  - [ ] `DELETE`
- [ ] Manual Journals
  - [ ] `GET`
  - [ ] `DELETE`
- [ ] Organisation
  - [ ] `GET`
- [ ] Overpayments
  - [ ] `GET`
- [ ] Payments
  - [ ] `GET`
- [ ] Prepayments
  - [ ] `GET`
- [ ] Purchase Orders
  - [ ] `GET`
- [ ] Receipts
  - [ ] `GET`
- [ ] Repeating Invoices
  - [ ] `GET`
- [ ] Reports
  - [ ] `GET`
- [ ] Tax Rates
  - [ ] `GET`
- [ ] Tracking Categories
  - [ ] `GET`
  - [ ] `DELETE`
- [ ] Users
  - [ ] `GET`
