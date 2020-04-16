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
<<<<<<< HEAD
	"testing"
)

func TestSimpleBase(t *testing.T) {
	th := NewKustTestHarness(t, "/app/base")
	th.writeK("/app/base", `
=======
	"strings"
	"testing"

	"sigs.k8s.io/kustomize/v3/pkg/kusttest"
)

func TestSimpleBase(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app/base")
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
<<<<<<< HEAD
  - deployment.yaml
  - networkpolicy.yaml
  - service.yaml
`)
	th.writeF("/app/base/service.yaml", `
=======
  - service.yaml
  - deployment.yaml
  - networkpolicy.yaml
`)
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
	th.writeF("/app/base/networkpolicy.yaml", `
=======
	th.WriteF("/app/base/networkpolicy.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: nginx
spec:
  podSelector:
    matchExpressions:
      - {key: app, operator: In, values: [test]}
  ingress:
    - from:
        - podSelector:
            matchLabels:
              app: nginx
`)
<<<<<<< HEAD
	th.writeF("/app/base/deployment.yaml", `
=======
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
kind: Service
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    org: example.com
    team: foo
  name: team-foo-nginx
spec:
  ports:
  - port: 80
  selector:
    app: mynginx
    org: example.com
    team: foo
---
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    org: example.com
    team: foo
  name: team-foo-nginx
spec:
  selector:
    matchLabels:
      app: mynginx
      org: example.com
      team: foo
  template:
    metadata:
      annotations:
        note: This is a test annotation
      labels:
        app: mynginx
        org: example.com
        team: foo
    spec:
      containers:
      - image: nginx
        name: nginx
---
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    org: example.com
    team: foo
  name: team-foo-nginx
spec:
  ingress:
  - from:
    - podSelector:
        matchLabels:
          app: mynginx
          org: example.com
          team: foo
  podSelector:
    matchExpressions:
    - key: app
      operator: In
      values:
      - test
`)
}

<<<<<<< HEAD
func makeBaseWithGenerators(th *KustTestHarness) {
	th.writeK("/app", `
=======
func makeBaseWithGenerators(th *kusttest_test.KustTestHarness) {
	th.WriteK("/app", `
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
secretGenerator:
- name: secret-in-base
  literals:
    - username=admin
    - password=somepw
`)
	th.writeF("/app/deployment.yaml", `
=======
- name: configmap-in-base
  literals:
  - foo=bar
secretGenerator:
- name: secret-in-base
  literals:
  - username=admin
  - password=somepw
`)
	th.WriteF("/app/deployment.yaml", `
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
      volumes:
      - name: nginx-persistent-storage
        emptyDir: {}
      - configMap:
          name: configmap-in-base
        name: configmap-in-base
`)
<<<<<<< HEAD
	th.writeF("/app/service.yaml", `
=======
	th.WriteF("/app/service.yaml", `
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
}

func TestBaseWithGeneratorsAlone(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/app")
	makeBaseWithGenerators(th)
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
    org: example.com
    team: foo
  name: team-foo-configmap-in-base-bbdmdh7m8t
---
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    org: example.com
    team: foo
  name: team-foo-secret-in-base-tkm7hhtf8d
type: Opaque
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    org: example.com
    team: foo
  name: team-foo-nginx
spec:
  ports:
  - port: 80
  selector:
    app: mynginx
    org: example.com
    team: foo
---
=======
	th := kusttest_test.NewKustTestHarness(t, "/app")
	makeBaseWithGenerators(th)
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
    org: example.com
    team: foo
  name: team-foo-nginx
spec:
  selector:
    matchLabels:
      app: mynginx
      org: example.com
      team: foo
  template:
    metadata:
      annotations:
        note: This is a test annotation
      labels:
        app: mynginx
        org: example.com
        team: foo
    spec:
      containers:
      - image: nginx
        name: nginx
        volumeMounts:
        - mountPath: /tmp/ps
          name: nginx-persistent-storage
      volumes:
      - emptyDir: {}
        name: nginx-persistent-storage
      - configMap:
          name: team-foo-configmap-in-base-bbdmdh7m8t
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
    org: example.com
    team: foo
  name: team-foo-nginx
spec:
  ports:
  - port: 80
  selector:
    app: mynginx
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
    org: example.com
    team: foo
  name: team-foo-configmap-in-base-bbdmdh7m8t
---
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    org: example.com
    team: foo
  name: team-foo-secret-in-base-tkm7hhtf8d
type: Opaque
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
`)
}

func TestMergeAndReplaceGenerators(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/overlay")
	makeBaseWithGenerators(th)
	th.writeF("/overlay/deployment.yaml", `
=======
	th := kusttest_test.NewKustTestHarness(t, "/overlay")
	makeBaseWithGenerators(th)
	th.WriteF("/overlay/deployment.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: apps/v1beta2
kind: Deployment
metadata:
  name: nginx
spec:
  template:
    spec:
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
	th.writeK("/overlay", `
=======
	th.WriteK("/overlay", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namePrefix: staging-
commonLabels:
  env: staging
  team: override-foo
patchesStrategicMerge:
<<<<<<< HEAD
  - deployment.yaml
bases:
  - ../app
configMapGenerator:
  - name: configmap-in-overlay
    literals:
      - hello=world
  - name: configmap-in-base
    behavior: replace
    literals:
      - foo=override-bar
=======
- deployment.yaml
resources:
- ../app
configMapGenerator:
- name: configmap-in-overlay
  literals:
  - hello=world
- name: configmap-in-base
  behavior: replace
  literals:
  - foo=override-bar
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
secretGenerator:
- name: secret-in-base
  behavior: merge
  literals:
<<<<<<< HEAD
   - proxy=haproxy
`)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.assertActualEqualsExpected(m, `
apiVersion: v1
data:
  foo: override-bar
kind: ConfigMap
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: override-foo
  name: staging-team-foo-configmap-in-base-gh9d7t85gb
---
apiVersion: v1
data:
  hello: world
kind: ConfigMap
metadata:
  labels:
    env: staging
    team: override-foo
  name: staging-configmap-in-overlay-k7cbc75tg8
---
apiVersion: v1
data:
  password: c29tZXB3
  proxy: aGFwcm94eQ==
  username: YWRtaW4=
kind: Secret
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: override-foo
  name: staging-team-foo-secret-in-base-c8db7gk2m2
type: Opaque
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
    team: override-foo
  name: staging-team-foo-nginx
spec:
  ports:
  - port: 80
  selector:
    app: mynginx
    env: staging
    org: example.com
    team: override-foo
---
=======
  - proxy=haproxy
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
    team: override-foo
  name: staging-team-foo-nginx
spec:
  selector:
    matchLabels:
      app: mynginx
      env: staging
      org: example.com
      team: override-foo
  template:
    metadata:
      annotations:
        note: This is a test annotation
      labels:
        app: mynginx
        env: staging
        org: example.com
        team: override-foo
    spec:
      containers:
      - image: nginx
        name: nginx
        volumeMounts:
        - mountPath: /tmp/ps
          name: nginx-persistent-storage
      volumes:
      - gcePersistentDisk:
          pdName: nginx-persistent-storage
        name: nginx-persistent-storage
      - configMap:
          name: staging-configmap-in-overlay-k7cbc75tg8
        name: configmap-in-overlay
      - configMap:
          name: staging-team-foo-configmap-in-base-gh9d7t85gb
        name: configmap-in-base
<<<<<<< HEAD
`)
}
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
    team: override-foo
  name: staging-team-foo-nginx
spec:
  ports:
  - port: 80
  selector:
    app: mynginx
    env: staging
    org: example.com
    team: override-foo
---
apiVersion: v1
data:
  foo: override-bar
kind: ConfigMap
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: override-foo
  name: staging-team-foo-configmap-in-base-gh9d7t85gb
---
apiVersion: v1
data:
  password: c29tZXB3
  proxy: aGFwcm94eQ==
  username: YWRtaW4=
kind: Secret
metadata:
  annotations:
    note: This is a test annotation
  labels:
    app: mynginx
    env: staging
    org: example.com
    team: override-foo
  name: staging-team-foo-secret-in-base-c8db7gk2m2
type: Opaque
---
apiVersion: v1
data:
  hello: world
kind: ConfigMap
metadata:
  labels:
    env: staging
    team: override-foo
  name: staging-configmap-in-overlay-k7cbc75tg8
`)
}

func TestGeneratingIntoNamespaces(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app")
	th.WriteK("/app", `
configMapGenerator:
- name: test
  namespace: default
  literals:
    - key=value
- name: test
  namespace: kube-system
  literals:
    - key=value
secretGenerator:
- name: test
  namespace: default
  literals:
  - username=admin
  - password=somepw
- name: test
  namespace: kube-system
  literals:
  - username=admin
  - password=somepw
`)
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
data:
  key: value
kind: ConfigMap
metadata:
  name: test-t5t4md8fdm
  namespace: default
---
apiVersion: v1
data:
  key: value
kind: ConfigMap
metadata:
  name: test-t5t4md8fdm
  namespace: kube-system
---
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  name: test-h65t9hg6kc
  namespace: default
type: Opaque
---
apiVersion: v1
data:
  password: c29tZXB3
  username: YWRtaW4=
kind: Secret
metadata:
  name: test-h65t9hg6kc
  namespace: kube-system
type: Opaque
`)
}

// Valid that conflict is detected is the name are identical
// and namespace left to default
func TestConfigMapGeneratingIntoSameNamespace(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app")
	th.WriteK("/app", `
configMapGenerator:
- name: test
  namespace: default
  literals:
  - key=value
- name: test
  literals:
  - key=value
`)
	_, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "must merge or replace") {
		t.Fatalf("unexpected error %v", err)
	}
}

// Valid that conflict is detected is the name are identical
// and namespace left to default
func TestSecretGeneratingIntoSameNamespace(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app")
	th.WriteK("/app", `
secretGenerator:
- name: test
  namespace: default
  literals:
  - username=admin
  - password=somepw
- name: test
  literals:
  - username=admin
  - password=somepw
`)
	_, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "must merge or replace") {
		t.Fatalf("unexpected error %v", err)
	}
}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
