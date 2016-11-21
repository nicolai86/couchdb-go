// package couchdb
//
// provides a wrapper around the couchdb HTTP API
package couchdb

import (
	"net/http"
	"time"
)

// New returns a configured couchdb client
func New(host string) *Client {
	return &Client{
		host: host,
		client: &http.Client{
			Timeout: time.Duration(60 * time.Second),
		},
	}
}
