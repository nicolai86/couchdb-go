package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Document contains basic document identifications
type Document struct {
	ID  string `json:"_id,omitempty"`
	Rev string `json:"_rev,omitempty"`
}

func revision(etag string) string {
	if etag == "" {
		return ""
	}
	return etag[1 : len(etag)-1]
}

// Get fetches a document identified by it's id. GET /{db}/{id}
// this results in couchdb automatically returning the latest revision of the document
func (d *Database) Get(id string, doc interface{}) error {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/%s", id), nil)
	resp, err := d.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("couchdb: GET %s returned %d", id, resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return err
	}
	return nil
}

// Put creates or updates a document, returning the new revision. PUT /{db}/{id}
func (d *Database) Put(id string, doc interface{}) (string, error) {
	bs, err := json.Marshal(doc)
	if err != nil {
		return "", err
	}
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/%s", id), bytes.NewReader(bs))
	resp, err := d.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("couchdb: PUT %s returned %d:\n%s\n", id, resp.StatusCode)
	}
	return revision(resp.Header.Get("Etag")), nil
}

func (d *Database) Delete(id, rev string) (string, error) {
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/%s", id), nil)
	values := req.URL.Query()
	values.Set("rev", rev)
	req.URL.RawQuery = values.Encode()

	resp, err := d.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("couchdb: DELETE %s returned %d:\n%s\n", id, resp.StatusCode)
	}
	return revision(resp.Header.Get("Etag")), nil
}

// Rev fetches the latest revision for a document. HEAD /{db}/{id}
func (d *Database) Rev(id string) (string, error) {
	req, _ := http.NewRequest("HEAD", fmt.Sprintf("/%s", id), nil)
	resp, err := d.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("couchdb: HEAD %s returned %d", id, resp.StatusCode)
	}
	return revision(resp.Header.Get("Etag")), nil
}
