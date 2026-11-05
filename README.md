# KubeVirt UI Operator

Kubernetes operator and install assets for the **OpenShift Virtualization console UI** stack: the dynamic plugin, apiserver-proxy, UI configuration, and everything the cluster needs to run the KubeVirt UI **without** owning the virtualization backend (KubeVirt, CDI, CNAO, and the rest of the CNV/HCO operands).

Today this repository ships a **Helm chart** and an **alpha operator controller** (`main.go`, `controllers/`) that reconcile the same stack from a `KubeVirtConsoleUI` CR. The chart remains a supported bootstrap path; the CR is the long-term day-2 interface.

## Why this repo exists

| Layer | Owned by | This repo |
|-------|----------|-----------|
| VM / storage / network control plane | CNV, HCO, KubeVirt, CDI, … | No |
| Console plugin + proxy + UI config | Today: HCO when CNV is installed | **Yes** (Helm + operator CR) |
| Plugin source (React) | [kubevirt-plugin](https://github.com/kubevirt-ui/kubevirt-plugin) | No (consumes published images) |
| Proxy implementation | [kubevirt-apiserver-proxy](https://github.com/kubevirt-ui/kubevirt-apiserver-proxy) | No (consumes published images) |

Use this project when you want the **same UI experience as OpenShift Virtualization** on a cluster that already has (or will have) a KubeVirt API, but you do not want HCO to own the console plugin deployment—for example custom distributions, dev clusters, or splitting UI lifecycle from CNV upgrades.

## What will be managed

The operator (and today’s Helm chart) target the **UI operands only**:

- **Deployments**: `kubevirt-console-plugin`, `kubevirt-apiserver-proxy`
- **Services** with OpenShift serving certificates
- **ConfigMaps**: `nginx-conf`, `kubevirt-user-settings`, `kubevirt-ui-features`
- **ConsolePlugin** CR: `kubevirt-plugin` (backend service + proxy to apiserver-proxy)
- **RBAC**: UI config reader Role/RoleBinding for authenticated console users
- **Optional**: NetworkPolicies, Route for off-cluster console, registration of the plugin on `consoles.operator.openshift.io/cluster`

The operator will **not** install virt-api, virt-controller, CDI, subscriptions, or `HyperConverged`—only the console-facing slice.

## Repository layout

```text
kubevirt-ui-operator/
├── cmd/main.go                 # Operator SDK manager entrypoint
├── api/v1alpha1/               # KubeVirtConsoleUI CRD types
├── controllers/                # Reconcile loop
├── internal/                   # defaults, labels, manifests, platform checks
├── config/                     # kustomize: crd, rbac, manager, OLM CSV/bundle
├── Dockerfile
├── Makefile                    # Operator SDK targets (build, deploy, bundle, …)
├── docs/
│   ├── ROADMAP.md
│   └── design/
└── helm/kubevirt-console-ui/
```

Planned additions: envtest/e2e on OpenShift, published operator image defaults, compatibility matrix.

## Quick start (Helm)

Prerequisites:

- **OpenShift Console:** Console Operator healthy, `ConsolePlugin` CRD (`console.openshift.io`), console UI running in `openshift-console`; cluster admin or equivalent if enabling automatic plugin registration on `consoles.operator.openshift.io/cluster`
- Namespace **`kubevirt-ui`** by default (same as operator `OPERAND_NAMESPACE`)—do **not** install on top of an existing HCO-managed UI in the same namespace
- KubeVirt API available to users (RBAC + running control plane)
- Container images for plugin and proxy (pin to your CNV/HCO versions when possible)

Full prerequisite list: [docs/ROADMAP.md#prerequisites-cluster-and-console](docs/ROADMAP.md#prerequisites-cluster-and-console).

```bash
helm upgrade --install kubevirt-ui ./helm/kubevirt-console-ui \
  --namespace kubevirt-ui \
  --create-namespace \
  --set plugin.image=quay.io/kubevirt-ui/kubevirt-plugin:<tag> \
  --set proxy.image=quay.io/kubevirt-ui/kubevirt-apiserver-proxy:<tag>
```

Image references used by HCO are defined in [`deploy/images.env`](https://github.com/kubevirt/hyperconverged-cluster-operator/blob/main/deploy/images.env) (`KUBEVIRT_CONSOLE_PLUGIN_IMAGE`, `KUBEVIRT_CONSOLE_PROXY_IMAGE`).

Full chart documentation: [`helm/kubevirt-console-ui/README.md`](helm/kubevirt-console-ui/README.md).

## Roadmap

Detailed plan (phases, prerequisites including **OpenShift Console**, milestones, and team workflow): **[docs/ROADMAP.md](docs/ROADMAP.md)**.

1. **Helm chart** (current): install and document UI-only stack; align names/labels with HCO for easier migration.
2. **API + controller**: reconcile the same resources from a custom resource (image overrides, replicas, TLS, network policy, console registration).
3. **OLM packaging**: ClusterServiceVersion for install via Operator Lifecycle Manager where appropriate.
4. **Integration**: optional coordination with [kubevirt-plugin](https://github.com/kubevirt-ui/kubevirt-plugin) release process for default image tags and compatibility matrix.

## Related projects

- [kubevirt-ui/kubevirt-plugin](https://github.com/kubevirt-ui/kubevirt-plugin) — OpenShift Console dynamic plugin (UI source)
- [kubevirt-ui/kubevirt-apiserver-proxy](https://github.com/kubevirt-ui/kubevirt-apiserver-proxy) — API proxy used by the plugin for large-cluster list/filter behavior
- [kubevirt/hyperconverged-cluster-operator](https://github.com/kubevirt/hyperconverged-cluster-operator) — upstream owner of the same UI operands in full CNV installs

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) and the operator delivery plan in [docs/ROADMAP.md](docs/ROADMAP.md). Reviewers are listed in [OWNERS](OWNERS) (same roster as kubevirt-plugin).

## License

See [LICENSE](LICENSE).
