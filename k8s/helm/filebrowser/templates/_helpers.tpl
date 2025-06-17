{{/*
Expand the name of the chart.
*/}}
{{- define "filebrowser.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "filebrowser.fullname" -}}
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
{{- define "filebrowser.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "filebrowser.labels" -}}
helm.sh/chart: {{ include "filebrowser.chart" . }}
{{ include "filebrowser.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "filebrowser.selectorLabels" -}}
app.kubernetes.io/name: {{ include "filebrowser.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "filebrowser.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "filebrowser.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Get the FileBrowser image repository
*/}}
{{- define "filebrowser.image" -}}
{{- if .Values.global.imageRegistry }}
{{- printf "%s/%s" .Values.global.imageRegistry .Values.image.repository }}
{{- else }}
{{- .Values.image.repository }}
{{- end }}
{{- end }}

{{/*
Get the Ubuntu image repository
*/}}
{{- define "filebrowser.ubuntuImage" -}}
{{- if .Values.global.imageRegistry }}
{{- printf "%s/%s" .Values.global.imageRegistry .Values.ubuntu.image.repository }}
{{- else }}
{{- .Values.ubuntu.image.repository }}
{{- end }}
{{- end }}

{{/*
Storage class for persistent volumes
*/}}
{{- define "filebrowser.storageClass" -}}
{{- if .Values.global.storageClass }}
{{- .Values.global.storageClass }}
{{- end }}
{{- end }}

{{/*
Common volumes for filebrowser pod
*/}}
{{- define "filebrowser.volumes" -}}
{{- if .Values.persistence.backup.enabled }}
- name: backup
  persistentVolumeClaim:
    claimName: {{ include "filebrowser.fullname" . }}-backup
{{- end }}
{{- if .Values.persistence.config.enabled }}
- name: config
  persistentVolumeClaim:
    claimName: {{ include "filebrowser.fullname" . }}-config
{{- end }}
{{- if .Values.volumes.dummyStorage.enabled }}
- name: dummy-storage
  emptyDir:
    {{- if .Values.volumes.dummyStorage.sizeLimit }}
    sizeLimit: {{ .Values.volumes.dummyStorage.sizeLimit }}
    {{- end }}
{{- end }}
{{- if .Values.volumes.elfbot.enabled }}
- name: elfbot
  emptyDir:
    {{- if .Values.volumes.elfbot.sizeLimit }}
    sizeLimit: {{ .Values.volumes.elfbot.sizeLimit }}
    {{- end }}
{{- end }}
{{- if .Values.configMaps.filebrowserElfbotScript.enabled }}
- name: elfbot-script
  configMap:
    name: {{ include "filebrowser.fullname" . }}-elfbot-script
    defaultMode: 493
- name: elfbot-script-ucfirst
  configMap:
    name: {{ include "filebrowser.fullname" . }}-elfbot-script
    defaultMode: 493
{{- end }}
{{- if .Values.volumes.elftermState.enabled }}
- name: elfterm-state
  emptyDir:
    {{- if .Values.volumes.elftermState.sizeLimit }}
    sizeLimit: {{ .Values.volumes.elftermState.sizeLimit }}
    {{- end }}
{{- end }}
{{- if .Values.persistence.logs.enabled }}
- name: logs
  persistentVolumeClaim:
    claimName: {{ include "filebrowser.fullname" . }}-logs
{{- end }}
{{- if .Values.persistence.rclone.enabled }}
- name: rclone
  persistentVolumeClaim:
    claimName: {{ include "filebrowser.fullname" . }}-rclone
{{- end }}
{{- if .Values.persistence.realdebridZurg.enabled }}
- name: rclonemountrealdebridzurg
  persistentVolumeClaim:
    claimName: {{ include "filebrowser.fullname" . }}-realdebrid-zurg
{{- end }}
{{- if .Values.configMaps.recyclarrConfig.enabled }}
- name: recyclarr-config
  configMap:
    name: {{ include "filebrowser.fullname" . }}-recyclarr-config
    defaultMode: 420
{{- end }}
{{- if .Values.persistence.symlinks.enabled }}
- name: symlinks
  persistentVolumeClaim:
    claimName: {{ include "filebrowser.fullname" . }}-symlinks
{{- end }}
{{- if .Values.volumes.tmp.enabled }}
- name: tmp
  emptyDir: {}
{{- end }}
{{- end }}

{{/*
FileBrowser commands as comma-separated string
*/}}
{{- define "filebrowser.commands" -}}
{{- join "," .Values.filebrowser.commands }}
{{- end }} 