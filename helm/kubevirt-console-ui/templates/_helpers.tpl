{{- define "kubevirt-console-ui.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "kubevirt-console-ui.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{- define "kubevirt-console-ui.namespace" -}}
{{- .Values.namespace }}
{{- end }}

{{- define "kubevirt-console-ui.baseLabels" -}}
helm.sh/chart: {{ include "kubevirt-console-ui.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .Values.hcoLabels.enabled }}
app: {{ .Values.hcoLabels.app }}
app.kubernetes.io/managed-by: {{ .Values.hcoLabels.managedBy }}
app.kubernetes.io/part-of: {{ .Values.hcoLabels.partOf }}
{{- with .Values.hcoLabels.version }}
app.kubernetes.io/version: {{ . | quote }}
{{- end }}
{{- else }}
app: kubevirt-console-ui
app.kubernetes.io/managed-by: kubevirt-ui-operator
app.kubernetes.io/part-of: kubevirt-ui
{{- end }}
{{- end }}

{{- define "kubevirt-console-ui.componentLabels" -}}
{{- $component := .component }}
{{- include "kubevirt-console-ui.baseLabels" .root }}
app.kubernetes.io/component: {{ $component }}
{{- end }}

{{- define "kubevirt-console-ui.pluginSelectorLabels" -}}
app.kubernetes.io/component: kubevirt-console-plugin
{{- if .Values.hcoLabels.enabled }}
app: {{ .Values.hcoLabels.app }}
{{- end }}
{{- end }}

{{- define "kubevirt-console-ui.proxySelectorLabels" -}}
app.kubernetes.io/component: kubevirt-apiserver-proxy
{{- if .Values.hcoLabels.enabled }}
app: {{ .Values.hcoLabels.app }}
{{- end }}
{{- end }}
