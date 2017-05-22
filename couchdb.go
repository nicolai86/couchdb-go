// Package couchdb provides a wrapper around the couchdb HTTP API
package couchdb

import (
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
)

// New returns a configured couchdb client
func New(host string, tracer opentracing.Tracer, configs ...func(*Client)) *Client {
	c := &Client{
		host: host,
		client: &http.Client{
			Timeout: time.Duration(60 * time.Second),
		},
		tracer: tracer,
	}
	for _, config := range configs {
		config(c)
	}
	c.DB = &db{
		c: c,
	}
	return c
}
