<<<<<<< HEAD
=======
// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
package transformers

import (
	"reflect"
	"testing"

<<<<<<< HEAD
	"sigs.k8s.io/kustomize/pkg/resid"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/transformers/config"
=======
	"sigs.k8s.io/kustomize/v3/pkg/gvk"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/resmaptest"
	"sigs.k8s.io/kustomize/v3/pkg/transformers/config"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

func TestVarRef(t *testing.T) {
	type given struct {
<<<<<<< HEAD
		varMap map[string]string
=======
		varMap map[string]interface{}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		fs     []config.FieldSpec
		res    resmap.ResMap
	}
	type expected struct {
		res    resmap.ResMap
		unused []string
	}
	testCases := []struct {
		description string
		given       given
		expected    expected
	}{
		{
			description: "var replacement in map[string]",
			given: given{
<<<<<<< HEAD
				varMap: map[string]string{
					"FOO": "replacementForFoo",
					"BAR": "replacementForBar",
				},
				fs: []config.FieldSpec{
					{Gvk: cmap, Path: "data"},
				},
				res: resmap.ResMap{
					resid.NewResId(cmap, "cm1"): rf.FromMap(
						map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name": "cm1",
							},
							"data": map[string]interface{}{
								"item1": "$(FOO)",
								"item2": "bla",
							},
						}),
				},
			},
			expected: expected{
				res: resmap.ResMap{
					resid.NewResId(cmap, "cm1"): rf.FromMap(
						map[string]interface{}{
							"apiVersion": "v1",
							"kind":       "ConfigMap",
							"metadata": map[string]interface{}{
								"name": "cm1",
							},
							"data": map[string]interface{}{
								"item1": "replacementForFoo",
								"item2": "bla",
							},
						}),
				},
=======
				varMap: map[string]interface{}{
					"FOO": "replacementForFoo",
					"BAR": "replacementForBar",
					"BAZ": int64(5),
					"BOO": true,
				},
				fs: []config.FieldSpec{
					{Gvk: gvk.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/map"},
					{Gvk: gvk.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/slice"},
					{Gvk: gvk.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/interface"},
					{Gvk: gvk.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/nil"},
					{Gvk: gvk.Gvk{Version: "v1", Kind: "ConfigMap"}, Path: "data/num"},
				},
				res: resmaptest_test.NewRmBuilder(t, rf).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"map": map[string]interface{}{
								"item1": "$(FOO)",
								"item2": "bla",
								"item3": "$(BAZ)",
								"item4": "$(BAZ)+$(BAZ)",
								"item5": "$(BOO)",
								"item6": "if $(BOO)",
								"item7": 2019,
							},
							"slice": []interface{}{
								"$(FOO)",
								"bla",
								"$(BAZ)",
								"$(BAZ)+$(BAZ)",
								"$(BOO)",
								"if $(BOO)",
							},
							"interface": "$(FOO)",
							"nil":       nil,
							"num":       2019,
						}}).ResMap(),
			},
			expected: expected{
				res: resmaptest_test.NewRmBuilder(t, rf).
					Add(map[string]interface{}{
						"apiVersion": "v1",
						"kind":       "ConfigMap",
						"metadata": map[string]interface{}{
							"name": "cm1",
						},
						"data": map[string]interface{}{
							"map": map[string]interface{}{
								"item1": "replacementForFoo",
								"item2": "bla",
								"item3": int64(5),
								"item4": "5+5",
								"item5": true,
								"item6": "if true",
								"item7": 2019,
							},
							"slice": []interface{}{
								"replacementForFoo",
								"bla",
								int64(5),
								"5+5",
								true,
								"if true",
							},
							"interface": "replacementForFoo",
							"nil":       nil,
							"num":       2019,
						}}).ResMap(),
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
				unused: []string{"BAR"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			// arrange
			tr := NewRefVarTransformer(tc.given.varMap, tc.given.fs)

			// act
			err := tr.Transform(tc.given.res)

			// assert
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			a, e := tc.given.res, tc.expected.res
			if !reflect.DeepEqual(a, e) {
<<<<<<< HEAD
				err = e.ErrorIfNotEqual(a)
=======
				err = e.ErrorIfNotEqualLists(a)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
				t.Fatalf("actual doesn't match expected: \nACTUAL:\n%v\nEXPECTED:\n%v\nERR: %v", a, e, err)
			}

		})
	}
}
