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

	DB            *db
	Authenticator Authentication
}

type db struct {
	c *Client
}

// Create creates a new database by calling PUT /{db}
func (d *db) Create(name string) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("/%s", name), nil)
	if err != nil {
		return err
	}

	_, err = d.c.Do(req)
	return err
}

// Delete removes a database
func (d *db) Delete(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", name), nil)
	if err != nil {
		return err
	}

	_, err = d.c.Do(req)
	return err
}

// Exists checks if the given database exists with a HEAD /{db} request
func (d *db) Exists(name string) (bool, error) {
	req, err := http.NewRequest("HEAD", fmt.Sprintf("/%s", name), nil)
	if err != nil {
		return false, err
	}
	resp, err := d.c.Do(req)
	if err != nil {
		return false, err
	}
	exists := resp.StatusCode == http.StatusOK
	resp.Body.Close()
	return exists, nil
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

	if c.Authenticator != nil {
		if err := c.Authenticator.Decorate(req); err != nil {
			return nil, err
		}
	}

	resp, err := c.client.Do(req)
	return resp, err
}
