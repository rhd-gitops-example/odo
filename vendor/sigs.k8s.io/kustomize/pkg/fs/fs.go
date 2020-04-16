<<<<<<< HEAD
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
=======
// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)

// Package fs provides a file system abstraction layer.
package fs

import (
	"io"
	"os"
<<<<<<< HEAD
=======
	"path/filepath"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

// FileSystem groups basic os filesystem methods.
type FileSystem interface {
<<<<<<< HEAD
	Create(name string) (File, error)
	Mkdir(name string) error
	MkdirAll(name string) error
	RemoveAll(name string) error
	Open(name string) (File, error)
	IsDir(name string) bool
	CleanedAbs(path string) (ConfirmedDir, string, error)
	Exists(name string) bool
	Glob(pattern string) ([]string, error)
	ReadFile(name string) ([]byte, error)
	WriteFile(name string, data []byte) error
=======
	// Create a file.
	Create(name string) (File, error)
	// MkDir makes a directory.
	Mkdir(path string) error
	// MkDir makes a directory path, creating intervening directories.
	MkdirAll(path string) error
	// RemoveAll removes path and any children it contains.
	RemoveAll(path string) error
	// Open opens the named file for reading.
	Open(path string) (File, error)
	// IsDir returns true if the path is a directory.
	IsDir(path string) bool
	// CleanedAbs converts the given path into a
	// directory and a file name, where the directory
	// is represented as a ConfirmedDir and all that implies.
	// If the entire path is a directory, the file component
	// is an empty string.
	CleanedAbs(path string) (ConfirmedDir, string, error)
	// Exists is true if the path exists in the file system.
	Exists(path string) bool
	// Glob returns the list of matching files
	Glob(pattern string) ([]string, error)
	// ReadFile returns the contents of the file at the given path.
	ReadFile(path string) ([]byte, error)
	// WriteFile writes the data to a file at the given path.
	WriteFile(path string, data []byte) error
	// Walk walks the file system with the given WalkFunc.
	Walk(path string, walkFn filepath.WalkFunc) error
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// File groups the basic os.File methods.
type File interface {
	io.ReadWriteCloser
	Stat() (os.FileInfo, error)
}
