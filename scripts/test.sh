#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
set -eu
echo "******************"
echo "Running unit tests"
cd $DIR/../gss
go test
echo "******************"
echo "Using gometalinter with misspell, vet, ineffassign, and gosec"
echo "Testing $DIR/../gss"
gometalinter --misspell-locale=US --disable-all --enable=misspell --enable=vet --enable=ineffassign --enable=gosec $DIR/../gss

echo "Testing $DIR/../plugins/gss"
gometalinter --misspell-locale=US --disable-all --enable=misspell --enable=vet --enable=ineffassign --enable=gosec $DIR/../plugins/gss

echo "Testing $DIR/../cmd/gss"
gometalinter --misspell-locale=US --disable-all --enable=misspell --enable=vet --enable=ineffassign --enable=gosec $DIR/../cmd/gss

echo "Testing $DIR/../cmd/gss.js"
gometalinter --misspell-locale=US --disable-all --enable=misspell --enable=vet --enable=ineffassign --enable=gosec $DIR/../cmd/gss.js