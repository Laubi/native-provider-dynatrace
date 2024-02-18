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

package autotag

import (
	"context"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/provider-dynatrace/apis/tags/v1alpha1"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api"
	autotaggingservice "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/tags/autotagging/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/crossplane/crossplane-runtime/pkg/test"
)

// Unlike many Kubernetes projects Crossplane does not use third party testing
// libraries, per the common Go test review comments. Crossplane encourages the
// use of table driven unit tests. The tests of the crossplane-runtime project
// are representative of the testing style Crossplane encourages.
//
// https://github.com/golang/go/wiki/TestComments
// https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md#contributing-code

type mockClient struct {
	list func() (api.Stubs, error)
	get  func(id string, v *autotaggingservice.Settings) error
}

func (m mockClient) List() (api.Stubs, error) {
	return m.list()
}

func (m mockClient) Get(id string, v *autotaggingservice.Settings) error {
	return m.get(id, v)
}

func (m mockClient) SchemaID() string {
	panic("not used")
}

func (m mockClient) Create(_ *autotaggingservice.Settings) (*api.Stub, error) {
	panic("not used")
}

func (m mockClient) Update(_ string, _ *autotaggingservice.Settings) error {
	panic("not used")
}

func (m mockClient) Delete(_ string) error {
	panic("not used")
}

func (m mockClient) Name() string {
	panic("not used")
}

var _ settings.CRUDService[*autotaggingservice.Settings] = mockClient{}

func TestObserve(t *testing.T) {
	type fields struct {
		service mockClient
	}

	type args struct {
		ctx context.Context
		mg  resource.Managed
	}

	type want struct {
		o   managed.ExternalObservation
		err error
	}

	cases := map[string]struct {
		reason string
		fields fields
		args   args
		want   want
	}{
		"SuccessNotExists": {
			reason: "We should not return an error if a resource does not exist",
			fields: fields{
				service: mockClient{
					get: func(_ string, _ *autotaggingservice.Settings) error {
						return rest.Error{
							Code: http.StatusNotFound,
						}
					},
				},
			},
			args: args{
				ctx: nil,
				mg: &v1alpha1.AutoTag{
					ObjectMeta: v1.ObjectMeta{
						Annotations: map[string]string{
							meta.AnnotationKeyExternalName: "generated-id",
						},
					},
				},
			},
			want: want{
				o: managed.ExternalObservation{
					ResourceExists: false,
				},
				err: nil,
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			e := external{client: tc.fields.service}
			got, err := e.Observe(tc.args.ctx, tc.args.mg)
			if diff := cmp.Diff(tc.want.err, err, test.EquateErrors()); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want error, +got error:\n%s\n", tc.reason, diff)
			}
			if diff := cmp.Diff(tc.want.o, got); diff != "" {
				t.Errorf("\n%s\ne.Observe(...): -want, +got:\n%s\n", tc.reason, diff)
			}
		})
	}
}
