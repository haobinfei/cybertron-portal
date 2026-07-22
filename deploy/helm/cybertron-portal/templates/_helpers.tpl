{{/*
Expand the name of the chart.
*/}}
{{- define "cybertron-portal.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cybertron-portal.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cybertron-portal.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cybertron-portal.labels" -}}
helm.sh/chart: {{ include "cybertron-portal.chart" . }}
{{ include "cybertron-portal.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cybertron-portal.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cybertron-portal.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the image pull secret to use
*/}}
{{- define "cybertron-portal.imagePullSecrets" -}}
{{- with .Values.imagePullSecrets }}
imagePullSecrets:
  {{- toYaml . | nindent 2 }}
{{- end }}
{{- end }}

{{/*
Backend fullname
*/}}
{{- define "cybertron-portal.backend.fullname" -}}
{{- printf "%s-backend" (include "cybertron-portal.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Frontend fullname
*/}}
{{- define "cybertron-portal.frontend.fullname" -}}
{{- printf "%s-frontend" (include "cybertron-portal.fullname" .) | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Generate base64-encoded docker config json for image pull secret.
*/}}
{{- define "cybertron-portal.dockerconfigjson" -}}
{{- $auth := printf "%s:%s" .Values.registry.username .Values.registry.password | b64enc -}}
{{- $config := dict "auths" (dict .Values.registry.server (dict "username" .Values.registry.username "password" .Values.registry.password "email" .Values.registry.email "auth" $auth)) -}}
{{- $config | toJson | b64enc -}}
{{- end }}

