package component

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/openshift/odo/pkg/devfile/adapters/common"
	"github.com/openshift/odo/pkg/devfile/adapters/kubernetes/utils"
	"github.com/openshift/odo/pkg/kclient"
	"github.com/openshift/odo/pkg/log"
	"github.com/pkg/errors"
	"github.com/shipwright-io/build/pkg/apis/build/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	ClusterBuildStrategyKind v1alpha1.BuildStrategyKind = "ClusterBuildStrategy"
)

func (a Adapter) runBuildpacksV3(parameters common.BuildParameters, isImageRegistryInternal bool) error {
	if isImageRegistryInternal {
		if err := a.createDockerCfgSecretForInternalRegistry(parameters.EnvSpecificInfo.GetNamespace()); err != nil {
			return errors.Wrap(err, "failed to create dockerconfig secret")
		}
	}

	labels := map[string]string{
		"component": a.ComponentName,
	}

	podSelector := fmt.Sprintf("component=%s", a.ComponentName)
	watchOptions := metav1.ListOptions{
		LabelSelector: podSelector,
	}

	// This should ideally be a field in common.BuildParameters
	sourceUrl := "https://github.com/sclorg/nodejs-ex"

	buildObj := a.build(labels, sourceUrl, parameters.Tag, regcredName)
	buildRunObj := a.buildRun(labels)

	// Create the build object on the cluster
	if _, err := a.Client.DynamicClient.Resource(secretGroupVersionResource).
		Namespace(parameters.EnvSpecificInfo.GetNamespace()).
		Create(context.Background(), buildObj); err != nil {
		return errors.Wrap(err, "failed to create build on cluster")
	}

	// Create the build Run object that references your build object on the cluster. This spins up the build pod
	if _, err := a.Client.DynamicClient.Resource(secretGroupVersionResource).
		Namespace(parameters.EnvSpecificInfo.GetNamespace()).
		Create(context.Background(), buildRunObj); err != nil {
		return errors.Wrap(err, "failed to create build on cluster")
	}

	log.Successf("Started builder pod using Buildpacks-v3 Build strategy")

	reader, _ := io.Pipe()
	controlC := make(chan os.Signal, 1)

	var cmdOutput string

	go utils.PipeStdOutput(cmdOutput, reader, controlC)

	s := log.Spinner("Waiting for builder pod to complete")

	if _, err := a.Client.WaitAndGetPod(watchOptions, corev1.PodSucceeded, "Waiting for builder pod to complete", false); err != nil {
		s.End(false)
		return errors.Wrapf(err, "unable to build image using Buildpacks-v3, error: %s", cmdOutput)
	}

	s.End(true)
	// Stop listening for a ^C so it doesnt perform terminateBuild during any later stages
	signal.Stop(controlC)
	log.Successf("Successfully built container image: %s", parameters.Tag)
	return nil

}

func (a Adapter) build(labels map[string]string, sourceURL, imageDestination, secretName string) *v1alpha1.Build {
	objectMeta := kclient.CreateObjectMeta("buildpack-"+a.ComponentName+"-build", a.Client.Namespace, labels, nil)
	typeMeta := metav1.TypeMeta{
		Kind:       "Build",
		APIVersion: "build.dev/v1alpha1",
	}
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

	build.ObjectMeta = objectMeta
	build.TypeMeta = typeMeta

	return build
}

func (a Adapter) buildRun(labels map[string]string) *v1alpha1.BuildRun {
	objectMeta := kclient.CreateObjectMeta("buildpack-"+a.ComponentName+"-buildrun", a.Client.Namespace, labels, nil)
	typeMeta := metav1.TypeMeta{
		Kind:       "BuildRun",
		APIVersion: "build.dev/v1alpha1",
	}
	buildRun := &v1alpha1.BuildRun{
		Spec: v1alpha1.BuildRunSpec{
			BuildRef: &v1alpha1.BuildRef{
				Name: "buildpack-" + a.ComponentName + "-build",
			},
		},
	}

	buildRun.ObjectMeta = objectMeta
	buildRun.TypeMeta = typeMeta
	return buildRun
}
