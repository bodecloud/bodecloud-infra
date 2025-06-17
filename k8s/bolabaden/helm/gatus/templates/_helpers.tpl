{{/*
Expand the name of the chart.
*/}}
{{- define "gatus.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "gatus.fullname" -}}
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
{{- define "gatus.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "gatus.labels" -}}
helm.sh/chart: {{ include "gatus.chart" . }}
{{ include "gatus.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "gatus.selectorLabels" -}}
app.kubernetes.io/name: {{ include "gatus.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "gatus.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "gatus.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get the image repository
*/}}
{{- define "gatus.image" -}}
{{- if .Values.global.imageRegistry }}
{{- printf "%s/%s" .Values.global.imageRegistry .Values.image.repository }}
{{- else }}
{{- .Values.image.repository }}
{{- end }}
{{- end }}

{{/*
Get the full image reference including tag and digest
*/}}
{{- define "gatus.imageReference" -}}
{{- $image := include "gatus.image" . }}
{{- if .Values.image.digest }}
{{- printf "%s@%s" $image .Values.image.digest }}
{{- else }}
{{- printf "%s:%s" $image .Values.image.tag }}
{{- end }}
{{- end }}

{{/*
Common volume mounts for gatus container
*/}}
{{- define "gatus.volumeMounts" -}}
{{- if .Values.persistence.config.enabled }}
- name: config
  mountPath: /data/
  subPath: gatus
{{- end }}
{{- if .Values.configMap.enabled }}
- name: gatus-config
  mountPath: /config
{{- end }}
{{- end }}

{{/*
Common volumes for gatus pod
*/}}
{{- define "gatus.volumes" -}}
{{- if .Values.persistence.config.enabled }}
- name: config
  persistentVolumeClaim:
    claimName: {{ include "gatus.fullname" . }}-config
{{- end }}
{{- if .Values.configMap.enabled }}
- name: gatus-config
  configMap:
    name: {{ include "gatus.fullname" . }}-config
    defaultMode: 420
{{- end }}
{{- end }}

{{/*
Environment variables for gatus container
*/}}
{{- define "gatus.env" -}}
{{- range $key, $value := .Values.env }}
- name: {{ $key }}
  value: {{ $value | quote }}
{{- end }}
{{- end }}

{{/*
Environment variables from secrets/configmaps
*/}}
{{- define "gatus.envFrom" -}}
{{- if .Values.envFrom }}
{{- toYaml .Values.envFrom }}
{{- end }}
{{- if and .Values.smtp.enabled (not .Values.envFrom) }}
- secretRef:
    name: {{ include "gatus.fullname" . }}-smtp
{{- end }}
{{- end }}

{{/*
Storage class for persistent volumes
*/}}
{{- define "gatus.storageClass" -}}
{{- if .Values.global.storageClass }}
{{- .Values.global.storageClass }}
{{- else if .Values.persistence.config.storageClass }}
{{- .Values.persistence.config.storageClass }}
{{- end }}
{{- end }} 