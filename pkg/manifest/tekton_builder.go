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
	elPatchDir      = "eventlistener_patches"
	rolebindingFile = "edit-rolebinding.yaml"
)

type patchStringValue struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type tektonBuilder struct {
	files    res.Resources
	cicdEnv  *config.Environment
	services []string // slice of all services
}

func buildEventlistenerResources(m *config.Manifest) (res.Resources, error) {
	files := make(res.Resources)
	cicdEnv, err := m.GetCICDEnvironment()
	if err != nil {
		return nil, err
	}
	tb := &tektonBuilder{files: files, cicdEnv: cicdEnv}
	err = m.Walk(tb)
	return tb.files, err
}

func (tk *tektonBuilder) Service(env *config.Environment, app *config.Application, svc *config.Service) error {

	tk.services = append(tk.services, svc.Name)
	svcFiles, err := getServiceFiles(tk.services, tk.cicdEnv, env, svc)
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
	return []patchStringValue{
		{
			Op:    "add",
			Path:  "/spec/triggers/-",
			Value: eventlisteners.CreateListenerTrigger(triggerName(svc.Name), eventlisteners.StageCIDryRunFilters, svcRepo, pipelines.Integration.Binding, pipelines.Integration.Template, svc.Webhook.Secret.Name, svc.Webhook.Secret.Namespace),
		},
	}, err
}

func triggerName(svc string) string {
	return fmt.Sprintf("app-ci-build-from-pr-%s", svc)
}

func getServiceFiles(services []string, cicdEnv *config.Environment, env *config.Environment, svc *config.Service) (res.Resources, error) {
	envFiles := res.Resources{}
	cicdPath := config.PathForEnvironment(cicdEnv)
	basePath := filepath.Join(cicdPath, "base")
	overlaysPath := filepath.Join(cicdPath, "overlays")
	overlaysFile := filepath.Join(overlaysPath, kustomization)
	overlayRel, err := filepath.Rel(overlaysPath, basePath)
	patchDir := filepath.Join(overlaysPath, elPatchDir)
	if err != nil {
		return nil, err
	}
	elPatch, err := eventlistenerPatch(env, svc)
	if err != nil {
		return nil, err
	}
	envFiles[filepath.Join(patchDir, patchFile(svc.Name))] = elPatch
	envFiles[overlaysFile] = elKustomiseTarget(cicdEnv.Name, overlayRel, services)

	return envFiles, nil
}

func patchFile(svc string) string {
	return fmt.Sprintf("%s_patch.yaml", svc)
}

func elKustomiseTarget(cicdNs string, base string, services []string) interface{} {

	GVK := gvk.Gvk{
		Group:   "tekton.dev",
		Version: "v1alpha1",
		Kind:    "EventListener",
	}
	target := &types.PatchTarget{
		Gvk:       GVK,
		Name:      "cicd-event-listener",
		Namespace: cicdNs,
	}
	Patches := []types.PatchJson6902{}
	for _, svc := range services {
		patch := types.PatchJson6902{
			Target: target,
			Path:   filepath.Join(elPatchDir, patchFile(svc)),
		}
		Patches = append(Patches, patch)
	}
	file := types.Kustomization{
		Bases:           []string{base},
		PatchesJson6902: Patches,
	}
	return file
}
