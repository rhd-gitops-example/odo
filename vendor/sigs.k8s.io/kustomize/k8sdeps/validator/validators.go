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

// Package validator provides functions to validate labels, annotations, namespace using apimachinery
=======
// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Package validator provides functions to validate labels, annotations,
// namespaces and configmap/secret keys using apimachinery functions.
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
package validator

import (
	"errors"
<<<<<<< HEAD
=======
	"fmt"
	"strings"

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	v1validation "k8s.io/apimachinery/pkg/apis/meta/v1/validation"
	"k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// KustValidator validates Labels and annotations by apimachinery
type KustValidator struct{}

// NewKustValidator returns a KustValidator object
func NewKustValidator() *KustValidator {
	return &KustValidator{}
}

<<<<<<< HEAD
=======
func (v *KustValidator) ErrIfInvalidKey(k string) error {
	if errs := validation.IsConfigMapKey(k); len(errs) != 0 {
		return fmt.Errorf(
			"%q is not a valid key name: %s",
			k, strings.Join(errs, ";"))
	}
	return nil
}

func (v *KustValidator) IsEnvVarName(k string) error {
	if errs := validation.IsEnvVarName(k); len(errs) != 0 {
		return fmt.Errorf(
			"%q is not a valid key name: %s",
			k, strings.Join(errs, ";"))
	}
	return nil
}

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
// MakeAnnotationValidator returns a MapValidatorFunc using apimachinery.
func (v *KustValidator) MakeAnnotationValidator() func(map[string]string) error {
	return func(x map[string]string) error {
		errs := apivalidation.ValidateAnnotations(x, field.NewPath("field"))
		if len(errs) > 0 {
			return errors.New(errs.ToAggregate().Error())
		}
		return nil
	}
}

<<<<<<< HEAD
=======
// MakeAnnotationNameValidator returns a MapValidatorFunc using apimachinery.
func (v *KustValidator) MakeAnnotationNameValidator() func([]string) error {
	return func(x []string) error {
		errs := field.ErrorList{}
		fldPath := field.NewPath("field")
		for _, k := range x {
			for _, msg := range validation.IsQualifiedName(strings.ToLower(k)) {
				errs = append(errs, field.Invalid(fldPath, k, msg))
			}
		}
		if len(errs) > 0 {
			return errors.New(errs.ToAggregate().Error())
		}
		return nil
	}
}

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
// MakeLabelValidator returns a MapValidatorFunc using apimachinery.
func (v *KustValidator) MakeLabelValidator() func(map[string]string) error {
	return func(x map[string]string) error {
		errs := v1validation.ValidateLabels(x, field.NewPath("field"))
		if len(errs) > 0 {
			return errors.New(errs.ToAggregate().Error())
		}
		return nil
	}
}

<<<<<<< HEAD
=======
// MakeLabelNameValidator returns a ArrayValidatorFunc using apimachinery.
func (v *KustValidator) MakeLabelNameValidator() func([]string) error {
	return func(x []string) error {
		errs := field.ErrorList{}
		fldPath := field.NewPath("field")
		for _, k := range x {
			errs = append(errs, v1validation.ValidateLabelName(k, fldPath)...)
		}
		if len(errs) > 0 {
			return errors.New(errs.ToAggregate().Error())
		}
		return nil
	}
}

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
// ValidateNamespace validates a string is a valid namespace using apimachinery.
func (v *KustValidator) ValidateNamespace(s string) []string {
	return validation.IsDNS1123Label(s)
}
