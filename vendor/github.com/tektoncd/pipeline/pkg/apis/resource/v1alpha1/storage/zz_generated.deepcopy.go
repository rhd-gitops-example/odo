// +build !ignore_autogenerated

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package storage

import (
	v1alpha1 "github.com/tektoncd/pipeline/pkg/apis/resource/v1alpha1"
	v1 "k8s.io/api/core/v1"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactBucket) DeepCopyInto(out *ArtifactBucket) {
	*out = *in
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]v1alpha1.SecretParam, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactBucket.
func (in *ArtifactBucket) DeepCopy() *ArtifactBucket {
	if in == nil {
		return nil
	}
	out := new(ArtifactBucket)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArtifactPVC) DeepCopyInto(out *ArtifactPVC) {
	*out = *in
	if in.PersistentVolumeClaim != nil {
		in, out := &in.PersistentVolumeClaim, &out.PersistentVolumeClaim
		*out = new(v1.PersistentVolumeClaim)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArtifactPVC.
func (in *ArtifactPVC) DeepCopy() *ArtifactPVC {
	if in == nil {
		return nil
	}
	out := new(ArtifactPVC)
	in.DeepCopyInto(out)
	return out
}
