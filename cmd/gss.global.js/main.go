// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

// gss.global.js is the package for go-simple-serializer (GSS) that adds GSS functions to the global scope under the "gss" variable.
//
// In Node, depending on where require is called and the build system used, the functions may need to be required at the top of each module file.
// In a web browser, gss can be made available to the entire web page.
// The functions are defined in the Exports variable in the gssjs package.
//
// Usage
//	// Below is a simple set of examples of how to use this package in a JavaScript application.
//
//	// load functions into global scope
//	// require('./dist/gss.global.min.js);
//
//	// Serialize an object to a string.
//	// Returns an object, which can be destructured to the formatted string and error as a string.
//	// If there is no error, then err will be null.
//	var { str, err } = gss.serialize(object, format, options);
//
//	// Deserialize a formatted string to an object.
//	// Returns an object, which can be destructured to the object and error as a string.
//	// If there is no error, then err will be null.
//	var { obj, err } = gss.deserialize(str, format, options);
//
//	// Convert a formatted string to a new format.
//	// Returns an object, which can be destructured to the string and error as a string.
//	// If there is no error, then err will be null.
//	var { str, err } = gss.convert(str, inputFormat, ouputFormat, inputOptions, outputOptions);
//
// References
//	- https://godoc.org/pkg/github.com/spatialcurrent/go-simple-serializer/pkg/gssjs/
//	- https://nodejs.org/api/globals.html#globals_global_objects
//	- https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects
package main

import (
	"github.com/gopherjs/gopherjs/js"

	"github.com/spatialcurrent/go-simple-serializer/pkg/gssjs"
)

func main() {
	js.Global.Set("gss", gssjs.Exports)
}
