package couchdb

import "testing"

func TestDatabase_Put(t *testing.T) {
	t.Parallel()

	var doc = testDoc{
		Document: Document{
			ID: "employee:martin",
		},
		Name: "Martin",
	}

	db := client.Database("put-test")
	client.DB.Create(db.Name)
	defer client.DB.Delete(db.Name)

	t.Run("insert", func(t *testing.T) {
		rev, err := db.Put(doc.ID, &doc)
		if err != nil {
			t.Fatal(err)
		}
		if rev == "" {
			t.Fatal("Expected to receive a document revision, but got nothing")
		}
		doc.Rev = rev
	})

	t.Run("update", func(t *testing.T) {
		doc.Name = "Klaus"
		rev, err := db.Put(doc.ID, &doc)
		if err != nil {
			t.Fatal(err)
		}
		if rev == doc.Rev {
			t.Fatalf("Expected update to succeed, but didn't")
		}
	})
}

func TestDatabase_Delete(t *testing.T) {
	t.Parallel()

	db := client.Database("delete-test")
	client.DB.Create(db.Name)
	defer client.DB.Delete(db.Name)

	rev, _ := db.Put("test", Document{
		ID: "test",
	})
	rmRev, err := db.Delete("test", rev)
	if err != nil {
		t.Fatal(err)
	}
	if rmRev == rev {
		t.Fatalf("Expected new revision, but got %q", rmRev)
	}
}

func TestDatabase_Get(t *testing.T) {
	t.Parallel()

	db := client.Database("_users")
	t.Run("known", func(t *testing.T) {
		t.Parallel()

		docID := "_design/_auth"
		var doc Document
		err := db.Get(docID, &doc)
		if err != nil {
			t.Fatal(err)
		}
		if doc.ID != docID {
			t.Fatalf("Expected doc %q, but got %q", docID, doc.ID)
		}
	})
	t.Run("unknown", func(t *testing.T) {
		t.Parallel()

		docID := "admin"
		var doc Document

		if err := db.Get(docID, &doc); err == nil {
			t.Fatal(err)
		}
	})
}

func TestDatabase_Rev(t *testing.T) {
	t.Parallel()
	db := client.Database("_users")

	t.Run("known", func(t *testing.T) {
		t.Parallel()

		docID := "_design/_auth"
		var doc Document
		err := db.Get(docID, &doc)
		if err != nil {
			t.Fatal(err)
		}

		rev, err := db.Rev(docID)
		if err != nil {
			t.Fatal(err)
		}
		if rev != doc.Rev {
			t.Fatalf("Expected revisions to match, but didn't: %q != %q", rev, doc.Rev)
		}
	})
	t.Run("unknown", func(t *testing.T) {
		t.Parallel()

		db := client.Database("_users")
		if _, err := db.Rev("admin"); err == nil {
			t.Fatal(err)
		}
	})
}
