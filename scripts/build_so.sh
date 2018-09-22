#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DEST=$(realpath ${1:-$DIR/../bin})

mkdir -p $DEST

echo "******************"
echo "Building Shared Object (*.so) for GSS"
cd $DEST
go build -o gss.so -buildmode=c-shared github.com/spatialcurrent/go-simple-serializer/plugins/gss
if [[ "$?" != 0 ]] ; then
    echo "Error Building Shared Object (*.so) for GSS"
    exit 1
fi
echo "Shared Object (*.so) built at $DEST"
