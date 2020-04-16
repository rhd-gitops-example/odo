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

package configmapandsecret

import (
	"fmt"
<<<<<<< HEAD
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/validation"
	"sigs.k8s.io/kustomize/k8sdeps/kv"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/types"
)

// SecretFactory makes Secrets.
type SecretFactory struct {
	ldr ifc.Loader
}

// NewSecretFactory returns a new SecretFactory.
func NewSecretFactory(ldr ifc.Loader) *SecretFactory {
	return &SecretFactory{ldr: ldr}
}

func (f *SecretFactory) makeFreshSecret(args *types.SecretArgs) *corev1.Secret {
=======

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/kustomize/v3/pkg/types"
)

func makeFreshSecret(
	args *types.SecretArgs) *corev1.Secret {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	s := &corev1.Secret{}
	s.APIVersion = "v1"
	s.Kind = "Secret"
	s.Name = args.Name
	s.Namespace = args.Namespace
	s.Type = corev1.SecretType(args.Type)
	if s.Type == "" {
		s.Type = corev1.SecretTypeOpaque
	}
	s.Data = map[string][]byte{}
	return s
}

// MakeSecret returns a new secret.
<<<<<<< HEAD
func (f *SecretFactory) MakeSecret(args *types.SecretArgs, options *types.GeneratorOptions) (*corev1.Secret, error) {
	var all []kv.Pair
	var err error
	s := f.makeFreshSecret(args)

	pairs, err := keyValuesFromEnvFile(f.ldr, args.EnvSource)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(
			"env source file: %s",
			args.EnvSource))
	}
	all = append(all, pairs...)

	pairs, err = keyValuesFromLiteralSources(args.LiteralSources)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(
			"literal sources %v", args.LiteralSources))
	}
	all = append(all, pairs...)

	pairs, err = keyValuesFromFileSources(f.ldr, args.FileSources)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf(
			"file sources: %v", args.FileSources))
	}
	all = append(all, pairs...)

	for _, p := range all {
		err = addKvToSecret(s, p.Key, p.Value)
=======
func (f *Factory) MakeSecret(
	args *types.SecretArgs) (*corev1.Secret, error) {
	all, err := f.ldr.LoadKvPairs(args.GeneratorArgs)
	if err != nil {
		return nil, err
	}
	s := makeFreshSecret(args)
	for _, p := range all {
		err = f.addKvToSecret(s, p.Key, p.Value)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		if err != nil {
			return nil, err
		}
	}
<<<<<<< HEAD
	if options != nil {
		s.SetLabels(options.Labels)
		s.SetAnnotations(options.Annotations)
=======
	if f.options != nil {
		s.SetLabels(f.options.Labels)
		s.SetAnnotations(f.options.Annotations)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
	return s, nil
}

<<<<<<< HEAD
func addKvToSecret(secret *corev1.Secret, keyName, data string) error {
	// Note, the rules for SecretKeys  keys are the exact same as the ones for ConfigMap.
	if errs := validation.IsConfigMapKey(keyName); len(errs) != 0 {
		return fmt.Errorf("%q is not a valid key name for a Secret: %s", keyName, strings.Join(errs, ";"))
	}
	if _, entryExists := secret.Data[keyName]; entryExists {
		return fmt.Errorf("cannot add key %s, another key by that name already exists", keyName)
=======
func (f *Factory) addKvToSecret(secret *corev1.Secret, keyName, data string) error {
	if err := f.ldr.Validator().ErrIfInvalidKey(keyName); err != nil {
		return err
	}
	if _, entryExists := secret.Data[keyName]; entryExists {
		return fmt.Errorf(keyExistsErrorMsg, keyName, secret.Data)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
	secret.Data[keyName] = []byte(data)
	return nil
}
