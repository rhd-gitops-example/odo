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

package target_test

import (
	"strings"
	"testing"
<<<<<<< HEAD
)

func writeCombinedOverlays(th *KustTestHarness) {
	// Base
	th.writeK("/app/base", `
resources:
- serviceaccount.yaml
- rolebinding.yaml
namePrefix: base-
nameSuffix: -suffix
`)
	th.writeF("/app/base/rolebinding.yaml", `
=======

	"sigs.k8s.io/kustomize/v3/pkg/kusttest"
)

func writeBase(th *kusttest_test.KustTestHarness) {
	th.WriteK("/app/base", `
resources:
- serviceaccount.yaml
- rolebinding.yaml
- clusterrolebinding.yaml
- clusterrole.yaml
namePrefix: pfx-
nameSuffix: -sfx
`)
	th.WriteF("/app/base/serviceaccount.yaml", `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: serviceaccount
`)
	th.WriteF("/app/base/rolebinding.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
<<<<<<< HEAD
  kind: Role
=======
  kind: ClusterRole
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
  name: role
subjects:
- kind: ServiceAccount
  name: serviceaccount
`)
<<<<<<< HEAD
	th.writeF("/app/base/serviceaccount.yaml", `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: serviceaccount
`)

	// Mid-level overlays
	th.writeK("/app/overlays/a", `
bases:
=======
	th.WriteF("/app/base/clusterrolebinding.yaml", `
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: rolebinding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: role
subjects:
- kind: ServiceAccount
  name: serviceaccount
`)
	th.WriteF("/app/base/clusterrole.yaml", `
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: role
rules:
- apiGroups: [""]
  resources: ["secrets"]
  verbs: ["get", "watch", "list"]
`)
}

func writeMidOverlays(th *kusttest_test.KustTestHarness) {
	// Mid-level overlays
	th.WriteK("/app/overlays/a", `
resources:
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
- ../../base
namePrefix: a-
nameSuffix: -suffixA
`)
<<<<<<< HEAD
	th.writeK("/app/overlays/b", `
bases:
=======
	th.WriteK("/app/overlays/b", `
resources:
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
- ../../base
namePrefix: b-
nameSuffix: -suffixB
`)
<<<<<<< HEAD

	// Top overlay, combining the mid-level overlays
	th.writeK("/app/combined", `
bases:
=======
}

func writeTopOverlay(th *kusttest_test.KustTestHarness) {
	// Top overlay, combining the mid-level overlays
	th.WriteK("/app/combined", `
resources:
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
- ../overlays/a
- ../overlays/b
`)
}

<<<<<<< HEAD
func TestMultibasesNoConflict(t *testing.T) {
	th := NewKustTestHarness(t, "/app/combined")
	writeCombinedOverlays(th)
	m, err := th.makeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	th.assertActualEqualsExpected(m, `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: a-base-serviceaccount-suffix-suffixA
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: b-base-serviceaccount-suffix-suffixB
=======
func TestBase(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app/base")
	writeBase(th)
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: pfx-serviceaccount-sfx
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: pfx-rolebinding-sfx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: pfx-role-sfx
subjects:
- kind: ServiceAccount
  name: pfx-serviceaccount-sfx
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: pfx-rolebinding-sfx
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: pfx-role-sfx
subjects:
- kind: ServiceAccount
  name: pfx-serviceaccount-sfx
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: pfx-role-sfx
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - watch
  - list
`)
}

func TestMidLevelA(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app/overlays/a")
	writeBase(th)
	writeMidOverlays(th)
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: a-pfx-serviceaccount-sfx-suffixA
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
<<<<<<< HEAD
  name: a-base-rolebinding-suffix-suffixA
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: role
subjects:
- kind: ServiceAccount
  name: a-base-serviceaccount-suffix-suffixA
=======
  name: a-pfx-rolebinding-sfx-suffixA
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: a-pfx-role-sfx-suffixA
subjects:
- kind: ServiceAccount
  name: a-pfx-serviceaccount-sfx-suffixA
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: a-pfx-rolebinding-sfx-suffixA
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: a-pfx-role-sfx-suffixA
subjects:
- kind: ServiceAccount
  name: a-pfx-serviceaccount-sfx-suffixA
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: a-pfx-role-sfx-suffixA
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - watch
  - list
`)
}

func TestMidLevelB(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app/overlays/b")
	writeBase(th)
	writeMidOverlays(th)
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: b-pfx-serviceaccount-sfx-suffixB
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
<<<<<<< HEAD
  name: b-base-rolebinding-suffix-suffixB
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: role
subjects:
- kind: ServiceAccount
  name: b-base-serviceaccount-suffix-suffixB
=======
  name: b-pfx-rolebinding-sfx-suffixB
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: b-pfx-role-sfx-suffixB
subjects:
- kind: ServiceAccount
  name: b-pfx-serviceaccount-sfx-suffixB
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: b-pfx-rolebinding-sfx-suffixB
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: b-pfx-role-sfx-suffixB
subjects:
- kind: ServiceAccount
  name: b-pfx-serviceaccount-sfx-suffixB
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: b-pfx-role-sfx-suffixB
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - watch
  - list
`)
}

func TestMultibasesNoConflict(t *testing.T) {
	th := kusttest_test.NewKustTestHarness(t, "/app/combined")
	writeBase(th)
	writeMidOverlays(th)
	writeTopOverlay(th)
	m, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}
	th.AssertActualEqualsExpected(m, `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: a-pfx-serviceaccount-sfx-suffixA
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: a-pfx-rolebinding-sfx-suffixA
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: a-pfx-role-sfx-suffixA
subjects:
- kind: ServiceAccount
  name: a-pfx-serviceaccount-sfx-suffixA
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: a-pfx-rolebinding-sfx-suffixA
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: a-pfx-role-sfx-suffixA
subjects:
- kind: ServiceAccount
  name: a-pfx-serviceaccount-sfx-suffixA
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: a-pfx-role-sfx-suffixA
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - watch
  - list
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: b-pfx-serviceaccount-sfx-suffixB
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: RoleBinding
metadata:
  name: b-pfx-rolebinding-sfx-suffixB
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: b-pfx-role-sfx-suffixB
subjects:
- kind: ServiceAccount
  name: b-pfx-serviceaccount-sfx-suffixB
---
apiVersion: rbac.authorization.k8s.io/v1beta1
kind: ClusterRoleBinding
metadata:
  name: b-pfx-rolebinding-sfx-suffixB
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: b-pfx-role-sfx-suffixB
subjects:
- kind: ServiceAccount
  name: b-pfx-serviceaccount-sfx-suffixB
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: b-pfx-role-sfx-suffixB
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - get
  - watch
  - list
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
`)
}

func TestMultibasesWithConflict(t *testing.T) {
<<<<<<< HEAD
	th := NewKustTestHarness(t, "/app/combined")
	writeCombinedOverlays(th)

	th.writeK("/app/overlays/a", `
bases:
- ../../base
=======
	th := kusttest_test.NewKustTestHarness(t, "/app/combined")
	writeBase(th)
	writeMidOverlays(th)
	writeTopOverlay(th)

	th.WriteK("/app/overlays/a", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
namePrefix: a-
nameSuffix: -suffixA
resources:
- serviceaccount.yaml
<<<<<<< HEAD
`)
	// Expect an error because this resource in the overlay
	// matches a resource in the base.
	th.writeF("/app/overlays/a/serviceaccount.yaml", `
=======
- ../../base
`)
	// Expect an error because this resource in the overlay
	// matches a resource in the base.
	th.WriteF("/app/overlays/a/serviceaccount.yaml", `
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
apiVersion: v1
kind: ServiceAccount
metadata:
  name: serviceaccount
`)

<<<<<<< HEAD
	_, err := th.makeKustTarget().MakeCustomizedResMap()
	if err == nil {
		t.Fatalf("Expected resource conflict.")
	}
	if !strings.Contains(
		err.Error(), "Multiple matches for name ~G_v1_ServiceAccount") {
		t.Fatalf("Unexpected err: %v", err)
=======
	_, err := th.MakeKustTarget().MakeCustomizedResMap()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "multiple matches for ~G_v1_ServiceAccount|~X|serviceaccount") {
		t.Fatalf("unexpected error %v", err)
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
}
