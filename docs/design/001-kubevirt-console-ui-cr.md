# RFC 001 — KubeVirtConsoleUI custom resource

**Status:** Draft  
**Authors:** kubevirt-ui team  
**Last updated:** 2026-03-24

## Summary

Introduce a namespaced custom resource that declares how the **KubeVirt console UI** (plugin, apiserver-proxy, config, `ConsolePlugin`, optional console registration) should run on a cluster. The operator reconciles cluster objects to match spec; status reports readiness for operators and the OpenShift console.

## Motivation

- Decouple UI lifecycle from HCO/CNV while keeping HCO-aligned resource names.
- Single supported day-2 interface: upgrade images, replicas, network policy, route, and console registration via CR.
- Enable OLM packaging and GitOps (Argo CD) without raw Helm as the only option.

## Prerequisites

Same as [ROADMAP — Prerequisites](../ROADMAP.md#prerequisites-cluster-and-console), including:

- Console Operator and `ConsolePlugin` CRD.
- Working console in `openshift-console` for validation.
- Cluster admin (or delegated RBAC) if `spec.console.registerPlugin` patches `consoles.operator.openshift.io/cluster`.

## API sketch (v1alpha1)

**Group:** `consoleui.kubevirt.io` (proposed)  
**Kind:** `KubeVirtConsoleUI`  
**Scope:** Cluster. A single instance named **`cluster`**. Operands deploy in **`OPERAND_NAMESPACE`** (operator env), not in CR metadata.

### Spec (align with Helm `values.yaml`)

| Field | Purpose |
|-------|---------|
| `plugin.image`, `plugin.replicas`, … | Console plugin deployment |
| `proxy.image`, `proxy.replicas`, `proxy.tlsArgs` | Apiserver proxy |
| `config.features` | `kubevirt-ui-features` data |
| `networkPolicy.enabled` | Plugin/proxy NetworkPolicies |
| `route.enabled`, `route.host` | Apiserver-proxy Route |
| `console.registerPlugin` | Patch Console cluster config `spec.plugins` |
| `priorityClass.create` | `kubevirt-ui-cluster-critical` when needed |
| `hcoLabels.enabled` | HCO-aligned labels for migration |

### Status

| Field | Purpose |
|-------|---------|
| `conditions[]` | `Ready`, `Progressing`, `Degraded` |
| `observedGeneration` | Spec generation applied |
| `pluginImage`, `proxyImage` | Running image digests or IDs if available |
| `consolePluginRegistered` | Whether cluster console lists `kubevirt-plugin` |

## Reconcile order (proposed)

1. PriorityClass (if `priorityClass.create`) — cluster-scoped, not owner-referenced
2. ServiceAccounts
3. ConfigMaps (`nginx-conf`, user settings, features)
4. Deployments (plugin, proxy)
5. Services (serving-cert annotations)
6. NetworkPolicies (if enabled)
7. Route (if enabled)
8. `ConsolePlugin` (`console.openshift.io`)
9. Console operator plugin registration (if enabled)
10. RBAC (config reader Role/RoleBinding)

## Open questions

1. Single cluster-wide CR vs one CR per namespace? **v1alpha1:** one CR per operand namespace; CR must be installed in that namespace.
2. Finalizer to remove console plugin entry on CR delete? **Not in v1alpha1** — `ConsolePlugin`, `PriorityClass`, and console `spec.plugins` are not garbage-collected on CR delete.
3. Validation webhook for conflicting HCO-managed UI in same namespace?
4. How to surface cluster TLS profile on proxy (`tlsArgs`) without coupling to HCO?

## Acceptance criteria (M1)

- [ ] CRD installed via `make install`
- [ ] Sample CR in `config/samples/` applies without Helm
- [ ] `ConsolePlugin` `kubevirt-plugin` exists
- [ ] With `console.registerPlugin: true`, plugin appears in console navigation (on a cluster meeting prerequisites)

## References

- Helm chart: `helm/kubevirt-console-ui/values.yaml`
- HCO: `controllers/handlers/kubevirtConsolePlugin.go`
- OpenShift: [ConsolePlugin](https://docs.openshift.com/container-platform/latest/web_console/web-console.html#dynamic-plugins-about_dynamic-plugins), Console Operator API
