package couchdb

import (
	"fmt"
	"net/http"
)

type Bulk struct {
	NewEdits bool `json:"new_edits,omitempty"`
}

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
