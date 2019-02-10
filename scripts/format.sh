#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR/..
echo "******************"
echo "Formatting"
go fmt $(go list ./... )
