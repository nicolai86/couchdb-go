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
		if exists, _ := client.DB.Exists("playground"); !exists {
			client.DB.Create("playground")
		}
		playground = client.Database("playground")

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
