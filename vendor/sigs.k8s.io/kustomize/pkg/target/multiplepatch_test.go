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
	"strings"
	"testing"
<<<<<<< HEAD
)

func makeCommonFileForMultiplePatchTest(th *KustTestHarness) {
	th.writeK("/app/base", `
=======

	"sigs.k8s.io/kustomize/v3/pkg/kusttest"
)

func makeCommonFileForMultiplePatchTest(th *kusttest_test.KustTestHarness) {
	th.WriteK("/app/base", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namePrefix: team-foo-
commonLabels:
  app: mynginx
  org: example.com
  team: foo
commonAnnotations:
  note: This is a test annotation
resources:
  - deployment.yaml
  - service.yaml
configMapGenerator:
<<<<<<< HEAD
  - name: configmap-in-base
    literals:
      - foo=bar
`)
	th.writeF("/app/base/deployment.yaml", `
=======
- name: configmap-in-base
  literals:
  - foo=bar
`)
	th.WriteF("/app/base/deployment.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        volumeMounts:
        - name: nginx-persistent-storage
          mountPath: /tmp/ps
<<<<<<< HEAD
=======
      - name: sidecar
        image: sidecar:latest
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
      volumes:
      - name: nginx-persistent-storage
        emptyDir: {}
      - configMap:
          name: configmap-in-base
        name: configmap-in-base
`)
<<<<<<< HEAD
	th.writeF("/app/base/service.yaml", `
=======
	th.WriteF("/app/base/service.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: v1
kind: Service
metadata:
  name: nginx
  labels:
    app: nginx
spec:
  ports:
    - port: 80
  selector:
    app: nginx
`)
<<<<<<< HEAD
	th.writeK("/app/overlay/staging", `
=======
	th.WriteK("/app/overlay/staging", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namePrefix: staging-
commonLabels:
  env: staging
patchesStrategicMerge:
  - deployment-patch1.yaml
  - deployment-patch2.yaml
<<<<<<< HEAD
bases:
  - ../../base
configMapGenerator:
  - name: configmap-in-overlay
    literals:
      - hello=world
=======
resources:
  - ../../base
configMapGenerator:
- name: configmap-in-overlay
  literals:
  - hello=world
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
`)
}

func TestMultiplePatchesNoConflict(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/app/overlay/staging")
	makeCommonFileForMultiplePatchTest(th)
	th.writeF("/app/overlay/staging/deployment-patch1.yaml", `
=======
	th := kusttest_test.NewKustTestHarness(t, "/app/overlay/staging")
	makeCommonFileForMultiplePatchTest(th)
	th.WriteF("/app/overlay/staging/deployment-patch1.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        env:
        - name: ENVKEY
          value: ENVVALUE
      volumes:
      - name: nginx-persistent-storage
        emptyDir: null
        gcePersistentDisk:
          pdName: nginx-persistent-storage
      - configMap:
          name: configmap-in-overlay
        name: configmap-in-overlay
`)
<<<<<<< HEAD
	th.writeF("/app/overlay/staging/deployment-patch2.yaml", `
=======
	th.WriteF("/app/overlay/staging/deployment-patch2.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - name: nginx
        env:
        - name: ANOTHERENV
          value: FOO
<<<<<<< HEAD
      - name: sidecar
        image: sidecar
      volumes:
      - name: nginx-persistent-storage
`)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.assertActualEqualsExpected(m, `
apiVersion: v1
data:
  foo: bar
kind: ConfigMap
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: foo
  name: staging-team-foo-configmap-in-base-g7k6gt2889
---
apiVersion: v1
data:
  hello: world
kind: ConfigMap
metadata:
  labels:
    env: staging
  name: staging-configmap-in-overlay-k7cbc75tg8
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: foo
  name: staging-team-foo-nginx
spec:
  ports:
  - port: 80
  selector:
    app: mynginx
    env: staging
    org: example.com
    team: foo
---
=======
      volumes:
      - name: nginx-persistent-storage
`)
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: foo
  name: staging-team-foo-nginx
spec:
  selector:
    matchLabels:
      app: mynginx
      env: staging
      org: example.com
      team: foo
  template:
    metadata:
      annotations:
        note: This is a test annotation
      labels:
        app: mynginx
        env: staging
        org: example.com
        team: foo
    spec:
      containers:
      - env:
        - name: ANOTHERENV
          value: FOO
        - name: ENVKEY
          value: ENVVALUE
        image: nginx:latest
        name: nginx
        volumeMounts:
        - mountPath: /tmp/ps
          name: nginx-persistent-storage
<<<<<<< HEAD
      - image: sidecar
=======
      - image: sidecar:latest
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
        name: sidecar
      volumes:
      - gcePersistentDisk:
          pdName: nginx-persistent-storage
        name: nginx-persistent-storage
      - configMap:
          name: staging-configmap-in-overlay-k7cbc75tg8
        name: configmap-in-overlay
      - configMap:
          name: staging-team-foo-configmap-in-base-g7k6gt2889
        name: configmap-in-base
<<<<<<< HEAD
=======
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: foo
  name: staging-team-foo-nginx
spec:
  ports:
  - port: 80
  selector:
    app: mynginx
    env: staging
    org: example.com
    team: foo
---
apiVersion: v1
data:
  foo: bar
kind: ConfigMap
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: foo
  name: staging-team-foo-configmap-in-base-g7k6gt2889
---
apiVersion: v1
data:
  hello: world
kind: ConfigMap
metadata:
  labels:
    env: staging
  name: staging-configmap-in-overlay-k7cbc75tg8
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
`)
}

func TestMultiplePatchesWithConflict(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/app/overlay/staging")
	makeCommonFileForMultiplePatchTest(th)
	th.writeF("/app/overlay/staging/deployment-patch1.yaml", `
=======
	th := kusttest_test.NewKustTestHarness(t, "/app/overlay/staging")
	makeCommonFileForMultiplePatchTest(th)
	th.WriteF("/app/overlay/staging/deployment-patch1.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - name: nginx
        env:
        - name: ENABLE_FEATURE_FOO
          value: TRUE
      volumes:
      - name: nginx-persistent-storage
        emptyDir: null
        gcePersistentDisk:
          pdName: nginx-persistent-storage
      - configMap:
          name: configmap-in-overlay
        name: configmap-in-overlay
`)
<<<<<<< HEAD
	th.writeF("/app/overlay/staging/deployment-patch2.yaml", `
=======
	th.WriteF("/app/overlay/staging/deployment-patch2.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - name: nginx
        env:
        - name: ENABLE_FEATURE_FOO
          value: FALSE
`)
<<<<<<< HEAD
	_, err := th.makeKustTarget().MakeCustomizedResMap()
=======
	_, err := th.MakeKustTarget().MakeCustomizedResMap()
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	if err == nil {
		t.Fatalf("expected conflict")
	}
	if !strings.Contains(
		err.Error(), "conflict between ") {
		t.Fatalf("Unexpected err: %v", err)
	}
}
<<<<<<< HEAD
=======

func TestMultiplePatchesWithOnePatchDeleteDirective(t *testing.T) {
	additivePatch := `apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - name: nginx
        env:
        - name: SOME_NAME
          value: somevalue
`
	deletePatch := `apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - $patch: delete
        name: sidecar
`
	cases := []struct {
		name        string
		patch1      string
		patch2      string
		expectError bool
	}{
		{
			name:   "Patch with delete directive first",
			patch1: deletePatch,
			patch2: additivePatch,
		},
		{
			name:   "Patch with delete directive second",
			patch1: additivePatch,
			patch2: deletePatch,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			th := kusttest_test.NewKustTestHarness(t, "/app/overlay/staging")

			makeCommonFileForMultiplePatchTest(th)
			th.WriteF("/app/overlay/staging/deployment-patch1.yaml", c.patch1)
			th.WriteF("/app/overlay/staging/deployment-patch2.yaml", c.patch2)
			m, err := th.MakeKustTarget().MakeCustomizedResMap()
			if err != nil {
				t.Fatalf("Err: %v", err)
			}
			th.AssertActualEqualsExpected(m, `apiVersion: apps/v1beta2
kind: Deployment
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: foo
  name: staging-team-foo-nginx
spec:
  selector:
    matchLabels:
      app: mynginx
      env: staging
      org: example.com
      team: foo
  template:
    metadata:
      annotations:
        note: This is a test annotation
      labels:
        app: mynginx
        env: staging
        org: example.com
        team: foo
    spec:
      containers:
      - env:
        - name: SOME_NAME
          value: somevalue
        image: nginx
        name: nginx
        volumeMounts:
        - mountPath: /tmp/ps
          name: nginx-persistent-storage
      volumes:
      - emptyDir: {}
        name: nginx-persistent-storage
      - configMap:
          name: staging-team-foo-configmap-in-base-g7k6gt2889
        name: configmap-in-base
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: foo
  name: staging-team-foo-nginx
spec:
  ports:
  - port: 80
  selector:
    app: mynginx
    env: staging
    org: example.com
    team: foo
---
apiVersion: v1
data:
  foo: bar
kind: ConfigMap
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: foo
  name: staging-team-foo-configmap-in-base-g7k6gt2889
---
apiVersion: v1
data:
  hello: world
kind: ConfigMap
metadata:
  labels:
    env: staging
  name: staging-configmap-in-overlay-k7cbc75tg8
`)
		})
	}
}

func TestMultiplePatchesBothWithPatchDeleteDirective(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app/overlay/staging")
	makeCommonFileForMultiplePatchTest(th)
	th.WriteF("/app/overlay/staging/deployment-patch1.yaml", `
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - $patch: delete
        name: sidecar
`)
	th.WriteF("/app/overlay/staging/deployment-patch2.yaml", `
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
      containers:
      - $patch: delete
        name: nginx
`)
	_, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err == nil {
		t.Fatalf("Expected error")
	}
	if !strings.Contains(
		err.Error(), "both containing ") {
		t.Fatalf("Unexpected err: %v", err)
	}
	return
}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
