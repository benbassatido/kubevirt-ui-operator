#!/usr/bin/env bash
# Create the operand namespace (OPERAND_NAMESPACE) and apply the cluster-scoped sample CR.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SAMPLE="${ROOT}/config/samples/consoleui_v1alpha1_kubevirtconsoleui.yaml"

if [[ ! -f "${SAMPLE}" ]]; then
  echo "sample not found: ${SAMPLE}" >&2
  exit 1
fi

namespace="${OPERAND_NAMESPACE:-kubevirt-ui}"

kubectl create namespace "${namespace}" --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f "${SAMPLE}"
echo "Applied ${SAMPLE} (cluster-scoped); operands use namespace ${namespace} (OPERAND_NAMESPACE)"
