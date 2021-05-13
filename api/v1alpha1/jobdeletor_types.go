/*
Copyright 2021.

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

package v1alpha1

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JobDeletorSpec defines the desired state of JobDeletor
type JobDeletorSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Format=string
	// +kubebuilder:validation:Pattern=succeeded|failed|all
	DeletionTargetStatus string `json:"deletionTargetStatus"`
}

// JobDeletorStatus defines the observed state of JobDeletor
type JobDeletorStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// +optional
	DeletedJobs []batchv1.Job `json:"deletedJobs"`
}

//+kubebuilder:object:root=true
//+kubebuilder:resource:shortName=jd;jds
//+kubebuilder:subresource:status

// JobDeletor is the Schema for the jobdeletors API
type JobDeletor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JobDeletorSpec   `json:"spec,omitempty"`
	Status JobDeletorStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JobDeletorList contains a list of JobDeletor
type JobDeletorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []JobDeletor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&JobDeletor{}, &JobDeletorList{})
}
