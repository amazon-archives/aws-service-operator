{{/*
Helm standard labels.
*/}}
{{- define "aws-service-operator.helmStandardLabels" -}}
app.kubernetes.io/name: {{ include "aws-service-operator.name" . }}
helm.sh/chart: {{ include "aws-service-operator.chart" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
{{- end -}}

{{/*
Helm pod labels.
*/}}
{{- define "aws-service-operator.helmPodLabels" -}}
app.kubernetes.io/name: {{ include "aws-service-operator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "aws-service-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Expand the name of the chart.
*/}}
{{- define "aws-service-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "aws-service-operator.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create the name of the service account to use.
*/}}
{{- define "aws-service-operator.serviceAccountName" -}}
{{- if .Values.rbac.create -}}
    {{ default (include "aws-service-operator.fullname" .) .Values.rbac.serviceAccountName }}
{{- else -}}
    {{ default "default" .Values.rbac.serviceAccountName }}
{{- end -}}
{{- end -}}
