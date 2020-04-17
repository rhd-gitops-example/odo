package manifest

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/openshift/odo/pkg/manifest/config"
	"github.com/openshift/odo/pkg/manifest/eventlisteners"
	res "github.com/openshift/odo/pkg/manifest/resources"
	"sigs.k8s.io/kustomize/pkg/gvk"
	"sigs.k8s.io/kustomize/pkg/types"
)

const (
	elPatchFile     = "eventlistener_patch.yaml"
	rolebindingFile = "edit-rolebinding.yaml"
)

type patchStringValue struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type tektonBuilder struct {
	files res.Resources
}

func buildEventlistener(m *config.Manifest) (res.Resources, error) {
	files := make(res.Resources)
	tb := &tektonBuilder{files: files}
	err := m.Walk(tb)
	return tb.files, err
}

func (tk *tektonBuilder) Service(env *config.Environment, app *config.Application, svc *config.Service) error {

	svcPath := config.PathForService(env, svc)
	svcFiles, err := getServiceFiles(svcPath, env, svc)
	if err != nil {
		return err
	}
	tk.files = res.Merge(svcFiles, tk.files)
	return nil
}

func getPipelines(env *config.Environment, svc *config.Service) *config.Pipelines {
	if svc.Pipelines != nil {
		return svc.Pipelines
	}
	if env.Pipelines != nil {
		return env.Pipelines
	}
	return defaultPipeline()
}

func defaultPipeline() *config.Pipelines {
	return &config.Pipelines{
		Integration: &config.TemplateBinding{
			Template: "app-ci-template",
			Binding:  "app-ci-binding",
		},
		Deployment: &config.TemplateBinding{
			Template: "app-cd-template",
			Binding:  "app-cd-binding",
		},
	}
}

func extractRepo(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	parts := strings.Split(parsed.Path, "/")
	return fmt.Sprintf("%s/%s", parts[1], strings.TrimSuffix(parts[2], ".git")), nil
}

func eventlistenerPatch(env *config.Environment, svc *config.Service) ([]patchStringValue, error) {
	pipelines := getPipelines(env, svc)
	svcRepo, err := extractRepo(svc.SourceURL)
	if err != nil {
		return nil, err
	}

	if env.IsArgoCD {
		return []patchStringValue{
			{
				Op:    "add",
				Path:  "/spec/triggers/-",
				Value: eventlisteners.CreateListenerTrigger("app-ci-build-from-pr", eventlisteners.StageCIDryRunFilters, svcRepo, pipelines.Integration.Binding, pipelines.Integration.Template, svc.Webhook.Secret.Name, svc.Webhook.Secret.Namespace),
			},
		}, nil
	}

	return []patchStringValue{
		{
			Op:    "add",
			Path:  "/spec/triggers/-",
			Value: eventlisteners.CreateListenerTrigger("app-ci-build-from-pr", eventlisteners.StageCIDryRunFilters, svcRepo, pipelines.Integration.Binding, pipelines.Integration.Template, svc.Webhook.Secret.Name, svc.Webhook.Secret.Namespace),
		},
		{
			Op:    "add",
			Path:  "/spec/triggers/-",
			Value: eventlisteners.CreateListenerTrigger("app-cd-deploy-from-master", eventlisteners.StageCDDeployFilters, svcRepo, pipelines.Deployment.Binding, pipelines.Deployment.Template, svc.Webhook.Secret.Name, svc.Webhook.Secret.Namespace),
		},
	}, nil
}

func getServiceFiles(svcPath string, env *config.Environment, svc *config.Service) (res.Resources, error) {
	envFiles := res.Resources{}
	basePath := filepath.Join(svcPath, "base")
	overlaysPath := filepath.Join(svcPath, "overlays")
	overlaysFile := filepath.Join(overlaysPath, kustomization)
	overlayRel, err := filepath.Rel(overlaysPath, basePath)
	if err != nil {
		return nil, err
	}
	elPatch, err := eventlistenerPatch(env, svc)
	if err != nil {
		return nil, err
	}
	envFiles[filepath.Join(overlaysPath, elPatchFile)] = elPatch

	envFiles[overlaysFile] = elKustomiseTarget(overlayRel)

	return envFiles, nil
}

func elKustomiseTarget(base string) interface{} {

	GVK := gvk.Gvk{
		Group:   "tekton.dev",
		Version: "v1alpha1",
		Kind:    "EventListener",
	}
	target := &types.PatchTarget{
		Gvk:  GVK,
		Name: "cicd-event-listener",
	}
	Patches := []types.PatchJson6902{
		{
			Target: target,
			Path:   elPatchFile,
		},
	}
	file := types.Kustomization{
		Bases:           []string{base},
		PatchesJson6902: Patches,
	}
	return file
}
