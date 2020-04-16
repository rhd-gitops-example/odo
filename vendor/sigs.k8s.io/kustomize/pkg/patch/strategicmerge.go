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

package patch

<<<<<<< HEAD
// StrategicMerge represents a relative path to a
// stategic merge patch with the format
// https://github.com/kubernetes/community/blob/master/contributors/devel/strategic-merge-patch.md
type StrategicMerge string

// Append appends a slice of patch paths to a StrategicMerge slice
func Append(patches []StrategicMerge, paths ...string) []StrategicMerge {
	for _, p := range paths {
		patches = append(patches, StrategicMerge(p))
=======
import "sigs.k8s.io/kustomize/v3/pkg/types"

// Append appends a slice of patch paths to a PatchStrategicMerge slice
func Append(patches []types.PatchStrategicMerge, paths ...string) []types.PatchStrategicMerge {
	for _, p := range paths {
		patches = append(patches, types.PatchStrategicMerge(p))
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
	}
	return patches
}

<<<<<<< HEAD
// Exist determines if a patch path exists in a slice of StrategicMerge
func Exist(patches []StrategicMerge, path string) bool {
	for _, p := range patches {
		if p == StrategicMerge(path) {
=======
// Exist determines if a patch path exists in a slice of PatchStrategicMerge
func Exist(patches []types.PatchStrategicMerge, path string) bool {
	for _, p := range patches {
		if p == types.PatchStrategicMerge(path) {
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
			return true
		}
	}
	return false
}
<<<<<<< HEAD
=======

// Delete deletes patches from a PatchStrategicMerge slice
func Delete(patches []types.PatchStrategicMerge, paths ...string) []types.PatchStrategicMerge {
	// Convert paths into PatchStrategicMerge slice
	convertedPath := make([]types.PatchStrategicMerge, len(paths))
	for i, p := range paths {
		convertedPath[i] = types.PatchStrategicMerge(p)
	}

	filteredPatches := make([]types.PatchStrategicMerge, 0, len(patches))
	for _, containedPatch := range patches {
		if !Exist(convertedPath, string(containedPatch)) {
			filteredPatches = append(filteredPatches, containedPatch)
		}
	}
	return filteredPatches
}
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
