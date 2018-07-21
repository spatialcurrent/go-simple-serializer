[![Build Status](https://travis-ci.org/spatialcurrent/go-simple-serializer.svg)](https://travis-ci.org/spatialcurrent/go-simple-serializer) [![GoDoc](https://godoc.org/github.com/spatialcurrent/go-simple-serializer?status.svg)](https://godoc.org/github.com/spatialcurrent/go-simple-serializer)

# go-simple-serializer

# Description

**go-simple-serializer** (aka GSS) is a simple library for serializing/deserializing objects.  GSS supports `csv`, `hcl`, `hcl2`, `json`, `jsonl`, `toml`, `yaml`.  `hcl` and `hcl2` implementation is fragile and very much in `alpha`.

# Usage

**CLI**

You can use the command line tool to convert between formats.

```
Usage: gss -i INPUT_FORMAT -o OUTPUT_FORMAT
Options:
  -help
    	Print help.
  -i string
    	The input format: csv, hcl, hcl2, json, jsonl, toml, yaml
  -o string
    	The output format: csv, hcl, hcl2, json, jsonl, toml, yaml
  -version
    	Prints version to stdout.
```

**Go**

You can import **go-simple-serializer** as a library with:

```go
import (
  "github.com/spatialcurrent/go-simple-serializer/gss"
)
...
  output_string, err := gss.Convert(input_string, input_format, output_format)
...
```

**JavaScript**

```html
<html>
  <head>
    <script src="https://...gss.js"></script>
  </head>
  <body>
    <script>
      var input = "{\"a\":1}";
      var output = gss.convert(input, "json", "yaml")
      ...
    </script>
  </body>
</html>
```

**Android**

The `go-simple-serializer` code is available under `com.spatialcurrent.gss`.  For example,

```java
import com.spatialcurrent.gss.Gss;
...
  String output_format = Gss.convert(input_string, input_format, output_format);
...
```

# Releases

**go-simple-serializer** is currently in **alpha**.  See releases at https://github.com/spatialcurrent/go-simple-serializer/releases.

# Building

**CLI**

The `build_cli.sh` script is used to build executables for Linux and Windows.

**JavaScript**

You can compile GSS to pure JavaScript with the `scripts/build_javascript.sh` script.

**Android**

The `build_android.sh` script is used to build an [Android Archive](https://developer.android.com/studio/projects/android-library) (AAR) file and associated Javadocs.

# Contributing

[Spatial Current, Inc.](https://spatialcurrent.io) is currently accepting pull requests for this repository.  We'd love to have your contributions!  Please see [Contributing.md](https://github.com/spatialcurrent/go-simple-serializer/blob/master/CONTRIBUTING.md) for how to get started.

# License

This work is distributed under the **MIT License**.  See **LICENSE** file.
