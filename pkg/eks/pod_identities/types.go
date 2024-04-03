/*
Copyright 2021 The Kubernetes Authors.

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

package pod_identities

import (
	"github.com/google/go-cmp/cmp"
)

// EKSPodIdentityAssociation represents an EKS pod identity association.
type EKSPodIdentityAssociation struct {
	ServiceAccountName      *string
	ServiceAccountNamespace *string
	RoleARN                 *string
	AssociationId           *string
}

// IsEqual determines if 2 EKSPodIdentityAssociation are equal.
func (e *EKSPodIdentityAssociation) IsEqual(other *EKSPodIdentityAssociation) bool {
	// NOTE: we don't compare the ARN as thats only for existing addons
	if e == other {
		return true
	}
	if !cmp.Equal(e.ServiceAccountName, other.ServiceAccountName) {
		return false
	}
	if !cmp.Equal(e.ServiceAccountNamespace, other.ServiceAccountNamespace) {
		return false
	}
	if !cmp.Equal(e.RoleARN, other.RoleARN) {
		return false
	}

	return true
}
