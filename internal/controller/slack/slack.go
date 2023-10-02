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

package slack

import (
	"context"
	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/provider-dynatrace/internal/credentials"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications"
	notificationSettings "github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/builtin/problem/notifications/slack/settings"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/settings"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/connection"
	"github.com/crossplane/crossplane-runtime/pkg/controller"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/crossplane/provider-dynatrace/apis/notification/v1alpha1"
	apisv1alpha1 "github.com/crossplane/provider-dynatrace/apis/v1alpha1"
	"github.com/crossplane/provider-dynatrace/internal/features"
)

const (
	errNotSlack     = "managed resource is not a Slack custom resource"
	errTrackPCUsage = "cannot track ProviderConfig usage"
	errGetPC        = "cannot get ProviderConfig"
	errGetCreds     = "cannot get credentials"

	errNewClient = "cannot create new Service"
)

type service struct {
	client settings.CRUDService[*notifications.Notification]
}

func newService(data []byte) (*service, error) {
	c, err := credentials.Unmarshal(data)
	if err != nil {
		return nil, err
	}

	return &service{
		client: notifications.Service(c, notifications.Types.Email),
	}, nil
}

// Setup adds a controller that reconciles Slack managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(v1alpha1.SlackGroupKind)

	cps := []managed.ConnectionPublisher{managed.NewAPISecretPublisher(mgr.GetClient(), mgr.GetScheme())}
	if o.Features.Enabled(features.EnableAlphaExternalSecretStores) {
		cps = append(cps, connection.NewDetailsManager(mgr.GetClient(), apisv1alpha1.StoreConfigGroupVersionKind))
	}

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.SlackGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube:         mgr.GetClient(),
			usage:        resource.NewProviderConfigUsageTracker(mgr.GetClient(), &apisv1alpha1.ProviderConfigUsage{}),
			newServiceFn: newService}),
		managed.WithLogger(o.Logger.WithValues("controller", name)),
		managed.WithPollInterval(o.PollInterval),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))),
		managed.WithConnectionPublishers(cps...))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		WithEventFilter(resource.DesiredStateChanged()).
		For(&v1alpha1.Slack{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

// A connector is expected to produce an ExternalClient when its Connect method
// is called.
type connector struct {
	kube         client.Client
	usage        resource.Tracker
	newServiceFn func(creds []byte) (*service, error)
}

// Connect typically produces an ExternalClient by:
// 1. Tracking that the managed resource is using a ProviderConfig.
// 2. Getting the managed resource's ProviderConfig.
// 3. Getting the credentials specified by the ProviderConfig.
// 4. Using the credentials to form a client.
func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*v1alpha1.Slack)
	if !ok {
		return nil, errors.New(errNotSlack)
	}

	if err := c.usage.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, errTrackPCUsage)
	}

	pc := &apisv1alpha1.ProviderConfig{}
	if err := c.kube.Get(ctx, types.NamespacedName{Name: cr.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, errGetPC)
	}

	cd := pc.Spec.Credentials
	data, err := resource.CommonCredentialExtractor(ctx, cd.Source, c.kube, cd.CommonCredentialSelectors)
	if err != nil {
		return nil, errors.Wrap(err, errGetCreds)
	}

	svc, err := c.newServiceFn(data)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	return &external{service: svc}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	// A 'client' used to connect to the external resource API. In practice this
	// would be something like an AWS SDK client.
	service *service
}

func (c *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Slack)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotSlack)
	}

	id := meta.GetExternalName(cr)
	var n notifications.Notification
	err := c.service.client.Get(id, &n)
	if err != nil {
		var restError rest.Error
		if errors.As(err, &restError) {
			if restError.Code == http.StatusNotFound {
				return managed.ExternalObservation{ResourceExists: false}, nil
			}
		}

		return managed.ExternalObservation{}, err
	}

	cr.Status.SetConditions(xpv1.Available())

	local := crdToDto(cr.Spec.ForProvider)
	if diff := cmp.Diff(n, local, cmpopts.IgnoreFields(notifications.Notification{}, "LegacyID"), cmpopts.EquateEmpty()); diff != "" {

		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: false,
			Diff:             diff,
		}, nil
	}

	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (c *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Slack)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotSlack)
	}

	cr.Status.SetConditions(xpv1.Creating())

	n := crdToDto(cr.Spec.ForProvider)
	apiResp, err := c.service.client.Create(&n)
	if err != nil {
		return managed.ExternalCreation{}, err
	}

	meta.SetExternalName(cr, apiResp.ID)

	return managed.ExternalCreation{}, nil
}

func (c *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.Slack)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotSlack)
	}

	id := meta.GetExternalName(cr)
	n := crdToDto(cr.Spec.ForProvider)
	err := c.service.client.Update(id, &n)
	if err != nil {
		return managed.ExternalUpdate{}, err
	}

	return managed.ExternalUpdate{}, nil
}

func (c *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Slack)
	if !ok {
		return errors.New(errNotSlack)
	}

	err := c.service.client.Delete(meta.GetExternalName(cr))
	if err != nil {
		return err
	}

	return nil
}

func crdToDto(v v1alpha1.SlackParameters) notifications.Notification {

	n := notifications.Notification{
		Type:      notifications.Types.Slack,
		Enabled:   v.Active,
		Name:      v.Name,
		ProfileID: *v.AlertingProfile,

		Slack: &notificationSettings.Slack{
			URL:     v.Url,
			Channel: v.Channel,
			Message: v.Message,
		},
	}
	return n
}
