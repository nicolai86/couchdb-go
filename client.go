package couchdb

import (
	"fmt"
	"net/http"
	"net/url"
)

// Client contains all state necessary to identify a specific couchdb server
type Client struct {
	host   string
	client *http.Client
}

// Database returns a database wrapper for a given db
func (c *Client) Database(name string) *Database {
	return &Database{
		c:    c,
		Name: name,
	}
}

// Do executes a http request against the specific couchdb, setting all required headers
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	uri := fmt.Sprintf("%s%s", c.host, req.URL)
	u, _ := url.Parse(uri)
	req.URL = u

	req.Header.Set("Content-Type", "application/json")

	return c.client.Do(req)
}
