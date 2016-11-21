#!/usr/bin/env bash
#
set -e

couchdb_id=$(docker run -p 5984:5984 -d couchdb)
function cleanup {
  docker stop $couchdb_id > /dev/null
  docker rm $couchdb_id > /dev/null
}
trap cleanup EXIT

echo Waiting for couchdb to start…
while ! curl http://localhost:5984 -s -q > /dev/null; do sleep 1; done

export COUCHDB_HOST_PORT=http://localhost:5984
echo Running tests…
go test ./... -v
