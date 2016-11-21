package couchdb

import (
	"flag"
	"os"
	"testing"
)

var (
	client     *Client
	playground *Database
)

type testDoc struct {
	Document
	Name string `json:"name"`
}

func TestMain(m *testing.M) {
	client = New(os.Getenv("COUCHDB_HOST_PORT"))

	func() {
		playground = client.Database("playground")
		if exists, _ := playground.Exists(); !exists {
			playground.Create()
		}

		playground.Put("employee:michael", testDoc{
			Name: "Michael",
		})
		playground.Put("employee:raphael", testDoc{
			Name: "Raphael",
		})
		playground.Put("pet:yumi", testDoc{
			Name: "Yumi",
		})
	}()

	flag.Parse()
	os.Exit(m.Run())
}
