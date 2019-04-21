#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
set -eu
cd $DIR/..
pkgs=$(go list ./... | grep -v /vendor/ | tr "\n" " ")
echo "******************"
echo "Running unit tests"
go test -p 1 -count 1 -short $pkgs
echo "******************"
echo "Running go vet"
go vet $pkgs
echo "******************"
echo "Running go vet with shadow"
go vet -vettool=$(which shadow) $pkgs
echo "******************"
echo "Running errcheck"
errcheck ${pkgs}
echo "******************"
echo "Running ineffassign"
find . -name '*.go' | xargs ineffassign
echo "******************"
echo "Running staticcheck"
staticcheck -checks all ${pkgs}
echo "******************"
echo "Running misspell"
misspell -locale US -error *.md *.go
