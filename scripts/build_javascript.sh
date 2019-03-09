#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DEST=$(realpath ${1:-$DIR/../bin})
LDFLAGS="-X main.gitBranch=$(git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(git rev-list -1 HEAD)"
mkdir -p $DEST
echo "******************"
echo "Building Javascript for GSS"
cd $DEST
gopherjs build -o gss.js github.com/spatialcurrent/go-simple-serializer/cmd/gss.js
if [[ "$?" != 0 ]] ; then
    echo "Error building Javascript artifacts for GSS"
    exit 1
fi
gopherjs build -m -o gss.min.js github.com/spatialcurrent/go-simple-serializer/cmd/gss.js
if [[ "$?" != 0 ]] ; then
    echo "Error building Javascript artifacts for GSS"
    exit 1
fi
echo "JavaScript artificats built at $DEST"
