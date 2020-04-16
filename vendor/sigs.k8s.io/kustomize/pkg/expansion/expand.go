/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package expansion provides functions find and replace $(FOO) style variables in strings.
package expansion

import (
	"bytes"
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

const (
	operator        = '$'
	referenceOpener = '('
	referenceCloser = ')'
)

// syntaxWrap returns the input string wrapped by the expansion syntax.
func syntaxWrap(input string) string {
	return string(operator) + string(referenceOpener) + input + string(referenceCloser)
}

// MappingFuncFor returns a mapping function for use with Expand that
// implements the expansion semantics defined in the expansion spec; it
// returns the input string wrapped in the expansion syntax if no mapping
// for the input is found.
func MappingFuncFor(
	counts map[string]int,
<<<<<<< HEAD
	context ...map[string]string) func(string) string {
	return func(input string) string {
=======
	context ...map[string]interface{}) func(string) interface{} {
	return func(input string) interface{} {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		for _, vars := range context {
			val, ok := vars[input]
			if ok {
				counts[input]++
<<<<<<< HEAD
				return val
=======
				switch typedV := val.(type) {
				case string, int64, float64, bool:
					return typedV
				default:
					return syntaxWrap(input)
				}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			}
		}
		return syntaxWrap(input)
	}
}

// Expand replaces variable references in the input string according to
// the expansion spec using the given mapping function to resolve the
// values of variables.
<<<<<<< HEAD
func Expand(input string, mapping func(string) string) string {
=======
func Expand(input string, mapping func(string) interface{}) interface{} {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	var buf bytes.Buffer
	checkpoint := 0
	for cursor := 0; cursor < len(input); cursor++ {
		if input[cursor] == operator && cursor+1 < len(input) {
			// Copy the portion of the input string since the last
			// checkpoint into the buffer
			buf.WriteString(input[checkpoint:cursor])

			// Attempt to read the variable name as defined by the
			// syntax from the input string
			read, isVar, advance := tryReadVariableName(input[cursor+1:])

			if isVar {
				// We were able to read a variable name correctly;
				// apply the mapping to the variable name and copy the
				// bytes into the buffer
<<<<<<< HEAD
				buf.WriteString(mapping(read))
=======
				mapped := mapping(read)
				if input == syntaxWrap(read) {
					// Preserve the type of variable
					return mapped
				}

				// Variable is used in a middle of a string
				buf.WriteString(fmt.Sprintf("%v", mapped))
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			} else {
				// Not a variable name; copy the read bytes into the buffer
				buf.WriteString(read)
			}

			// Advance the cursor in the input string to account for
			// bytes consumed to read the variable name expression
			cursor += advance

			// Advance the checkpoint in the input string
			checkpoint = cursor + 1
		}
	}

	// Return the buffer and any remaining unwritten bytes in the
	// input string.
	return buf.String() + input[checkpoint:]
}

// tryReadVariableName attempts to read a variable name from the input
// string and returns the content read from the input, whether that content
// represents a variable name to perform mapping on, and the number of bytes
// consumed in the input string.
//
// The input string is assumed not to contain the initial operator.
func tryReadVariableName(input string) (string, bool, int) {
	switch input[0] {
	case operator:
		// Escaped operator; return it.
		return input[0:1], false, 1
	case referenceOpener:
		// Scan to expression closer
		for i := 1; i < len(input); i++ {
			if input[i] == referenceCloser {
				return input[1:i], true, i + 1
			}
		}

		// Incomplete reference; return it.
		return string(operator) + string(referenceOpener), false, 1
	default:
		// Not the beginning of an expression, ie, an operator
		// that doesn't begin an expression.  Return the operator
		// and the first rune in the string.
		return string(operator) + string(input[0]), false, 1
	}
}
