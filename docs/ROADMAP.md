# KubeVirt UI Operator — roadmap and delivery plan

This document is the team plan for evolving **kubevirt-ui-operator** from a Helm-only install to a **production-grade operator**, while keeping the same squad practices as [kubevirt-plugin](https://github.com/kubevirt-ui/kubevirt-plugin).

## North star

| Goal | Definition of done |
|------|-------------------|
| **Functional** | Operator (and optional OLM) reconciles the same UI stack as the Helm chart and HCO’s [`kubevirtConsolePlugin`](https://github.com/kubevirt/hyperconverged-cluster-operator/blob/main/controllers/handlers/kubevirtConsolePlugin.go) handlers |
| **Team** | Any plugin-team engineer can clone, run locally, open a PR, and get CI signal without a private cluster for basic changes |
| **Production** | Upgrades, failure visibility, security/RBAC review, release artifacts, and a documented compatibility matrix with plugin/proxy images |

**Non-goals:** KubeVirt, CDI, CNAO, HCO, subscriptions, or `HyperConverged` — **UI operands only**.

## Prerequisites (cluster and console)

Install and development assume an **OpenShift** (or OpenShift-compatible) cluster with a working **web console**. The UI stack plugs into the console; it does not replace it.

### Console platform

- **OpenShift Console Operator** installed and healthy (`consoles.operator.openshift.io/cluster` exists and reports available).
- **`ConsolePlugin` CRD** (`console.openshift.io/v1`) available on the cluster API.
- **Console UI workloads** running in `openshift-console` (pods labeled for the console UI) so the plugin backend is reachable and NetworkPolicy ingress rules can target them.
- **Authenticated console access** for smoke testing (users need RBAC to open Virtualization views once the plugin is registered).
- **Plugin registration surface:** ability to patch `consoles.operator.openshift.io/cluster` `spec.plugins` when `console.registerPlugin` is enabled (Helm hook today; operator equivalent later). Requires cluster-admin or equivalent for that API.

### Virtualization API (not installed by this repo)

- **KubeVirt API** (or CNV-managed KubeVirt) reachable with appropriate user RBAC.
- Operand namespace (default **`openshift-cnv`**) available; do not duplicate HCO-managed UI in the same namespace.

### Images and TLS

- Published **plugin** and **apiserver-proxy** images (pin tags; avoid `latest` in production). Reference versions: HCO [`deploy/images.env`](https://github.com/kubevirt/hyperconverged-cluster-operator/blob/main/deploy/images.env) (`KUBEVIRT_CONSOLE_PLUGIN_IMAGE`, `KUBEVIRT_CONSOLE_PROXY_IMAGE`).
- **Serving certificates:** OpenShift service serving-cert controller (Services annotated with `service.beta.openshift.io/serving-cert-secret-name`).
- **Optional Route:** cluster ingress / apps domain for off-cluster apiserver-proxy access (`route.enabled`).

### Tooling (developers)

- `kubectl` / `oc`, Go 1.22+ (when operator lands), Helm 3, `make`, container build tooling for operator image.
- For local controller development: cluster with CRD install permissions (`make install`) or a shared dev cluster.

See also: [Helm chart prerequisites](../helm/kubevirt-console-ui/README.md#prerequisites) and [TLS vs NetworkPolicy](../helm/kubevirt-console-ui/README.md#tls-and-networkpolicy).

## Current baseline

| In tree | Not in tree yet |
|---------|-----------------|
| Helm chart `helm/kubevirt-console-ui/` | Go operator, CI, OLM bundle |
| Root and chart README | CONTRIBUTING, OWNERS, design RFCs, e2e |

**In-org references:**

- [kube-gateway-operator](https://github.com/kubevirt-ui/kube-gateway-operator) — Kubebuilder v3 layout (`api/`, `controllers/`, `config/`, `Makefile`, `deploy/`).
- [kubevirt-apiserver-proxy](https://github.com/kubevirt-ui/kubevirt-apiserver-proxy) — simple GHA (`go build`, `go test`).
- [kubevirt-plugin](https://github.com/kubevirt-ui/kubevirt-plugin) — team ownership and review culture.

## Target architecture

```text
Install surfaces                Cluster runtime
─────────────────              ───────────────────────────────────
Helm chart (supported)    →    KubeVirtConsoleUI CR
OLM CSV (later)           →         ↓
config/deploy YAML (dev)         kubevirt-ui-operator
                                      ↓
                    Deployments, Services, ConfigMaps,
                    ConsolePlugin, RBAC, NetworkPolicy, Route,
                    optional Console operator plugin registration
```

**Reconcile implementation (recommended):** Go operand builders (HCO-style), shared semantics with the Helm chart so values and CR spec do not drift. Helm remains a supported bootstrap path until the operator is the default.

**Authoritative behavior reference:** HCO `controllers/handlers/kubevirtConsolePlugin.go`.

## Phased delivery

### Phase 0 — Repo foundations (≈1–2 weeks)

Unblock parallel contributions before the controller is large.

| Work item | Deliverable |
|-----------|-------------|
| Governance | `OWNERS` (kubevirt-plugin roster), `CONTRIBUTING.md`, `SECURITY.md`, PR and issue templates |
| CI (Helm) | GHA: `helm lint`, `helm template`, optional manifest schema check |
| CI (Go) | Scaffold `go.mod` / `make test`; expand when operator merges |
| Docs | `docs/ARCHITECTURE.md`, `docs/DEVELOPMENT.md`; link from root README |
| Design | [RFC 001 — KubeVirtConsoleUI CR](design/001-kubevirt-console-ui-cr.md) |
| Backlog | Labels (`area/helm`, `area/operator`, `area/olm`), milestones M0–M3 |

**Exit:** Required checks on `main`, team can land chart-only PRs with CI.

### Phase 1 — Scaffold + minimal reconcile (≈2–4 weeks)

| Work item | Notes |
|-----------|--------|
| Kubebuilder v3 init | Align with kube-gateway-operator |
| API | e.g. `consoleui.kubevirt.io/v1alpha1` `KubeVirtConsoleUI` (name in RFC) |
| Spec | Mirror Helm `values.yaml`: images, replicas, features, network policy, route, console registration, priority class |
| Status | Conditions: `Ready`, `Progressing`, `Degraded`; observed image versions |
| Reconcile v0 | SAs → ConfigMaps → Deployments → Services → `ConsolePlugin` → RBAC |
| Console integration | Reconcile `ConsolePlugin`; optional cluster Console `spec.plugins` when enabled |
| Samples | `config/samples/` + `make run` / `make deploy` docs |

**Parallel streams:** plugin+nginx CM; proxy+route+NP; ConsolePlugin+registration; CRD+validation.

**Exit:** Sample CR → ready pods, `ConsolePlugin` present, Virtualization visible in console on a test cluster.

### Phase 2 — Helm ↔ operator parity (≈2–3 weeks)

| Work item | Notes |
|-----------|--------|
| Parity matrix | [docs/parity.md](parity.md) — every chart template ↔ controller operand |
| Migration | Runbook: Helm uninstall → operator CR (no double ownership) |
| Upgrade | Spec changes roll workloads; document immutable fields |
| E2E (light) | Optional CI on OpenShift: apply CR, assert Deployments + `ConsolePlugin` |

**Exit:** Parity matrix complete; documented migration.

### Phase 3 — Production hardening (≈3–5 weeks)

| Area | Expectations |
|------|----------------|
| Reliability | Leader election, sensible requeue; finalizers only if console deregistration on delete is required |
| Security | Minimal operator RBAC; NetworkPolicy defaults documented; SCC `restricted-v2` |
| TLS | Serving certs; optional proxy `tlsArgs` aligned with cluster TLS profile |
| Observability | Structured logs, controller metrics; runbook for common failures |
| Release | `quay.io/kubevirtui/kubevirt-ui-operator` semver tags, changelog |
| Compatibility | [docs/COMPATIBILITY.md](COMPATIBILITY.md) — OpenShift versions, plugin/proxy pins |

**Exit:** `docs/RUNBOOK.md` covers console plugin not listed, NP blocking console traffic, serving cert not ready.

### Phase 4 — OLM (optional, after M2)

Operator SDK bundle/CSV, `relatedImages`, subscription docs. Stays **UI-only**; no CNV backend operands.

## Milestones

| Milestone | Target | Outcome |
|-----------|--------|---------|
| **M0** | ~2 weeks | CI, contributing docs, RFC merged, labeled backlog |
| **M1** | ~4–6 weeks | Scaffold + CR + core operands + `ConsolePlugin` |
| **M2** | ~8–10 weeks | Helm parity + migration + basic e2e |
| **M3** | ~12–14 weeks | Hardening, releases, compatibility matrix, runbook |
| **M4** | Optional | OLM bundle |

## Team operating model

- **Weekly operator sync** (30m): HCO drift, PR splits, review load.
- **RFC before API changes** — CRD is a public contract.
- **Release train:** operator tag + default image bumps with plugin/proxy releases; update compatibility doc.
- **HCO drift watch:** changes to `kubevirtConsolePlugin.go` → tracking issue in this repo.

```text
kubevirt-plugin release ──┐
kubevirt-apiserver-proxy ─┼──► kubevirt-ui-operator release + COMPATIBILITY.md
HCO images.env (reference)┘
```

## Target repository layout

```text
kubevirt-ui-operator/
├── api/v1alpha1/
├── controllers/
├── internal/                  # Shared manifest builders (Helm parity)
├── config/                    # manager, rbac, crd, samples
├── deploy/
├── helm/kubevirt-console-ui/
├── hack/
├── docs/
│   ├── ROADMAP.md             # this file
│   ├── ARCHITECTURE.md
│   ├── DEVELOPMENT.md
│   ├── COMPATIBILITY.md
│   ├── RUNBOOK.md
│   ├── parity.md
│   └── design/
├── .github/workflows/
├── OWNERS
├── CONTRIBUTING.md
├── Makefile
├── PROJECT
└── main.go
```

## Decisions (RFC 001)

1. CRD name and API group.
2. Operator namespace vs operand namespace (`openshift-cnv` default).
3. Console registration: default on vs opt-in.
4. Helm long-term: co-maintained vs deprecated after OLM.
5. Production image pins (no `latest` in samples).

## Related documents

- [Design RFC 001 — KubeVirtConsoleUI CR](design/001-kubevirt-console-ui-cr.md)
- [Helm chart](../helm/kubevirt-console-ui/README.md)
- [Root README](../README.md)
