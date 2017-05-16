# Contacts Example

This example shows the basics of making an iterative paginated call to the `/Contacts`
Xero API endpoont to return a list of `Contact`'s

```
$ go get github.com/thisissoon/go-xero
$ go get github.com/thisissoon/go-xero/authorizers/go-oauth
$ go run examples/contacts/main.go -pemfile=/path/to.pem -token=TOKEN
```
