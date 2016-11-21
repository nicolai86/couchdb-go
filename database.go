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

// Exists checks if the given database exists with a HEAD /{db} request
func (d *Database) Exists() (bool, error) {
	req, err := http.NewRequest("HEAD", "", nil)
	if err != nil {
		return false, err
	}
	resp, err := d.Do(req)
	if err != nil {
		return false, err
	}
	exists := resp.StatusCode == http.StatusOK
	resp.Body.Close()
	return exists, nil
}

// Create creates a new database by calling PUT /{db}
func (d *Database) Create() error {
	req, err := http.NewRequest("PUT", "", nil)
	if err != nil {
		return err
	}

	_, err = d.Do(req)
	return err
}

// Destroy deletes a database
func (d *Database) Destroy() error {
	req, err := http.NewRequest("DELETE", "", nil)
	if err != nil {
		return err
	}

	_, err = d.Do(req)
	return err
}

// Do forwards requests to the http client, including the database in the URL path
func (d *Database) Do(req *http.Request) (*http.Response, error) {
	req.URL.Path = fmt.Sprintf("/%s%s", d.Name, req.URL.EscapedPath())
	return d.c.Do(req)
}
