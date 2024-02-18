# Dynatrace Crossplane Provider

`provider-dynatrace` is a minimal [Crossplane](https://crossplane.io/) Provider to configure Dynatrace environments.

Currently, only a small subset of configurations are supported:
* Alerting Profiles
* Notifications for Email and Slack
* Auto-Tags

## Developing & Contributing

Refer to [our development documentation](DEVELOPING.MD) for more details on how to create a new resource.
Furtheremore, refer to Crossplane's [CONTRIBUTING.md] file for more information on how the
Crossplane community prefers to work. The [Provider Development][provider-dev]
guide may also be of use.

## Setting up the Dynatrace Crossplane Provider

All necessary configs are available in [the examples](./examples) directory. 

[CONTRIBUTING.md]: https://github.com/crossplane/crossplane/blob/master/CONTRIBUTING.md
[provider-dev]: https://github.com/crossplane/crossplane/blob/master/contributing/guide-provider-development.md
