# Development

## Helm chart

```bash
helm lint ./helm/kubevirt-console-ui
helm template kubevirt-ui ./helm/kubevirt-console-ui \
  --set plugin.image=quay.io/kubevirt-ui/kubevirt-plugin:<tag> \
  --set proxy.image=quay.io/kubevirt-ui/kubevirt-apiserver-proxy:<tag>
```

Smoke install on OpenShift: see [helm/kubevirt-console-ui/README.md](../helm/kubevirt-console-ui/README.md).

## Operator (Operator SDK + Kubebuilder)

This project uses the [Operator SDK](https://sdk.operatorframework.io/) Go layout (`go.kubebuilder.io/v4`):

- **Entrypoint:** `cmd/main.go`
- **API:** `api/v1alpha1/`
- **Controller:** `controllers/`
- **Operands:** `internal/manifests/` (Helm parity via `internal/defaults`, `internal/labels`)
- **Deploy / OLM:** `config/` (kustomize, CSV under `config/manifests/`)

Install the CLI (match `OPERATOR_SDK_VERSION` in the Makefile, currently v1.39.1):

```bash
# macOS example
brew install operator-sdk
operator-sdk version
```

### Local development

```bash
make help
make manifests        # CRD + RBAC from kubebuilder markers (needs controller-gen)
make verify-manifests # same as CI: fails if config/crd/bases or role.yaml differ or are untracked
make generate    # deepcopy (needs controller-gen)
make test
make build       # bin/manager
make run         # controller against current kubeconfig (uses your user RBAC)
```

Install the CRD, create sample namespaces, run the controller, then apply the sample:

```bash
make install
./scripts/apply-sample-cr.sh   # creates namespaces + applies config/samples/…
make run                       # separate terminal
```

Or create namespaces manually before `kubectl apply`:

```bash
kubectl create namespace openshift-cnv-test --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f config/samples/consoleui_v1alpha1_kubevirtconsoleui.yaml
```

Create the operand namespace before applying the CR (`scripts/apply-sample-cr.sh` uses `metadata.namespace` from the sample).

### Deploy manager to the cluster

```bash
make install     # CRD only (kustomize config/crd)
export IMG=quay.io/<your-user>/kubevirt-ui-operator:dev
make docker-build docker-push
make deploy      # openshift-cnv Namespace (if missing), manager Deployment, RBAC, SA
./scripts/apply-sample-cr.sh
```

`make deploy` applies **`config/deploy`** (namespace, CRD, manager RBAC/Deployment). **`make undeploy`** deletes only **`config/default`**—the operator manager and its RBAC, not **`openshift-cnv`** and not the CRD—so existing **`KubeVirtConsoleUI`** objects stay on the cluster until you remove them or run **`make uninstall`**.

The manager runs in **`openshift-cnv`** (HCO convention). The sample CR may use another namespace (e.g. `openshift-cnv-test`) for an isolated UI stack. Do not install UI operands on top of an existing HCO-managed stack in the same namespace.

Or manually:

```bash
kubectl create namespace openshift-cnv-test --dry-run=client -o yaml | kubectl apply -f -
kubectl apply -f config/samples/consoleui_v1alpha1_kubevirtconsoleui.yaml
```

`make deploy` runs `kustomize edit set image` under `config/manager/` (local change to that file is normal; do not commit the generated `images:` block unless you intend to pin a default image). The default `IMG=controller:latest` only works if that image exists on the cluster (build and push first).

If `make deploy` fails building `controller-gen`, ensure Go 1.22+ and use the Makefile’s pinned `CONTROLLER_TOOLS_VERSION` (v0.16.5), or install `controller-gen` on your `PATH`.

Undeploy the operator manager (CRD and CRs remain):

```bash
make undeploy    # Deployment, RBAC, metrics Service in openshift-cnv — not the namespace or CRD
```

Remove the API (deletes **all** `KubeVirtConsoleUI` instances cluster-wide):

```bash
make uninstall   # CRD only; custom resources are garbage-collected with the CRD
```

### OLM bundle (optional)

```bash
VERSION=0.1.0
IMG=quay.io/<your-user>/kubevirt-ui-operator:v${VERSION}
BUNDLE_IMG=quay.io/<your-user>/kubevirt-ui-operator-bundle:v${VERSION}

make bundle VERSION=${VERSION} IMG=${IMG}
make bundle-build bundle-push VERSION=${VERSION} BUNDLE_IMG=${BUNDLE_IMG}
```

Use the same `VERSION`, `IMG`, and `BUNDLE_IMG` on every target. A second `make` without those variables falls back to `VERSION=0.0.1` and `quay.io/kubevirt-ui/kubevirt-ui-operator-bundle:v0.0.1`.

See [Operator SDK bundle tutorial](https://sdk.operatorframework.io/docs/olm-integration/tutorial-bundle/).

The reconciler sets `PlatformRequirementsMet` when OpenShift console APIs are available, then reconciles UI operands. With `spec.console.registerPlugin: true`, the manager patches `consoles.operator.openshift.io/cluster` (requires ClusterRole rules on the operator SA).

### CI: operator image

Workflow [`.github/workflows/operator-image.yaml`](../.github/workflows/operator-image.yaml):

| Event | Behavior |
|-------|----------|
| Pull request | `docker build` only (validates the Dockerfile) |
| Push to `main` | Build and **push** to `quay.io/kubevirt-ui/kubevirt-ui-operator` |

Tags on `main`: `main`, `latest`, and `sha-<short>`.

Repository secrets required for publish (Quay robot account):

- `QUAY_USERNAME`
- `QUAY_PASSWORD`

Deploy the published manager to a cluster:

```bash
export IMG=quay.io/kubevirt-ui/kubevirt-ui-operator:sha-$(git rev-parse --short HEAD)
# or :main / :latest after a main-branch publish
make deploy
```

## Cluster prerequisites

Console and cluster requirements: [ROADMAP.md — Prerequisites](ROADMAP.md#prerequisites-cluster-and-console).
