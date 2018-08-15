#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

mkdir -p $DIR/../bin

echo "******************"
echo "Formatting $DIR/gss"
cd $DIR/../gss
go fmt
echo "Formatting $DIR/../cmd/gss"
cd $DIR/../cmd/gss
go fmt
echo "Done formatting."
echo "******************"
echo "Building Shared Object (*.so) for GSS"
cd $DIR/../bin
go build -o gss.so -buildmode=c-shared github.com/spatialcurrent/go-simple-serializer/plugins/gss
if [[ "$?" != 0 ]] ; then
    echo "Error Building Shared Object (*.so) for GSS"
    exit 1
fi
echo "Executable built at $(realpath $DIR/../bin)"
