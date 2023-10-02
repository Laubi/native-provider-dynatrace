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

// EmailParameters are the configurable fields of a Email.
type EmailParameters struct {
	Enabled     bool   `json:"enabled"`
	DisplayName string `json:"displayName"`

	To []string `json:"to"`
	// +optional
	Cc []string `json:"cc"`
	// +optional
	Bcc []string `json:"bcc"`

	Subject                      string `json:"subject"`
	SendEmailWhenProblemIsClosed bool   `json:"sendEmailWhenProblemIsClosed"`
	Body                         string `json:"body"`

	// +crossplane:generate:reference:type=github.com/crossplane/provider-dynatrace/apis/alerting/v1alpha1.Profile
	// +crossplane:generate:reference:extractor=github.com/crossplane/provider-dynatrace/apis/alerting/v1alpha1.ProfileID()
	// +optional
	AlertingProfile *string `json:"alertingProfile"`

	// +optional
	// +immutable
	AlertingProfileRef *xpv1.Reference `json:"alertingProfileRef,omitempty"`

	// +optional
	// +immutable
	AlertingProfileSelector *xpv1.Selector `json:"alertingProfileSelector,omitempty"`
}

// EmailObservation are the observable fields of a Email.
type EmailObservation struct {
}

// A EmailSpec defines the desired state of a Email.
type EmailSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       EmailParameters `json:"forProvider"`
}

// A EmailStatus represents the observed state of a Email.
type EmailStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          EmailObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Email is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,dynatrace}
type Email struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSpec   `json:"spec"`
	Status EmailStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// EmailList contains a list of Email
type EmailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Email `json:"items"`
}

// Email type metadata.
var (
	EmailKind             = reflect.TypeOf(Email{}).Name()
	EmailGroupKind        = schema.GroupKind{Group: Group, Kind: EmailKind}.String()
	EmailKindAPIVersion   = EmailKind + "." + SchemeGroupVersion.String()
	EmailGroupVersionKind = SchemeGroupVersion.WithKind(EmailKind)
)

func init() {
	SchemeBuilder.Register(&Email{}, &EmailList{})
}
