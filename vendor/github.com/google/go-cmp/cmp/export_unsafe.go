// Copyright 2017, The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE.md file.

// +build !purego

package cmp

import (
	"reflect"
	"unsafe"
)

<<<<<<< HEAD
const supportAllowUnexported = true
=======
const supportExporters = true
>>>>>>> Start of generating Secrets.

// retrieveUnexportedField uses unsafe to forcibly retrieve any field from
// a struct such that the value has read-write permissions.
//
// The parent struct, v, must be addressable, while f must be a StructField
// describing the field to retrieve.
func retrieveUnexportedField(v reflect.Value, f reflect.StructField) reflect.Value {
<<<<<<< HEAD
	return reflect.NewAt(f.Type, unsafe.Pointer(v.UnsafeAddr()+f.Offset)).Elem()
=======
	// See https://github.com/google/go-cmp/issues/167 for discussion of the
	// following expression.
	return reflect.NewAt(f.Type, unsafe.Pointer(uintptr(unsafe.Pointer(v.UnsafeAddr()))+f.Offset)).Elem()
>>>>>>> Start of generating Secrets.
}
