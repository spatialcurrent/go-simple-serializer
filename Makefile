# =================================================================
#
# Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

ifdef GOPATH
GCFLAGS=-trimpath=$(shell printenv GOPATH)/src
else
GCFLAGS=-trimpath=$(shell go env GOPATH)/src
endif

LDFLAGS=-X main.gitBranch=$(shell git branch | grep \* | cut -d ' ' -f2) -X main.gitCommit=$(shell git rev-list -1 HEAD)

ifndef DEST
DEST=bin
endif

deps:
	go get -d -t ./...

fmt:
	go fmt $$(go list ./... )

vet:
	go vet $$(go list ./...)

test:
	bash scripts/test.sh

bin/gss_darwin_amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(DEST)/gss_darwin_amd64 -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" github.com/spatialcurrent/go-simple-serializer/cmd/gss

bin/gss_linux_amd64:
	GOOS=linux GOARCH=amd64 go build -o $(DEST)/gss_linux_amd64 -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" github.com/spatialcurrent/go-simple-serializer/cmd/gss

bin/gss_windows_amd64.exe:
	GOOS=windows GOARCH=amd64 go build -o $(DEST)/gss_windows_amd64.exe -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" github.com/spatialcurrent/go-simple-serializer/cmd/gss

bin/gss_linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o $(DEST)/gss_linux_arm64 -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" github.com/spatialcurrent/go-simple-serializer/cmd/gss

# Build Shared Object
bin/gss_linux_amd64.so:
	# https://golang.org/cmd/link/
	# CGO Enabled : https://github.com/golang/go/issues/24068
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $(DEST)/gss_linux_amd64.so -buildmode=c-shared -ldflags "$(LDFLAGS)" -gcflags="$(GCFLAGS)" github.com/spatialcurrent/go-simple-serializer/plugins/gss

bin/gss_linux_armv7.so:
	# LDFLAGS - https://golang.org/cmd/link/
	# CGO Enabled  - https://github.com/golang/go/issues/24068
	# GOARM/GOARCH Compatability Table - https://github.com/golang/go/wiki/GoArm
	# ARM Cross Compiler Required - https://www.acmesystems.it/arm9_toolchain
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc go build -ldflags "-linkmode external -extldflags -static" -o $(DEST)/gss_linux_armv7.so -buildmode=c-shared -ldflags "$(LDFLAGS)" -gcflags="$(GCFLAGS)" github.com/spatialcurrent/go-simple-serializer/plugins/gss

# Build JavaScript Library
bin/gss.js:
	gopherjs build -o $(DEST)/gss.js github.com/spatialcurrent/go-simple-serializer/cmd/gss.js

# Build Minified JavaScript Library
bin/gss.min.js:
	gopherjs build -m -o $(DEST)/gss.min.js github.com/spatialcurrent/go-simple-serializer/cmd/gss.js

# Build Android Archive Library
bin/gss.aar:
	gomobile bind -target android -javapkg=com.spatialcurrent -o $(DEST)/gss.aar -gcflags="$(GCFLAGS)" github.com/spatialcurrent/go-simple-serializer/pkg/gss

build_cli: bin/gss_darwin_amd64 bin/gss_linux_amd64 bin/gss_windows_amd64.exe bin/gss_linux_arm64

build_javascript: bin/gss.js bin/gss.min.js

build_android: bin/gss.arr

build_so: bin/gss_linux_amd64.so bin/gss_linux_armv7.so

build: build_cli build_javascript build_so build_android

install:
	go install -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" github.com/spatialcurrent/go-simple-serializer/cmd/gss

clean:
	rm -fr bin
