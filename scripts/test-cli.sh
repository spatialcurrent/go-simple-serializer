#!/bin/bash

# =================================================================
#
# Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
# Released as open source under the MIT License.  See LICENSE file.
#
# =================================================================

set -euo pipefail

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

expectedFormats="bson,csv,fmt,go,gob,hcl,json,jsonl,properties,tags,toml,tsv,yaml"

testFormats() {
  formats=$("${DIR}/../bin/gss" formats -f csv)
  assertEquals "unexpected formats" "${expectedFormats}" "${formats}"
}

testCSVJSON() {
  local expected='[{"hello":"world"}]'
  local output=$(printf 'hello\nworld' | "${DIR}/../bin/gss" -i csv -o json)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

testCSVJSONL() {
  local expected='{"hello":"world"}'
  local output=$(printf 'hello\nworld' | "${DIR}/../bin/gss" -i csv -o jsonl)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

testJSONCSV() {
  local expected='hello\nworld'
  local output=$(echo '{"hello":"world"}' | "${DIR}/../bin/gss" -i json -o csv)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

testJSONJSON() {
  local expected='{"hello":"world"}'
  local output=$(echo '{"hello":"world"}' | "${DIR}/../bin/gss" -i  json -o json)
  assertEquals "unexpected output" "${expected}" "${output}"
}

testJSONLGOBJSONL() {
  local input='{"hello":"world"}'
  local expected='{"hello":"world"}'
  local output=$(echo "${input}" | gss -i jsonl -o gob | gss -i gob -t -o jsonl)
  assertEquals "unexpected output" "${expected}" "${output}"
}

testJSONArrayCSV() {
  local input='[{"a":"x"},{"b":"y"},{"c":"z"}]'
  local expected='a=x\nb=y\nc=z'
  local output=$(echo -e "${input}" | "${DIR}/../bin/gss" -i json -o tags)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

testJSONLTags() {
  local input='{"a":"x"}\n{"b":"y"}\n{"c":"z"}'
  local expected='a=x\nb=y\nc=z'
  local output=$(echo -e "${input}" | "${DIR}/../bin/gss" -i jsonl -o tags)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

testJSONLFmt() {
  local input='{"a":"x"}\n{"b":"y"}\n{"c":"z"}'
  local expected='map[a:x]\nmap[b:y]\nmap[c:z]'
  local output=$(echo -e "${input}" | "${DIR}/../bin/gss" -i jsonl -o fmt --output-format-specifier "%v")
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

testJSONGo() {
  local input='["a", "x"]'
  local expected='[]interface {}{"a", "x"}'
  local output=$(echo -e "${input}" | "${DIR}/../bin/gss" -i json -o go)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}

testJSONGoFit() {
  local input='["a", "x"]'
  local expected='[]string{"a", "x"}'
  local output=$(echo -e "${input}" | "${DIR}/../bin/gss" -i json -o go --output-fit)
  assertEquals "unexpected output" "$(echo -e "${expected}")" "${output}"
}


testHCLJSON() {
  local expected='{"data":[{"aws_caller_identity":[{"current":[{}]}]}]}'
  local output=$(echo 'data "aws_caller_identity" "current" {}' | "${DIR}/../bin/gss" -i  hcl -o json)
  assertEquals "unexpected output" "${expected}" "${output}"
}

testYAMLJSON() {
  local expected='[{"a":"x"},{"b":"y"},"foo",null]'
  local output=$(echo -e '---\na: x\n---\nb: "y"\n---\nfoo\n---\n' | "${DIR}/../bin/gss" -i  yaml -o json)
  assertEquals "unexpected output" "${expected}" "${output}"
}

oneTimeSetUp() {
  echo "Setting up"
  echo "Using temporary directory at ${SHUNIT_TMPDIR}"
}

oneTimeTearDown() {
  echo "Tearing Down"
}

# Load shUnit2.
. "${DIR}/shunit2"
