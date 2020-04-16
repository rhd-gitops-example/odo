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

// Package validators defines a FakeValidator that can be used in tests
package validators

import (
	"errors"
	"regexp"
	"testing"
)

// FakeValidator can be used in tests.
type FakeValidator struct {
	happy  bool
	called bool
	t      *testing.T
}

// SAD is an error string.
const SAD = "i'm not happy Bob, NOT HAPPY"

// MakeHappyMapValidator makes a FakeValidator that always passes.
func MakeHappyMapValidator(t *testing.T) *FakeValidator {
	return &FakeValidator{happy: true, t: t}
}

// MakeSadMapValidator makes a FakeValidator that always fails.
func MakeSadMapValidator(t *testing.T) *FakeValidator {
	return &FakeValidator{happy: false, t: t}
}

// MakeFakeValidator makes an empty Fake Validator.
func MakeFakeValidator() *FakeValidator {
	return &FakeValidator{}
}

<<<<<<< HEAD
=======
// ErrIfInvalidKey returns nil
func (v *FakeValidator) ErrIfInvalidKey(k string) error {
	return nil
}

// IsEnvVarName returns nil
func (v *FakeValidator) IsEnvVarName(k string) error {
	return nil
}

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
// MakeAnnotationValidator returns a nil function
func (v *FakeValidator) MakeAnnotationValidator() func(map[string]string) error {
	return nil
}

<<<<<<< HEAD
=======
// MakeAnnotationNameValidator returns a nil function
func (v *FakeValidator) MakeAnnotationNameValidator() func([]string) error {
	return nil
}

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
// MakeLabelValidator returns a nil function
func (v *FakeValidator) MakeLabelValidator() func(map[string]string) error {
	return nil
}

<<<<<<< HEAD
=======
// MakeLabelNameValidator returns a nil function
func (v *FakeValidator) MakeLabelNameValidator() func([]string) error {
	return nil
}

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
// ValidateNamespace validates namespace by regexp
func (v *FakeValidator) ValidateNamespace(s string) []string {
	pattern := regexp.MustCompile(`^[a-zA-Z].*`)
	if pattern.MatchString(s) {
		return nil
	}
	return []string{"doesn't match"}
}

// Validator replaces apimachinery validation in tests.
// Can be set to fail or succeed to test error handling.
// Can confirm if run or not run by surrounding code.
func (v *FakeValidator) Validator(_ map[string]string) error {
	v.called = true
	if v.happy {
		return nil
	}
	return errors.New(SAD)
}

<<<<<<< HEAD
=======
func (v *FakeValidator) ValidatorArray(_ []string) error {
	v.called = true
	if v.happy {
		return nil
	}
	return errors.New(SAD)
}

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
// VerifyCall returns true if Validator was used.
func (v *FakeValidator) VerifyCall() {
	if !v.called {
		v.t.Errorf("should have called Validator")
	}
}

// VerifyNoCall returns true if Validator was not used.
func (v *FakeValidator) VerifyNoCall() {
	if v.called {
		v.t.Errorf("should not have called Validator")
	}
}
