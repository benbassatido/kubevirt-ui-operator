# Helm chart ↔ operator parity matrix

Track reconciliation parity between `helm/kubevirt-console-ui` templates and the operator controller.

| Helm template | Resource(s) | Operator | Notes |
|---------------|-------------|----------|-------|
| `namespace.yaml` | Namespace | No | Create namespace before applying the CR (`scripts/apply-sample-cr.sh`); operator does not create namespaces |
| `priorityclass.yaml` | PriorityClass | Yes | Cluster-scoped, no owner ref |
| `serviceaccount-plugin.yaml` | ServiceAccount | Yes | |
| `serviceaccount-proxy.yaml` | ServiceAccount | Yes | |
| `configmap-nginx.yaml` | ConfigMap | Yes | `nginx-conf` |
| `configmap-user-settings.yaml` | ConfigMap | Yes | |
| `configmap-ui-features.yaml` | ConfigMap | Yes | |
| `deployment-plugin.yaml` | Deployment | Yes | SCC `restricted-v2` |
| `deployment-proxy.yaml` | Deployment | Yes | |
| `service-plugin.yaml` | Service | Yes | Serving cert annotation |
| `service-proxy.yaml` | Service | Yes | Serving cert annotation |
| `networkpolicy-plugin.yaml` | NetworkPolicy | Yes | When `networkPolicy.enabled` |
| `networkpolicy-proxy.yaml` | NetworkPolicy | Yes | When `networkPolicy.enabled` |
| `route-proxy.yaml` | Route | Yes | When `route.enabled` |
| `consoleplugin.yaml` | ConsolePlugin | Yes | Cluster-scoped |
| `console-register-plugin.yaml` | Console cluster patch | Partial | In-process patch (no Job/ClusterRole) |
| `role-config-reader.yaml` | Role | Yes | |
| `rolebinding-config-reader.yaml` | RoleBinding | Yes | |

**Legend:** Partial → behavior differs from Helm or needs hardening (RBAC for registration, status checks).
