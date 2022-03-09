#!/bin/bash

# =================================================================
#
# Copyright (C) 2022 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
set -eu

# move up a directory
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
go vet -vettool="bin/shadow" $pkgs
echo "******************"
echo "Running errcheck"
bin/errcheck ${pkgs}
echo "******************"
echo "Running staticcheck"
bin/staticcheck -checks all ${pkgs}
echo "******************"
echo "Running misspell"
bin/misspell -locale US -error *.md *.go
