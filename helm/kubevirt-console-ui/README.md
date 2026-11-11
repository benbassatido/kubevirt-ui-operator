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
- Operand namespace **`kubevirt-ui`** by default (same as operator **`OPERAND_NAMESPACE`**). The chart creates it when **`createNamespace: true`** (default).
- **KubeVirt control plane** reachable via the Kubernetes API (this chart does not install it)
- Do not install a second UI stack in **`openshift-cnv`** while CNV/HCO still owns the console plugin there. This chart (and the operator) target **`kubevirt-ui`**; **`openshift-cnv`** is only relevant to the **operator** via **`HCO_NAMESPACE`** for optional **`hcoVersion`** on the features ConfigMap when HyperConverged exists.

Full list: [docs/ROADMAP.md#prerequisites-cluster-and-console](../../docs/ROADMAP.md#prerequisites-cluster-and-console).

## Install

Default (operand namespace **`kubevirt-ui`**, operator-parity labels; images on `quay.io/kubevirt-ui`):

```bash
helm upgrade --install kubevirt-ui ./helm/kubevirt-console-ui \
  --namespace kubevirt-ui \
  --set plugin.image=quay.io/kubevirt-ui/kubevirt-plugin:<tag> \
  --set proxy.image=quay.io/kubevirt-ui/kubevirt-apiserver-proxy:<tag>
```

For the **operator** path instead of Helm, install the cluster-scoped **`KubeVirtConsoleUI`** named **`cluster`** after deploying the manager—see the root [README](../../README.md#quick-start-operator).

Pin images to the same versions shipped with your CNV release (see HCO `deploy/images.env`: `KUBEVIRT_CONSOLE_PLUGIN_IMAGE`, `KUBEVIRT_CONSOLE_PROXY_IMAGE`).

### Common values

| Value | Purpose |
|-------|---------|
| `namespace` | Operand namespace for Deployments, Services, ConfigMaps (**default `kubevirt-ui`**) |
| `createNamespace` | Render a labeled **`Namespace`** for operands (**default `true`**) |
| `plugin.image` / `proxy.image` | Container images (`quay.io/kubevirt-ui/...`; pin tags in production) |
| `plugin.replicas` / `proxy.replicas` | HA replica count (HCO uses 2 on highly available infra) |
| `networkPolicy.enabled` | Network segmentation for plugin/proxy (**default `true`**); not required for TLS—set `false` if rules block your console traffic |
| `route.enabled` + `route.host` | Public Route for apiserver-proxy (local/off-cluster console) |
| `console.registerPlugin` | Helm hook patches `consoles.operator.openshift.io/cluster` |
| `hcoLabels.enabled` | HCO-style labels (`app: kubevirt-hyperconverged`); default **`false`** (operator-owned labels) |
| `config.features.hcoVersion` | Leave empty in Helm; the operator sets **`hcoVersion`** from **`HCO_NAMESPACE`** (default **`openshift-cnv`**) when HyperConverged is installed |
| `priorityClass.create` | Install `kubevirt-ui-cluster-critical` (operator-owned; set `false` if you manage the class elsewhere) |

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
