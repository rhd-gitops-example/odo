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

package resmap_test

import (
	"encoding/base64"
<<<<<<< HEAD
	"fmt"
	"reflect"
	"testing"

	"sigs.k8s.io/kustomize/pkg/fs"
	"sigs.k8s.io/kustomize/pkg/gvk"
	"sigs.k8s.io/kustomize/pkg/ifc"
	"sigs.k8s.io/kustomize/pkg/internal/loadertest"
	"sigs.k8s.io/kustomize/pkg/loader"
	"sigs.k8s.io/kustomize/pkg/resid"
	. "sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/types"
)

func TestFromFiles(t *testing.T) {
=======
	"reflect"
	"testing"

	"sigs.k8s.io/kustomize/v3/internal/loadertest"
	"sigs.k8s.io/kustomize/v3/pkg/fs"
	"sigs.k8s.io/kustomize/v3/pkg/gvk"
	"sigs.k8s.io/kustomize/v3/pkg/ifc"
	"sigs.k8s.io/kustomize/v3/pkg/loader"
	. "sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/resmaptest"
	"sigs.k8s.io/kustomize/v3/pkg/types"
	"sigs.k8s.io/kustomize/v3/pkg/validators"
)

func TestFromFile(t *testing.T) {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)

	resourceStr := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: dply1
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dply2
---
# some comment
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: dply2
  namespace: test
---
`
<<<<<<< HEAD

=======
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	l := loadertest.NewFakeLoader("/whatever/project")
	if ferr := l.AddFile("/whatever/project/deployment.yaml", []byte(resourceStr)); ferr != nil {
		t.Fatalf("Error adding fake file: %v\n", ferr)
	}
<<<<<<< HEAD
	expected := ResMap{resid.NewResId(deploy, "dply1"): rf.FromMap(
		map[string]interface{}{
=======
	expected := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "dply1",
<<<<<<< HEAD
			},
		}),
		resid.NewResId(deploy, "dply2"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "dply2",
				},
			}),
		resid.NewResIdWithPrefixNamespace(deploy, "dply2", "", "test"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "apps/v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name":      "dply2",
					"namespace": "test",
				},
			}),
	}

	m, _ := rmF.FromFiles(
		l, []string{"deployment.yaml"})
	if len(m) != 3 {
		t.Fatalf("%#v should contain 3 appResource, but got %d", m, len(m))
	}

	if err := expected.ErrorIfNotEqual(m); err != nil {
=======
			}}).
		Add(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "dply2",
			}}).
		Add(map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name":      "dply2",
				"namespace": "test",
			}}).ResMap()

	m, _ := rmF.FromFile(l, "deployment.yaml")
	if m.Size() != 3 {
		t.Fatalf("result should contain 3, but got %d", m.Size())
	}
	if err := expected.ErrorIfNotEqualLists(m); err != nil {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

func TestFromBytes(t *testing.T) {
	encoded := []byte(`apiVersion: v1
kind: ConfigMap
metadata:
  name: cm1
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cm2
`)
<<<<<<< HEAD
	expected := ResMap{
		resid.NewResId(cmap, "cm1"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "cm1",
				},
			}),
		resid.NewResId(cmap, "cm2"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "cm2",
				},
			}),
	}
	m, err := rmF.NewResMapFromBytes(encoded)
	fmt.Printf("%v\n", m)
=======
	expected := resmaptest_test.NewRmBuilder(t, rf).
		Add(map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": "cm1",
			}}).
		Add(map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": "cm2",
			}}).ResMap()
	m, err := rmF.NewResMapFromBytes(encoded)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(m, expected) {
		t.Fatalf("%#v doesn't match expected %#v", m, expected)
	}
}

var cmap = gvk.Gvk{Version: "v1", Kind: "ConfigMap"}

func TestNewFromConfigMaps(t *testing.T) {
	type testCase struct {
		description string
		input       []types.ConfigMapArgs
		filepath    string
		content     string
		expected    ResMap
	}

	l := loadertest.NewFakeLoader("/whatever/project")
	testCases := []testCase{
		{
			description: "construct config map from env",
			input: []types.ConfigMapArgs{
				{
					GeneratorArgs: types.GeneratorArgs{
						Name: "envConfigMap",
						DataSources: types.DataSources{
<<<<<<< HEAD
							EnvSource: "app.env",
=======
							EnvSources: []string{"app.env"},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
						},
					},
				},
			},
			filepath: "/whatever/project/app.env",
			content:  "DB_USERNAME=admin\nDB_PASSWORD=somepw",
<<<<<<< HEAD
			expected: ResMap{
				resid.NewResId(cmap, "envConfigMap"): rf.FromMapAndOption(
					map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "envConfigMap",
						},
						"data": map[string]interface{}{
							"DB_USERNAME": "admin",
							"DB_PASSWORD": "somepw",
						},
					}, &types.GeneratorArgs{}, nil),
			},
		},
=======
			expected: resmaptest_test.NewRmBuilder(t, rf).Add(
				map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name": "envConfigMap",
					},
					"data": map[string]interface{}{
						"DB_USERNAME": "admin",
						"DB_PASSWORD": "somepw",
					}}).ResMap(),
		},

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		{
			description: "construct config map from file",
			input: []types.ConfigMapArgs{{
				GeneratorArgs: types.GeneratorArgs{
					Name: "fileConfigMap",
					DataSources: types.DataSources{
						FileSources: []string{"app-init.ini"},
					},
				},
			},
			},
			filepath: "/whatever/project/app-init.ini",
			content:  "FOO=bar\nBAR=baz\n",
<<<<<<< HEAD
			expected: ResMap{
				resid.NewResId(cmap, "fileConfigMap"): rf.FromMapAndOption(
					map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "fileConfigMap",
						},
						"data": map[string]interface{}{
							"app-init.ini": `FOO=bar
BAR=baz
`,
						},
					}, &types.GeneratorArgs{}, nil),
			},
=======
			expected: resmaptest_test.NewRmBuilder(t, rf).Add(
				map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name": "fileConfigMap",
					},
					"data": map[string]interface{}{
						"app-init.ini": `FOO=bar
BAR=baz
`,
					},
				}).ResMap(),
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		},
		{
			description: "construct config map from literal",
			input: []types.ConfigMapArgs{
				{
					GeneratorArgs: types.GeneratorArgs{
						Name: "literalConfigMap",
						DataSources: types.DataSources{
							LiteralSources: []string{"a=x", "b=y", "c=\"Good Morning\"", "d=\"false\""},
						},
					},
				},
			},
<<<<<<< HEAD
			expected: ResMap{
				resid.NewResId(cmap, "literalConfigMap"): rf.FromMapAndOption(
					map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "literalConfigMap",
						},
						"data": map[string]interface{}{
							"a": "x",
							"b": "y",
							"c": "Good Morning",
							"d": "false",
						},
					}, &types.GeneratorArgs{}, nil),
			},
		},
		// TODO: add testcase for data coming from multiple sources like
		// files/literal/env etc.
	}
	rmF.Set(l)
=======
			expected: resmaptest_test.NewRmBuilder(t, rf).Add(
				map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata": map[string]interface{}{
						"name": "literalConfigMap",
					},
					"data": map[string]interface{}{
						"a": "x",
						"b": "y",
						"c": "Good Morning",
						"d": "false",
					},
				}).ResMap(),
		},

		// TODO: add testcase for data coming from multiple sources like
		// files/literal/env etc.
	}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	for _, tc := range testCases {
		if ferr := l.AddFile(tc.filepath, []byte(tc.content)); ferr != nil {
			t.Fatalf("Error adding fake file: %v\n", ferr)
		}
<<<<<<< HEAD
		r, err := rmF.NewResMapFromConfigMapArgs(tc.input, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(r, tc.expected) {
			t.Fatalf("in testcase: %q got:\n%+v\n expected:\n%+v\n", tc.description, r, tc.expected)
=======
		r, err := rmF.NewResMapFromConfigMapArgs(l, nil, tc.input)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if err = tc.expected.ErrorIfNotEqualLists(r); err != nil {
			t.Fatalf("testcase: %q, err: %v", tc.description, err)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		}
	}
}

<<<<<<< HEAD
var secret = gvk.Gvk{Version: "v1", Kind: "Secret"}

=======
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
func TestNewResMapFromSecretArgs(t *testing.T) {
	secrets := []types.SecretArgs{
		{
			GeneratorArgs: types.GeneratorArgs{
				Name: "apple",
				DataSources: types.DataSources{
					LiteralSources: []string{
						"DB_USERNAME=admin",
						"DB_PASSWORD=somepw",
					},
				},
			},
			Type: ifc.SecretTypeOpaque,
		},
	}
<<<<<<< HEAD
	fakeFs := fs.MakeFakeFS()
	fakeFs.Mkdir(".")
	rmF.Set(loader.NewFileLoaderAtRoot(fakeFs))
	actual, err := rmF.NewResMapFromSecretArgs(secrets, nil)

=======
	fSys := fs.MakeFsInMemory()
	fSys.Mkdir(".")
	actual, err := rmF.NewResMapFromSecretArgs(
		loader.NewFileLoaderAtRoot(validators.MakeFakeValidator(), fSys), nil, secrets)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

<<<<<<< HEAD
	expected := ResMap{
		resid.NewResId(secret, "apple"): rf.FromMapAndOption(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Secret",
				"metadata": map[string]interface{}{
					"name": "apple",
				},
				"type": ifc.SecretTypeOpaque,
				"data": map[string]interface{}{
					"DB_USERNAME": base64.StdEncoding.EncodeToString([]byte("admin")),
					"DB_PASSWORD": base64.StdEncoding.EncodeToString([]byte("somepw")),
				},
			}, &types.GeneratorArgs{}, nil),
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("%#v\ndoesn't match expected:\n%#v", actual, expected)
=======
	expected := resmaptest_test.NewRmBuilder(t, rf).Add(
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name": "apple",
			},
			"type": ifc.SecretTypeOpaque,
			"data": map[string]interface{}{
				"DB_USERNAME": base64.StdEncoding.EncodeToString([]byte("admin")),
				"DB_PASSWORD": base64.StdEncoding.EncodeToString([]byte("somepw")),
			},
		}).ResMap()
	if err = expected.ErrorIfNotEqualLists(actual); err != nil {
		t.Fatalf("error: %s", err)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
}
