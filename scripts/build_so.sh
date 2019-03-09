#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DEST=$(realpath ${1:-$DIR/../bin})
LDFLAGS="-X main.gitBranch=$(git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(git rev-list -1 HEAD)"
mkdir -p $DEST
echo "******************"
echo "Building Shared Object (*.so) for GSS"
cd $DEST
go build -o gss.so -buildmode=c-shared -ldflags "$LDFLAGS" github.com/spatialcurrent/go-simple-serializer/plugins/gss
if [[ "$?" != 0 ]] ; then
    echo "Error Building Shared Object (*.so) for GSS"
    exit 1
fi
echo "Shared Object (*.so) built at $DEST"
