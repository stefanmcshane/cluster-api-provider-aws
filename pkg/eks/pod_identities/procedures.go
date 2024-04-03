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
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
)

var (
	// ErrNilPodIdentityAssociation defines an error for when a nil eks pod identity is returned.
	ErrNilPodIdentityAssociation = errors.New("nil eks pod identity association returned from create")
	// ErrPodIdentityAssociationNotFound defines an error for when an eks pod identity is not found.
	ErrPodIdentityAssociationNotFound = errors.New("eks pod identity association not found")
	// ErrPodIdentityAssociationAlreadyExists defines an error for when an eks pod identity already exists.
	ErrPodIdentityAssociationAlreadyExists = errors.New("eks pod identity association already exists")
)

// DeletePodIdentityAssociationProcedure is a procedure that will delete an EKS eks pod identity.
type DeletePodIdentityAssociationProcedure struct {
	eksClient             eksiface.EKSAPI
	clusterName           string
	existingAssociationId string
}

// Do implements the logic for the procedure.
func (p *DeletePodIdentityAssociationProcedure) Do(_ context.Context) error {
	input := &eks.DeletePodIdentityAssociationInput{
		AssociationId: aws.String(p.existingAssociationId),
		ClusterName:   aws.String(p.clusterName),
	}

	if _, err := p.eksClient.DeletePodIdentityAssociation(input); err != nil {
		return fmt.Errorf("deleting eks eks pod identity %s: %w", p.existingAssociationId, err)
	}

	return nil
}

// Name is the name of the procedure.
func (p *DeletePodIdentityAssociationProcedure) Name() string {
	return "eks_pod_identity_delete"
}

// CreatePodIdentityAssociationProcedure is a procedure that will create an EKS eks pod identity for a cluster.
type CreatePodIdentityAssociationProcedure struct {
	eksClient      eksiface.EKSAPI
	clusterName    string
	newAssociation *EKSPodIdentityAssociation
}

// Do implements the logic for the procedure.
func (p *CreatePodIdentityAssociationProcedure) Do(_ context.Context) error {
	if p.newAssociation == nil {
		return fmt.Errorf("getting desired eks pod identity for cluster %s: %w", p.clusterName, ErrPodIdentityAssociationNotFound)
	}

	input := &eks.CreatePodIdentityAssociationInput{
		ClusterName:    aws.String(p.clusterName),
		Namespace:      p.newAssociation.ServiceAccountNamespace,
		RoleArn:        p.newAssociation.RoleARN,
		ServiceAccount: p.newAssociation.ServiceAccountName,
	}

	output, err := p.eksClient.CreatePodIdentityAssociation(input)
	if err != nil {
		return fmt.Errorf("creating desired eks pod identity for cluster %s: %w", p.clusterName, err)
	}

	if output == nil {
		return ErrNilPodIdentityAssociation
	}
	if output.Association == nil {
		return ErrNilPodIdentityAssociation
	}

	return nil
}

// Name is the name of the procedure.
func (p *CreatePodIdentityAssociationProcedure) Name() string {
	return "eks_pod_identity_create"
}
