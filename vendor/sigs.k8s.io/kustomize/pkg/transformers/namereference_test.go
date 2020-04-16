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

package transformers

import (
<<<<<<< HEAD
	"reflect"
	"strings"
	"testing"

	"sigs.k8s.io/kustomize/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/pkg/resid"
	"sigs.k8s.io/kustomize/pkg/resmap"
	"sigs.k8s.io/kustomize/pkg/resource"
=======
	"strings"
	"testing"

	"sigs.k8s.io/kustomize/v3/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/v3/pkg/gvk"
	"sigs.k8s.io/kustomize/v3/pkg/resid"
	"sigs.k8s.io/kustomize/v3/pkg/resmap"
	"sigs.k8s.io/kustomize/v3/pkg/resmaptest"
	"sigs.k8s.io/kustomize/v3/pkg/resource"
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
)

func TestNameReferenceHappyRun(t *testing.T) {
	rf := resource.NewFactory(
		kunstruct.NewKunstructuredFactoryImpl())
<<<<<<< HEAD
	m := resmap.ResMap{
		resid.NewResId(cmap, "cm1"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "someprefix-cm1-somehash",
				},
			}),
		resid.NewResId(cmap, "cm2"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name": "someprefix-cm2-somehash",
				},
			}),
		resid.NewResId(secret, "secret1"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "Secret",
				"metadata": map[string]interface{}{
					"name": "someprefix-secret1-somehash",
				},
			}),
		resid.NewResId(pvc, "claim1"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "PersistentVolumeClaim",
				"metadata": map[string]interface{}{
					"name": "someprefix-claim1",
				},
			}),
		resid.NewResId(ingress, "ingress1"): rf.FromMap(
			map[string]interface{}{
				"group":      "extensions",
				"apiVersion": "v1beta1",
				"kind":       "Ingress",
				"metadata": map[string]interface{}{
					"name": "ingress1",
					"annotations": map[string]interface{}{
						"ingress.kubernetes.io/auth-secret":       "secret1",
						"nginx.ingress.kubernetes.io/auth-secret": "secret1",
					},
				},
				"spec": map[string]interface{}{
					"backend": map[string]interface{}{
						"serviceName": "testsvc",
						"servicePort": "80",
					},
				},
			},
		),
		resid.NewResId(deploy, "deploy1"): rf.FromMap(
			map[string]interface{}{
				"group":      "apps",
				"apiVersion": "v1",
				"kind":       "Deployment",
				"metadata": map[string]interface{}{
					"name": "deploy1",
				},
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name":  "nginx",
									"image": "nginx:1.7.9",
									"env": []interface{}{
										map[string]interface{}{
											"name": "CM_FOO",
											"valueFrom": map[string]interface{}{
												"configMapKeyRef": map[string]interface{}{
													"name": "cm1",
													"key":  "somekey",
												},
											},
										},
										map[string]interface{}{
											"name": "SECRET_FOO",
											"valueFrom": map[string]interface{}{
												"secretKeyRef": map[string]interface{}{
													"name": "secret1",
													"key":  "somekey",
												},
											},
										},
									},
									"envFrom": []interface{}{
										map[string]interface{}{
											"configMapRef": map[string]interface{}{
=======
	m := resmaptest_test.NewRmBuilder(t, rf).AddWithName(
		"cm1",
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": "someprefix-cm1-somehash",
			},
		}).AddWithName(
		"cm2",
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": "someprefix-cm2-somehash",
			},
		}).AddWithName(
		"secret1",
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name": "someprefix-secret1-somehash",
			},
		}).AddWithName(
		"claim1",
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolumeClaim",
			"metadata": map[string]interface{}{
				"name": "someprefix-claim1",
			},
		}).Add(
		map[string]interface{}{
			"group":      "networking.k8s.io",
			"apiVersion": "v1beta1",
			"kind":       "Ingress",
			"metadata": map[string]interface{}{
				"name": "ingress1",
				"annotations": map[string]interface{}{
					"ingress.kubernetes.io/auth-secret":           "secret1",
					"nginx.ingress.kubernetes.io/auth-secret":     "secret1",
					"nginx.ingress.kubernetes.io/auth-tls-secret": "secret1",
				},
			},
			"spec": map[string]interface{}{
				"backend": map[string]interface{}{
					"serviceName": "testsvc",
					"servicePort": "80",
				},
			},
		},
	).Add(
		map[string]interface{}{
			"group":      "apps",
			"apiVersion": "v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:1.7.9",
								"env": []interface{}{
									map[string]interface{}{
										"name": "CM_FOO",
										"valueFrom": map[string]interface{}{
											"configMapKeyRef": map[string]interface{}{
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
												"name": "cm1",
												"key":  "somekey",
											},
										},
<<<<<<< HEAD
										map[string]interface{}{
											"secretRef": map[string]interface{}{
=======
									},
									map[string]interface{}{
										"name": "SECRET_FOO",
										"valueFrom": map[string]interface{}{
											"secretKeyRef": map[string]interface{}{
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
												"name": "secret1",
												"key":  "somekey",
											},
										},
									},
								},
<<<<<<< HEAD
							},
							"imagePullSecrets": []interface{}{
								map[string]interface{}{
									"name": "secret1",
								},
							},
							"volumes": map[string]interface{}{
								"configMap": map[string]interface{}{
									"name": "cm1",
								},
								"projected": map[string]interface{}{
									"sources": map[string]interface{}{
										"configMap": map[string]interface{}{
											"name": "cm2",
										},
										"secret": map[string]interface{}{
											"name": "secret1",
										},
									},
								},
								"secret": map[string]interface{}{
									"secretName": "secret1",
								},
								"persistentVolumeClaim": map[string]interface{}{
									"claimName": "claim1",
								},
							},
						},
					},
				},
			}),
		resid.NewResId(statefulset, "statefulset1"): rf.FromMap(
			map[string]interface{}{
				"group":      "apps",
				"apiVersion": "v1",
				"kind":       "StatefulSet",
				"metadata": map[string]interface{}{
					"name": "statefulset1",
				},
				"spec": map[string]interface{}{
					"template": map[string]interface{}{
						"spec": map[string]interface{}{
							"containers": []interface{}{
								map[string]interface{}{
									"name":  "nginx",
									"image": "nginx:1.7.9",
								},
							},
							"volumes": map[string]interface{}{
								"projected": map[string]interface{}{
									"sources": map[string]interface{}{
										"configMap": map[string]interface{}{
											"name": "cm2",
										},
										"secret": map[string]interface{}{
											"name": "secret1",
										},
=======
								"envFrom": []interface{}{
									map[string]interface{}{
										"configMapRef": map[string]interface{}{
											"name": "cm1",
											"key":  "somekey",
										},
									},
									map[string]interface{}{
										"secretRef": map[string]interface{}{
											"name": "secret1",
											"key":  "somekey",
										},
									},
								},
							},
						},
						"imagePullSecrets": []interface{}{
							map[string]interface{}{
								"name": "secret1",
							},
						},
						"volumes": map[string]interface{}{
							"configMap": map[string]interface{}{
								"name": "cm1",
							},
							"projected": map[string]interface{}{
								"sources": map[string]interface{}{
									"configMap": map[string]interface{}{
										"name": "cm2",
									},
									"secret": map[string]interface{}{
										"name": "secret1",
									},
								},
							},
							"secret": map[string]interface{}{
								"secretName": "secret1",
							},
							"persistentVolumeClaim": map[string]interface{}{
								"claimName": "claim1",
							},
						},
					},
				},
			},
		}).Add(
		map[string]interface{}{
			"group":      "apps",
			"apiVersion": "v1",
			"kind":       "StatefulSet",
			"metadata": map[string]interface{}{
				"name": "statefulset1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:1.7.9",
							},
						},
						"volumes": map[string]interface{}{
							"projected": map[string]interface{}{
								"sources": map[string]interface{}{
									"configMap": map[string]interface{}{
										"name": "cm2",
									},
									"secret": map[string]interface{}{
										"name": "secret1",
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
									},
								},
							},
						},
					},
				},
<<<<<<< HEAD
			}),
		resid.NewResIdWithPrefixNamespace(sa, "sa", "", "test"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "ServiceAccount",
				"metadata": map[string]interface{}{
					"name":      "someprefix-sa",
					"namespace": "test",
				},
			}),
		resid.NewResId(crb, "crb"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "ClusterRoleBinding",
				"metadata": map[string]interface{}{
					"name": "crb",
				},
				"subjects": []interface{}{
					map[string]interface{}{
						"kind":      "ServiceAccount",
						"name":      "sa",
						"namespace": "test",
					},
				},
			}),
		resid.NewResId(cr, "cr"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "ClusterRole",
				"metadata": map[string]interface{}{
					"name": "cr",
				},
				"rules": []interface{}{
					map[string]interface{}{
						"resources": []interface{}{
							"secrets",
						},
						"resourceNames": []interface{}{
							"secret1",
							"secret1",
							"secret2",
						},
					},
				},
			}),
	}

	expected := resmap.ResMap{}
	for k, v := range m {
		expected[k] = v
	}

	expected[resid.NewResId(deploy, "deploy1")] = rf.FromMap(
=======
			},
		}).AddWithName("sa",
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name":      "someprefix-sa",
				"namespace": "test",
			},
		}).Add(
		map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRoleBinding",
			"metadata": map[string]interface{}{
				"name": "crb",
			},
			"subjects": []interface{}{
				map[string]interface{}{
					"kind":      "ServiceAccount",
					"name":      "sa",
					"namespace": "test",
				},
			},
		}).Add(
		map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRole",
			"metadata": map[string]interface{}{
				"name": "cr",
			},
			"rules": []interface{}{
				map[string]interface{}{
					"resources": []interface{}{
						"secrets",
					},
					"resourceNames": []interface{}{
						"secret1",
						"secret1",
						"secret2",
					},
				},
			},
		}).Add(
		map[string]interface{}{
			"apiVersion": "batch/v1beta1",
			"kind":       "CronJob",
			"metadata": map[string]interface{}{
				"name": "cronjob1",
			},
			"spec": map[string]interface{}{
				"schedule": "0 14 * * *",
				"jobTemplate": map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "main",
										"image": "myimage",
									},
								},
								"volumes": map[string]interface{}{
									"projected": map[string]interface{}{
										"sources": map[string]interface{}{
											"configMap": map[string]interface{}{
												"name": "cm2",
											},
											"secret": map[string]interface{}{
												"name": "secret1",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}).ResMap()

	expected := resmaptest_test.NewSeededRmBuilder(t, rf, m.ShallowCopy()).ReplaceResource(
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		map[string]interface{}{
			"group":      "apps",
			"apiVersion": "v1",
			"kind":       "Deployment",
			"metadata": map[string]interface{}{
				"name": "deploy1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:1.7.9",
								"env": []interface{}{
									map[string]interface{}{
										"name": "CM_FOO",
										"valueFrom": map[string]interface{}{
											"configMapKeyRef": map[string]interface{}{
												"name": "someprefix-cm1-somehash",
												"key":  "somekey",
											},
										},
									},
									map[string]interface{}{
										"name": "SECRET_FOO",
										"valueFrom": map[string]interface{}{
											"secretKeyRef": map[string]interface{}{
												"name": "someprefix-secret1-somehash",
												"key":  "somekey",
											},
										},
									},
								},
								"envFrom": []interface{}{
									map[string]interface{}{
										"configMapRef": map[string]interface{}{
											"name": "someprefix-cm1-somehash",
											"key":  "somekey",
										},
									},
									map[string]interface{}{
										"secretRef": map[string]interface{}{
											"name": "someprefix-secret1-somehash",
											"key":  "somekey",
										},
									},
								},
							},
						},
						"imagePullSecrets": []interface{}{
							map[string]interface{}{
								"name": "someprefix-secret1-somehash",
							},
						},
						"volumes": map[string]interface{}{
							"configMap": map[string]interface{}{
								"name": "someprefix-cm1-somehash",
							},
							"projected": map[string]interface{}{
								"sources": map[string]interface{}{
									"configMap": map[string]interface{}{
										"name": "someprefix-cm2-somehash",
									},
									"secret": map[string]interface{}{
										"name": "someprefix-secret1-somehash",
									},
								},
							},
							"secret": map[string]interface{}{
								"secretName": "someprefix-secret1-somehash",
							},
							"persistentVolumeClaim": map[string]interface{}{
								"claimName": "someprefix-claim1",
							},
						},
					},
				},
			},
<<<<<<< HEAD
		})
	expected[resid.NewResId(statefulset, "statefulset1")] = rf.FromMap(
=======
		}).ReplaceResource(
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		map[string]interface{}{
			"group":      "apps",
			"apiVersion": "v1",
			"kind":       "StatefulSet",
			"metadata": map[string]interface{}{
				"name": "statefulset1",
			},
			"spec": map[string]interface{}{
				"template": map[string]interface{}{
					"spec": map[string]interface{}{
						"containers": []interface{}{
							map[string]interface{}{
								"name":  "nginx",
								"image": "nginx:1.7.9",
							},
						},
						"volumes": map[string]interface{}{
							"projected": map[string]interface{}{
								"sources": map[string]interface{}{
									"configMap": map[string]interface{}{
										"name": "someprefix-cm2-somehash",
									},
									"secret": map[string]interface{}{
										"name": "someprefix-secret1-somehash",
									},
								},
							},
						},
					},
				},
			},
<<<<<<< HEAD
		})
	expected[resid.NewResId(ingress, "ingress1")] = rf.FromMap(
		map[string]interface{}{
			"group":      "extensions",
=======
		}).ReplaceResource(
		map[string]interface{}{
			"group":      "networking.k8s.io",
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			"apiVersion": "v1beta1",
			"kind":       "Ingress",
			"metadata": map[string]interface{}{
				"name": "ingress1",
				"annotations": map[string]interface{}{
<<<<<<< HEAD
					"ingress.kubernetes.io/auth-secret":       "someprefix-secret1-somehash",
					"nginx.ingress.kubernetes.io/auth-secret": "someprefix-secret1-somehash",
=======
					"ingress.kubernetes.io/auth-secret":           "someprefix-secret1-somehash",
					"nginx.ingress.kubernetes.io/auth-secret":     "someprefix-secret1-somehash",
					"nginx.ingress.kubernetes.io/auth-tls-secret": "someprefix-secret1-somehash",
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
				},
			},
			"spec": map[string]interface{}{
				"backend": map[string]interface{}{
					"serviceName": "testsvc",
					"servicePort": "80",
				},
			},
<<<<<<< HEAD
		},
	)
	expected[resid.NewResId(crb, "crb")] = rf.FromMap(
=======
		}).ReplaceResource(
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRoleBinding",
			"metadata": map[string]interface{}{
				"name": "crb",
			},
			"subjects": []interface{}{
				map[string]interface{}{
					"kind":      "ServiceAccount",
					"name":      "someprefix-sa",
					"namespace": "test",
				},
			},
<<<<<<< HEAD
		})
	expected[resid.NewResId(cr, "cr")] = rf.FromMap(
=======
		}).ReplaceResource(
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRole",
			"metadata": map[string]interface{}{
				"name": "cr",
			},
			"rules": []interface{}{
				map[string]interface{}{
					"resources": []interface{}{
						"secrets",
					},
					"resourceNames": []interface{}{
						"someprefix-secret1-somehash",
						"someprefix-secret1-somehash",
						"secret2",
					},
				},
			},
<<<<<<< HEAD
		})
=======
		}).ReplaceResource(
		map[string]interface{}{
			"apiVersion": "batch/v1beta1",
			"kind":       "CronJob",
			"metadata": map[string]interface{}{
				"name": "cronjob1",
			},
			"spec": map[string]interface{}{
				"schedule": "0 14 * * *",
				"jobTemplate": map[string]interface{}{
					"spec": map[string]interface{}{
						"template": map[string]interface{}{
							"spec": map[string]interface{}{
								"containers": []interface{}{
									map[string]interface{}{
										"name":  "main",
										"image": "myimage",
									},
								},
								"volumes": map[string]interface{}{
									"projected": map[string]interface{}{
										"sources": map[string]interface{}{
											"configMap": map[string]interface{}{
												"name": "someprefix-cm2-somehash",
											},
											"secret": map[string]interface{}{
												"name": "someprefix-secret1-somehash",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		}).ResMap()

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	nrt := NewNameReferenceTransformer(defaultTransformerConfig.NameReference)
	err := nrt.Transform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
<<<<<<< HEAD
	if !reflect.DeepEqual(m, expected) {
		err = expected.ErrorIfNotEqual(m)
=======

	if err = expected.ErrorIfNotEqualLists(m); err != nil {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

func TestNameReferenceUnhappyRun(t *testing.T) {
	rf := resource.NewFactory(
		kunstruct.NewKunstructuredFactoryImpl())
	tests := []struct {
		resMap      resmap.ResMap
		expectedErr string
	}{
		{
<<<<<<< HEAD
			resMap: resmap.ResMap{
				resid.NewResId(cr, "cr"): rf.FromMap(
					map[string]interface{}{
						"apiVersion": "rbac.authorization.k8s.io/v1",
						"kind":       "ClusterRole",
						"metadata": map[string]interface{}{
							"name": "cr",
						},
						"rules": []interface{}{
							map[string]interface{}{
								"resources": []interface{}{
									"secrets",
								},
								"resourceNames": []interface{}{
									[]interface{}{},
								},
							},
						},
					}),
			},
			expectedErr: "is expected to be string"},
		{resMap: resmap.ResMap{
			resid.NewResId(cr, "cr"): rf.FromMap(
=======
			resMap: resmaptest_test.NewRmBuilder(t, rf).Add(
				map[string]interface{}{
					"apiVersion": "rbac.authorization.k8s.io/v1",
					"kind":       "ClusterRole",
					"metadata": map[string]interface{}{
						"name": "cr",
					},
					"rules": []interface{}{
						map[string]interface{}{
							"resources": []interface{}{
								"secrets",
							},
							"resourceNames": []interface{}{
								[]interface{}{},
							},
						},
					},
				}).ResMap(),
			expectedErr: "is expected to be"},
		{
			resMap: resmaptest_test.NewRmBuilder(t, rf).Add(
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
				map[string]interface{}{
					"apiVersion": "rbac.authorization.k8s.io/v1",
					"kind":       "ClusterRole",
					"metadata": map[string]interface{}{
						"name": "cr",
					},
					"rules": []interface{}{
						map[string]interface{}{
							"resources": []interface{}{
								"secrets",
							},
							"resourceNames": map[string]interface{}{
								"foo": "bar",
							},
						},
					},
<<<<<<< HEAD
				}),
		},
			expectedErr: "is expected to be either a string or a []interface{}"},
=======
				}).ResMap(),
			expectedErr: "is expected to contain a name field"},
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}

	nrt := NewNameReferenceTransformer(defaultTransformerConfig.NameReference)
	for _, test := range tests {
		err := nrt.Transform(test.resMap)
		if err == nil {
			t.Fatalf("expected error to happen")
		}

		if !strings.Contains(err.Error(), test.expectedErr) {
			t.Fatalf("Incorrect error.\nExpected: %s, but got %v",
				test.expectedErr, err)
		}
	}
}

func TestNameReferencePersistentVolumeHappyRun(t *testing.T) {
	rf := resource.NewFactory(
		kunstruct.NewKunstructuredFactoryImpl())
<<<<<<< HEAD
	m := resmap.ResMap{
		resid.NewResId(pv, "volume1"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "PersistentVolume",
				"metadata": map[string]interface{}{
					"name": "someprefix-volume1",
				},
			}),

		resid.NewResId(pvc, "claim1"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "PersistentVolumeClaim",
				"metadata": map[string]interface{}{
					"name":      "someprefix-claim1",
					"namespace": "some-namespace",
				},
				"spec": map[string]interface{}{
					"volumeName": "volume1",
				},
			}),
	}

	expected := resmap.ResMap{
		resid.NewResId(pv, "volume1"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "PersistentVolume",
				"metadata": map[string]interface{}{
					"name": "someprefix-volume1",
				},
			}),

		resid.NewResId(pvc, "claim1"): rf.FromMap(
			map[string]interface{}{
				"apiVersion": "v1",
				"kind":       "PersistentVolumeClaim",
				"metadata": map[string]interface{}{
					"name":      "someprefix-claim1",
					"namespace": "some-namespace",
				},
				"spec": map[string]interface{}{
					"volumeName": "someprefix-volume1",
				},
			}),
	}
=======

	v1 := rf.FromMapWithName(
		"volume1",
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolume",
			"metadata": map[string]interface{}{
				"name": "someprefix-volume1",
			},
		})
	c1 := rf.FromMapWithName(
		"claim1",
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolumeClaim",
			"metadata": map[string]interface{}{
				"name":      "someprefix-claim1",
				"namespace": "some-namespace",
			},
			"spec": map[string]interface{}{
				"volumeName": "volume1",
			},
		})

	v2 := v1.DeepCopy()
	c2 := rf.FromMapWithName(
		"claim1",
		map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolumeClaim",
			"metadata": map[string]interface{}{
				"name":      "someprefix-claim1",
				"namespace": "some-namespace",
			},
			"spec": map[string]interface{}{
				"volumeName": "someprefix-volume1",
			},
		})

	m1 := resmaptest_test.NewRmBuilder(t, rf).AddR(v1).AddR(c1).ResMap()

	nrt := NewNameReferenceTransformer(defaultTransformerConfig.NameReference)
	if err := nrt.Transform(m1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	m2 := resmaptest_test.NewRmBuilder(t, rf).AddR(v2).AddR(c2).ResMap()
	v2.AppendRefBy(c2.CurId())

	if err := m1.ErrorIfNotEqualLists(m2); err != nil {
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

// utility map to create a deployment object
// with (metadatanamespace, metadataname) as key
// and pointing to "refname" secret and configmap
func deploymentMap(metadatanamespace string, metadataname string,
	configmapref string, secretref string) map[string]interface{} {
	deployment := map[string]interface{}{
		"group":      "apps",
		"apiVersion": "v1",
		"kind":       "Deployment",
		"metadata": map[string]interface{}{
			"name": metadataname,
		},
		"spec": map[string]interface{}{
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{
							"name":  "nginx",
							"image": "nginx:1.7.9",
							"env": []interface{}{
								map[string]interface{}{
									"name": "CM_FOO",
									"valueFrom": map[string]interface{}{
										"configMapKeyRef": map[string]interface{}{
											"name": configmapref,
											"key":  "somekey",
										},
									},
								},
								map[string]interface{}{
									"name": "SECRET_FOO",
									"valueFrom": map[string]interface{}{
										"secretKeyRef": map[string]interface{}{
											"name": secretref,
											"key":  "somekey",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	if metadatanamespace != "" {
		metadata := deployment["metadata"].(map[string]interface{})
		metadata["namespace"] = metadatanamespace
	}
	return deployment
}

const (
	defaultNs = "default"
	ns1       = "ns1"
	ns2       = "ns2"
	ns3       = "ns3"
	ns4       = "ns4"

	orgname      = "uniquename"
	prefixedname = "prefix-uniquename"
	suffixedname = "uniquename-suffix"
	modifiedname = "modifiedname"
)

// TestNameReferenceNamespace creates serviceAccount and clusterRoleBinding
// object with the same original names (uniquename) in different namespaces
// and with different current Id.
func TestNameReferenceNamespace(t *testing.T) {
	rf := resource.NewFactory(
		kunstruct.NewKunstructuredFactoryImpl())
	m := resmaptest_test.NewRmBuilder(t, rf).
		// Add ConfigMap with the same org name in noNs, "ns1" and "ns2" namespaces
		AddWithName(orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": modifiedname,
			}}).
		AddWithNsAndName(ns1, orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      prefixedname,
				"namespace": ns1,
			}}).
		AddWithNsAndName(ns2, orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name":      suffixedname,
				"namespace": ns2,
			}}).
		// Add Secret with the same org name in noNs, "ns1" and "ns2" namespaces
		AddWithNsAndName(defaultNs, orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      modifiedname,
				"namespace": defaultNs,
			}}).
		AddWithNsAndName(ns1, orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      prefixedname,
				"namespace": ns1,
			}}).
		AddWithNsAndName(ns2, orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      suffixedname,
				"namespace": ns2,
			}}).
		// Add Deployment with the same org name in noNs, "ns1" and "ns2" namespaces
		AddWithNsAndName(defaultNs, orgname, deploymentMap(defaultNs, modifiedname, modifiedname, modifiedname)).
		AddWithNsAndName(ns1, orgname, deploymentMap(ns1, prefixedname, orgname, orgname)).
		AddWithNsAndName(ns2, orgname, deploymentMap(ns2, suffixedname, orgname, orgname)).ResMap()

	expected := resmaptest_test.NewSeededRmBuilder(t, rf, m.ShallowCopy()).
		ReplaceResource(deploymentMap(defaultNs, modifiedname, modifiedname, modifiedname)).
		ReplaceResource(deploymentMap(ns1, prefixedname, prefixedname, prefixedname)).
		ReplaceResource(deploymentMap(ns2, suffixedname, suffixedname, suffixedname)).ResMap()

	nrt := NewNameReferenceTransformer(defaultTransformerConfig.NameReference)
	err := nrt.Transform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err = expected.ErrorIfNotEqualLists(m); err != nil {
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

// TestNameReferenceNamespace creates serviceAccount and clusterRoleBinding
// object with the same original names (uniquename) in different namespaces
// and with different current Id.
func TestNameReferenceClusterWide(t *testing.T) {
	rf := resource.NewFactory(
		kunstruct.NewKunstructuredFactoryImpl())
	m := resmaptest_test.NewRmBuilder(t, rf).
		// Add ServiceAccount with the same org name in noNs, "ns1" and "ns2" namespaces
		AddWithName(orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name": modifiedname,
			}}).
		AddWithNsAndName(ns1, orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name":      prefixedname,
				"namespace": ns1,
			}}).
		AddWithNsAndName(ns2, orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name":      suffixedname,
				"namespace": ns2,
			}}).
		// Add a PersistentVolume to have a clusterwide resource
		AddWithName(orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "PersistentVolume",
			"metadata": map[string]interface{}{
				"name": modifiedname,
			}}).
		AddWithName(orgname, map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRole",
			"metadata": map[string]interface{}{
				"name": modifiedname,
			},
			"rules": []interface{}{
				map[string]interface{}{
					"resources": []interface{}{
						"persistentvolumes",
					},
					"resourceNames": []interface{}{
						orgname,
					},
				},
			}}).
		AddWithName(orgname, map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRoleBinding",
			"metadata": map[string]interface{}{
				"name": modifiedname,
			},
			"roleRef": map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "ClusterRole",
				"name":       orgname,
			},
			"subjects": []interface{}{
				map[string]interface{}{
					"kind":      "ServiceAccount",
					"name":      orgname,
					"namespace": defaultNs,
				},
				map[string]interface{}{
					"kind":      "ServiceAccount",
					"name":      orgname,
					"namespace": ns1,
				},
				map[string]interface{}{
					"kind":      "ServiceAccount",
					"name":      orgname,
					"namespace": ns2,
				},
				map[string]interface{}{
					"kind":      "ServiceAccount",
					"name":      orgname,
					"namespace": "random",
				},
			}}).ResMap()

	expected := resmaptest_test.NewSeededRmBuilder(t, rf, m.ShallowCopy()).
		ReplaceResource(
			map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "ClusterRole",
				"metadata": map[string]interface{}{
					"name": modifiedname,
				},
				// Behavior of the transformer is still imperfect
				// It should use the (resources,apigroup,resourceNames) as
				// combination to select the candidates.
				"rules": []interface{}{
					map[string]interface{}{
						"resources": []interface{}{
							"persistentvolumes",
						},
						"resourceNames": []interface{}{
							modifiedname,
						},
					},
				}}).
		ReplaceResource(
			map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "ClusterRoleBinding",
				"metadata": map[string]interface{}{
					"name": modifiedname,
				},
				"roleRef": map[string]interface{}{
					"apiVersion": "rbac.authorization.k8s.io/v1",
					"kind":       "ClusterRole",
					"name":       modifiedname,
				},
				// The following tests required a change in
				// getNameFunc implementation in order to leverage
				// the namespace field.
				"subjects": []interface{}{
					map[string]interface{}{
						"kind":      "ServiceAccount",
						"name":      modifiedname,
						"namespace": defaultNs,
					},
					map[string]interface{}{
						"kind":      "ServiceAccount",
						"name":      prefixedname,
						"namespace": ns1,
					},
					map[string]interface{}{
						"kind":      "ServiceAccount",
						"name":      suffixedname,
						"namespace": ns2,
					},
					map[string]interface{}{
						"kind":      "ServiceAccount",
						"name":      orgname,
						"namespace": "random",
					},
				},
			}).ResMap()

	clusterRoleId := resid.NewResId(
		gvk.Gvk{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"}, modifiedname)
	clusterRoleBindingId := resid.NewResId(
		gvk.Gvk{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"}, modifiedname)
	clusterRole, _ := expected.GetByCurrentId(clusterRoleId)
	clusterRole.AppendRefBy(clusterRoleBindingId)

	nrt := NewNameReferenceTransformer(defaultTransformerConfig.NameReference)
	err := nrt.Transform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err = expected.ErrorIfNotEqualLists(m); err != nil {
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

// TestNameReferenceNamespaceTransformation creates serviceAccount and clusterRoleBinding
// object with the same original names (uniquename) in different namespaces
// and with different current Id.
func TestNameReferenceNamespaceTransformation(t *testing.T) {
	rf := resource.NewFactory(
		kunstruct.NewKunstructuredFactoryImpl())
	m := resmaptest_test.NewRmBuilder(t, rf).
		AddWithNsAndName(ns4, orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      orgname,
				"namespace": ns4,
			}}).
		// Add ServiceAccount with the same org name in "ns1" namespaces
		AddWithNsAndName(ns1, orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name":      prefixedname,
				"namespace": ns1,
			}}).
		// Simulate NamespaceTransformer effect (ns3 transformed in ns2)
		AddWithNsAndName(ns3, orgname, map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ServiceAccount",
			"metadata": map[string]interface{}{
				"name":      suffixedname,
				"namespace": ns2,
			}}).
		AddWithName(orgname, map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRole",
			"metadata": map[string]interface{}{
				"name": modifiedname,
			}}).
		AddWithName(orgname, map[string]interface{}{
			"apiVersion": "rbac.authorization.k8s.io/v1",
			"kind":       "ClusterRoleBinding",
			"metadata": map[string]interface{}{
				"name": modifiedname,
			},
			"roleRef": map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "ClusterRole",
				"name":       orgname,
			},
			"subjects": []interface{}{
				map[string]interface{}{
					"kind":      "ServiceAccount",
					"name":      orgname,
					"namespace": ns1,
				},
				map[string]interface{}{
					"kind":      "ServiceAccount",
					"name":      orgname,
					"namespace": ns3,
				},
				map[string]interface{}{
					"kind":      "ServiceAccount",
					"name":      orgname,
					"namespace": "random",
				},
				map[string]interface{}{
					"kind":      "ServiceAccount",
					"name":      orgname,
					"namespace": ns4,
				},
			}}).ResMap()

	expected := resmaptest_test.NewSeededRmBuilder(t, rf, m.ShallowCopy()).
		ReplaceResource(
			map[string]interface{}{
				"apiVersion": "rbac.authorization.k8s.io/v1",
				"kind":       "ClusterRoleBinding",
				"metadata": map[string]interface{}{
					"name": modifiedname,
				},
				"roleRef": map[string]interface{}{
					"apiVersion": "rbac.authorization.k8s.io/v1",
					"kind":       "ClusterRole",
					"name":       modifiedname,
				},
				// The following tests required a change in
				// getNameFunc implementation in order to leverage
				// the namespace field.
				"subjects": []interface{}{
					map[string]interface{}{
						"kind":      "ServiceAccount",
						"name":      prefixedname,
						"namespace": ns1,
					},
					map[string]interface{}{
						"kind":      "ServiceAccount",
						"name":      suffixedname,
						"namespace": ns2,
					},
					map[string]interface{}{
						"kind":      "ServiceAccount",
						"name":      orgname,
						"namespace": "random",
					},
					map[string]interface{}{
						"kind":      "ServiceAccount",
						"name":      orgname,
						"namespace": ns4,
					},
				},
			}).ResMap()

	clusterRoleId := resid.NewResId(
		gvk.Gvk{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRole"}, modifiedname)
	clusterRoleBindingId := resid.NewResId(
		gvk.Gvk{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: "ClusterRoleBinding"}, modifiedname)
	clusterRole, _ := expected.GetByCurrentId(clusterRoleId)
	clusterRole.AppendRefBy(clusterRoleBindingId)

	nrt := NewNameReferenceTransformer(defaultTransformerConfig.NameReference)
	err := nrt.Transform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err = expected.ErrorIfNotEqualLists(m); err != nil {
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}

// TestNameReferenceNamespace creates configmap, secret, deployment
// It validates the change done is IsSameFuzzyNamespace which
// uses the IsNsEquals method instead of the simple == operator.
func TestNameReferenceCandidateSelection(t *testing.T) {
	rf := resource.NewFactory(
		kunstruct.NewKunstructuredFactoryImpl())
	m := resmaptest_test.NewRmBuilder(t, rf).
		AddWithName("cm1", map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "ConfigMap",
			"metadata": map[string]interface{}{
				"name": "p1-cm1-hash",
			}}).
		AddWithNsAndName("default", "secret1", map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":      "p1-secret1-hash",
				"namespace": "default",
			}}).
		AddWithName("deploy1", deploymentMap("", "p1-deploy1", "cm1", "secret1")).
		ResMap()

	expected := resmaptest_test.NewSeededRmBuilder(t, rf, m.ShallowCopy()).
		ReplaceResource(deploymentMap("", "p1-deploy1", "p1-cm1-hash", "p1-secret1-hash")).
		ResMap()

>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	nrt := NewNameReferenceTransformer(defaultTransformerConfig.NameReference)
	err := nrt.Transform(m)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
<<<<<<< HEAD
	if !reflect.DeepEqual(m, expected) {
		err = expected.ErrorIfNotEqual(m)
=======

	if err = expected.ErrorIfNotEqualLists(m); err != nil {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
		t.Fatalf("actual doesn't match expected: %v", err)
	}
}
