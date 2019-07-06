// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package escaper

import (
	"fmt"
)

// This example shows that when no prefix is set, then Escape is a no-op.
func ExampleEscaper_Prefix_off() {
	in := "Hello*World"
	e := New().Prefix("") // set the prefix to a blank string
	out := e.Escape(in)   // escape has no effect since their is no prefix set
	fmt.Println(out)
	// Output: Hello*World
}

// This example shows that any string can be used as an escape prefix.
func ExampleEscaper_Prefix_asterix() {
	in := "Hello*World"
	e := New().Prefix("*") // set the prefix to a custom string (asterix)
	out := e.Escape(in)    // escape escapes the prefix itself resulting in 2 asterixis
	fmt.Println(out)
	// Output: Hello**World
}

// This example shows that sub can be called to set substrings.
func ExampleEscaper_Sub_simple() {
	in := "Hello Beautiful=World:Again"
	e := New().Prefix("\\").Sub("=").Sub(":")
	out := e.Escape(in)
	fmt.Println(out)
	// Output: Hello Beautiful\=World\:Again
}

// This example shows that sub is a variadic function and accepts multiple substrings.
func ExampleEscaper_Sub_multiple() {
	in := "Hello Beautiful=World:Again"
	e := New().Prefix("\\").Sub(" ", "=", ":")
	out := e.Escape(in)
	fmt.Println(out)
	// Output: Hello\ Beautiful\=World\:Again
}

// This examples shows an simple use case for Escape.
func ExampleEscaper_Escape_simple() {
	in := "Hello Beautiful\\World Again"
	e := New().Prefix("\\")
	out := e.Escape(in)
	fmt.Println(out)
	// Output: Hello Beautiful\\World Again
}

// This example shows that you can specify the substrings to escape.
func ExampleEscaper_Escape_equal() {
	in := "Hello=World"
	e := New().Prefix("\\").Sub("=")
	out := e.Escape(in)
	fmt.Println(out)
	// Output: Hello\=World
}

// This example shows you can escape multiple substrings.
func ExampleEscaper_Escape_multiple() {
	in := "Hello Beautiful=World:Again"
	e := New().Prefix("\\").Sub(" ", "=", ":")
	out := e.Escape(in)
	fmt.Println(out)
	// Output: Hello\ Beautiful\=World\:Again
}

// This exmaple shows that you can unescape text.
func ExampleEscaper_Unescape_simple() {
	in := "Hello Beautiful\\\\World Again"
	e := New().Prefix("\\")
	out := e.Unescape(in)
	fmt.Println(out)
	// Output: Hello Beautiful\World Again
}

// This example shows you can unescape new line characters.
func ExampleEscaper_Unescape_newline() {
	// You can escape new line characters
	in := "Hello Beautiful\\\nWorld Again"
	e := New().Prefix("\\").Sub("\n")
	out := e.Unescape(in)
	fmt.Println(out)
	// Output: Hello Beautiful
	//World Again
}

// This example shows you can unescape multiple substrings.
func ExampleEscaper_Unescape_multiple() {
	in := "Hello\\ Beautiful\\=World\\:Again"
	e := New().Prefix("\\").Sub(" ", "=", ":")
	out := e.Unescape(in)
	fmt.Println(out)
	// Output: Hello Beautiful=World:Again
}
