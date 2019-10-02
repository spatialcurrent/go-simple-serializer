# Formats

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
| tsv | ✓ | ✓ | ✓ |[ Tab-Separated Values](https://en.wikipedia.org/wiki/Tab-separated_values) |
| yaml | ✓ | ✓ | - | [YAML](https://yaml.org/) |