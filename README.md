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

## Task List

- [x] Base API Request Client
- [x] `go-oauth` Authorizer
- [x] Simple `GET|POST|PUT` support
- [ ] Base Test Suite / CI
- [ ] PUT/POST Error Handling
- [ ] Attchments
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Accounts
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
  - [ ] `DELETE`
- [ ] Bank Transactions
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Bank Transfers
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Branding Themes
  - [ ] `GET`
- [ ] Contacts
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Contact Groups
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
  - [ ] `DELETE`
- [ ] Credit Notes
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Currencies
  - [ ] `GET`
- [ ] Employees
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Expense Claims
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Invoices
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Invoice Reminders
  - [ ] `GET`
- [ ] Items
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
  - [ ] `DELETE`
- [ ] Journals
  - [ ] `GET`
- [ ] Linked Transactions
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
  - [ ] `DELETE`
- [ ] Manual Journals
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
  - [ ] `DELETE`
- [ ] Organisation
  - [ ] `GET`
- [ ] Overpayments
  - [ ] `GET`
  - [ ] `PUT`
- [ ] Payments
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Prepayments
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Purchase Orders
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Receipts
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Repeating Invoices
  - [ ] `GET`
- [ ] Reports
  - [ ] `GET`
- [ ] Tax Rates
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
- [ ] Tracking Categories
  - [ ] `GET`
  - [ ] `POST`
  - [ ] `PUT`
  - [ ] `DELETE`
- [ ] Users
  - [ ] `GET`
