[![CircleCI](https://circleci.com/gh/spatialcurrent/go-simple-serializer/tree/master.svg?style=svg)](https://circleci.com/gh/spatialcurrent/go-simple-serializer/tree/master) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/go-simple-serializer)](https://goreportcard.com/report/spatialcurrent/go-simple-serializer)  [![GoDoc](https://godoc.org/github.com/spatialcurrent/go-simple-serializer?status.svg)](https://godoc.org/github.com/spatialcurrent/go-simple-serializer) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/go-simple-serializer/blob/master/LICENSE)

# go-simple-serializer

# Description

**go-simple-serializer** (aka GSS) is a simple library for serializing/deserializing objects.  GSS supports `bson`, `csv`, `tsv`, `hcl`, `hcl2`, `json`, `jsonl`, `properties`, `toml`, `yaml`.  `hcl` and `hcl2` implementation is fragile and very much in `alpha`.

Using cross compilers, this library can also be called by other languages.  This library is cross compiled into a Shared Object file (`*.so`).  The Shared Object file can be called by `C`, `C++`, and `Python` on Linux machines.  See the examples folder for patterns that you can use.  This library is also compiled to pure `JavaScript` using [GopherJS](https://github.com/gopherjs/gopherjs).

**Packages**

The main public api for GSS is the `gss` package.  However, the library does ship with lower-level packages that can be imported directly as well.

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

You can use the command line tool to convert between serialization formats.

```
gss is a simple program for serializing/deserializing data.

Usage:
  gss [flags]
  gss [command]

Available Commands:
  completion  Generates bash completion scripts
  help        Help about any command
  version     print version information to stdout

Flags:
  -a, --async                   async processing
  -h, --help                    help for gss
  -c, --input-comment string    The input comment character, e.g., #.  Commented lines are not sent to output.
  -i, --input-format string     The input format: bson, csv, tsv, hcl, hcl2, json, jsonl, properties, toml, yaml
      --input-header strings    The input header if the stdin input has no header.
      --input-lazy-quotes       allows lazy quotes for CSV and TSV
  -l, --input-limit int         The input limit (default -1)
      --input-skip-lines int    The number of lines to skip before processing
  -t, --input-trim              trim input lines
  -o, --output-format string    The output format: bson, csv, tsv, hcl, hcl2, json, jsonl, properties, toml, yaml
      --output-header strings   The output header if the stdout output has no header.
  -n, --output-limit int        the output limit (default -1)
  -p, --output-pretty           print pretty output
  -s, --output-sorted           sort output
      --verbose                 Print debug info to stdout

Use "gss [command] --help" for more information about a command.
```

**Go**

You can import **go-simple-serializer** as a library with:

```go
import (
  "github.com/spatialcurrent/go-simple-serializer/gss"
)
```

The `Convert`, `Deserialize`, and `Serialize` functions are the core functions to use.

```go
...
  output_string, err := gss.Convert(input_string, input_format, input_header, input_comment, output_format, verbose)
...
  input_type, err := GetType(input_string, input_format)
	if err != nil {
		return "", errors.Wrap(err, "error creating new object for format "+input_format)
	}
  output, err := gss.Deserialize(input, format, input_header, input_comment, input_type, verbose)
...
  output_string, err := gss.Serialize(input, format)
...
```

See [gss](https://godoc.org/github.com/spatialcurrent/go-simple-serializer/gss) in GoDoc for information on how to use Go API.

**JavaScript**

```html
<html>
  <head>
    <script src="https://...gss.js"></script>
  </head>
  <body>
    <script>
      var input = "{\"a\":1}";
      var output = gss.convert(input, "json", "yaml", )
      ...
      // You can also pass the input header for a csv/tsv that has none
      var output = gss.convert(input, "csv", "json", {"header": ["a","b"]})
    </script>
  </body>
</html>
```

**Android**

The `go-simple-serializer` code is available for use in Android applications under `com.spatialcurrent.gss`.  For example,

```java
import com.spatialcurrent.gss.Gss;
...
  String output_format = Gss.convert(input_string, input_format, input_header, input_comment, output_format, verbose);
...
```

**C**

A variant of the `Convert` function is exported in a Shared Object file (`*.so`), which can be called by `C`, `C++`, and `Python` programs on Linux machines.  For example:

```
char *input_string = "<YOUR INPUT>";
char *output_string;
err = Convert(input_string, input_format, input_header_csv, input_comment, output_format, &output_string);
```

The Go function definition defined in `plugins/gss/main.go` uses `*C.char` for all input except `output_string` which uses a double pointer (`**C.char`) to write to the output.

```
func Convert(input_string *C.char, input_format *C.char, input_header *C.char, input_comment *C.char, output_format *C.char, output_string **C.char) *C.char
```

For complete patterns for `C`, `C++`, and `Python`, see the `go-simpler-serializer/examples` folder.

# Releases

**go-simple-serializer** is currently in **alpha**.  See releases at https://github.com/spatialcurrent/go-simple-serializer/releases.

# Examples

`.gitignore` file to jsonl

```
cat .gitignore | ./gss -i csv -h pattern -o jsonl
```

Get language from `.travis.yml` and set to variable

```
language=$(cat .travis.yml | ./gss_linux_amd64 -i yaml -o json -c '#' | jq .language -r)
```

# Building

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

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-simple-serializer/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
