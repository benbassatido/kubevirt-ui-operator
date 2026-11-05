# Compatibility matrix

Document supported OpenShift versions and pinned **kubevirt-plugin** / **kubevirt-apiserver-proxy** image tags per operator release. Operand image tags follow HCO [`deploy/images.env`](https://github.com/kubevirt/hyperconverged-cluster-operator/blob/main/deploy/images.env) unless overridden in the CR or Helm values.

## Operator build (tracking `main`)

| Component | Version |
|-----------|---------|
| Go | 1.27 |
| Kubernetes API (`k8s.io/client-go`) | 1.37 (`v0.37.x`) |
| controller-runtime | 0.25.x |

## OpenShift / OKD

| Platform | Support status |
|----------|----------------|
| OpenShift 4.14+ with Console Operator and `ConsolePlugin` CRD | Target platform; validate on test clusters before claiming a release |

Console integration requires `consoles.operator.openshift.io` and `consoleplugins.console.openshift.io`. See [docs/ROADMAP.md](ROADMAP.md#prerequisites-cluster-and-console).

## Operand images

| Workload | Image reference |
|----------|-----------------|
| Console plugin | Helm `plugin.image` / operator env `KUBEVIRT_CONSOLE_PLUGIN_IMAGE` — default `quay.io/kubevirt-ui/kubevirt-plugin` |
| Apiserver proxy | Helm `proxy.image` / operator env `KUBEVIRT_APISERVER_PROXY_IMAGE` — default `quay.io/kubevirt-ui/kubevirt-apiserver-proxy` |

Example tags from HCO (change over time; verify in HCO for your release):

- `KUBEVIRT_CONSOLE_PLUGIN_IMAGE`
- `KUBEVIRT_CONSOLE_PROXY_IMAGE`

## Operator image

| Artifact | Registry |
|----------|----------|
| Manager | `quay.io/kubevirtui/kubevirt-ui-operator` |

Release tags: semver git tags (`v*`). Continuous builds from `main` publish `latest` and `sha-*` (see `.github/workflows/operator-image.yaml`).

## Release rows (fill in when tagging)

| Operator version | OpenShift tested | Plugin image tag | Proxy image tag |
|------------------|------------------|------------------|-----------------|
| _unreleased_ | — | — | — |
