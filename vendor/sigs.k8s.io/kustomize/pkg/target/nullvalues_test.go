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

func TestNullValues(t *testing.T) {
	th := NewKustTestHarness(t, "/app")
	th.writeF("/app/deployment.yaml", `
=======

	"sigs.k8s.io/kustomize/v3/pkg/kusttest"
)

func TestNullValues(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app")
	th.WriteF("/app/deployment.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: example
  name: example
spec:
  selector:
    matchLabels:
      app: example
  template:
    metadata:
      labels:
        app: example
    spec:
      containers:
      - args: null
        image: image
        name: example
`)
<<<<<<< HEAD
	th.writeF("/app/kustomization.yaml", `
=======
	th.WriteF("/app/kustomization.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
- deployment.yaml
`)
<<<<<<< HEAD
	m, err := th.makeKustTarget().MakeCustomizedResMap()
=======
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

<<<<<<< HEAD
	th.assertActualEqualsExpected(m, `
=======
	th.AssertActualEqualsExpected(m, `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: example
  name: example
spec:
  selector:
    matchLabels:
      app: example
  template:
    metadata:
      labels:
        app: example
    spec:
      containers:
      - args: null
        image: image
        name: example
`)
}
