package v1alpha1

import (
	"github.com/crossplane/crossplane-runtime/pkg/reference"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
)

func ProfileID() reference.ExtractValueFn {
	return func(mg resource.Managed) string {
		r, ok := mg.(*Profile)
		if !ok {
			return ""
		}
		return r.Status.AtProvider.Id
	}
}
