# Contributing

Thanks for helping build **kubevirt-ui-operator**. This repo ships the Helm chart today and the KubeVirt console UI operator; see [docs/ROADMAP.md](docs/ROADMAP.md) for phases and scope.

## Before you open a PR

1. Read [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) and [docs/ROADMAP.md](docs/ROADMAP.md#prerequisites-cluster-and-console).
2. For Helm-only changes: run `helm lint` and `helm template` (CI runs the same).
3. For operator changes: run `make test`, `make lint`, and `make verify-manifests` if you change API markers; run `make verify-generate` if you change API types. Commit generated files under `config/crd/bases/`, `config/rbac/role.yaml`, and `api/v1alpha1/zz_generated.deepcopy.go`.

## Pull requests

- One logical change per PR when possible (scaffold vs reconcile vs chart).
- Link related issues or Jira tickets in the description.
- Update [docs/parity.md](docs/parity.md) when an operand is implemented in the controller.
- API changes need an update to [docs/design/001-kubevirt-console-ui-cr.md](docs/design/001-kubevirt-console-ui-cr.md) or a new RFC.

## Code style

- Go: `make fmt vet test lint`
- YAML/Helm: match existing chart conventions and HCO resource names.

## Reviews

`OWNERS` lists approvers and reviewers (same roster as [kubevirt-plugin](https://github.com/kubevirt-ui/kubevirt-plugin/blob/main/OWNERS)). Operator and RBAC changes should get a second pair of eyes.
