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

// Package ifc holds miscellaneous interfaces used by kustomize.
package ifc

import (
<<<<<<< HEAD
	"sigs.k8s.io/kustomize/pkg/gvk"
	"sigs.k8s.io/kustomize/pkg/types"
=======
	"sigs.k8s.io/kustomize/v3/pkg/gvk"
	"sigs.k8s.io/kustomize/v3/pkg/types"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

// Validator provides functions to validate annotations and labels
type Validator interface {
	MakeAnnotationValidator() func(map[string]string) error
<<<<<<< HEAD
	MakeLabelValidator() func(map[string]string) error
	ValidateNamespace(string) []string
=======
	MakeAnnotationNameValidator() func([]string) error
	MakeLabelValidator() func(map[string]string) error
	MakeLabelNameValidator() func([]string) error
	ValidateNamespace(string) []string
	ErrIfInvalidKey(string) error
	IsEnvVarName(k string) error
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// Loader interface exposes methods to read bytes.
type Loader interface {
	// Root returns the root location for this Loader.
	Root() string
	// New returns Loader located at newRoot.
	New(newRoot string) (Loader, error)
	// Load returns the bytes read from the location or an error.
	Load(location string) ([]byte, error)
	// Cleanup cleans the loader
	Cleanup() error
<<<<<<< HEAD
=======
	// Validator validates data for use in various k8s fields.
	Validator() Validator
	// Loads pairs.
	LoadKvPairs(args types.GeneratorArgs) ([]types.Pair, error)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// Kunstructured allows manipulation of k8s objects
// that do not have Golang structs.
type Kunstructured interface {
	Map() map[string]interface{}
	SetMap(map[string]interface{})
	Copy() Kunstructured
<<<<<<< HEAD
	GetFieldValue(string) (string, error)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	GetGvk() gvk.Gvk
	GetKind() string
	GetName() string
	SetName(string)
=======
	GetFieldValue(string) (interface{}, error)
	GetString(string) (string, error)
	GetStringSlice(string) ([]string, error)
	GetBool(path string) (bool, error)
	GetFloat64(path string) (float64, error)
	GetInt64(path string) (int64, error)
	GetSlice(path string) ([]interface{}, error)
	GetStringMap(path string) (map[string]string, error)
	GetMap(path string) (map[string]interface{}, error)
	MarshalJSON() ([]byte, error)
	UnmarshalJSON([]byte) error
	GetGvk() gvk.Gvk
	SetGvk(gvk.Gvk)
	GetKind() string
	GetName() string
	SetName(string)
	SetNamespace(string)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	GetLabels() map[string]string
	SetLabels(map[string]string)
	GetAnnotations() map[string]string
	SetAnnotations(map[string]string)
<<<<<<< HEAD
=======
	MatchesLabelSelector(selector string) (bool, error)
	MatchesAnnotationSelector(selector string) (bool, error)
	Patch(Kunstructured) error
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// KunstructuredFactory makes instances of Kunstructured.
type KunstructuredFactory interface {
	SliceFromBytes([]byte) ([]Kunstructured, error)
	FromMap(m map[string]interface{}) Kunstructured
<<<<<<< HEAD
	MakeConfigMap(args *types.ConfigMapArgs, options *types.GeneratorOptions) (Kunstructured, error)
	MakeSecret(args *types.SecretArgs, options *types.GeneratorOptions) (Kunstructured, error)
	Set(ldr Loader)
=======
	Hasher() KunstructuredHasher
	MakeConfigMap(
		ldr Loader,
		options *types.GeneratorOptions,
		args *types.ConfigMapArgs) (Kunstructured, error)
	MakeSecret(
		ldr Loader,
		options *types.GeneratorOptions,
		args *types.SecretArgs) (Kunstructured, error)
}

// KunstructuredHasher returns a hash of the argument
// or an error.
type KunstructuredHasher interface {
	Hash(Kunstructured) (string, error)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// See core.v1.SecretTypeOpaque
const SecretTypeOpaque = "Opaque"
