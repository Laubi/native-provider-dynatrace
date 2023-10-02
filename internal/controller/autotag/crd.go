package autotag

import (
	"github.com/crossplane/provider-dynatrace/apis/tags/v1alpha1"
	autotagging "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tags/autotagging/settings"
)

func crdToDto(v v1alpha1.AutoTagParameters) autotagging.Settings {

	return autotagging.Settings{
		Description: v.Description,
		Name:        v.Name,
		Rules:       convertRules(v.Rules),
	}
}

func convertRules(rules []v1alpha1.Rule) autotagging.Rules {
	result := make(autotagging.Rules, len(rules))

	for i, r := range rules {
		rule := autotagging.Rule{
			Type:               autotagging.RuleType(r.Type),
			Enabled:            r.Enabled,
			EntitySelector:     r.EntitySelector,
			ValueFormat:        r.Value,
			ValueNormalization: autotagging.Normalization(r.TagValueNormalization),
		}

		if rule.Type == autotagging.RuleTypes.Me {
			rule.AttributeRule = &autotagging.AutoTagAttributeRule{
				EntityType:                autotagging.AutoTagMeType(*r.AppliesTo),
				Conditions:                convertConditions(r.Conditions),
				AzureToPGPropagation:      r.AzureToPgPropagation,
				AzureToServicePropagation: r.AzureToServicePropagation,
				HostToPGPropagation:       r.HostToPgPropagation,
				PGToHostPropagation:       r.PgToHostPropagation,
				PGToServicePropagation:    r.PgToServicePropagation,
				ServiceToHostPropagation:  r.ServiceToHostPropagation,
				ServiceToPGPropagation:    r.ServiceToPGPropagation,
			}
		}

		result[i] = &rule
	}

	return result
}

func convertConditions(conditions []v1alpha1.Condition) autotagging.AttributeConditions {
	result := make(autotagging.AttributeConditions, len(conditions))

	for i, r := range conditions {
		result[i] = &autotagging.AttributeCondition{
			CaseSensitive:    r.CaseSensitive,
			DynamicKey:       r.DynamicKey,
			DynamicKeySource: r.DynamicKeySource,
			EntityID:         r.EntityId,
			EnumValue:        r.EnumValue,
			IntegerValue:     r.IntegerValue,
			Key:              autotagging.Attribute(r.Property),
			Operator:         autotagging.Operator(r.Operator),
			StringValue:      r.StringValue,
			Tag:              r.Tag,
		}
	}

	return result
}
