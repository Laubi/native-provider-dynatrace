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

// AutoTagParameters are the configurable fields of a AutoTag.
type AutoTagParameters struct {
	Name string `json:"name"`

	// +optional
	Description *string `json:"description"`

	// +optional
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Enabled bool `json:"enabled"`

	// +optional
	Value *string `json:"value"`

	// +kubebuilder:validation:Enum=ME;SELECTOR
	Type string `json:"type"`

	// +kubebuilder:validation:Enum=Leave text as-is;To lower case;To upper case
	TagValueNormalization string `json:"tagValueNormalization"`

	// +optional
	EntitySelector *string `json:"entitySelector"`

	// +kubebuilder:validation:Enum=APPLICATION;AWS_APPLICATION_LOAD_BALANCER;AWS_CLASSIC_LOAD_BALANCER;AWS_NETWORK_LOAD_BALANCER;AWS_RELATIONAL_DATABASE_SERVICE;AZURE;CUSTOM_APPLICATION;CUSTOM_DEVICE;DCRUM_APPLICATION;ESXI_HOST;EXTERNAL_SYNTHETIC_TEST;HOST;HTTP_CHECK;MOBILE_APPLICATION;PROCESS_GROUP;SERVICE;SYNTHETIC_TEST
	// +optional
	AppliesTo *string `json:"appliesTo"`

	// +optional
	Conditions []Condition `json:"conditions"`

	// +optional
	ServiceToHostPropagation *bool `json:"serviceToHostPropagation"`

	// +optional
	ServiceToPGPropagation *bool `json:"serviceToPGPropagation"`

	// +optional
	HostToPgPropagation *bool `json:"hostToPgPropagation"`

	// +optional
	PgToHostPropagation *bool `json:"pgToHostPropagation"`

	// +optional
	PgToServicePropagation *bool `json:"pgToServicePropagation"`

	// +optional
	ServiceToPgPropagation *bool `json:"serviceToPgPropagation"`

	// +optional
	AzureToPgPropagation *bool `json:"azureToPgPropagation"`

	// +optional
	AzureToServicePropagation *bool `json:"azureToServicePropagation"`
}

type Condition struct {
	Property string `json:"property"`

	// +kubebuilder:validation:Enum=BEGINS_WITH;CONTAINS;ENDS_WITH;EQUALS;EXISTS;GREATER_THAN;GREATER_THAN_OR_EQUAL;IS_IP_IN_RANGE;LOWER_THAN;LOWER_THAN_OR_EQUAL;NOT_BEGINS_WITH;NOT_CONTAINS;NOT_ENDS_WITH;NOT_EQUALS;NOT_EXISTS;NOT_GREATER_THAN;NOT_GREATER_THAN_OR_EQUAL;NOT_IS_IP_IN_RANGE;NOT_LOWER_THAN;NOT_LOWER_THAN_OR_EQUAL;NOT_REGEX_MATCHES;NOT_TAG_KEY_EQUALS;REGEX_MATCHES;TAG_KEY_EQUALS
	Operator string `json:"operator"`

	// +optional
	CaseSensitive *bool `json:"caseSensitive"`

	// +optional
	DynamicKey *string `json:"dynamicKey"`

	// +optional
	DynamicKeySource *string `json:"dynamicKeySource"`

	// +optional
	EntityId *string `json:"entityId"`

	// +optional
	EnumValue *string `json:"enumValue"`

	// +optional
	IntegerValue *int `json:"integerValue"`

	// +optional
	StringValue *string `json:"stringValue"`

	// +optional
	Tag *string `json:"tag"`
}

// AutoTagObservation are the observable fields of a AutoTag.
type AutoTagObservation struct {
}

// A AutoTagSpec defines the desired state of a AutoTag.
type AutoTagSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       AutoTagParameters `json:"forProvider"`
}

// A AutoTagStatus represents the observed state of a AutoTag.
type AutoTagStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          AutoTagObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A AutoTag is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,dynatrace}
type AutoTag struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AutoTagSpec   `json:"spec"`
	Status AutoTagStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AutoTagList contains a list of AutoTag
type AutoTagList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AutoTag `json:"items"`
}

// AutoTag type metadata.
var (
	AutoTagKind             = reflect.TypeOf(AutoTag{}).Name()
	AutoTagGroupKind        = schema.GroupKind{Group: Group, Kind: AutoTagKind}.String()
	AutoTagKindAPIVersion   = AutoTagKind + "." + SchemeGroupVersion.String()
	AutoTagGroupVersionKind = SchemeGroupVersion.WithKind(AutoTagKind)
)

func init() {
	SchemeBuilder.Register(&AutoTag{}, &AutoTagList{})
}
