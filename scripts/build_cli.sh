#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
DEST=$(realpath ${1:-$DIR/../bin})

mkdir -p $DEST

echo "******************"
echo "Building program gss"
cd $DEST
####################################################
#echo "Building program for darwin"
#GOTAGS= CGO_ENABLED=1 GOOS=${GOOS} GOARCH=amd64 go build --tags "darwin" -o "gss_darwin_amd64" github.com/spatialcurrent/go-simple-serializer/cmd/gss
#if [[ "$?" != 0 ]] ; then
#    echo "Error building gss for Darwin"
#    exit 1
#fi
#echo "Executable built at $(realpath $DIR/../bin/gss_darwin_amd64)"
####################################################
echo "Building program for Linux"
GOTAGS= CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build --tags "linux" -o "gss_linux_amd64" github.com/spatialcurrent/go-simple-serializer/cmd/gss
if [[ "$?" != 0 ]] ; then
    echo "Error building gss for Linux"
    exit 1
fi
echo "Executable built at $DEST"
####################################################
echo "Building program for Windows"
GOTAGS= CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CXX=x86_64-w64-mingw32-g++ CC=x86_64-w64-mingw32-gcc go build -o "gss_windows_amd64.exe" github.com/spatialcurrent/go-simple-serializer/cmd/gss
if [[ "$?" != 0 ]] ; then
    echo "Error building gss for Windows"
    exit 1
fi
echo "Executable built at $DEST"
