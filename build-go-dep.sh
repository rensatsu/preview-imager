#!/bin/bash
set -o errexit
set -o pipefail
set -o nounset

go get -v -t -d ./...
if [ ! -d dist ]; then
    mkdir dist
fi
go build -v -o dist/preview-imager ./go
