// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// gss.mod.js is the package for go-simple-serializer (GSS) that is built as a JavaScript module.
// In modern JavaScript, the module can be imported using destructuring assignment.
// The functions are defined in the Exports variable in the gssjs package.
//
// Usage
//	// Below is a simple set of examples of how to use this package in a JavaScript application.
//
//	// load functions into current scope
//	const { serialize, deserialize, convert, formats } = require('./dist/gss.global.min.js);
//
//	// Serialize an object to a string.
//	// Returns an object, which can be destructured to the formatted string and error as a string.
//	// If there is no error, then err will be null.
//	var { str, err } = serialize(object, format, options);
//
//	// Deserialize a formatted string to an object.
//	// Returns an object, which can be destructured to the object and error as a string.
//	// If there is no error, then err will be null.
//	var { obj, err } = deserialize(str, format, options);
//
//	// Convert a formatted string to a new format.
//	// Returns an object, which can be destructured to the string and error as a string.
//	// If there is no error, then err will be null.
//	var { str, err } = convert(str, inputFormat, ouputFormat, inputOptions, outputOptions);
//
// References
//	- https://godoc.org/pkg/github.com/spatialcurrent/go-simple-serializer/pkg/gssjs/
//	- https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/Destructuring_assignment
package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/spatialcurrent/go-simple-serializer/pkg/gssjs"
)

func main() {
	jsModuleExports := js.Module.Get("exports")
	for name, value := range gssjs.Exports {
		jsModuleExports.Set(name, value)
	}
}
