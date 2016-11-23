package couchdb

import (
	"fmt"
	"net/http"
)

// Database is a client for a specific couchdb server & database
type Database struct {
	c    *Client
	Name string
}

// Do forwards requests to the http client, prefixing the URL path with the database name
func (d *Database) Do(req *http.Request) (*http.Response, error) {
	req.URL.Path = fmt.Sprintf("/%s%s", d.Name, req.URL.EscapedPath())
	return d.c.Do(req)
}
