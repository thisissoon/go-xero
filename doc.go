/*
This package provides an API wrapper around the Xero HTTP API.
Warning: This package is in early development and should not be used in production.

Authorizarion

The Xero API uses OAuth to authorize API requests. This package on it's own makes
no assumptions about how you want to implement this flow. You can either use the
provided go-oauth (which uses http://github.com/garyburd/go-oauth/oauth) authorizer
or implement your own using any other OAuth library of your choosing,
you just need to implement the Authorizer interface.

This means the package itself does not enforce lock-in to any external 3rd party library.
*/
package xero
