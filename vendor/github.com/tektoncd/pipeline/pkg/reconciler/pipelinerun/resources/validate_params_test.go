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

package resources

import (
	"testing"

	tb "github.com/tektoncd/pipeline/internal/builder/v1beta1"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1beta1"
)

func TestValidateParamTypesMatching_Valid(t *testing.T) {
	tcs := []struct {
		name          string
		p             *v1beta1.Pipeline
		pr            *v1beta1.PipelineRun
		errorExpected bool
	}{{
		name: "proper param types",
		p: tb.Pipeline("a-pipeline", tb.PipelineSpec(
			tb.PipelineParamSpec("correct-type-1", v1beta1.ParamTypeString),
			tb.PipelineParamSpec("mismatching-type", v1beta1.ParamTypeString),
			tb.PipelineParamSpec("correct-type-2", v1beta1.ParamTypeArray))),
		pr: tb.PipelineRun("a-pipelinerun", tb.PipelineRunSpec(
			"test-pipeline",
			tb.PipelineRunParam("correct-type-1", "somestring"),
			tb.PipelineRunParam("mismatching-type", "astring"),
			tb.PipelineRunParam("correct-type-2", "another", "array"))),
		errorExpected: false,
	}, {
		name:          "no params to get wrong",
		p:             tb.Pipeline("a-pipeline"),
		pr:            tb.PipelineRun("a-pipelinerun"),
		errorExpected: false,
	}, {
		name: "string-array mismatch",
		p: tb.Pipeline("a-pipeline", tb.PipelineSpec(
			tb.PipelineParamSpec("correct-type-1", v1beta1.ParamTypeString),
			tb.PipelineParamSpec("mismatching-type", v1beta1.ParamTypeString),
			tb.PipelineParamSpec("correct-type-2", v1beta1.ParamTypeArray))),
		pr: tb.PipelineRun("a-pipelinerun",
			tb.PipelineRunSpec("test-pipeline",
				tb.PipelineRunParam("correct-type-1", "somestring"),
				tb.PipelineRunParam("mismatching-type", "an", "array"),
				tb.PipelineRunParam("correct-type-2", "another", "array"))),
		errorExpected: true,
	}, {
		name: "array-string mismatch",
		p: tb.Pipeline("a-pipeline", tb.PipelineSpec(
			tb.PipelineParamSpec("correct-type-1", v1beta1.ParamTypeString),
			tb.PipelineParamSpec("mismatching-type", v1beta1.ParamTypeArray),
			tb.PipelineParamSpec("correct-type-2", v1beta1.ParamTypeArray))),
		pr: tb.PipelineRun("a-pipelinerun",
			tb.PipelineRunSpec("test-pipeline",
				tb.PipelineRunParam("correct-type-1", "somestring"),
				tb.PipelineRunParam("mismatching-type", "astring"),
				tb.PipelineRunParam("correct-type-2", "another", "array"))),
		errorExpected: true,
	}}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateParamTypesMatching(&tc.p.Spec, tc.pr)
			if (!tc.errorExpected) && (err != nil) {
				t.Errorf("Pipeline.Validate() returned error: %v", err)
			}

			if tc.errorExpected && (err == nil) {
				t.Error("Pipeline.Validate() did not return error, wanted error")
			}
		})
	}
}

func TestValidateParamTypesMatching_Invalid(t *testing.T) {
	tcs := []struct {
		name string
		p    *v1beta1.Pipeline
		pr   *v1beta1.PipelineRun
	}{{
		name: "string-array mismatch",
		p: tb.Pipeline("a-pipeline", tb.PipelineSpec(
			tb.PipelineParamSpec("correct-type-1", v1beta1.ParamTypeString),
			tb.PipelineParamSpec("mismatching-type", v1beta1.ParamTypeString),
			tb.PipelineParamSpec("correct-type-2", v1beta1.ParamTypeArray))),
		pr: tb.PipelineRun("a-pipelinerun",
			tb.PipelineRunSpec("test-pipeline",
				tb.PipelineRunParam("correct-type-1", "somestring"),
				tb.PipelineRunParam("mismatching-type", "an", "array"),
				tb.PipelineRunParam("correct-type-2", "another", "array"))),
	}, {
		name: "array-string mismatch",
		p: tb.Pipeline("a-pipeline", tb.PipelineSpec(
			tb.PipelineParamSpec("correct-type-1", v1beta1.ParamTypeString),
			tb.PipelineParamSpec("mismatching-type", v1beta1.ParamTypeArray),
			tb.PipelineParamSpec("correct-type-2", v1beta1.ParamTypeArray))),
		pr: tb.PipelineRun("a-pipelinerun",
			tb.PipelineRunSpec("test-pipeline",
				tb.PipelineRunParam("correct-type-1", "somestring"),
				tb.PipelineRunParam("mismatching-type", "astring"),
				tb.PipelineRunParam("correct-type-2", "another", "array"))),
	}}
	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			if err := ValidateParamTypesMatching(&tc.p.Spec, tc.pr); err == nil {
				t.Errorf("Expected to see error when validating PipelineRun/Pipeline param types but saw none")
			}
		})
	}
}
