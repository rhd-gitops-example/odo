/*
Copyright 2018 The Kubernetes Authors.
<<<<<<< HEAD
 Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
     http://www.apache.org/licenses/LICENSE-2.0
 Unless required by applicable law or agreed to in writing, software
=======

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types_test

import (
	"testing"

<<<<<<< HEAD
	. "sigs.k8s.io/kustomize/pkg/types"
=======
	. "sigs.k8s.io/kustomize/v3/pkg/types"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

func TestGenArgs_String(t *testing.T) {
	tests := []struct {
		ga       *GenArgs
		expected string
	}{
		{
			ga:       nil,
			expected: "{nilGenArgs}",
		},
		{
			ga:       &GenArgs{},
			expected: "{nsfx:false,beh:unspecified}",
		},
		{
			ga: NewGenArgs(
				&GeneratorArgs{Behavior: "merge"},
				&GeneratorOptions{DisableNameSuffixHash: false}),
			expected: "{nsfx:true,beh:merge}",
		},
	}
	for _, test := range tests {
		if test.ga.String() != test.expected {
			t.Fatalf("Expected '%s', got '%s'", test.expected, test.ga.String())
		}
	}
}
