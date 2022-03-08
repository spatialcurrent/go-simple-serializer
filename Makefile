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

deps_arm:  ## Install dependencies to cross-compile to ARM
	# ARMv7
	apt-get install -y libc6-armel-cross libc6-dev-armel-cross binutils-arm-linux-gnueabi libncurses5-dev gcc-arm-linux-gnueabi g++-arm-linux-gnueabi
  # ARMv8
	apt-get install gcc-aarch64-linux-gnu g++-aarch64-linux-gnu

#
# Go building, formatting, testing, and installing
#

fmt:  ## Format Go source code
	go fmt $$(go list ./... )

.PHONY: imports
imports: bin/goimports ## Update imports in Go source code
	bin/goimports -w -local github.com/spatialcurrent/go-simple-serializer,github.com/spatialcurrent $$(find . -iname '*.go')

vet: ## Vet Go source code
	go vet github.com/spatialcurrent/go-simple-serializer/pkg/... # vet packages
	go vet github.com/spatialcurrent/go-simple-serializer/cmd/... # vet commands

tidy: ## Tidy Go source code
	go mod tidy

.PHONY: test_go
test_go: bin/errcheck bin/misspell bin/staticcheck bin/shadow ## Run Go tests
	bash scripts/test.sh

.PHONY: test_cli
test_cli: bin/gss ## Run CLI tests
	bash scripts/test-cli.sh

install:  ## Install the CLI on current platform
	go install github.com/spatialcurrent/go-simple-serializer/cmd/gss

#
# Command line Programs
#

bin/errcheck:
	go build -o bin/errcheck github.com/kisielk/errcheck

bin/goimports:
	go build -o bin/goimports golang.org/x/tools/cmd/goimports

bin/gox:
	go build -o bin/gox github.com/mitchellh/gox

bin/misspell:
	go build -o bin/misspell github.com/client9/misspell/cmd/misspell

bin/staticcheck:
	go build -o bin/staticcheck honnef.co/go/tools/cmd/staticcheck

bin/shadow:
	go build -o bin/shadow golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow

bin/gss: ## Build CLI for Darwin / amd64
	go build -o bin/gss github.com/spatialcurrent/go-simple-serializer/cmd/gss

bin/gss_linux_amd64: bin/gox ## Build CLI for Darwin / amd64
	scripts/build-release linux amd64

.PHONY: build
build: bin/gss

.PHONY: build_release
build_release: bin/gox
	scripts/build-release

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


## Clean

clean:  ## Clean artifacts
	rm -fr bin
	rm -fr dist
