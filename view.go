package couchdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Results is a struct meant to be embedded in a couchdb request struct with correct
// rows, e.g.
//
//   type UserResults struct {
//       couchdb.Results
//       Users []user `json:"rows"`
//   }
type Results struct {
	Offset    int `json:"offset"`
	TotalRows int `json:"total_rows"`
}

// View defines map & reduce functions for a single view
type View struct {
	MapFn    string `json:"map,omitempty"`
	ReduceFn string `json:"reduce,omitempty"`
}

// DesignDocument describes a language and all associated views
type DesignDocument struct {
	Document
	Language string          `json:"language"`
	Views    map[string]View `json:"views"`
}

func (d *Database) Results(design, view string, opts AllDocOpts, result interface{}) error {
	req, _ := http.NewRequest("GET", fmt.Sprintf("/_design/%s/_view/%s", design, view), nil)
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
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("couchdb: GET %s returned %d", req.URL.Path, resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}
	return nil
}
