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
echo "Building Javascript for GSS"
cd $DIR/../bin
gopherjs build -o gss.js github.com/spatialcurrent/go-simple-serializer/cmd/gss.js
if [[ "$?" != 0 ]] ; then
    echo "Error building Javascript for GSS"
    exit 1
fi
echo "Executable built at $(realpath $DIR/../bin)"
