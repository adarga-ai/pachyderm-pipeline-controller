<!-- Copyright 2021 Adarga Limited -->
<!-- SPDX-License-Identifier: Apache-2.0 -->

# pachyderm-pipeline-controller

### A Kubernetes Operator for controlling Pachyderm Pipeline objects.

[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.0-4baaaa.svg)](CODE_OF_CONDUCT.md)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
[![REUSE status](https://api.reuse.software/badge/git.fsfe.org/reuse/api)](https://api.reuse.software/info/git.fsfe.org/reuse/api)

This controller introduces a custom `PachydermPipeline` object type
to Kubernetes. It holds the same information as a normal pachyderm
pipeline definition file, and this controller monitors these objects,
reflecting any changes into Pachyderm itself using the pachd
API.

On the surface, this doesn't appear to add much value. The process
for managing Pachyderm pipelines is largely the same; the only
tangible difference is what format you store them in and what
command you run to enact the change (e.g. `pachctl create pipeline`
vs `kubectl create`).

However, this subtle difference allows you to leverage your existing
Kubernetes-based CICD tooling to apply GitOps principles to the
pipelines. A tool such as ArgoCD or Flux - if it is already managing
your Kubernetes estate - can start to manage your Pachyderm pipelines
also.

## Getting Started

We currently do not have a publically accessible container image
for this controller, so you will need to create your own. You can
run `make image` to create a local docker image that can then be
pushed to the registry of your choice.

Before running the controller, deploy the CRD to your Kubernetes
cluster using:

```bash
kubectl create -f ./manifests/adarga.ai_pachydermpipelines.yaml
```

And, once you have modified the deployment file to point to your
image, deploy the controller into the same namespace as pachyderm
using:

```bash
kubectl create -n pachyderm -f ./manifests/controller.yaml
```

## Supported Constructs

The Pachyderm pipeline spec is currently only partially supported.
More features will be added, but currently the only supported input
types are:

- [X] PFS
- [ ] Union
- [X] Cross
- [ ] Cron
- [X] Join
- [ ] Group
- [ ] Git

If in doubt, look at the CRD definition.

## Examples

Examples are kept in the [examples](examples/) directory.

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md).

