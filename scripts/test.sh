#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
set -eu
cd $DIR/..
pkgs=$(go list ./... )
echo "******************"
echo "Running unit tests for: "
for pkg in "${pkgs[@]}"; do
   go test -p 1 -count 1 -short $pkgs
done
echo "******************"
echo "Using gometalinter with misspell, vet, ineffassign, and gosec"
gometalinter \
--misspell-locale=US \
--disable-all \
--enable=misspell \
--enable=vet \
--enable=ineffassign \
--enable=gosec \
./...
