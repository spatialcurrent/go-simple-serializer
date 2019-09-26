[![CircleCI](https://circleci.com/gh/spatialcurrent/go-simple-serializer/tree/master.svg?style=svg)](https://circleci.com/gh/spatialcurrent/go-simple-serializer/tree/master) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/go-simple-serializer)](https://goreportcard.com/report/spatialcurrent/go-simple-serializer)  [![GoDoc](https://godoc.org/github.com/spatialcurrent/go-simple-serializer?status.svg)](https://godoc.org/github.com/spatialcurrent/go-simple-serializer) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/go-simple-serializer/blob/master/LICENSE)

# go-simple-serializer

# Description

**go-simple-serializer** (aka GSS) is a simple library for serializing/deserializing objects that aims to decrease the burden on developers to support multiple serialization formats in their applications.  GSS supports a variety of operating systems, architectures, and use cases.  A CLI is released for Microsoft Windows, Linux distributions, and [Darwin](https://en.wikipedia.org/wiki/Darwin_%28operating_system%29) platforms.

Using cross compilers, this library can also be called by other languages, including `C`, `C++`, `Python`, and `JavaScript`.  This library is cross compiled into a Shared Object file (`*.so`), which can be called by `C`, `C++`, and `Python` on Linux machines.  This library is also compiled to pure `JavaScript` using [GopherJS](https://github.com/gopherjs/gopherjs), which can be called by [Node.js](https://nodejs.org) and loaded in the browser.  See the examples folder for patterns that you can use.

**Formats**

GSS supports the following formats.

| Format | Description |
| ---- | ------ |
| bson | [Binary JSON](https://en.wikipedia.org/wiki/BSON) |
| csv | [Comma-Separated Values](https://en.wikipedia.org/wiki/Comma-separated_values) |
| hcl | HashiCorp Configuration Language |
| hcl2 | HashiCorp Configuration Language (version 2.x) |
| json | [JSON](http://json.org/) |
| jsonl | [JSON Lines](http://jsonlines.org/) |
| properties | [Properties](https://en.wikipedia.org/wiki/.properties) |
| tags | single-line key=value tags |
| toml | [TOML](https://github.com/toml-lang/toml) |
| tsv | Tab-Separated Values |
| yaml | [YAML](https://yaml.org/) |

`hcl` and `hcl2` implementation is fragile and very much in `alpha`.  The other formats are well-supported.

**Packages**

The main public api for GSS is in the `gss` package.  However, this library does ship with lower-level packages that can be imported directly as well.

| Package | Purpose |
| ---- | ------ |
| bson | Binary JSON |
| escaper | Escape/unescape strings |
| gss | The main public API |
| inspector | Reusable functions for inspecting objects |
| iterator | Wrapper for iterable formats |
| json | JSON |
| jsonl | JSON Lines |
| properties | Properties Files |
| scanner | Scanning through a stream of bytes |
| splitter | Creating custom bufio.SplitFunc |
| sv | Separated-Values formats, i.e., CSV and TSV. |
| toml | TOML |
| yaml | YAML |

# Usage

**CLI**

The command line tool, `gss`, can be used to easily covert data between formats.  We currently support the following platforms.

| GOOS | GOARCH |
| ---- | ------ |
| darwin | amd64 |
| linux | amd64 |
| windows | amd64 |
| linux | arm64 |

Pull requests to support other platforms are welcome!  See the [examples](#examples) section below for usage.

**Go**

You can install the go-simple-serializer packages with.


```shell
go get -u -d github.com/spatialcurrent/go-simple-serializer/...
```

You can then import the main public API with `import "github.com/spatialcurrent/go-simple-serializer/pkg/gss"` or one of the underlying packages, e.g., `import "github.com/spatialcurrent/go-simple-serializer/pkg/tags"`.

See [go-simple-serializer](https://godoc.org/github.com/spatialcurrent/go-simple-serializer) in GoDoc for API documentation and examples.

**Node**

GSS is built as a module.  In modern JavaScript, the module can be imported using [destructuring assignment](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Destructuring_assignment).

```javascript
const { serialize, deserialize, convert, formats } = require('./dist/gss.mod.min.js');
```

In legacy JavaScript, you can use the `gss.global.js` file that simply adds `gss` to the global scope.

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

# Releases

**go-simple-serializer** is currently in **alpha**.  See releases at https://github.com/spatialcurrent/go-simple-serializer/releases.  See the **Building** section below to build from scratch.

**JavaScript**

- `gss.global.js`, `gss.global.js.map` - JavaScript global build  with source map
- `gss.global.min.js`, `gss.global.min.js.map` - Minified JavaScript global build with source map
- `gss.mod.js`, `gss.mod.js.map` - JavaScript module build  with source map
- `gss.mod.min.js`, `gss.mod.min.js.map` - Minified JavaScript module with source map

**Darwin**

- `gss_darwin_amd64` - CLI for Darwin on amd64 (includes `macOS` and `iOS` platforms)

**Linux**

- `gss_linux_amd64` - CLI for Linux on amd64
- `gss_linux_amd64` - CLI for Linux on arm64
- `gss_linux_amd64.h`, `gss_linuxamd64.so` - Shared Object for Linux on amd64
- `gss_linux_armv7.h`, `gss_linux_armv7.so` - Shared Object for Linux on ARMv7
- `gss_linux_armv8.h`, `gss_linux_armv8.so` - Shared Object for Linux on ARMv8

**Windows**

- `gss_windows_amd64.exe` - CLI for Windows on amd64

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

See the examples in [GoDoc](https://godoc.org/github.com/spatialcurrent/go-simple-serializer).

**C**

See the `examples/c/main.c` file.  You can run the example with `make run_example_c`.

**C++**

See the `examples/cpp/main.cpp` file.  You can run the example with `make run_example_cpp`.

**Python**

See the `examples/python/test.py` file.  You can run the example with `make run_example_python`.

**JavaScript**

See the `examples/js/index.js` file.  You can run the example with `make run_example_javascript`.

# Building

Use `make help` to see help information for each target.

**CLI**

The `make build_cli` script is used to build executables for Linux and Windows.

**JavaScript**

You can compile GSS to pure JavaScript with the `make build_javascript` script.

**Android**

The `make build_android` script is used to build an [Android Archive](https://developer.android.com/studio/projects/android-library) (AAR) file and associated Javadocs.

**Shared Object**

The `make build_so` script is used to build a Shared Object (`*.go`), which can be called by `C`, `C++`, and `Python` on Linux machines.

**Changing Destination**

The default destination for build artifacts is `go-simple-serializer/bin`, but you can change the destination with an environment variable.  For building on a Chromebook consider saving the artifacts in `/usr/local/go/bin`, e.g., `DEST=/usr/local/go/bin make build_cli`

# Testing

**CLI**

To run CLI testes use `make test_cli`, which uses [shUnit2](https://github.com/kward/shunit2).  If you recive a `shunit2:FATAL Please declare TMPDIR with path on partition with exec permission.` error, you can modify the `TMPDIR` environment variable in line or with `export TMPDIR=<YOUR TEMP DIRECTORY HERE>`. For example:

```
TMPDIR="/usr/local/tmp" make test_cli
```

**Go**

To run Go tests use `make test_go` (or `bash scripts/test.sh`), which runs unit tests, `go vet`, `go vet with shadow`, [errcheck](https://github.com/kisielk/errcheck), [ineffassign](https://github.com/gordonklaus/ineffassign), [staticcheck](https://staticcheck.io/), and [misspell](https://github.com/client9/misspell).

**JavaScript**

To run JavaScript tests, first install [Jest](https://jestjs.io/) using `make deps_javascript`, use [Yarn](https://yarnpkg.com/en/), or another method.  Then, build the JavaScript module with `make build_javascript`.  To run tests, use `make test_javascript`.  You can also use the scripts in the `package.json`.

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-simple-serializer/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
