#!/bin/bash
set -e

trap 'killall distribkv' SIGINT

cd $(dirname $0)

killall distribkv || true
sleep 0.1

go build distributed-kv

go run distributed-kv -db-location=mayur.db -http-addr=127.0.0.1:8080 -config-file=sharding.toml -shard=Mayur &
go run distributed-kv -db-location=fahad.db -http-addr=127.0.0.1:8081 -config-file=sharding.toml -shard=Fahad &
go run distributed-kv -db-location=farhan.db -http-addr=127.0.0.1:8082 -config-file=sharding.toml -shard=Farhan &

wait