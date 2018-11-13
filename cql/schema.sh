#!/bin/sh

DIR=$(dirname "$0")

until cqlsh -f "$DIR/schema.cql"; do
  echo "Failed to create schema. Retry in 1 second"
  sleep 1
done &

exec /docker-entrypoint.sh "$@"
