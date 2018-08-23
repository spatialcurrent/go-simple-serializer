#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DEST=$(realpath ${1:-$DIR/../bin})

mkdir -p $DEST

echo "******************"
echo "Formatting $(realpath $DIR/../gss)"
cd $DIR/../gss
go fmt
echo "Done formatting."
echo "******************"
echo "Building AAR for GSS"
cd $DEST
gomobile bind -target android -javapkg=com.spatialcurrent -o gss.aar github.com/spatialcurrent/go-simple-serializer/gss
if [[ "$?" != 0 ]] ; then
    echo "Error building AAR for Android"
    exit 1
fi
echo "AAR built at $DEST"
