# =================================================================
#
# Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

build_so:
	go build -o gss.so -buildmode=c-shared github.com/spatialcurrent/go-simple-serializer/plugins/gss
build_cpp:
	mkdir -p bin && g++ -o bin/gss_test_cpp test.cpp -L. -l:gss.so
build: build_so build_cpp
run:
	LD_LIBRARY_PATH=. bin/gss_test_cpp
clean:
	rm -f gss.h gss.so bin/gss_test_cpp
