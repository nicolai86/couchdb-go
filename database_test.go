package couchdb

import "testing"

func TestDatabase_NotExisting(t *testing.T) {
	t.Parallel()

	db := client.Database("foobar")
	exists, err := db.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("Expected database %q to not exist, but does.", db.Name)
	}
}

func TestDatabase_Exists(t *testing.T) {
	t.Parallel()

	db := client.Database("_replicator")
	exists, err := db.Exists()
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("Expected database %q to exist, but didn't.", db.Name)
	}
}

func TestDatabase_Create(t *testing.T) {
	db := client.Database("new-db")
	if err := db.Create(); err != nil {
		t.Fatal(err)
	}
}

func TestDatabase_Destroy(t *testing.T) {
	db := client.Database("new-db")
	if err := db.Destroy(); err != nil {
		t.Fatal(err)
	}
}
