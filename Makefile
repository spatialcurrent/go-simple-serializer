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

.PHONY: help

help:  ## Print the help documentation
	@grep -E '^[a-zA-Z0-9_-\]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#
# Dependencies
#

deps_go:  ## Install Go dependencies
	go get -d -t ./...

.PHONY: deps_go_test
deps_go_test: ## Download Go dependencies for tests
	go get golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow # download shadow
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow # install shadow
	go get -u github.com/kisielk/errcheck # download and install errcheck
	go get -u github.com/client9/misspell/cmd/misspell # download and install misspell
	go get -u github.com/gordonklaus/ineffassign # download and install ineffassign
	go get -u honnef.co/go/tools/cmd/staticcheck # download and instal staticcheck
	go get -u golang.org/x/tools/cmd/goimports # download and install goimports

deps_arm:  ## Install dependencies to cross-compile to ARM
	# ARMv7
	apt-get install -y libc6-armel-cross libc6-dev-armel-cross binutils-arm-linux-gnueabi libncurses5-dev gcc-arm-linux-gnueabi g++-arm-linux-gnueabi
  # ARMv8
	apt-get install gcc-aarch64-linux-gnu g++-aarch64-linux-gnu

deps_gopherjs:  ## Install GopherJS
	go get -u github.com/gopherjs/gopherjs

deps_javascript:  ## Install dependencies for JavaScript tests
	npm install .

#
# Go building, formatting, testing, and installing
#

fmt:  ## Format Go source code
	go fmt $$(go list ./... )

imports: ## Update imports in Go source code
	# If missing, install goimports with: go get golang.org/x/tools/cmd/goimports
	goimports -w -local github.com/spatialcurrent/go-simple-serializer,github.com/spatialcurrent/ $$(find . -iname '*.go')

vet: ## Vet Go source code
	go vet $$(go list ./...)

test_go: ## Run Go tests
	bash scripts/test.sh

build: build_cli build_javascript build_so build_android  ## Build CLI, Shared Objects (.so), JavaScript, and Android

install:  ## Install GSS CLI on current platform
	go install -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" github.com/spatialcurrent/go-simple-serializer/cmd/gss

#
# Command line Programs
#

bin/gss_darwin_amd64: ## Build GSS CLI for Darwin / amd64
	GOOS=darwin GOARCH=amd64 go build -o $(DEST)/gss_darwin_amd64 -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" github.com/spatialcurrent/go-simple-serializer/cmd/gss

bin/gss_linux_amd64: ## Build GSS CLI for Linux / amd64
	GOOS=linux GOARCH=amd64 go build -o $(DEST)/gss_linux_amd64 -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" github.com/spatialcurrent/go-simple-serializer/cmd/gss

bin/gss_windows_amd64.exe:  ## Build GSS CLI for Windows / amd64
	GOOS=windows GOARCH=amd64 go build -o $(DEST)/gss_windows_amd64.exe -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" github.com/spatialcurrent/go-simple-serializer/cmd/gss

bin/gss_linux_arm64: ## Build GSS CLI for Linux / arm64
	GOOS=linux GOARCH=arm64 go build -o $(DEST)/gss_linux_arm64 -gcflags="$(GCFLAGS)" -ldflags="$(LDFLAGS)" github.com/spatialcurrent/go-simple-serializer/cmd/gss

build_cli: bin/gss_darwin_amd64 bin/gss_linux_amd64 bin/gss_windows_amd64.exe bin/gss_linux_arm64  ## Build command line programs

#
# Shared Objects
#

bin/gss.so:  ## Compile Shared Object for current platform
	# https://golang.org/cmd/link/
	# CGO Enabled : https://github.com/golang/go/issues/24068
	CGO_ENABLED=1 go build -o $(DEST)/gss.so -buildmode=c-shared -ldflags "$(LDFLAGS)" -gcflags="$(GCFLAGS)" github.com/spatialcurrent/go-simple-serializer/plugins/gss

bin/gss_linux_amd64.so:  ## Compile Shared Object for Linux / amd64
	# https://golang.org/cmd/link/
	# CGO Enabled : https://github.com/golang/go/issues/24068
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o $(DEST)/gss_linux_amd64.so -buildmode=c-shared -ldflags "$(LDFLAGS)" -gcflags="$(GCFLAGS)" github.com/spatialcurrent/go-simple-serializer/plugins/gss

bin/gss_linux_armv7.so:  ## Compile Shared Object for Linux / ARMv7
	# LDFLAGS - https://golang.org/cmd/link/
	# CGO Enabled  - https://github.com/golang/go/issues/24068
	# GOARM/GOARCH Compatability Table - https://github.com/golang/go/wiki/GoArm
	# ARM Cross Compiler Required - https://www.acmesystems.it/arm9_toolchain
	GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc go build -ldflags "-linkmode external -extldflags -static" -o $(DEST)/gss_linux_armv7.so -buildmode=c-shared -ldflags "$(LDFLAGS)" -gcflags="$(GCFLAGS)" github.com/spatialcurrent/go-simple-serializer/plugins/gss

bin/gss_linux_armv8.so:   ## Compile Shared Object for Linux / ARMv8
	# LDFLAGS - https://golang.org/cmd/link/
	# CGO Enabled  - https://github.com/golang/go/issues/24068
	# GOARM/GOARCH Compatability Table - https://github.com/golang/go/wiki/GoArm
	# ARM Cross Compiler Required - https://www.acmesystems.it/arm9_toolchain
	# Dependencies - https://www.96boards.org/blog/cross-compile-files-x86-linux-to-96boards/
	GOOS=linux GOARCH=arm64 CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc go build -ldflags "-linkmode external -extldflags -static" -o $(DEST)/gss_linux_armv8.so -buildmode=c-shared -ldflags "$(LDFLAGS)" -gcflags="$(GCFLAGS)" github.com/spatialcurrent/go-simple-serializer/plugins/gss

build_so: bin/gss_linux_amd64.so bin/gss_linux_armv7.so bin/gss_linux_armv8.so  ## Build Shared Objects (.so)

#
# Android
#

bin/gss.aar:  ## Build Android Archive Library
	gomobile bind -target android -javapkg=com.spatialcurrent -o $(DEST)/gss.aar -gcflags="$(GCFLAGS)" github.com/spatialcurrent/go-simple-serializer/pkg/gss

build_android: bin/gss.arr  ## Build artifacts for Android

#
# JavaScript
#

dist/gss.mod.js:  ## Build JavaScript module
	gopherjs build -o dist/gss.mod.js github.com/spatialcurrent/go-simple-serializer/cmd/gss.mod.js

dist/gss.mod.min.js:  ## Build minified JavaScript module
	gopherjs build -m -o dist/gss.mod.min.js github.com/spatialcurrent/go-simple-serializer/cmd/gss.mod.js

dist/gss.global.js:  ## Build JavaScript library that attaches to global or window.
	gopherjs build -o dist/gss.global.js github.com/spatialcurrent/go-simple-serializer/cmd/gss.global.js

dist/gss.global.min.js:  ## Build minified JavaScript library that attaches to global or window.
	gopherjs build -m -o dist/gss.global.min.js github.com/spatialcurrent/go-simple-serializer/cmd/gss.global.js

build_javascript: dist/gss.mod.js dist/gss.mod.min.js dist/gss.global.js dist/gss.global.min.js  ## Build artifacts for JavaScript

test_javascript:  ## Run JavaScript tests
	npm run test

lint:  ## Lint JavaScript source code
	npm run lint

#
# Examples
#

bin/gss_example_c: bin/gss.so  ## Build C example
	mkdir -p bin && cd bin && gcc -o gss_example_c -I. ./../examples/c/main.c -L. -l:gss.so

bin/gss_example_cpp: bin/gss.so  ## Build C++ example
	mkdir -p bin && cd bin && g++ -o gss_example_cpp -I . ./../examples/cpp/main.cpp -L. -l:gss.so

run_example_c: bin/gss.so bin/gss_example_c  ## Run C example
	cd bin && LD_LIBRARY_PATH=. ./gss_example_c

run_example_cpp: bin/gss.so bin/gss_example_cpp  ## Run C++ example
	cd bin && LD_LIBRARY_PATH=. ./gss_example_cpp

run_example_python: bin/gss.so  ## Run Python example
	LD_LIBRARY_PATH=bin python examples/python/test.py

run_example_javascript: dist/gss.mod.min.js  ## Run JavaScript module example
	node examples/js/index.mod.js

## Clean

clean:  ## Clean artifacts
	rm -fr bin
	rm -fr dist
