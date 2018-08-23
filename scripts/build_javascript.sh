#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DEST=$(realpath ${1:-$DIR/../bin})

mkdir -p $DEST

echo "******************"
echo "Formatting $(realpath $DIR/../gss)"
cd $DIR/../gss
go fmt
echo "Formatting $(realpath $DIR/../cmd/gss.js)"
cd $DIR/../cmd/gss.js
go fmt
echo "Done formatting."
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
