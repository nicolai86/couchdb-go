package couchdb

import "testing"

func TestDatabase_NotExisting(t *testing.T) {
	t.Parallel()

	db := client.Database("foobar")
	exists, err := client.DB.Exists(db.Name)
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
	exists, err := client.DB.Exists(db.Name)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatalf("Expected database %q to exist, but didn't.", db.Name)
	}
}

func TestClient_Create(t *testing.T) {
	if err := client.DB.Create("new-db"); err != nil {
		t.Fatal(err)
	}
}

func TestClient_Delete(t *testing.T) {
	if err := client.DB.Delete("new-db"); err != nil {
		t.Fatal(err)
	}
}
