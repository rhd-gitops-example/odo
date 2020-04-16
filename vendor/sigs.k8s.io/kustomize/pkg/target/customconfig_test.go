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
)

func makeBaseReferencingCustomConfig(th *KustTestHarness) {
	th.writeK("/app/base", `
=======

	"sigs.k8s.io/kustomize/v3/pkg/kusttest"
)

func makeBaseReferencingCustomConfig(th *kusttest_test.KustTestHarness) {
	th.WriteK("/app/base", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
namePrefix: x-
commonLabels:
  app: myApp
vars:
- name: APRIL_DIET
  objref:
    kind: Giraffe
    name: april
  fieldref:
    fieldpath: spec.diet
- name: KOKO_DIET
  objref:
    kind: Gorilla
    name: koko
  fieldref:
    fieldpath: spec.diet
resources:
<<<<<<< HEAD
- giraffes.yaml
- gorilla.yaml
- animalPark.yaml
=======
- animalPark.yaml
- giraffes.yaml
- gorilla.yaml
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
configurations:
- config/defaults.yaml
- config/custom.yaml
`)
<<<<<<< HEAD
	th.writeF("/app/base/giraffes.yaml", `
kind: Giraffe
metadata:
  name: may
spec:
  diet: acacia
  location: SE
---
kind: Giraffe
metadata:
  name: april
spec:
  diet: mimosa
  location: NE
`)
	th.writeF("/app/base/gorilla.yaml", `
=======
	th.WriteF("/app/base/giraffes.yaml", `
kind: Giraffe
metadata:
  name: april
spec:
  diet: mimosa
  location: NE
---
kind: Giraffe
metadata:
  name: may
spec:
  diet: acacia
  location: SE
`)
	th.WriteF("/app/base/gorilla.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
kind: Gorilla
metadata:
  name: koko
spec:
  diet: bambooshoots
  location: SW
`)
<<<<<<< HEAD
	th.writeF("/app/base/animalPark.yaml", `
=======
	th.WriteF("/app/base/animalPark.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
kind: AnimalPark
metadata:
  name: sandiego
spec:
  gorillaRef:
    name: koko
  giraffeRef:
    name: april
  food:
  - "$(APRIL_DIET)"
  - "$(KOKO_DIET)"
`)
}

func TestCustomConfig(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/app/base")
	makeBaseReferencingCustomConfig(th)
	th.writeDefaultConfigs("/app/base/config/defaults.yaml")
	th.writeF("/app/base/config/custom.yaml", `
=======
	th := kusttest_test.NewKustTestHarness(t, "/app/base")
	makeBaseReferencingCustomConfig(th)
	th.WriteDefaultConfigs("/app/base/config/defaults.yaml")
	th.WriteF("/app/base/config/custom.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
nameReference:
- kind: Gorilla
  fieldSpecs:
  - kind: AnimalPark
    path: spec/gorillaRef/name
- kind: Giraffe
  fieldSpecs:
  - kind: AnimalPark
    path: spec/giraffeRef/name
varReference:
- path: spec/food
  kind: AnimalPark
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
kind: AnimalPark
metadata:
  labels:
    app: myApp
  name: x-sandiego
spec:
  food:
  - mimosa
  - bambooshoots
  giraffeRef:
    name: x-april
  gorillaRef:
    name: x-koko
---
kind: Giraffe
metadata:
  labels:
    app: myApp
  name: x-april
spec:
  diet: mimosa
  location: NE
---
kind: Giraffe
metadata:
  labels:
    app: myApp
  name: x-may
spec:
  diet: acacia
  location: SE
---
kind: Gorilla
metadata:
  labels:
    app: myApp
  name: x-koko
spec:
  diet: bambooshoots
  location: SW
`)
}

func TestCustomConfigWithDefaultOverspecification(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/app/base")
	makeBaseReferencingCustomConfig(th)
	th.writeDefaultConfigs("/app/base/config/defaults.yaml")
	// Specifying namePrefix here conflicts with (is the same as)
	// the defaults written above.  This is intentional in the
	// test to assure duplicate config doesn't cause problems.
	th.writeF("/app/base/config/custom.yaml", `
=======
	th := kusttest_test.NewKustTestHarness(t, "/app/base")
	makeBaseReferencingCustomConfig(th)
	th.WriteDefaultConfigs("/app/base/config/defaults.yaml")
	// Specifying namePrefix here conflicts with (is the same as)
	// the defaults written above.  This is intentional in the
	// test to assure duplicate config doesn't cause problems.
	th.WriteF("/app/base/config/custom.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
namePrefix:
- path: metadata/name
nameReference:
- kind: Gorilla
  fieldSpecs:
  - kind: AnimalPark
    path: spec/gorillaRef/name
- kind: Giraffe
  fieldSpecs:
  - kind: AnimalPark
    path: spec/giraffeRef/name
varReference:
- path: spec/food
  kind: AnimalPark
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
kind: AnimalPark
metadata:
  labels:
    app: myApp
  name: x-sandiego
spec:
  food:
  - mimosa
  - bambooshoots
  giraffeRef:
    name: x-april
  gorillaRef:
    name: x-koko
---
kind: Giraffe
metadata:
  labels:
    app: myApp
  name: x-april
spec:
  diet: mimosa
  location: NE
---
kind: Giraffe
metadata:
  labels:
    app: myApp
  name: x-may
spec:
  diet: acacia
  location: SE
---
kind: Gorilla
metadata:
  labels:
    app: myApp
  name: x-koko
spec:
  diet: bambooshoots
  location: SW
`)
}

func TestFixedBug605_BaseCustomizationAvailableInOverlay(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/app/overlay")
	makeBaseReferencingCustomConfig(th)
	th.writeDefaultConfigs("/app/base/config/defaults.yaml")
	th.writeF("/app/base/config/custom.yaml", `
=======
	th := kusttest_test.NewKustTestHarness(t, "/app/overlay")
	makeBaseReferencingCustomConfig(th)
	th.WriteDefaultConfigs("/app/base/config/defaults.yaml")
	th.WriteF("/app/base/config/custom.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
nameReference:
- kind: Gorilla
  fieldSpecs:
  - kind: AnimalPark
    path: spec/gorillaRef/name
- kind: Giraffe
  fieldSpecs:
  - kind: AnimalPark
    path: spec/giraffeRef/name
varReference:
- path: spec/food
  kind: AnimalPark
`)
<<<<<<< HEAD
	th.writeK("/app/overlay", `
=======
	th.WriteK("/app/overlay", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
namePrefix: o-
commonLabels:
  movie: planetOfTheApes
patchesStrategicMerge:
- animalPark.yaml
resources:
<<<<<<< HEAD
- ursus.yaml
bases:
- ../base
`)
	th.writeF("/app/overlay/ursus.yaml", `
=======
- ../base
- ursus.yaml
`)
	th.WriteF("/app/overlay/ursus.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
kind: Gorilla
metadata:
  name: ursus
spec:
  diet: heston
  location: Arizona
`)
	// The following replaces the gorillaRef in the AnimalPark.
<<<<<<< HEAD
	th.writeF("/app/overlay/animalPark.yaml", `
=======
	th.WriteF("/app/overlay/animalPark.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
kind: AnimalPark
metadata:
  name: sandiego
spec:
  gorillaRef:
    name: ursus
`)

<<<<<<< HEAD
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	// TODO(#669): The name of AnimalPark should be x-o-sandiego,
	// not o-sandiego, since AnimalPark appears in the base.
	th.assertActualEqualsExpected(m, `
=======
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	th.AssertActualEqualsExpected(m, `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
kind: AnimalPark
metadata:
  labels:
    app: myApp
    movie: planetOfTheApes
<<<<<<< HEAD
  name: o-sandiego
=======
  name: o-x-sandiego
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
spec:
  food:
  - mimosa
  - bambooshoots
  giraffeRef:
    name: o-x-april
  gorillaRef:
    name: o-ursus
---
kind: Giraffe
metadata:
  labels:
    app: myApp
    movie: planetOfTheApes
  name: o-x-april
spec:
  diet: mimosa
  location: NE
---
kind: Giraffe
metadata:
  labels:
    app: myApp
    movie: planetOfTheApes
  name: o-x-may
spec:
  diet: acacia
  location: SE
---
kind: Gorilla
metadata:
  labels:
    app: myApp
    movie: planetOfTheApes
  name: o-x-koko
spec:
  diet: bambooshoots
  location: SW
---
kind: Gorilla
metadata:
  labels:
    movie: planetOfTheApes
  name: o-ursus
spec:
  diet: heston
  location: Arizona
`)
}
