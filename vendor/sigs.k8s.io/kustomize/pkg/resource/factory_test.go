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

package resource_test

import (
	"reflect"
	"testing"

<<<<<<< HEAD
	"sigs.k8s.io/kustomize/pkg/internal/loadertest"
	"sigs.k8s.io/kustomize/pkg/patch"
	. "sigs.k8s.io/kustomize/pkg/resource"
=======
	"sigs.k8s.io/kustomize/v3/internal/loadertest"
	. "sigs.k8s.io/kustomize/v3/pkg/resource"
	"sigs.k8s.io/kustomize/v3/pkg/types"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

func TestSliceFromPatches(t *testing.T) {

<<<<<<< HEAD
	patchGood1 := patch.StrategicMerge("patch1.yaml")
=======
	patchGood1 := types.PatchStrategicMerge("patch1.yaml")
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	patch1 := `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pooh
`
<<<<<<< HEAD
	patchGood2 := patch.StrategicMerge("patch2.yaml")
=======
	patchGood2 := types.PatchStrategicMerge("patch2.yaml")
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	patch2 := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: winnie
  namespace: hundred-acre-wood
---
# some comment
---
---
`
<<<<<<< HEAD
	patchBad := patch.StrategicMerge("patch3.yaml")
	patch3 := `
WOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOT: woot
`
	patchList := patch.StrategicMerge("patch4.yaml")
=======
	patchBad := types.PatchStrategicMerge("patch3.yaml")
	patch3 := `
WOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOOT: woot
`
	patchList := types.PatchStrategicMerge("patch4.yaml")
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	patch4 := `
apiVersion: v1
kind: List
items:
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: pooh
- apiVersion: v1
  kind: ConfigMap
  metadata:
    name: winnie
    namespace: hundred-acre-wood
`
<<<<<<< HEAD
	patchList2 := patch.StrategicMerge("patch5.yaml")
=======
	patchList2 := types.PatchStrategicMerge("patch5.yaml")
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	patch5 := `
apiVersion: v1
kind: DeploymentList
items:
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: deployment-a
  spec: &hostAliases
    template:
      spec:
        hostAliases:
        - hostnames:
          - a.example.com
          ip: 8.8.8.8
- apiVersion: apps/v1
  kind: Deployment
  metadata:
    name: deployment-b
  spec:
    <<: *hostAliases
`
<<<<<<< HEAD
	patchList3 := patch.StrategicMerge("patch6.yaml")
=======
	patchList3 := types.PatchStrategicMerge("patch6.yaml")
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	patch6 := `
apiVersion: v1
kind: List
items:
`
<<<<<<< HEAD
	patchList4 := patch.StrategicMerge("patch7.yaml")
=======
	patchList4 := types.PatchStrategicMerge("patch7.yaml")
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	patch7 := `
apiVersion: v1
kind: List
`
	testDeploymentSpec := map[string]interface{}{
		"template": map[string]interface{}{
			"spec": map[string]interface{}{
				"hostAliases": []interface{}{
					map[string]interface{}{
						"hostnames": []interface{}{
							"a.example.com",
						},
						"ip": "8.8.8.8",
					},
				},
			},
		},
	}
	testDeploymentA := factory.FromMap(
		map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deployment-a",
			},
			"spec": testDeploymentSpec,
		})
	testDeploymentB := factory.FromMap(
		map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deployment-b",
			},
			"spec": testDeploymentSpec,
		})
	l := loadertest.NewFakeLoader("/")
	l.AddFile("/"+string(patchGood1), []byte(patch1))
	l.AddFile("/"+string(patchGood2), []byte(patch2))
	l.AddFile("/"+string(patchBad), []byte(patch3))
	l.AddFile("/"+string(patchList), []byte(patch4))
	l.AddFile("/"+string(patchList2), []byte(patch5))
	l.AddFile("/"+string(patchList3), []byte(patch6))
	l.AddFile("/"+string(patchList4), []byte(patch7))

	tests := []struct {
		name        string
<<<<<<< HEAD
		input       []patch.StrategicMerge
=======
		input       []types.PatchStrategicMerge
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		expectedOut []*Resource
		expectedErr bool
	}{
		{
			name:        "happy",
<<<<<<< HEAD
			input:       []patch.StrategicMerge{patchGood1, patchGood2},
=======
			input:       []types.PatchStrategicMerge{patchGood1, patchGood2},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			expectedOut: []*Resource{testDeployment, testConfigMap},
			expectedErr: false,
		},
		{
			name:        "badFileName",
<<<<<<< HEAD
			input:       []patch.StrategicMerge{patchGood1, "doesNotExist"},
=======
			input:       []types.PatchStrategicMerge{patchGood1, "doesNotExist"},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			expectedOut: []*Resource{},
			expectedErr: true,
		},
		{
			name:        "badData",
<<<<<<< HEAD
			input:       []patch.StrategicMerge{patchGood1, patchBad},
=======
			input:       []types.PatchStrategicMerge{patchGood1, patchBad},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			expectedOut: []*Resource{},
			expectedErr: true,
		},
		{
			name:        "listOfPatches",
<<<<<<< HEAD
			input:       []patch.StrategicMerge{patchList},
=======
			input:       []types.PatchStrategicMerge{patchList},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			expectedOut: []*Resource{testDeployment, testConfigMap},
			expectedErr: false,
		},
		{
			name:        "listWithAnchorReference",
<<<<<<< HEAD
			input:       []patch.StrategicMerge{patchList2},
=======
			input:       []types.PatchStrategicMerge{patchList2},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			expectedOut: []*Resource{testDeploymentA, testDeploymentB},
			expectedErr: false,
		},
		{
			name:        "listWithNoEntries",
<<<<<<< HEAD
			input:       []patch.StrategicMerge{patchList3},
=======
			input:       []types.PatchStrategicMerge{patchList3},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			expectedOut: []*Resource{},
			expectedErr: false,
		},
		{
			name:        "listWithNo'items:'",
<<<<<<< HEAD
			input:       []patch.StrategicMerge{patchList4},
=======
			input:       []types.PatchStrategicMerge{patchList4},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			expectedOut: []*Resource{},
			expectedErr: false,
		},
	}
	for _, test := range tests {
		rs, err := factory.SliceFromPatches(l, test.input)
		if test.expectedErr && err == nil {
			t.Fatalf("%v: should return error", test.name)
		}
		if !test.expectedErr && err != nil {
			t.Fatalf("%v: unexpected error: %s", test.name, err)
		}
		if len(rs) != len(test.expectedOut) {
			t.Fatalf("%s: length mismatch %d != %d",
				test.name, len(rs), len(test.expectedOut))
		}
		for i := range rs {
			if !reflect.DeepEqual(test.expectedOut[i], rs[i]) {
				t.Fatalf("%s: Got: %v\nexpected:%v",
					test.name, test.expectedOut[i], rs[i])
			}
		}
	}
}
