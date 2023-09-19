/*
Copyright 2022 The Crossplane Authors.

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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// ProfileParameters are the configurable fields of a Profile.
type ProfileParameters struct {
	Name string `json:"name"`
}

// ProfileObservation are the observable fields of a Profile.
type ProfileObservation struct {
	Id string `json:"id"`
}

// A ProfileSpec defines the desired state of a Profile.
type ProfileSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       ProfileParameters `json:"forProvider"`
}

// A ProfileStatus represents the observed state of a Profile.
type ProfileStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          ProfileObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Profile is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,dynatrace}
type Profile struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProfileSpec   `json:"spec"`
	Status ProfileStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProfileList contains a list of Profile
type ProfileList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Profile `json:"items"`
}

// Profile type metadata.
var (
	ProfileKind             = reflect.TypeOf(Profile{}).Name()
	ProfileGroupKind        = schema.GroupKind{Group: Group, Kind: ProfileKind}.String()
	ProfileKindAPIVersion   = ProfileKind + "." + SchemeGroupVersion.String()
	ProfileGroupVersionKind = SchemeGroupVersion.WithKind(ProfileKind)
)

func init() {
	SchemeBuilder.Register(&Profile{}, &ProfileList{})
}
