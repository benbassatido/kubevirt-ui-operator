#!/usr/bin/env bash
# Create the namespace required by the sample CR, then apply it.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SAMPLE="${ROOT}/config/samples/consoleui_v1alpha1_kubevirtconsoleui.yaml"

if [[ ! -f "${SAMPLE}" ]]; then
  echo "sample not found: ${SAMPLE}" >&2
  exit 1
fi

namespace="$(kubectl create -f "${SAMPLE}" --dry-run=client -o jsonpath='{.metadata.namespace}')"

if [[ -z "${namespace}" ]]; then
  echo "could not read metadata.namespace from ${SAMPLE}" >&2
  exit 1
fi

kubectl create namespace "${namespace}" --dry-run=client -o yaml | kubectl apply -f -

kubectl apply -f "${SAMPLE}"
echo "Applied ${SAMPLE} in namespace ${namespace}"
