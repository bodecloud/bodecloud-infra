{{/*
Expand the name of the chart.
*/}}
{{- define "riven.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "riven.fullname" -}}
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
{{- define "riven.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "riven.labels" -}}
helm.sh/chart: {{ include "riven.chart" . }}
{{ include "riven.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "riven.selectorLabels" -}}
app.kubernetes.io/name: {{ include "riven.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "riven.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "riven.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get the Riven image repository
*/}}
{{- define "riven.image" -}}
{{- if .Values.global.imageRegistry }}
{{- printf "%s/%s" .Values.global.imageRegistry .Values.image.repository }}
{{- else }}
{{- .Values.image.repository }}
{{- end }}
{{- end }}

{{/*
Get the PostgreSQL image repository
*/}}
{{- define "riven.postgresqlImage" -}}
{{- if .Values.global.imageRegistry }}
{{- printf "%s/%s" .Values.global.imageRegistry .Values.postgresql.image.repository }}
{{- else }}
{{- .Values.postgresql.image.repository }}
{{- end }}
{{- end }}

{{/*
Get the Bootstrap image repository
*/}}
{{- define "riven.bootstrapImage" -}}
{{- if .Values.global.imageRegistry }}
{{- printf "%s/%s" .Values.global.imageRegistry .Values.bootstrap.image.repository }}
{{- else }}
{{- .Values.bootstrap.image.repository }}
{{- end }}
{{- end }}

{{/*
Storage class for persistent volumes
*/}}
{{- define "riven.storageClass" -}}
{{- if .Values.global.storageClass }}
{{- .Values.global.storageClass }}
{{- end }}
{{- end }}

{{/*
Common volumes for riven pod
*/}}
{{- define "riven.volumes" -}}
{{- if .Values.persistence.backup.enabled }}
- name: backup
  persistentVolumeClaim:
    claimName: {{ include "riven.fullname" . }}-backup
{{- end }}
{{- if .Values.persistence.config.enabled }}
- name: config
  persistentVolumeClaim:
    claimName: {{ include "riven.fullname" . }}-config
{{- end }}
{{- if .Values.volumes.devShm.enabled }}
- name: dev-shm
  emptyDir:
    {{- if .Values.volumes.devShm.sizeLimit }}
    sizeLimit: {{ .Values.volumes.devShm.sizeLimit }}
    {{- end }}
{{- end }}
{{- if .Values.configMaps.elfbotRiven.enabled }}
- name: elfbot
  configMap:
    name: {{ include "riven.fullname" . }}-elfbot-riven
    defaultMode: 420
    optional: true
{{- end }}
{{- if .Values.persistence.logs.enabled }}
- name: logs
  persistentVolumeClaim:
    claimName: {{ include "riven.fullname" . }}-logs
{{- end }}
{{- if .Values.persistence.rclone.enabled }}
- name: rclone
  persistentVolumeClaim:
    claimName: {{ include "riven.fullname" . }}-rclone
{{- end }}
{{- if .Values.persistence.realdebridZurg.enabled }}
- name: rclonemountrealdebridzurg
  persistentVolumeClaim:
    claimName: {{ include "riven.fullname" . }}-realdebrid-zurg
{{- end }}
{{- if .Values.configMaps.rivenSetup.enabled }}
- name: setup
  configMap:
    name: {{ include "riven.fullname" . }}-riven-setup
    defaultMode: 420
{{- end }}
{{- if .Values.persistence.symlinks.enabled }}
- name: symlinks
  persistentVolumeClaim:
    claimName: {{ include "riven.fullname" . }}-symlinks
{{- end }}
{{- if .Values.volumes.tmp.enabled }}
- name: tmp
  emptyDir:
    {{- if .Values.volumes.tmp.sizeLimit }}
    sizeLimit: {{ .Values.volumes.tmp.sizeLimit }}
    {{- end }}
{{- end }}
{{- end }}

{{/*
Environment variables from configmaps
*/}}
{{- define "riven.envFrom" -}}
{{- if .Values.envFrom }}
{{- toYaml .Values.envFrom }}
{{- end }}
{{- end }} 