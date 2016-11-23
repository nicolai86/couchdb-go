package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Document contains basic document identifications
type Document struct {
	ID      string `json:"_id,omitempty"`
	Rev     string `json:"_rev,omitempty"`
	Deleted bool   `json:"_deleted,omitempty"`
}

func revision(etag string) string {
	if etag == "" {
		return ""
	}
	return etag[1 : len(etag)-1]
}

type AllDocOpts struct {
	Limit       int
	IncludeDocs bool
	StartKey    string
	EndKey      string
}

func (d *Database) AllDocs(opts AllDocOpts, results interface{}) error {
	if opts.Limit == 0 {
		opts.Limit = 100
	}

	req, _ := http.NewRequest("GET", "/_all_docs", nil)
	values := req.URL.Query()
	if opts.Limit == 0 {
		opts.Limit = 100
	}
	values.Set("limit", strconv.Itoa(opts.Limit))
	values.Set("include_docs", strconv.FormatBool(opts.IncludeDocs))
	if opts.StartKey != "" {
		values.Set("startkey", fmt.Sprintf("%q", opts.StartKey))
	}
	if opts.EndKey != "" {
		values.Set("endkey", fmt.Sprintf("%q", opts.EndKey))
	}
	req.URL.RawQuery = values.Encode()

	resp, err := d.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &results); err != nil {
		return err
	}
	return nil
}

// Get fetches a document identified by it's id. GET /{db}/{id}
// this results in couchdb automatically returning the latest revision of the document
//
//  var doc couchdb.Document
//  db.Get("some-id", &doc)
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
//
//  var doc = couchdb.Document{
//    ID: "whatever",
//    Rev: "1-62bc3c4d01e43ee9d0cead8cd7c76041",
//  }
//  db.Put(doc.ID, &doc)
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

// Delete removes a document from a database
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
