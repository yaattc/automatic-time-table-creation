#!/bin/sh

echo "Migrating to $DB_CONN_STR"
cd /srv/migrations
$(go env GOPATH)/bin/goose postgres $DB_CONN_STR up

echo "Running application"
export LOCATION=/db
/go/build/app $@
