// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

require('./../../dist/gss.global.min.js');
const { serialize, deserialize, convert, formats } = global.gss;
console.log('Formats:', formats);
console.log();
console.log("************************************");
console.log();

const obj = {
  "a": "1",
  "b": "c",
  "x": ["foo", "bar"]
}

const arr = [
  {
    "a": "x",
    "b": "y",
    "c": "z"
  },
  {
    "d": "g",
    "e": "h",
    "f": "i"
  }
];

console.log('Input:');
console.log(obj);
console.log();
// Destructure return value
var { str, err } = serialize(arr, "json", {"pretty": true});
console.log('Output:');
console.log(str);
console.log();
console.log('Error:');
console.log(err);
console.log();
console.log("************************************");
console.log();

console.log('Input:');
console.log(obj);
console.log();
// Destructure return value
var { str, err } = serialize(obj, "blah");
console.log('Output:');
console.log(str);
console.log();
console.log('Error:');
console.log(err);
console.log();
