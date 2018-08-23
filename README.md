[![Build Status](https://travis-ci.org/spatialcurrent/go-simple-serializer.svg)](https://travis-ci.org/spatialcurrent/go-simple-serializer) [![Go Report Card](https://goreportcard.com/badge/spatialcurrent/go-simple-serializer)](https://goreportcard.com/report/spatialcurrent/go-simple-serializer)  [![GoDoc](https://godoc.org/github.com/spatialcurrent/go-simple-serializer?status.svg)](https://godoc.org/github.com/spatialcurrent/go-simple-serializer) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://github.com/spatialcurrent/go-simple-serializer/blob/master/LICENSE)

# go-simple-serializer

# Description

**go-simple-serializer** (aka GSS) is a simple library for serializing/deserializing objects.

GSS supports `bson`, `csv`, `tsv`, `hcl`, `hcl2`, `json`, `jsonl`, `properties`, `toml`, `yaml`.  `hcl` and `hcl2` implementation is fragile and very much in `alpha`.

Using cross compilers, this library can also be called by other languages.  This library is cross compiled into a Shared Object file (`*.so`).  The Shared Object file can be called by `C`, `C++`, and `Python` on Linux machines.  See the examples folder for patterns that you can use.  This library is also compiled to pure `JavaScript` using [GopherJS](https://github.com/gopherjs/gopherjs).

# Usage

**CLI**

You can use the command line tool to convert between formats.

```
Usage: gss -i INPUT_FORMAT -o OUTPUT_FORMAT [-h HEADER] [-c COMMENT]
Options:
  -c string
    	The input comment character, e.g., #.  Commented lines are not sent to output.
  -h string
    	The input header if the stdin input has no header.
  -help
    	Print help.
  -i string
    	The input format: bson, csv, tsv, hcl, hcl2, json, jsonl, properties, toml, yaml
  -o string
    	The output format: bson, csv, tsv, hcl, hcl2, json, jsonl, properties, toml, yaml
  -version
    	Prints version to stdout.
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
  output_string, err := gss.Convert(input_string, input_format, input_header, input_comment, output_format)
...
  output = map[string]interface{}{}
  err := gss.Deserialize(input, format, input_header, input_comment, &output)
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
  String output_format = Gss.convert(input_string, input_format, input_header, input_comment, output_format);
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

The `build_cli.sh` script is used to build executables for Linux and Windows.

**JavaScript**

You can compile GSS to pure JavaScript with the `scripts/build_javascript.sh` script.

**Android**

The `build_android.sh` script is used to build an [Android Archive](https://developer.android.com/studio/projects/android-library) (AAR) file and associated Javadocs.

**Shared Object**

The `build_so.sh` script is used to build a Shared Object (`*.go`), which can be called by `C`, `C++`, and `Python` on Linux machines.

**Changing Destination**

The default destination for build artifacts is `go-simple-serializer/bin`, but you can change the destination with a CLI argument.  For building on a Chromebook consider saving the artifacts in `/usr/local/go/bin`, e.g., `bash scripts/build_cli.sh /usr/local/go/bin`

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-simple-serializer/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
