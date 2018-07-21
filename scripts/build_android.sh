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
echo "Building AAR for GSS"
cd $DIR/../bin
gomobile bind -target android -javapkg=com.spatialcurrent -o gss.aar github.com/spatialcurrent/go-simple-serializer/gss
if [[ "$?" != 0 ]] ; then
    echo "Error building AAR for Android"
    exit 1
fi
echo "Executable built at $(realpath $DIR/../bin)"
