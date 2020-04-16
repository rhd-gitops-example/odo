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

package resmap

import (
<<<<<<< HEAD
	"fmt"

	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/pkg/ifc"
	internal "sigs.k8s.io/kustomize/pkg/internal/error"
	"sigs.k8s.io/kustomize/pkg/resource"
	"sigs.k8s.io/kustomize/pkg/types"
=======
	"github.com/pkg/errors"
	"sigs.k8s.io/kustomize/v3/internal/kusterr"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/resource"
	"sigs.k8s.io/kustomize/v3/pkg/types"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

// Factory makes instances of ResMap.
type Factory struct {
	resF *resource.Factory
<<<<<<< HEAD
}

// NewFactory returns a new resmap.Factory.
func NewFactory(rf *resource.Factory) *Factory {
	return &Factory{resF: rf}
=======
	tf   PatchFactory
}

// NewFactory returns a new resmap.Factory.
func NewFactory(rf *resource.Factory, tf PatchFactory) *Factory {
	return &Factory{resF: rf, tf: tf}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// RF returns a resource.Factory.
func (rmF *Factory) RF() *resource.Factory {
	return rmF.resF
}

<<<<<<< HEAD
// FromFiles returns a ResMap given a resource path slice.
func (rmF *Factory) FromFiles(
	loader ifc.Loader, paths []string) (ResMap, error) {
	var result []ResMap
	for _, path := range paths {
		content, err := loader.Load(path)
		if err != nil {
			return nil, errors.Wrap(err, "Load from path "+path+" failed")
		}
		res, err := rmF.NewResMapFromBytes(content)
		if err != nil {
			return nil, internal.Handler(err, path)
		}
		result = append(result, res)
	}
	return MergeWithErrorOnIdCollision(result...)
}

// newResMapFromBytes decodes a list of objects in byte array format.
func (rmF *Factory) NewResMapFromBytes(b []byte) (ResMap, error) {
	resources, err := rmF.resF.SliceFromBytes(b)
	if err != nil {
		return nil, err
	}

	result := ResMap{}
	for _, res := range resources {
		id := res.Id()
		if _, found := result[id]; found {
			return result, fmt.Errorf("GroupVersionKindName: %#v already exists b the map", id)
		}
		result[id] = res
	}
	return result, nil
=======
func New() ResMap {
	return newOne()
}

// FromResource returns a ResMap with one entry.
func (rmF *Factory) FromResource(res *resource.Resource) ResMap {
	m, err := newResMapFromResourceSlice([]*resource.Resource{res})
	if err != nil {
		panic(err)
	}
	return m
}

// FromFile returns a ResMap given a resource path.
func (rmF *Factory) FromFile(
	loader ifc.Loader, path string) (ResMap, error) {
	content, err := loader.Load(path)
	if err != nil {
		return nil, err
	}
	m, err := rmF.NewResMapFromBytes(content)
	if err != nil {
		return nil, kusterr.Handler(err, path)
	}
	return m, nil
}

// NewResMapFromBytes decodes a list of objects in byte array format.
func (rmF *Factory) NewResMapFromBytes(b []byte) (ResMap, error) {
	resources, err := rmF.resF.SliceFromBytes(b)
	if err != nil {
		return nil, err
	}
	return newResMapFromResourceSlice(resources)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// NewResMapFromConfigMapArgs returns a Resource slice given
// a configmap metadata slice from kustomization file.
<<<<<<< HEAD
func (rmF *Factory) NewResMapFromConfigMapArgs(argList []types.ConfigMapArgs, options *types.GeneratorOptions) (ResMap, error) {
	var resources []*resource.Resource
	for _, args := range argList {
		res, err := rmF.resF.MakeConfigMap(&args, options)
=======
func (rmF *Factory) NewResMapFromConfigMapArgs(
	ldr ifc.Loader,
	options *types.GeneratorOptions,
	argList []types.ConfigMapArgs) (ResMap, error) {
	var resources []*resource.Resource
	for _, args := range argList {
		res, err := rmF.resF.MakeConfigMap(ldr, options, &args)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		if err != nil {
			return nil, errors.Wrap(err, "NewResMapFromConfigMapArgs")
		}
		resources = append(resources, res)
	}
	return newResMapFromResourceSlice(resources)
}

<<<<<<< HEAD
// NewResMapFromSecretArgs takes a SecretArgs slice, generates
// secrets from each entry, and accumulates them in a ResMap.
func (rmF *Factory) NewResMapFromSecretArgs(argsList []types.SecretArgs, options *types.GeneratorOptions) (ResMap, error) {
	var resources []*resource.Resource
	for _, args := range argsList {
		res, err := rmF.resF.MakeSecret(&args, options)
=======
func (rmF *Factory) FromConfigMapArgs(
	ldr ifc.Loader,
	options *types.GeneratorOptions,
	args types.ConfigMapArgs) (ResMap, error) {
	res, err := rmF.resF.MakeConfigMap(ldr, options, &args)
	if err != nil {
		return nil, err
	}
	return rmF.FromResource(res), nil
}

// NewResMapFromSecretArgs takes a SecretArgs slice, generates
// secrets from each entry, and accumulates them in a ResMap.
func (rmF *Factory) NewResMapFromSecretArgs(
	ldr ifc.Loader,
	options *types.GeneratorOptions,
	argsList []types.SecretArgs) (ResMap, error) {
	var resources []*resource.Resource
	for _, args := range argsList {
		res, err := rmF.resF.MakeSecret(ldr, options, &args)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		if err != nil {
			return nil, errors.Wrap(err, "NewResMapFromSecretArgs")
		}
		resources = append(resources, res)
	}
	return newResMapFromResourceSlice(resources)
}

<<<<<<< HEAD
// Set sets the loader for the underlying factory
func (rmF *Factory) Set(ldr ifc.Loader) {
	rmF.resF.Set(ldr)
}

func newResMapFromResourceSlice(resources []*resource.Resource) (ResMap, error) {
	result := ResMap{}
	for _, res := range resources {
		id := res.Id()
		if _, found := result[id]; found {
			return nil, fmt.Errorf("duplicated %#v is not allowed", id)
		}
		result[id] = res
=======
func (rmF *Factory) FromSecretArgs(
	ldr ifc.Loader,
	options *types.GeneratorOptions,
	args types.SecretArgs) (ResMap, error) {
	res, err := rmF.resF.MakeSecret(ldr, options, &args)
	if err != nil {
		return nil, err
	}
	return rmF.FromResource(res), nil
}

func (rmF *Factory) MergePatches(patches []*resource.Resource) (
	ResMap, error) {
	return rmF.tf.MergePatches(patches, rmF.resF)
}

func newResMapFromResourceSlice(resources []*resource.Resource) (ResMap, error) {
	result := New()
	for _, res := range resources {
		err := result.Append(res)
		if err != nil {
			return nil, err
		}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
	return result, nil
}
