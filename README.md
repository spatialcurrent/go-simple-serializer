[![CircleCI](https://circleci.com/gh/spatialcurrent/spatialcurrent/go-simple-serializer/tree/main.svg?style=svg)](https://circleci.com/gh/spatialcurrent/spatialcurrent/go-simple-serializer/tree/main)
[![Go Report Card](https://goreportcard.com/badge/spatialcurrent/spatialcurrent/go-simple-serializer?style=flat-square)](https://goreportcard.com/report/github.com/spatialcurrent/spatialcurrent/go-simple-serializer)
[![PkgGoDev](https://pkg.go.dev/badge/mod/github.com/spatialcurrent/spatialcurrent/go-simple-serializer)](https://pkg.go.dev/github.com/spatialcurrent/spatialcurrent/go-simple-serializer)
[![License](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/spatialcurrent/go-simple-serializer/blob/master/LICENSE)

# go-simple-serializer

# Description

**go-simple-serializer** (aka GSS) is a simple library to easily convert data between formats that aims to decrease the burden on developers to support multiple serialization formats in their applications.  GSS supports a variety of operating systems, architectures, and use cases.  A CLI is released for Microsoft Windows, Linux distributions, and [Darwin](https://en.wikipedia.org/wiki/Darwin_%28operating_system%29) platforms.

Using cross compilers, this library can also be called by other languages, including `C`, `C++`, and `Python`.  This library is cross compiled into a Shared Object file (`*.so`), which can be called by `C`, `C++`, and `Python` on Linux machines.  See the examples folder for patterns that you can use.

**Formats**

GSS supports many common formats, including CSV, JSON, and YAML.  Pull requests to support other formats are welcome!  See the [Formats.md](docs/Formats.md) document for a full list of supported formats.

**Packages**

The main public api for GSS is in the `gss` package.  However, this library does ship with internal packages under `/pkg/...` that can be imported and used directly.

## Platforms

The following platforms are supported.  Pull requests to support other platforms are welcome!

| GOOS | 386 | amd64 | arm | arm64 |
| ---- | --- | ----- | --- | ----- |
| darwin | - | ✓ | - | - |
| freebsd | ✓ | ✓ | ✓ | - |
| linux | ✓ | ✓ | ✓ | ✓ |
| openbsd | ✓ | ✓ | - | - |
| solaris | - | ✓ | - | - |
| windows | ✓ | ✓ | - | - |

## Releases

Find releases for the supported platforms at [https://github.com/spatialcurrent/go-simple-serializer/releases](https://github.com/spatialcurrent/go-simple-serializer/releases).  See the **Building** section below to build for another platform from source.

# Usage

**CLI**

See the [CLI.md](docs/CLI.md) document for detailed usage and examples.

**Go**

You can import the public *gss** package with:

```go
import (
  "github.com/spatialcurrent/go-simple-serializer/pkg/gss"
)
```

You can also import one of the internal packages such as **tags** with:

```go
import (
  "github.com/spatialcurrent/go-simple-serializer/pkg/tags"
)
```

See [gss](https://pkg.go.dev/github.com/spatialcurrent/go-simple-serializer/pkg/gss) in the docs for information on how to use Go API.

**Android**

The `go-simple-serializer` code is available for use in Android applications under `com.spatialcurrent.gss`.  For example,

```java
import com.spatialcurrent.gss.Gss;
...
  String output_format = Gss.convert(input_string, input_format, input_header, input_comment, output_format, verbose);
...
```

**C**

A variant of the `Convert` function is exported in a Shared Object file (`*.so`), which can be called by `C`, `C++`, and `Python` programs on Linux machines.  For complete patterns for `C`, `C++`, and `Python`, see the `examples` folder in this repo.

# Examples

**CLI**

`.gitignore` file to jsonl

```shell
cat .gitignore | gss -i csv --input-header path -o json
```

Get language from [CircleCI](https://circleci.com/) config.

```shell
cat .circleci/config.yml | gss -i yaml -o json -c '#' | jq -r .version
```

Convert list of files to JSON Lines

```shell
find . -name '*.go' | gss -i csv --input-header path -o jsonl
```

**Go**

See the examples in the [docs](https://godoc.org/github.com/spatialcurrent/go-simple-serializer).

**C**

See the `examples/c/main.c` file.  You can run the example with `make run_example_c`.

**C++**

See the `examples/cpp/main.cpp` file.  You can run the example with `make run_example_cpp`.

**Python**

See the `examples/python/test.py` file.  You can run the example with `make run_example_python`.

# Building

Use `make help` to see help information for each target.

**CLI**

The `make build_cli` script is used to build executables for Linux and Windows.

**Android**

The `make build_android` script is used to build an [Android Archive](https://developer.android.com/studio/projects/android-library) (AAR) file and associated Javadocs.

**Shared Object**

The `make build_so` script is used to build a Shared Object (`*.go`), which can be called by `C`, `C++`, and `Python` on Linux machines.

**Changing Destination**

The default destination for build artifacts is `go-simple-serializer/bin`, but you can change the destination with an environment variable.  For building on a Chromebook consider saving the artifacts in `/usr/local/go/bin`, e.g., `DEST=/usr/local/go/bin make build_cli`

# Testing

**CLI**

To run CLI tests use `make test_cli`, which uses [shUnit2](https://github.com/kward/shunit2).  If you recive a `shunit2:FATAL Please declare TMPDIR with path on partition with exec permission.` error, you can modify the `TMPDIR` environment variable in line or with `export TMPDIR=<YOUR TEMP DIRECTORY HERE>`. For example:

```
TMPDIR="/usr/local/tmp" make test_cli
```

**Go**

To run Go tests using `make test_go` or (`bash scripts/test.sh`), which runs unit tests, `go vet`, `go vet with shadow`, [errcheck](https://github.com/kisielk/errcheck), [staticcheck](https://staticcheck.io/), and [misspell](https://github.com/client9/misspell).

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/spatialcurrent/go-simple-serializer/blob/main/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
