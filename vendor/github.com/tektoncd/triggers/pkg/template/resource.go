/*
Copyright 2019 The Tekton Authors

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

package template

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	pipelinev1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	triggersv1 "github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/rand"
)

// uidMatch determines the uid variable within the resource template
var uidMatch = []byte(`$(uid)`)

// ResolvedTrigger contains the dereferenced TriggerBindings and
// TriggerTemplate after resolving the k8s ObjectRef.
type ResolvedTrigger struct {
	TriggerBindings []*triggersv1.TriggerBinding
	TriggerTemplate *triggersv1.TriggerTemplate
}

type getTriggerBinding func(name string, options metav1.GetOptions) (*triggersv1.TriggerBinding, error)
type getTriggerTemplate func(name string, options metav1.GetOptions) (*triggersv1.TriggerTemplate, error)

// ResolveTrigger takes in a trigger containing object refs to bindings and
// templates and resolves them to their underlying values.
func ResolveTrigger(trigger triggersv1.EventListenerTrigger, getTB getTriggerBinding, getTT getTriggerTemplate) (ResolvedTrigger, error) {
	tb := make([]*triggersv1.TriggerBinding, 0, len(trigger.Bindings))
	for _, b := range trigger.Bindings {
		tb2, err := getTB(b.Name, metav1.GetOptions{})
		if err != nil {
			return ResolvedTrigger{}, fmt.Errorf("error getting TriggerBinding %s: %w", b.Name, err)
		}
		tb = append(tb, tb2)
	}

	ttName := trigger.Template.Name
	tt, err := getTT(ttName, metav1.GetOptions{})
	if err != nil {
		return ResolvedTrigger{}, fmt.Errorf("error getting TriggerTemplate %s: %w", ttName, err)
	}
	return ResolvedTrigger{TriggerBindings: tb, TriggerTemplate: tt}, nil
}

// MergeInDefaultParams returns the params with the addition of all
// paramSpecs that have default values and are already in the params list
func MergeInDefaultParams(params []pipelinev1.Param, paramSpecs []pipelinev1.ParamSpec) []pipelinev1.Param {
	allParamsMap := map[string]pipelinev1.ArrayOrString{}
	for _, paramSpec := range paramSpecs {
		if paramSpec.Default != nil {
			allParamsMap[paramSpec.Name] = *paramSpec.Default
		}
	}
	for _, param := range params {
		allParamsMap[param.Name] = param.Value
	}
	return convertParamMapToArray(allParamsMap)
}

// ApplyParamsToResourceTemplate returns the TriggerResourceTemplate with the
// param values substituted for all matching param variables in the template
func ApplyParamsToResourceTemplate(params []pipelinev1.Param, rt json.RawMessage) json.RawMessage {
	// Assume the params are valid
	for _, param := range params {
		rt = applyParamToResourceTemplate(param, rt)
	}
	return rt
}

// applyParamToResourceTemplate returns the TriggerResourceTemplate with the
// param value substituted for all matching param variables in the template
func applyParamToResourceTemplate(param pipelinev1.Param, rt json.RawMessage) json.RawMessage {
	// Assume the param is valid
	paramVariable := fmt.Sprintf("$(params.%s)", param.Name)
	// Escape quotes so that that JSON strings can be appended to regular strings.
	// See #257 for discussion on this behavior.
	paramValue := strings.Replace(param.Value.StringVal, `"`, `\"`, -1)
	return bytes.Replace(rt, []byte(paramVariable), []byte(paramValue), -1)
}

// UID generates a random string like the Kubernetes apiserver generateName metafield postfix.
var UID = func() string { return rand.String(5) }

// ApplyUIDToResourceTemplate returns the TriggerResourceTemplate after uid replacement
// The same uid should be used per trigger to properly address resources throughout the TriggerTemplate.
func ApplyUIDToResourceTemplate(rt json.RawMessage, uid string) json.RawMessage {
	return bytes.Replace(rt, uidMatch, []byte(uid), -1)
}

func convertParamMapToArray(paramMap map[string]pipelinev1.ArrayOrString) []pipelinev1.Param {
	params := []pipelinev1.Param{}
	for name, value := range paramMap {
		params = append(params, pipelinev1.Param{Name: name, Value: value})
	}
	return params
}

// MergeBindingParams merges params across multiple bindings.
func MergeBindingParams(bindings []*triggersv1.TriggerBinding) ([]pipelinev1.Param, error) {
	params := []pipelinev1.Param{}
	for _, b := range bindings {
		params = append(params, b.Spec.Params...)
	}
	seen := make(map[string]bool, len(params))
	for _, p := range params {
		if seen[p.Name] {
			return nil, fmt.Errorf("duplicate param name: %s", p.Name)
		}
		seen[p.Name] = true
	}
	return params, nil
}
