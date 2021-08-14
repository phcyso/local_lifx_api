#!/usr/bin/env bash

docker build --target builder -t lifx_builder .

mkdir -p dist
rm -rf ./dist/local_lifx_api
docker run --rm  -v "$(pwd)/dist/:/tmp/" -t lifx_builder sh -c 'cp /app/local_lifx_api /tmp/local_lifx_api'
