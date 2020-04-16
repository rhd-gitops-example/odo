<<<<<<< HEAD
=======
// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
package resid

import (
	"testing"

<<<<<<< HEAD
	"sigs.k8s.io/kustomize/pkg/gvk"
=======
	"sigs.k8s.io/kustomize/v3/pkg/gvk"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

var stringTests = []struct {
	x ResId
	s string
}{
	{
		ResId{
<<<<<<< HEAD
			namespace: "ns",
			gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "p",
			suffix:    "s",
		},
		"g_v_k|ns|p|nm|s",
	},
	{
		ResId{
			namespace: "ns",
			gvKind:    gvk.Gvk{Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "p",
			suffix:    "s",
		},
		"~G_v_k|ns|p|nm|s",
	},
	{
		ResId{
			namespace: "ns",
			gvKind:    gvk.Gvk{Kind: "k"},
			name:      "nm",
			prefix:    "p",
			suffix:    "s",
		},
		"~G_~V_k|ns|p|nm|s",
	},
	{
		ResId{
			namespace: "ns",
			gvKind:    gvk.Gvk{},
			name:      "nm",
			prefix:    "p",
			suffix:    "s",
		},
		"~G_~V_~K|ns|p|nm|s",
	},
	{
		ResId{
			gvKind: gvk.Gvk{},
			name:   "nm",
			prefix: "p",
			suffix: "s",
		},
		"~G_~V_~K|~X|p|nm|s",
	},
	{
		ResId{
			gvKind: gvk.Gvk{},
			name:   "nm",
			suffix: "s",
		},
		"~G_~V_~K|~X|~P|nm|s",
	},
	{
		ResId{
			gvKind: gvk.Gvk{},
			suffix: "s",
		},
		"~G_~V_~K|~X|~P|~N|s",
	},
	{
		ResId{
			gvKind: gvk.Gvk{},
		},
		"~G_~V_~K|~X|~P|~N|~S",
	},
	{
		ResId{},
		"~G_~V_~K|~X|~P|~N|~S",
=======
			Namespace: "ns",
			Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
			Name:      "nm",
		},
		"g_v_k|ns|nm",
	},
	{
		ResId{
			Namespace: "ns",
			Gvk:       gvk.Gvk{Version: "v", Kind: "k"},
			Name:      "nm",
		},
		"~G_v_k|ns|nm",
	},
	{
		ResId{
			Namespace: "ns",
			Gvk:       gvk.Gvk{Kind: "k"},
			Name:      "nm",
		},
		"~G_~V_k|ns|nm",
	},
	{
		ResId{
			Namespace: "ns",
			Gvk:       gvk.Gvk{},
			Name:      "nm",
		},
		"~G_~V_~K|ns|nm",
	},
	{
		ResId{
			Gvk:  gvk.Gvk{},
			Name: "nm",
		},
		"~G_~V_~K|~X|nm",
	},
	{
		ResId{
			Gvk:  gvk.Gvk{},
			Name: "nm",
		},
		"~G_~V_~K|~X|nm",
	},
	{
		ResId{
			Gvk: gvk.Gvk{},
		},
		"~G_~V_~K|~X|~N",
	},
	{
		ResId{
			Gvk: gvk.Gvk{},
		},
		"~G_~V_~K|~X|~N",
	},
	{
		ResId{},
		"~G_~V_~K|~X|~N",
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	},
}

func TestString(t *testing.T) {
	for _, hey := range stringTests {
		if hey.x.String() != hey.s {
			t.Fatalf("Actual: %v,  Expected: '%s'", hey.x, hey.s)
		}
	}
}

var gvknStringTests = []struct {
	x ResId
	s string
}{
	{
		ResId{
<<<<<<< HEAD
			namespace: "ns",
			gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "p",
			suffix:    "s",
=======
			Namespace: "ns",
			Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
			Name:      "nm",
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		},
		"g_v_k|nm",
	},
	{
		ResId{
<<<<<<< HEAD
			namespace: "ns",
			gvKind:    gvk.Gvk{Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "p",
			suffix:    "s",
=======
			Namespace: "ns",
			Gvk:       gvk.Gvk{Version: "v", Kind: "k"},
			Name:      "nm",
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		},
		"~G_v_k|nm",
	},
	{
		ResId{
<<<<<<< HEAD
			namespace: "ns",
			gvKind:    gvk.Gvk{Kind: "k"},
			name:      "nm",
			prefix:    "p",
			suffix:    "s",
=======
			Namespace: "ns",
			Gvk:       gvk.Gvk{Kind: "k"},
			Name:      "nm",
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		},
		"~G_~V_k|nm",
	},
	{
		ResId{
<<<<<<< HEAD
			namespace: "ns",
			gvKind:    gvk.Gvk{},
			name:      "nm",
			prefix:    "p",
			suffix:    "s",
=======
			Namespace: "ns",
			Gvk:       gvk.Gvk{},
			Name:      "nm",
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		},
		"~G_~V_~K|nm",
	},
	{
		ResId{
<<<<<<< HEAD
			gvKind: gvk.Gvk{},
			name:   "nm",
			prefix: "p",
			suffix: "s",
=======
			Gvk:  gvk.Gvk{},
			Name: "nm",
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		},
		"~G_~V_~K|nm",
	},
	{
		ResId{
<<<<<<< HEAD
			gvKind: gvk.Gvk{},
			name:   "nm",
			suffix: "s",
=======
			Gvk:  gvk.Gvk{},
			Name: "nm",
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		},
		"~G_~V_~K|nm",
	},
	{
		ResId{
<<<<<<< HEAD
			gvKind: gvk.Gvk{},
			suffix: "s",
=======
			Gvk: gvk.Gvk{},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		},
		"~G_~V_~K|",
	},
	{
		ResId{
<<<<<<< HEAD
			gvKind: gvk.Gvk{},
=======
			Gvk: gvk.Gvk{},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		},
		"~G_~V_~K|",
	},
	{
		ResId{},
		"~G_~V_~K|",
	},
}

func TestGvknString(t *testing.T) {
	for _, hey := range gvknStringTests {
		if hey.x.GvknString() != hey.s {
			t.Fatalf("Actual: %s,  Expected: '%s'", hey.x.GvknString(), hey.s)
		}
	}
}

<<<<<<< HEAD
var GvknEqualsTest = []struct {
	id1          ResId
	id2          ResId
	gVknResult   bool
	nSgVknResult bool
}{
	{
		id1: ResId{
			namespace: "X",
			gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "AA",
			suffix:    "aa",
		},
		id2: ResId{
			namespace: "X",
			gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "BB",
			suffix:    "bb",
		},
		gVknResult:   true,
		nSgVknResult: true,
	},
	{
		id1: ResId{
			namespace: "X",
			gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "AA",
			suffix:    "aa",
		},
		id2: ResId{
			namespace: "Z",
			gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "BB",
			suffix:    "bb",
		},
		gVknResult:   true,
		nSgVknResult: false,
	},
	{
		id1: ResId{
			namespace: "X",
			gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "AA",
			suffix:    "aa",
		},
		id2: ResId{
			gvKind: gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
			name:   "nm",
			prefix: "BB",
			suffix: "bb",
		},
		gVknResult:   true,
		nSgVknResult: false,
	},
	{
		id1: ResId{
			namespace: "X",
			gvKind:    gvk.Gvk{Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "AA",
			suffix:    "aa",
		},
		id2: ResId{
			namespace: "Z",
			gvKind:    gvk.Gvk{Version: "v", Kind: "k"},
			name:      "nm",
			prefix:    "BB",
			suffix:    "bb",
		},
		gVknResult:   true,
		nSgVknResult: false,
	},
	{
		id1: ResId{
			namespace: "X",
			gvKind:    gvk.Gvk{Kind: "k"},
			name:      "nm",
			prefix:    "AA",
			suffix:    "aa",
		},
		id2: ResId{
			namespace: "Z",
			gvKind:    gvk.Gvk{Kind: "k"},
			name:      "nm",
			prefix:    "BB",
			suffix:    "bb",
		},
		gVknResult:   true,
		nSgVknResult: false,
	},
	{
		id1: ResId{
			namespace: "X",
			name:      "nm",
			prefix:    "AA",
			suffix:    "aa",
		},
		id2: ResId{
			namespace: "Z",
			name:      "nm",
			prefix:    "BB",
			suffix:    "bb",
		},
		gVknResult:   true,
		nSgVknResult: false,
	},
}

func TestEquals(t *testing.T) {
=======
func TestEquals(t *testing.T) {

	var GvknEqualsTest = []struct {
		id1        ResId
		id2        ResId
		gVknResult bool
		nsEquals   bool
		equals     bool
	}{
		{
			id1: ResId{
				Namespace: "X",
				Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
				Name:      "nm",
			},
			id2: ResId{
				Namespace: "X",
				Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
				Name:      "nm",
			},
			gVknResult: true,
			nsEquals:   true,
			equals:     true,
		},
		{
			id1: ResId{
				Namespace: "X",
				Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
				Name:      "nm",
			},
			id2: ResId{
				Namespace: "Z",
				Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
				Name:      "nm",
			},
			gVknResult: true,
			nsEquals:   false,
			equals:     false,
		},
		{
			id1: ResId{
				Namespace: "X",
				Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
				Name:      "nm",
			},
			id2: ResId{
				Gvk:  gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
				Name: "nm",
			},
			gVknResult: true,
			nsEquals:   false,
			equals:     false,
		},
		{
			id1: ResId{
				Namespace: "X",
				Gvk:       gvk.Gvk{Version: "v", Kind: "k"},
				Name:      "nm",
			},
			id2: ResId{
				Namespace: "Z",
				Gvk:       gvk.Gvk{Version: "v", Kind: "k"},
				Name:      "nm",
			},
			gVknResult: true,
			nsEquals:   false,
			equals:     false,
		},
		{
			id1: ResId{
				Namespace: "X",
				Gvk:       gvk.Gvk{Kind: "k"},
				Name:      "nm",
			},
			id2: ResId{
				Namespace: "Z",
				Gvk:       gvk.Gvk{Kind: "k"},
				Name:      "nm",
			},
			gVknResult: true,
			nsEquals:   false,
			equals:     false,
		},
		{
			id1: ResId{
				Gvk:  gvk.Gvk{Kind: "k"},
				Name: "nm",
			},
			id2: ResId{
				Gvk:  gvk.Gvk{Kind: "k"},
				Name: "nm2",
			},
			gVknResult: false,
			nsEquals:   true,
			equals:     false,
		},
		{
			id1: ResId{
				Gvk:  gvk.Gvk{Kind: "k"},
				Name: "nm",
			},
			id2: ResId{
				Gvk:  gvk.Gvk{Kind: "Node"},
				Name: "nm",
			},
			gVknResult: false,
			nsEquals:   false,
			equals:     false,
		},
		{
			id1: ResId{
				Gvk:  gvk.Gvk{Kind: "Node"},
				Name: "nm1",
			},
			id2: ResId{
				Gvk:  gvk.Gvk{Kind: "Node"},
				Name: "nm2",
			},
			gVknResult: false,
			nsEquals:   true,
			equals:     false,
		},
		{
			id1: ResId{
				Namespace: "default",
				Gvk:       gvk.Gvk{Kind: "k"},
				Name:      "nm1",
			},
			id2: ResId{
				Gvk:  gvk.Gvk{Kind: "k"},
				Name: "nm2",
			},
			gVknResult: false,
			nsEquals:   true,
			equals:     false,
		},
		{
			id1: ResId{
				Namespace: "X",
				Name:      "nm",
			},
			id2: ResId{
				Namespace: "Z",
				Name:      "nm",
			},
			gVknResult: true,
			nsEquals:   false,
			equals:     false,
		},
	}

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	for _, tst := range GvknEqualsTest {
		if tst.id1.GvknEquals(tst.id2) != tst.gVknResult {
			t.Fatalf("GvknEquals(\n%v,\n%v\n) should be %v",
				tst.id1, tst.id2, tst.gVknResult)
		}
<<<<<<< HEAD
		if tst.id1.NsGvknEquals(tst.id2) != tst.nSgVknResult {
			t.Fatalf("NsGvknEquals(\n%v,\n%v\n) should be %v",
				tst.id1, tst.id2, tst.nSgVknResult)
=======
		if tst.id1.IsNsEquals(tst.id2) != tst.nsEquals {
			t.Fatalf("IsNsEquals(\n%v,\n%v\n) should be %v",
				tst.id1, tst.id2, tst.equals)
		}
		if tst.id1.Equals(tst.id2) != tst.equals {
			t.Fatalf("Equals(\n%v,\n%v\n) should be %v",
				tst.id1, tst.id2, tst.equals)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		}
	}
}

<<<<<<< HEAD
func TestCopyWithNewPrefixSuffix(t *testing.T) {
	r1 := ResId{
		namespace: "X",
		gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
		name:      "nm",
		prefix:    "a",
		suffix:    "b",
	}
	r2 := r1.CopyWithNewPrefixSuffix("p-", "-s")
	expected := ResId{
		namespace: "X",
		gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
		name:      "nm",
		prefix:    "p-a",
		suffix:    "b-s",
	}
	if !r2.GvknEquals(expected) {
		t.Fatalf("%v should equal %v", r2, expected)
	}
}

func TestCopyWithNewNamespace(t *testing.T) {
	r1 := ResId{
		namespace: "X",
		gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
		name:      "nm",
		prefix:    "a",
		suffix:    "b",
	}
	r2 := r1.CopyWithNewNamespace("zzz")
	expected := ResId{
		namespace: "zzz",
		gvKind:    gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
		name:      "nm",
		prefix:    "a",
		suffix:    "b",
	}
	if !r2.GvknEquals(expected) {
		t.Fatalf("%v should equal %v", r2, expected)
=======
var ids = []ResId{
	{
		Namespace: "ns",
		Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
		Name:      "nm",
	},
	{
		Namespace: "ns",
		Gvk:       gvk.Gvk{Version: "v", Kind: "k"},
		Name:      "nm",
	},
	{
		Namespace: "ns",
		Gvk:       gvk.Gvk{Kind: "k"},
		Name:      "nm",
	},
	{
		Namespace: "ns",
		Gvk:       gvk.Gvk{},
		Name:      "nm",
	},
	{
		Gvk:  gvk.Gvk{},
		Name: "nm",
	},
	{
		Gvk:  gvk.Gvk{},
		Name: "nm",
	},
	{
		Gvk: gvk.Gvk{},
	},
}

func TestFromString(t *testing.T) {
	for _, id := range ids {
		newId := FromString(id.String())
		if newId != id {
			t.Fatalf("Actual: %v,  Expected: '%s'", newId, id)
		}
	}
}

func TestEffectiveNamespace(t *testing.T) {
	var test = []struct {
		id       ResId
		expected string
	}{
		{
			id: ResId{
				Gvk:  gvk.Gvk{Group: "g", Version: "v", Kind: "Node"},
				Name: "nm",
			},
			expected: TotallyNotANamespace,
		},
		{
			id: ResId{
				Namespace: "foo",
				Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "Node"},
				Name:      "nm",
			},
			expected: TotallyNotANamespace,
		},
		{
			id: ResId{
				Namespace: "foo",
				Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
				Name:      "nm",
			},
			expected: "foo",
		},
		{
			id: ResId{
				Namespace: "",
				Gvk:       gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
				Name:      "nm",
			},
			expected: DefaultNamespace,
		},
		{
			id: ResId{
				Gvk:  gvk.Gvk{Group: "g", Version: "v", Kind: "k"},
				Name: "nm",
			},
			expected: DefaultNamespace,
		},
	}

	for _, tst := range test {
		if actual := tst.id.EffectiveNamespace(); actual != tst.expected {
			t.Fatalf("EffectiveNamespace was %s, expected %s",
				actual, tst.expected)
		}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
}
