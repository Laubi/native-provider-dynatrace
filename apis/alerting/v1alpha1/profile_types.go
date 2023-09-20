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

// +kubebuilder:validation:Enum=PREDEFINED;CUSTOM
type EventFilterType string

const (
	EventFilterTypePredefined EventFilterType = "PREDEFINED"
	EventFilterTypeCustom     EventFilterType = "CUSTOM"
)

// +kubebuilder:validation:Enum=AVAILABILITY;CUSTOM_ALERT;ERRORS;MONITORING_UNAVAILABLE;PERFORMANCE;RESOURCE_CONTENTION
type SeverityLevel string

const (
	SeverityLevelAvailability          SeverityLevel = "AVAILABILITY"
	SeverityLevelCustom                SeverityLevel = "CUSTOM_ALERT"
	SeverityLevelError                 SeverityLevel = "ERRORS"
	SeverityLevelMonitoringUnavailable SeverityLevel = "MONITORING_UNAVAILABLE"
	SeverityLevelSlowdown              SeverityLevel = "PERFORMANCE"
	SeverityLevelResource              SeverityLevel = "RESOURCE_CONTENTION"
)

// +kubebuilder:validation:Enum=EC2_HIGH_CPU;OSI_HIGH_CPU;ELB_HIGH_BACKEND_ERROR_RATE;PROCESS_NA_HIGH_CONN_FAIL_RATE;CUSTOM_APP_CRASH_RATE_INCREASED;CUSTOM_APPLICATION_ERROR_RATE_INCREASED;CUSTOM_APPLICATION_SLOWDOWN;CUSTOM_APPLICATION_UNEXPECTED_LOW_LOAD;CUSTOM_APPLICATION_UNEXPECTED_HIGH_LOAD;DCRUM_SVC_PERFORMANCE_DEGRADATION;DCRUM_SVC_LOW_AVAILABILITY;ESXI_GUEST_CPU_LIMIT_REACHED;ESXI_GUEST_ACTIVE_SWAP_WAIT;ESXI_HOST_CPU_SATURATION;ESXI_HOST_MEMORY_SATURATION
type EventType string

const (
	EventTypeAWSCPUSaturation                        EventType = "EC2_HIGH_CPU"
	EventTypeCPUSaturation                           EventType = "OSI_HIGH_CPU"
	EventTypeELBHighBackendErrorRate                 EventType = "ELB_HIGH_BACKEND_ERROR_RATE"
	EventTypeConnectivityProblem                     EventType = "PROCESS_NA_HIGH_CONN_FAIL_RATE"
	EventTypeCustomAppCrashRateIncrease              EventType = "CUSTOM_APP_CRASH_RATE_INCREASED"
	EventTypeCustomAppErrorRateIncrease              EventType = "CUSTOM_APPLICATION_ERROR_RATE_INCREASED"
	EventTypeCustomAppSlowdown                       EventType = "CUSTOM_APPLICATION_SLOWDOWN"
	EventTypeCustomAppUnexpectedLowLoad              EventType = "CUSTOM_APPLICATION_UNEXPECTED_LOW_LOAD"
	EventTypeCustomAppUnexpectedHighLoad             EventType = "CUSTOM_APPLICATION_UNEXPECTED_HIGH_LOAD"
	EventTypeDataCenterServicePerformanceDegredation EventType = "DCRUM_SVC_PERFORMANCE_DEGRADATION"
	EventTypeDataCenterServiceUnvailable             EventType = "DCRUM_SVC_LOW_AVAILABILITY"
	EventTypeESXiGuestCPUSaturation                  EventType = "ESXI_GUEST_CPU_LIMIT_REACHED"
	EventTypeESXiGuestMemorySaturation               EventType = "ESXI_GUEST_ACTIVE_SWAP_WAIT"
	EventTypeESXiHostCPUSaturation                   EventType = "ESXI_HOST_CPU_SATURATION"
	EventTypeESXiHostMemorySaturation                EventType = "ESXI_HOST_MEMORY_SATURATION"
)

// +kubebuilder:validation:Enum=BEGINS_WITH;ENDS_WITH;CONTAINS;REGEX_MATCHES;STRING_EQUALS
type Operator string

const (
	OperatorBeginsWith   Operator = "BEGINS_WITH"
	OperatorEndsWith     Operator = "ENDS_WITH"
	OperatorContains     Operator = "CONTAINS"
	OperatorRegexMatches Operator = "REGEX_MATCHES"
	OperatorStringEquals Operator = "STRING_EQUALS"
)

type MetadataFilterItem struct {
	MetadataKey   string `json:"metadataKey"`   // GET /api/v2/eventProperties for list of available keys
	MetadataValue string `json:"metadataValue"` // Value
}

type MetadataFilter struct {
	MetadataFilterItems []MetadataFilterItem `json:"metadataFilterItems"` // Define filters for event properties. A maximum of 20 properties is allowed.
}
type CustomEventFilter struct {
	// +optional
	Title *TextFilter `json:"titleFilter,omitempty"` // Title filter
	// +optional
	Description *TextFilter `json:"descriptionFilter,omitempty"` // Description filter
	// +optional
	MetadataFilter *MetadataFilter `json:"metadataFilter,omitempty"` // Property filters
}

// +kubebuilder:validation:Enum=NONE;INCLUDE_ANY;INCLUDE_ALL
type TagFilterIncludeMode string

const (
	None       TagFilterIncludeMode = "NONE"
	IncludeAny TagFilterIncludeMode = "INCLUDE_ANY"
	IncludeAll TagFilterIncludeMode = "INCLUDE_ALL"
)

type TextFilter struct {
	Operator      Operator `json:"operator"`      // Operator of the comparison
	Value         string   `json:"value"`         // The value to compare with
	Negate        bool     `json:"negate"`        // Negate the operator
	Enabled       bool     `json:"enabled"`       // Enable this filter
	CaseSensitive bool     `json:"caseSensitive"` // Case Sensitive comparison of text
}

type PredefinedEventFilter struct {
	EventType EventType `json:"eventType"` // Filter problems by a Dynatrace event type
	Negate    bool      `json:"negate"`    // Negate the given event type
}

type EventFilter struct {
	Type EventFilterType `json:"type"` // The type of event to filter by
	// +optional
	Predefined *PredefinedEventFilter `json:"predefinedFilter,omitempty"` // The predefined filter. Only valid if `type` is `PREDEFINED`
	// +optional
	Custom *CustomEventFilter `json:"customFilter,omitempty"` // The custom filter. Only valid if `type` is `CUSTOM`
}

type SeverityRule struct {
	SeverityLevel        SeverityLevel        `json:"severityLevel"`        // Problem severity level
	DelayInMinutes       int32                `json:"delayInMinutes"`       // Send a notification if a problem remains open longer than X minutes. Must be between 0 and 10000.
	TagFilterIncludeMode TagFilterIncludeMode `json:"tagFilterIncludeMode"` // Possible values are `NONE`, `INCLUDE_ANY` and `INCLUDE_ALL`
	// +optional
	Tags []string `json:"tagFilter,omitempty"`
}

// ProfileParameters are the configurable fields of a Profile.
type ProfileParameters struct {
	Name string `json:"name"`
	// +optional
	ManagementZone *string `json:"managementZone,omitempty"`
	// +optional
	SeverityRules []SeverityRule `json:"severityRules,omitempty"` // Define severity rules for profile. A maximum of 100 severity rules is allowed.
	// +optional
	EventFilters []EventFilter `json:"eventFilters,omitempty"` // Define event filters for profile. A maximum of 100 event filters is allowed.

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
