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

// SlackParameters are the configurable fields of a Slack.
type SlackParameters struct {
	// The name of the Slack notification.
	Name string `json:"name"`

	// Whether this Slack notification is enabled.
	// +kubebuilder:default=true
	// +optional
	Enable bool `json:"enable"`

	// The URL to the Slack webhook.
	Url string `json:"url"`

	// Channel contains which channel the notification should be posted to.
	Channel string `json:"channel"`

	// The content of the message that will be postet.
	Message string `json:"message"`

	// ID of the associated alerting profile.
	// +crossplane:generate:reference:type=github.com/crossplane/provider-dynatrace/apis/alerting/v1alpha1.Profile
	// +crossplane:generate:reference:extractor=github.com/crossplane/provider-dynatrace/apis/alerting/v1alpha1.ProfileID()
	// +optional
	AlertingProfile *string `json:"alertingProfile"`

	// A referencer to retrieve the ID of an alerting profile.
	// +optional
	// +immutable
	AlertingProfileRef *xpv1.Reference `json:"alertingProfileRef,omitempty"`

	// A selector to select a referencer to retrieve the ID of an alerting profile.
	// +optional
	// +immutable
	AlertingProfileSelector *xpv1.Selector `json:"alertingProfileSelector,omitempty"`
}

// SlackObservation are the observable fields of a Slack.
type SlackObservation struct {
	ID            string  `json:"id,omitempty"`
	ObfuscatedUrl *string `json:"obfuscatedUrl,omitempty"`
}

// A SlackSpec defines the desired state of a Slack.
type SlackSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       SlackParameters `json:"forProvider"`
}

// A SlackStatus represents the observed state of a Slack.
type SlackStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          SlackObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Slack is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,dynatrace}
type Slack struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SlackSpec   `json:"spec"`
	Status SlackStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SlackList contains a list of Slack
type SlackList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Slack `json:"items"`
}

// Slack type metadata.
var (
	SlackKind             = reflect.TypeOf(Slack{}).Name()
	SlackGroupKind        = schema.GroupKind{Group: Group, Kind: SlackKind}.String()
	SlackKindAPIVersion   = SlackKind + "." + SchemeGroupVersion.String()
	SlackGroupVersionKind = SchemeGroupVersion.WithKind(SlackKind)
)

func init() {
	SchemeBuilder.Register(&Slack{}, &SlackList{})
}
