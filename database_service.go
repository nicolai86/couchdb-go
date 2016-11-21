package couchdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// DatabaseService exposes database management apis
type DatabaseService struct {
	c *Client
}

// Create creates a new database by calling PUT /{db}
func (d *DatabaseService) Create(name string) error {
	req, err := http.NewRequest("PUT", fmt.Sprintf("/%s", name), nil)
	if err != nil {
		return err
	}

	_, err = d.c.Do(req)
	return err
}

// Delete removes a database
func (d *DatabaseService) Delete(name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/%s", name), nil)
	if err != nil {
		return err
	}

	_, err = d.c.Do(req)
	return err
}

// DatabaseMeta contains ever changing meta data about a single database
type DatabaseMeta struct {
	Name                          string `json:"db_name"`
	DocumentCount                 int    `json:"doc_count"`
	DocumentDeletionCount         int    `json:"doc_del_count"`
	UpdateSequenceNumber          int    `json:"update_seq"`
	PurgeSequenceNumber           int    `json:"purge_seq"`
	CompactRunning                bool   `json:"compact_running"`
	DiskSize                      int    `json:"disk_size"`
	DataSize                      int    `json:"data_size"`
	InstanceStartTime             string `json:"instance_start_time"`
	DiskFormatVersion             int    `json:"disk_format_version"`
	CommittedUpdateSequenceNumber int    `json:"committed_update_seq"`
}

// Meta looks up database metadata
func (d *DatabaseService) Meta(name string) (DatabaseMeta, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("/%s", name), nil)
	if err != nil {
		return DatabaseMeta{}, err
	}
	resp, err := d.c.Do(req)
	if err != nil {
		return DatabaseMeta{}, err
	}
	defer resp.Body.Close()
	meta := DatabaseMeta{}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return DatabaseMeta{}, err
	}
	err = json.Unmarshal(bs, &meta)
	return meta, err
}

// Exists checks if the given database exists with a HEAD /{db} request
func (d *DatabaseService) Exists(name string) (bool, error) {
	req, err := http.NewRequest("HEAD", fmt.Sprintf("/%s", name), nil)
	if err != nil {
		return false, err
	}
	resp, err := d.c.Do(req)
	if err != nil {
		return false, err
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK, nil
}
