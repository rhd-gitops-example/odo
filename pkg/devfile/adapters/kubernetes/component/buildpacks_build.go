package component

import (
	"github.com/openshift/odo/pkg/devfile/adapters/common"
	"github.com/pkg/errors"
	"github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ClusterBuildStrategyKind v1alpha1.BuildStrategyKind = "ClusterBuildStrategy"
	registrySecret           *corev1.LocalObjectReference
)

func (a Adapter) runBuildpacksV3(parameters common.BuildParameters, isImageRegistryInternal bool) error {
	if isImageRegistryInternal {
		return errors.Errorf("Pushing to openshift internal image registry is currently unsupported")
	}
	registrySecret = &corev1.LocalObjectReference{
		Name: regcredName,
	}
}

func (a Adapter) build(sourceURL, imageDestination, secretName string) *v1alpha1.Build {
	// objectMeta := kclient.CreateObjectMeta("buildpack-"+a.ComponentName+"-build", a.Client.Namespace, labels, nil)
	// typeMeta := metav1.TypeMeta{
	// 	Kind: "Build",
	// 	APIVersion: "build.dev/v1alpha1",
	// }
	build := &v1alpha1.Build{
		Spec: v1alpha1.BuildSpec{
			Source: v1alpha1.GitSource{
				URL: sourceURL,
			},
			StrategyRef: &v1alpha1.StrategyRef{
				Name: "buildpacks-v3",
				Kind: &ClusterBuildStrategyKind,
			},
			Output: v1alpha1.Image{
				ImageURL: imageDestination,
				SecretRef: &corev1.LocalObjectReference{
					Name: regcredName,
				},
			},
		},
	}

}
