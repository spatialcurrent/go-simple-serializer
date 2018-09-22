#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
set -eu
echo "******************"
echo "Formatting"
cd $DIR/../gss
go fmt
cd ./../gssjs
go fmt
cd ./../cmd/gss
go fmt
cd ./../gss.js
go fmt
cd ./../../plugins/gss
go fmt