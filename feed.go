package couchdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// Feed allows users to consume the change feed of couchdb
type Feed struct {
	Ch  <-chan Change
	err error
}

// Error returns the last error when trying to access this feed
func (f *Feed) Error() error {
	return f.err
}

// FeedType describes supported couchdb change feeds
type FeedType string

// defines supported FeedType s
const (
	ContinuousFeed FeedType = "continuous"
)

// FeedOpts are passed along when subscribing for changes
type FeedOpts struct {
	Type        FeedType
	Since       uint64
	Timeout     int
	Heartbeat   int
	IncludeDocs bool
}

// Change identifies a single even from the _changes endpoint
type Change struct {
	Seq int64            `json:"seq"`
	ID  string           `json:"id"`
	Doc *json.RawMessage `json:"doc"`
}

// Subscribe spawns a goroutine used to continously pull data off a couchdb HTTP feed
func (d *Database) Subscribe(ctx context.Context, opts FeedOpts) *Feed {
	if opts.Type == "" {
		opts.Type = ContinuousFeed
	}
	if opts.Timeout == 0 {
		opts.Timeout = 60 * 1000
	}
	if opts.Heartbeat == 0 {
		opts.Heartbeat = 10 * 1000
	}
	ch := make(chan Change)
	var f = Feed{
		Ch: ch,
	}

	go func() {
		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/%s/_changes", d.c.host, d.Name), nil)

		values := req.URL.Query()
		values.Set("timeout", strconv.Itoa(opts.Timeout))
		values.Set("heartbeat", strconv.Itoa(opts.Heartbeat))
		if opts.Since != 0 {
			values.Set("since", strconv.Itoa(int(opts.Since)))
		}
		values.Set("feed", string(opts.Type))
		if opts.IncludeDocs {
			values.Set("include_docs", "true")
		}
		req.URL.RawQuery = values.Encode()

		resp, err := d.c.client.Get(req.URL.String())
		if err != nil {
			f.err = err
			return
		}
		if resp.StatusCode != http.StatusOK {
			f.err = fmt.Errorf("couchdb: GET %s returned %d", req.URL, resp.StatusCode)
			return
		}

		d := json.NewDecoder(resp.Body)

		go func() {
			select {
			case <-ctx.Done():
				close(ch)
				resp.Body.Close()
				return
			}
		}()

		for {
			var c Change
			if err := d.Decode(&c); err != nil {
				f.err = err
				resp.Body.Close()
				close(ch)
				return
			}
			ch <- c
		}
	}()

	return &f
}
