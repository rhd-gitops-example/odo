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

package transformers

import (
	"errors"
	"fmt"

<<<<<<< HEAD
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/transformers/config"
=======
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/transformers/config"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

// mapTransformer applies a string->string map to fieldSpecs.
type mapTransformer struct {
	m          map[string]string
	fieldSpecs []config.FieldSpec
}

<<<<<<< HEAD
var _ Transformer = &mapTransformer{}

// NewLabelsMapTransformer constructs a mapTransformer.
func NewLabelsMapTransformer(
	m map[string]string, fs []config.FieldSpec) (Transformer, error) {
=======
var _ resmap.Transformer = &mapTransformer{}

// NewLabelsMapTransformer constructs a mapTransformer.
func NewLabelsMapTransformer(
	m map[string]string, fs []config.FieldSpec) (resmap.Transformer, error) {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	return NewMapTransformer(fs, m)
}

// NewAnnotationsMapTransformer construct a mapTransformer.
func NewAnnotationsMapTransformer(
<<<<<<< HEAD
	m map[string]string, fs []config.FieldSpec) (Transformer, error) {
=======
	m map[string]string, fs []config.FieldSpec) (resmap.Transformer, error) {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	return NewMapTransformer(fs, m)
}

// NewMapTransformer construct a mapTransformer.
func NewMapTransformer(
<<<<<<< HEAD
	pc []config.FieldSpec, m map[string]string) (Transformer, error) {
=======
	pc []config.FieldSpec, m map[string]string) (resmap.Transformer, error) {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	if m == nil {
		return NewNoOpTransformer(), nil
	}
	if pc == nil {
		return nil, errors.New("fieldSpecs is not expected to be nil")
	}
	return &mapTransformer{fieldSpecs: pc, m: m}, nil
}

// Transform apply each <key, value> pair in the mapTransformer to the
// fields specified in mapTransformer.
func (o *mapTransformer) Transform(m resmap.ResMap) error {
<<<<<<< HEAD
	for id := range m {
		objMap := m[id].Map()
		for _, path := range o.fieldSpecs {
			if !id.Gvk().IsSelected(&path.Gvk) {
				continue
			}
			err := mutateField(objMap, path.PathSlice(), path.CreateIfNotPresent, o.addMap)
=======
	for _, r := range m.Resources() {
		for _, path := range o.fieldSpecs {
			if !r.OrgId().IsSelected(&path.Gvk) {
				continue
			}
			err := MutateField(
				r.Map(), path.PathSlice(),
				path.CreateIfNotPresent, o.addMap)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o *mapTransformer) addMap(in interface{}) (interface{}, error) {
	m, ok := in.(map[string]interface{})
<<<<<<< HEAD
	if !ok {
=======
	if in == nil {
		m = map[string]interface{}{}
	} else if !ok {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		return nil, fmt.Errorf("%#v is expected to be %T", in, m)
	}
	for k, v := range o.m {
		m[k] = v
	}
	return m, nil
}
