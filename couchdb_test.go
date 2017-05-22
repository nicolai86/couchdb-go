package couchdb

import (
	"flag"
	"os"
	"testing"

	"github.com/opentracing/opentracing-go"
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
	if os.Getenv("COUCHDB_HOST_PORT") == "" {
		fmt.Println("Skipping couchdb tests as COUCHDB_HOST_PORT is not configured")
		os.Exit(0)
	}
	client = New(os.Getenv("COUCHDB_HOST_PORT"), opentracing.NoopTracer{})

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
