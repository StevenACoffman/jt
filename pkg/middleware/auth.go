// package middleware - HTTP client round trippers
// sometimes called "tripperware" to disambiguate from
// HTTP server middleware.

package middleware

import (
	"encoding/base64"
	"net/http"
	"os"
	"time"
)

// HeaderRoundTripper is a client middleware for adding headers on every request.
type HeaderRoundTripper struct {
	next   http.RoundTripper
	Header http.Header
}

func NewHeaderRoundTripper(next http.RoundTripper, header http.Header) *HeaderRoundTripper {
	if next == nil {
		next = http.DefaultTransport
	}
	return &HeaderRoundTripper{
		next:   next,
		Header: header,
	}
}

func (rt *HeaderRoundTripper) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	if rt.Header != nil {
		for k, v := range rt.Header {
			req.Header[k] = v
		}
	}
	return rt.next.RoundTrip(req)
}

// BasicAuth is a convenience method
func (rt *HeaderRoundTripper) BasicAuth(username, password string) {
	if rt.Header == nil {
		rt.Header = make(http.Header)
	}

	auth := username + ":" + password
	base64Auth := base64.StdEncoding.EncodeToString([]byte(auth))
	rt.Header.Set("Authorization", "Basic "+base64Auth)
}

func (rt *HeaderRoundTripper) SetHeader(key, value string) {
	if rt.Header == nil {
		rt.Header = make(http.Header)
	}
	rt.Header.Set(key, value)
}

// NewBasicAuthHTTPClient configures an HTTP Client
// that adds basic auth header and json
// as well as a generous 60-second timeout.
func NewBasicAuthHTTPClient(user, token string) *http.Client {
	header := make(http.Header)
	header.Set("Content-Type", "application/json; charset=utf-8")
	header.Set("Accept", "application/json; charset=utf-8")
	rt := NewLoggingRoundTripper(http.DefaultTransport, os.Stdout)
	hrt := NewHeaderRoundTripper(rt, header)
	hrt.BasicAuth(user, token)

	return &http.Client{
		Transport: hrt,
		Timeout:   60 * time.Second,
	}
}
