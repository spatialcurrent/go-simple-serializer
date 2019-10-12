#!/bin/bash

# =================================================================
#
# Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
set -euo pipefail
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
staticcheck -checks all $(go list ./... | grep -v /vendor/ | grep -v /plugins/ | tr "\n" " ")
echo "******************"
echo "Running misspell"
misspell -locale US -error *.md *.go
echo "******************"
