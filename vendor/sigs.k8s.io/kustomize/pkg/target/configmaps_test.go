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

package target_test

import (
	"testing"
<<<<<<< HEAD
=======

	"sigs.k8s.io/kustomize/v3/pkg/kusttest"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

// Generate a Secret and a ConfigMap from the same data
// to compare the result.
func TestGeneratorBasics(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/app")
	th.writeK("/app", `
=======
	th := kusttest_test.NewKustTestHarness(t, "/app")
	th.WriteK("/app", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
namePrefix: blah-
configMapGenerator:
- name: bob
  literals:
<<<<<<< HEAD
    - fruit=apple
    - vegetable=broccoli
  env: foo.env
  files:
    - passphrase=phrase.dat
    - forces.txt
- name: json
  literals:
    - 'v2=[{"path": "var/druid/segment-cache"}]'
secretGenerator:
- name: bob
  literals:
    - fruit=apple
    - vegetable=broccoli
  env: foo.env
  files:
    - passphrase=phrase.dat
    - forces.txt
`)
	th.writeF("/app/foo.env", `
MOUNTAIN=everest
OCEAN=pacific
`)
	th.writeF("/app/phrase.dat", `
=======
  - fruit=apple
  - vegetable=broccoli
  envs:
  - foo.env
  files:
  - passphrase=phrase.dat
  - forces.txt
- name: json
  literals:
  - 'v2=[{"path": "var/druid/segment-cache"}]'
  - >- 
    druid_segmentCache_locations=[{"path": 
    "var/druid/segment-cache", 
    "maxSize": 32000000000, 
    "freeSpacePercent": 1.0}]
secretGenerator:
- name: bob
  literals:
  - fruit=apple
  - vegetable=broccoli
  envs:
  - foo.env
  files:
  - passphrase=phrase.dat
  - forces.txt
`)
	th.WriteF("/app/foo.env", `
MOUNTAIN=everest
OCEAN=pacific
`)
	th.WriteF("/app/phrase.dat", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
Life is short.
But the years are long.
Not while the evil days come not.
`)
<<<<<<< HEAD
	th.writeF("/app/forces.txt", `
=======
	th.WriteF("/app/forces.txt", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
gravitational
electromagnetic
strong nuclear
weak nuclear
`)
<<<<<<< HEAD
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.assertActualEqualsExpected(m, `
=======
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: v1
data:
  MOUNTAIN: everest
  OCEAN: pacific
  forces.txt: |2

    gravitational
    electromagnetic
    strong nuclear
    weak nuclear
  fruit: apple
  passphrase: |2

    Life is short.
    But the years are long.
    Not while the evil days come not.
  vegetable: broccoli
kind: ConfigMap
metadata:
  name: blah-bob-k772g5db55
---
apiVersion: v1
data:
<<<<<<< HEAD
  v2: '[{"path": "var/druid/segment-cache"}]'
kind: ConfigMap
metadata:
  name: blah-json-tkh79m5tbc
=======
  druid_segmentCache_locations: '[{"path":  "var/druid/segment-cache",  "maxSize":
    32000000000,  "freeSpacePercent": 1.0}]'
  v2: '[{"path": "var/druid/segment-cache"}]'
kind: ConfigMap
metadata:
  name: blah-json-9gtcc2fgb4
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
---
apiVersion: v1
data:
  MOUNTAIN: ZXZlcmVzdA==
  OCEAN: cGFjaWZpYw==
  forces.txt: CmdyYXZpdGF0aW9uYWwKZWxlY3Ryb21hZ25ldGljCnN0cm9uZyBudWNsZWFyCndlYWsgbnVjbGVhcgo=
  fruit: YXBwbGU=
  passphrase: CkxpZmUgaXMgc2hvcnQuCkJ1dCB0aGUgeWVhcnMgYXJlIGxvbmcuCk5vdCB3aGlsZSB0aGUgZXZpbCBkYXlzIGNvbWUgbm90Lgo=
  vegetable: YnJvY2NvbGk=
kind: Secret
metadata:
  name: blah-bob-gmc2824f4b
type: Opaque
`)
}

// TODO: These should be errors instead.
func TestGeneratorRepeatsInKustomization(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/app")
	th.writeK("/app", `
=======
	th := kusttest_test.NewKustTestHarness(t, "/app")
	th.WriteK("/app", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
namePrefix: blah-
configMapGenerator:
- name: bob
  behavior: create
  literals:
<<<<<<< HEAD
    - bean=pinto
    - star=wolf-rayet
  literals:
    - fruit=apple
    - vegetable=broccoli
  files:
    - forces.txt
  files:
    - nobles=nobility.txt
`)
	th.writeF("/app/forces.txt", `
=======
  - bean=pinto
  - star=wolf-rayet
  literals:
  - fruit=apple
  - vegetable=broccoli
  files:
  - forces.txt
  files:
  - nobles=nobility.txt
`)
	th.WriteF("/app/forces.txt", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
gravitational
electromagnetic
strong nuclear
weak nuclear
`)
<<<<<<< HEAD
	th.writeF("/app/nobility.txt", `
=======
	th.WriteF("/app/nobility.txt", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
helium
neon
argon
krypton
xenon
radon
`)
<<<<<<< HEAD
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.assertActualEqualsExpected(m, `
=======
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: v1
data:
  fruit: apple
  nobles: |2

    helium
    neon
    argon
    krypton
    xenon
    radon
  vegetable: broccoli
kind: ConfigMap
metadata:
  name: blah-bob-gfkcbk5ckf
`)
}

func TestGeneratorOverlays(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/app/overlay")
	th.writeK("/app/base1", `
=======
	th := kusttest_test.NewKustTestHarness(t, "/app/overlay")
	th.WriteK("/app/base1", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
namePrefix: p1-
configMapGenerator:
- name: com1
  behavior: create
  literals:
<<<<<<< HEAD
    - from=base
`)
	th.writeK("/app/base2", `
=======
  - from=base
`)
	th.WriteK("/app/base2", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
namePrefix: p2-
configMapGenerator:
- name: com2
  behavior: create
  literals:
<<<<<<< HEAD
    - from=base
`)
	th.writeK("/app/overlay/o1", `
bases:
=======
  - from=base
`)
	th.WriteK("/app/overlay/o1", `
resources:
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
- ../../base1
configMapGenerator:
- name: com1
  behavior: merge
  literals:
<<<<<<< HEAD
    - from=overlay
`)
	th.writeK("/app/overlay/o2", `
bases:
=======
  - from=overlay
`)
	th.WriteK("/app/overlay/o2", `
resources:
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
- ../../base2
configMapGenerator:
- name: com2
  behavior: merge
  literals:
<<<<<<< HEAD
    - from=overlay
`)
	th.writeK("/app/overlay", `
bases:
=======
  - from=overlay
`)
	th.WriteK("/app/overlay", `
resources:
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
- o1
- o2
configMapGenerator:
- name: com1
  behavior: merge
  literals:
<<<<<<< HEAD
    - foo=bar
    - baz=qux
`)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.assertActualEqualsExpected(m, `
=======
  - foo=bar
  - baz=qux
`)
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: v1
data:
  baz: qux
  foo: bar
  from: overlay
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: p1-com1-dhbbm922gd
---
apiVersion: v1
data:
  from: overlay
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: p2-com2-c4b8md75k9
`)
}
<<<<<<< HEAD
=======

func TestConfigMapGeneratorMergeNamePrefix(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app")
	th.WriteK("/app/base", `
configMapGenerator:
- name: cm
  behavior: create
  literals:
  - foo=bar
`)
	th.WriteK("/app/o1", `
resources:
- ../base
namePrefix: o1-
`)
	th.WriteK("/app/o2", `
resources:
- ../base
nameSuffix: -o2
`)
	th.WriteK("/app", `
resources:
- o1
- o2
configMapGenerator:
- name: o1-cm
  behavior: merge
  literals:
  - big=bang
- name: cm-o2
  behavior: merge
  literals:
  - big=crunch
`)
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
data:
  big: bang
  foo: bar
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: o1-cm-28g596k77k
---
apiVersion: v1
data:
  big: crunch
  foo: bar
kind: ConfigMap
metadata:
  annotations: {}
  labels: {}
  name: cm-o2-gfcc59fg5m
`)
}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
