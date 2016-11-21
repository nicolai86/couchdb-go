package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// HTTPExecutor wraps http client interactions. This allows users to pass in HTTP clients with tracing support.
type HTTPExecutor interface {
	Do(*http.Request) (*http.Response, error)
}

// Client contains all state necessary to identify a specific couchdb server
type Client struct {
	Host   string
	client HTTPExecutor

	Databases     *DatabaseService
	Users         *UserService
	Admins        *AdminUserService
	Replications  *ReplicationService
	Sessions      *SessionService
	Authenticator Authentication
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
	uri := fmt.Sprintf("%s%s", c.Host, req.URL)
	u, _ := url.Parse(uri)
	req.URL = u

	req.Header.Set("Content-Type", "application/json")

	if c.Authenticator != nil {
		if err := c.Authenticator.Decorate(req); err != nil {
			return nil, err
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		bs, err := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return resp, err
		}
		apiErr := ErrorResponse{}
		if err := json.Unmarshal(bs, &apiErr); err == nil {
			resp.Body = ioutil.NopCloser(bytes.NewBuffer(bs))
			return resp, apiErr
		}
	}

	return resp, err
}
