package couchdb

import (
	"fmt"
	"net/http"
)

// Bulk is a struct which should be embedded into user supplied structs to support
// editting multiple documents at once
type Bulk struct {
	NewEdits bool `json:"new_edits,omitempty"`
}

// BulkPut executes a bulk request. It assumes that the bulk parameter is a struct
// embedding the Bulk struct
func (d *Database) BulkPut(bulk interface{}) error {
	req, _ := http.NewRequest("POST", "/_bulk_docs", nil)
	resp, err := d.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("couchdb: POST %s returned %d", req.URL.Path, resp.StatusCode)
	}
	return nil
}
