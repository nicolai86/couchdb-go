package couchdb

import (
	"context"
	"testing"
	"time"
)

func TestDatabase_Subscribe(t *testing.T) {
	t.Run("since = 0", func(t *testing.T) {
		feed := playground.Subscribe(context.Background(), FeedOpts{})
		select {
		case change := <-feed.Ch:
			if change.ID != "employee:michael" {
				t.Fatalf("Inserted employee:michael first, but got %q", change.ID)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Expected change within 2 seconds, but got none.")
		}

	})

	t.Run("since != 0", func(t *testing.T) {
		feed := playground.Subscribe(context.Background(), FeedOpts{
			Since: 2,
		})

		select {
		case change := <-feed.Ch:
			if change.ID != "pet:yumi" {
				t.Fatalf("Inserted pet:yumi, but got %q", change.ID)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Expected change within 2 seconds, but got none.")
		}

	})
}
