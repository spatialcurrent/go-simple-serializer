# CLI

- [Formats](#formats) - list of supported file formats
- [Platforms](#platforms) - list of supported platforms
- [Releases](#releases) - where to find an executable
- [Examples](#examples)  - detailed usage exampels
- [Examples](#building) - how to build the CLI
- [Testing](#testing) - test the CLI
- [Troubleshooting](#Troubleshooting) - how to troubleshoot common errors

## Usage

The command line tool, `gss`, can be used to easily read and write compressed resources by uri.

### Formats

The following file formats are supported.  Pull requests to support other formats are welcome!

| Format | Read |  Write | Stream | Description |
| ---- | ------ |  ------ | ------ | ------ |
| bson | ✓ | ✓ | - | [Binary JSON](https://en.wikipedia.org/wiki/BSON) |
| csv | ✓ | ✓ | ✓ | [Comma-Separated Values](https://en.wikipedia.org/wiki/Comma-separated_values) |
| fmt | - | ✓ | ✓ | [fmt](https://godoc.org/fmt) |
| go | - | ✓ | ✓ | Go (format specifier: "%#v") |
| gob | ✓ | ✓ | ✓ | [gob](https://godoc.org/encoding/gob) |
| hcl | ✓ | - | - | [HashiCorp Configuration Language](https://github.com/hashicorp/hcl) |
| json | ✓ | ✓ | - | [JSON](http://json.org/) |
| jsonl | ✓ | ✓ | ✓ | [JSON Lines](http://jsonlines.org/) |
| properties | ✓ | ✓ | - |[Properties](https://en.wikipedia.org/wiki/.properties) |
| tags | ✓ | ✓ | ✓ | single-line series of key=value tags |
| toml | ✓ | ✓ | - | [TOML](https://github.com/toml-lang/toml) |
| tsv | ✓ | ✓ | - |[ Tab-Separated Values](https://en.wikipedia.org/wiki/Tab-separated_values) |
| yaml | ✓ | ✓ | ✓ | [YAML](https://yaml.org/) |


### Platforms

The following platforms are supported.  Pull requests to support other platforms are welcome!

| GOOS | GOARCH |
| ---- | ------ |
| darwin | amd64 |
| linux | amd64 |
| windows | amd64 |
| linux | arm64 |

## Releases

**gss** is currently in **alpha**.  Find releases at [https://github.com/spatialcurrent/go-reader-writer/releases](https://github.com/spatialcurrent/go-reader-writer/releases).  See the **Building** section below to build from scratch.

**Darwin**

- `gss_darwin_amd64` - CLI for Darwin on amd64 (includes `macOS` and `iOS` platforms)

**Linux**

- `gss_linux_amd64` - CLI for Linux on amd64
- `gss_linux_amd64` - CLI for Linux on arm64

**Windows**

- `gss_windows_amd64.exe` - CLI for Windows on amd64

# Examples

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

# Building

Use `make build_cli` to build executables for Linux and Windows.

**Changing Destination**

The default destination for build artifacts is `gss/bin`, but you can change the destination with an environment variable.  For building on a Chromebook consider saving the artifacts in `/usr/local/go/bin`, e.g., `DEST=/usr/local/go/bin make build_cli`

## Testing

To run CLI testes use `make test_cli`, which uses [shUnit2](https://github.com/kward/shunit2).  If you recive a `shunit2:FATAL Please declare TMPDIR with path on partition with exec permission.` error, you can modify the `TMPDIR` environment variable in line or with `export TMPDIR=<YOUR TEMP DIRECTORY HERE>`. For example:

```
TMPDIR="/usr/local/tmp" make test_cli
```

## Troubleshooting

### no such file or directory

#### Example

```text
error opening resource at uri %q: error opening file for writing at path %q: open %s: no such file or directory
```

#### Solution

This error typically occurs when a parent directory of an output file does not exist.  Use the `--output-mkdirs` command line flag to allow gss to create parent directories for output files as needed.

