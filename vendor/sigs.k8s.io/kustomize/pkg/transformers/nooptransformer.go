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

<<<<<<< HEAD
import "sigs.k8s.io/kustomize/pkg/resmap"
=======
import "sigs.k8s.io/kustomize/v3/pkg/resmap"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)

// noOpTransformer contains a no-op transformer.
type noOpTransformer struct{}

<<<<<<< HEAD
var _ Transformer = &noOpTransformer{}

// NewNoOpTransformer constructs a noOpTransformer.
func NewNoOpTransformer() Transformer {
=======
var _ resmap.Transformer = &noOpTransformer{}

// NewNoOpTransformer constructs a noOpTransformer.
func NewNoOpTransformer() resmap.Transformer {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	return &noOpTransformer{}
}

// Transform does nothing.
func (o *noOpTransformer) Transform(_ resmap.ResMap) error {
	return nil
}
