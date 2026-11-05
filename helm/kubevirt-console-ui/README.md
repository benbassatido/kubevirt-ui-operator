# kubevirt-console-ui Helm chart

Deploys the **OpenShift Virtualization console UI stack** without the HCO/CNV backend operands (KubeVirt, CDI, and so on). Resource names and wiring match **kubevirt-ui-operator** (`internal/manifests/builders`) and HCO’s [`kubevirtConsolePlugin`](https://github.com/kubevirt/hyperconverged-cluster-operator/blob/main/controllers/handlers/kubevirtConsolePlugin.go) handler.

## What gets installed

| Resource | Name (defaults) |
|----------|-------------------|
| Deployment | `kubevirt-console-plugin`, `kubevirt-apiserver-proxy` |
| Service | `kubevirt-console-plugin-service`, `kubevirt-apiserver-proxy-service` |
| ConfigMap | `nginx-conf`, `kubevirt-user-settings`, `kubevirt-ui-features` |
| ConsolePlugin | `kubevirt-plugin` |
| Role / RoleBinding | `kubevirt-ui-config-reader` (authenticated users → UI configmaps) |
| NetworkPolicy | `kubevirt-console-plugin-np`, `kubevirt-apiserver-proxy-np` (when `networkPolicy.enabled`, **on by default**) |

A Route for apiserver-proxy is **on by default** (`route.enabled`); disable it if you only use in-cluster ConsolePlugin proxy. Post-install console registration: `console.registerPlugin`.

### TLS and NetworkPolicy

**TLS does not depend on NetworkPolicy.** Encryption and serving certificates are configured separately:

- **Plugin and proxy Services** use the OpenShift annotation `service.beta.openshift.io/serving-cert-secret-name`. The serving-cert controller writes TLS material into Secrets; pods mount them (nginx on the plugin uses `/var/serving-cert`).
- **Optional `proxy.tlsArgs`** pass TLS protocol/cipher flags to the proxy container (HCO aligns these with the cluster TLS profile).
- **Route** (`route.enabled`) adds edge TLS when exposing apiserver-proxy outside the cluster.

**NetworkPolicy** is L3/L4 segmentation only: who may connect to the plugin/proxy pods and where the proxy may egress (DNS, Kubernetes API). It does not issue or terminate TLS. Defaults match HCO when network policies are deployed (`networkPolicy.enabled: true`). Set `networkPolicy.enabled: false` on clusters where these rules block legitimate console or API traffic (for example non-standard console labels/namespaces).

| Policy | Effect (summary) |
|--------|------------------|
| `kubevirt-console-plugin-np` | Ingress to the plugin only from OpenShift Console UI pods in `openshift-console` on the plugin port |
| `kubevirt-apiserver-proxy-np` | Ingress on the proxy port; egress to DNS and apiserver (443/6443) |

## Prerequisites

- **OpenShift Console:** Console Operator available; `ConsolePlugin` CRD (`console.openshift.io`); console UI pods in `openshift-console` (required for default NetworkPolicy ingress). If `console.registerPlugin` is enabled, permissions to update `consoles.operator.openshift.io/cluster` `spec.plugins`.
- Namespace (default **`kubevirt-ui`**, same as operator `OPERAND_NAMESPACE`) — set `createNamespace: true` or create it before install
- **KubeVirt control plane** reachable via the Kubernetes API (this chart does not install it)
- Avoid duplicate installs: if CNV/HCO already manages the UI in `openshift-cnv`, do not install a second stack there; use `--set namespace=openshift-cnv` only when replacing/migrating deliberately

Full list: [docs/ROADMAP.md#prerequisites-cluster-and-console](../../docs/ROADMAP.md#prerequisites-cluster-and-console).

## Install

Default (operator-parity namespace and labels; images on `quay.io/kubevirt-ui`):

```bash
helm upgrade --install kubevirt-ui ./helm/kubevirt-console-ui \
  --namespace kubevirt-ui \
  --create-namespace \
  --set plugin.image=quay.io/kubevirt-ui/kubevirt-plugin:<tag> \
  --set proxy.image=quay.io/kubevirt-ui/kubevirt-apiserver-proxy:<tag>
```

CNV namespace (HCO-aligned labels optional via `hcoLabels.enabled=true`):

```bash
helm upgrade --install kubevirt-ui ./helm/kubevirt-console-ui \
  --namespace openshift-cnv \
  --set namespace=openshift-cnv \
  --set createNamespace=false \
  --set hcoLabels.enabled=true \
  --set plugin.image=quay.io/kubevirt-ui/kubevirt-plugin:<tag> \
  --set proxy.image=quay.io/kubevirt-ui/kubevirt-apiserver-proxy:<tag>
```

Pin images to the same versions shipped with your CNV release (see HCO `deploy/images.env`: `KUBEVIRT_CONSOLE_PLUGIN_IMAGE`, `KUBEVIRT_CONSOLE_PROXY_IMAGE`).

### Common values

| Value | Purpose |
|-------|---------|
| `plugin.image` / `proxy.image` | Container images (`quay.io/kubevirt-ui/...`; pin tags in production) |
| `plugin.replicas` / `proxy.replicas` | HA replica count (HCO uses 2 on highly available infra) |
| `networkPolicy.enabled` | Network segmentation for plugin/proxy (**default `true`**); not required for TLS—set `false` if rules block your console traffic |
| `route.enabled` + `route.host` | Public Route for apiserver-proxy (local/off-cluster console) |
| `console.registerPlugin` | Helm hook patches `consoles.operator.openshift.io/cluster` |
| `hcoLabels.enabled` | HCO-style labels (`app: kubevirt-hyperconverged`); default **`false`** (operator-owned labels) |
| `config.features.hcoVersion` | Leave empty to omit from ConfigMap (operator sets from HCO when installed) |
| `priorityClass.create` | Install `kubevirt-cluster-critical` (required for UI-only clusters; set `false` if KubeVirt already created it) |

## Uninstall

```bash
helm uninstall kubevirt-ui -n kubevirt-ui
```

Remove `kubevirt-plugin` from the cluster console plugins list if no other installation provides it:

```bash
oc patch consoles.operator.openshift.io cluster --type=json \
  -p='[{"op": "remove", "path": "/spec/plugins/<index>"}]'
```

(Use `oc get consoles.operator.openshift.io cluster -o yaml` to find the plugin index.)
