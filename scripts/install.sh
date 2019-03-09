#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
set -eu
LDFLAGS="-X main.gitBranch=$(git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(git rev-list -1 HEAD)"
cd $DIR/..
pkgs=$(go list ./... | grep cmd | grep -v '.js')
echo "******************"
echo "Installing programs"
for pkg in "${pkgs[@]}"; do
  echo "Installing $(basename $pkg) from $pkg"
  go install -ldflags "$LDFLAGS" $pkg
done
