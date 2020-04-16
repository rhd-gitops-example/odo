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

package resource

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

<<<<<<< HEAD
	"sigs.k8s.io/kustomize/pkg/ifc"
	internal "sigs.k8s.io/kustomize/pkg/internal/error"
	"sigs.k8s.io/kustomize/pkg/patch"
	"sigs.k8s.io/kustomize/pkg/types"
=======
	"sigs.k8s.io/kustomize/v3/internal/kusterr"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/types"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

// Factory makes instances of Resource.
type Factory struct {
	kf ifc.KunstructuredFactory
}

// NewFactory makes an instance of Factory.
func NewFactory(kf ifc.KunstructuredFactory) *Factory {
	return &Factory{kf: kf}
}

<<<<<<< HEAD
// FromMap returns a new instance of Resource.
func (rf *Factory) FromMap(m map[string]interface{}) *Resource {
	return &Resource{
		Kunstructured: rf.kf.FromMap(m),
		options:       types.NewGenArgs(nil, nil),
	}
}

// FromMapAndOption returns a new instance of Resource with given options.
func (rf *Factory) FromMapAndOption(m map[string]interface{}, args *types.GeneratorArgs, option *types.GeneratorOptions) *Resource {
	return &Resource{
		Kunstructured: rf.kf.FromMap(m),
		options:       types.NewGenArgs(args, option),
	}
}

// FromKunstructured returns a new instance of Resource.
func (rf *Factory) FromKunstructured(
	u ifc.Kunstructured) *Resource {
	if u == nil {
		log.Fatal("unstruct ifc must not be null")
	}
	return &Resource{
		Kunstructured: u,
		options:       types.NewGenArgs(nil, nil),
	}
=======
func (rf *Factory) Hasher() ifc.KunstructuredHasher {
	return rf.kf.Hasher()
}

// FromMap returns a new instance of Resource.
func (rf *Factory) FromMap(m map[string]interface{}) *Resource {
	return rf.makeOne(rf.kf.FromMap(m), nil)
}

// FromMapWithName returns a new instance with the given "original" name.
func (rf *Factory) FromMapWithName(n string, m map[string]interface{}) *Resource {
	return rf.makeOne(rf.kf.FromMap(m), nil).setOriginalName(n)
}

// FromMapWithNamespace returns a new instance with the given "original" namespace.
func (rf *Factory) FromMapWithNamespace(n string, m map[string]interface{}) *Resource {
	return rf.makeOne(rf.kf.FromMap(m), nil).setOriginalNs(n)
}

// FromMapWithNamespaceAndName returns a new instance with the given "original" namespace.
func (rf *Factory) FromMapWithNamespaceAndName(ns string, n string, m map[string]interface{}) *Resource {
	return rf.makeOne(rf.kf.FromMap(m), nil).setOriginalNs(ns).setOriginalName(n)
}

// FromMapAndOption returns a new instance of Resource with given options.
func (rf *Factory) FromMapAndOption(
	m map[string]interface{}, args *types.GeneratorArgs, option *types.GeneratorOptions) *Resource {
	return rf.makeOne(rf.kf.FromMap(m), types.NewGenArgs(args, option))
}

// FromKunstructured returns a new instance of Resource.
func (rf *Factory) FromKunstructured(u ifc.Kunstructured) *Resource {
	return rf.makeOne(u, nil)
}

// makeOne returns a new instance of Resource.
func (rf *Factory) makeOne(
	u ifc.Kunstructured, o *types.GenArgs) *Resource {
	if u == nil {
		log.Fatal("unstruct ifc must not be null")
	}
	if o == nil {
		o = types.NewGenArgs(nil, nil)
	}
	r := &Resource{
		Kunstructured: u,
		options:       o,
	}
	return r.setOriginalName(r.GetName()).setOriginalNs(r.GetNamespace())
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}

// SliceFromPatches returns a slice of resources given a patch path
// slice from a kustomization file.
func (rf *Factory) SliceFromPatches(
<<<<<<< HEAD
	ldr ifc.Loader, paths []patch.StrategicMerge) ([]*Resource, error) {
=======
	ldr ifc.Loader, paths []types.PatchStrategicMerge) ([]*Resource, error) {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	var result []*Resource
	for _, path := range paths {
		content, err := ldr.Load(string(path))
		if err != nil {
			return nil, err
		}
		res, err := rf.SliceFromBytes(content)
		if err != nil {
<<<<<<< HEAD
			return nil, internal.Handler(err, string(path))
=======
			return nil, kusterr.Handler(err, string(path))
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		}
		result = append(result, res...)
	}
	return result, nil
}

<<<<<<< HEAD
// SliceFromBytes unmarshalls bytes into a Resource slice.
=======
// FromBytes unmarshals bytes into one Resource.
func (rf *Factory) FromBytes(in []byte) (*Resource, error) {
	result, err := rf.SliceFromBytes(in)
	if err != nil {
		return nil, err
	}
	if len(result) != 1 {
		return nil, fmt.Errorf(
			"expected 1 resource, found %d in %v", len(result), in)
	}
	return result[0], nil
}

// SliceFromBytes unmarshals bytes into a Resource slice.
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
func (rf *Factory) SliceFromBytes(in []byte) ([]*Resource, error) {
	kunStructs, err := rf.kf.SliceFromBytes(in)
	if err != nil {
		return nil, err
	}
	var result []*Resource
	for len(kunStructs) > 0 {
		u := kunStructs[0]
		kunStructs = kunStructs[1:]
		if strings.HasSuffix(u.GetKind(), "List") {
			items := u.Map()["items"]
			itemsSlice, ok := items.([]interface{})
			if !ok {
				if items == nil {
					// an empty list
					continue
				}
				return nil, fmt.Errorf("items in List is type %T, expected array", items)
			}
			for _, item := range itemsSlice {
				itemJSON, err := json.Marshal(item)
				if err != nil {
					return nil, err
				}
				innerU, err := rf.kf.SliceFromBytes(itemJSON)
				if err != nil {
					return nil, err
				}
				// append innerU to kunStructs so nested Lists can be handled
				kunStructs = append(kunStructs, innerU...)
			}
		} else {
			result = append(result, rf.FromKunstructured(u))
		}
	}
	return result, nil
}

<<<<<<< HEAD
// Set sets the loader for the underlying factory
func (rf *Factory) Set(ldr ifc.Loader) {
	rf.kf.Set(ldr)
}

// MakeConfigMap makes an instance of Resource for ConfigMap
func (rf *Factory) MakeConfigMap(args *types.ConfigMapArgs, options *types.GeneratorOptions) (*Resource, error) {
	u, err := rf.kf.MakeConfigMap(args, options)
	if err != nil {
		return nil, err
	}
	return &Resource{Kunstructured: u, options: types.NewGenArgs(&types.GeneratorArgs{Behavior: args.Behavior}, options)}, nil
}

// MakeSecret makes an instance of Resource for Secret
func (rf *Factory) MakeSecret(args *types.SecretArgs, options *types.GeneratorOptions) (*Resource, error) {
	u, err := rf.kf.MakeSecret(args, options)
	if err != nil {
		return nil, err
	}
	return &Resource{Kunstructured: u, options: types.NewGenArgs(&types.GeneratorArgs{Behavior: args.Behavior}, options)}, nil
=======
// MakeConfigMap makes an instance of Resource for ConfigMap
func (rf *Factory) MakeConfigMap(
	ldr ifc.Loader,
	options *types.GeneratorOptions,
	args *types.ConfigMapArgs) (*Resource, error) {
	u, err := rf.kf.MakeConfigMap(ldr, options, args)
	if err != nil {
		return nil, err
	}
	return rf.makeOne(
		u,
		types.NewGenArgs(
			&types.GeneratorArgs{Behavior: args.Behavior},
			options)), nil
}

// MakeSecret makes an instance of Resource for Secret
func (rf *Factory) MakeSecret(
	ldr ifc.Loader,
	options *types.GeneratorOptions,
	args *types.SecretArgs) (*Resource, error) {
	u, err := rf.kf.MakeSecret(ldr, options, args)
	if err != nil {
		return nil, err
	}
	return rf.makeOne(
		u,
		types.NewGenArgs(
			&types.GeneratorArgs{Behavior: args.Behavior},
			options)), nil
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
}
