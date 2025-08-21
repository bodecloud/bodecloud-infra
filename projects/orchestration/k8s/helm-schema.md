# Helm Chart Schema & Syntax Reference

## Overview

Helm is the package manager for Kubernetes. It allows you to define, install, and upgrade even the most complex Kubernetes applications. This document provides comprehensive schema and syntax information for Helm charts.

## Chart Structure

### Directory Structure

```shell
mychart/
  Chart.yaml          # Chart metadata
  values.yaml         # Default configuration values
  values.schema.json  # JSON schema for values validation
  charts/             # Directory containing subcharts
  templates/          # Directory containing template files
  templates/NOTES.txt # Usage notes displayed after installation
```

### Chart.yaml

```yaml
apiVersion: v2
name: mychart
description: A Helm chart for Kubernetes
type: application
version: 0.1.0
appVersion: "1.16.0"
keywords:
  - myapp
  - web
home: https://github.com/example/mychart
sources:
  - https://github.com/example/mychart
maintainers:
  - name: maintainer
    email: maintainer@example.com
    url: https://github.com/maintainer
icon: https://example.com/icon.png
deprecated: false
annotations:
  category: Infrastructure
kubeVersion: ">=1.19.0-0"
```

## Values Configuration

### Basic values.yaml

```yaml
# Default values for mychart
replicaCount: 1

image:
  repository: nginx
  pullPolicy: IfNotPresent
  tag: ""

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""

serviceAccount:
  create: true
  annotations: {}
  name: ""

rbac:
  create: true

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  className: ""
  annotations: {}
  hosts:
    - host: chart-example.local
      paths:
        - path: /
          pathType: ImplementationSpecific
  tls: []

resources: {}
nodeSelector: {}
tolerations: []
affinity: {}
```

### Values Schema (values.schema.json)

```json
{
  "$schema": "https://json-schema.org/draft-07/schema#",
  "properties": {
    "image": {
      "description": "Container Image",
      "properties": {
        "repo": {
          "type": "string"
        },
        "tag": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "name": {
      "description": "Service name",
      "type": "string"
    },
    "port": {
      "description": "Port",
      "minimum": 0,
      "type": "integer"
    },
    "protocol": {
      "type": "string"
    }
  },
  "required": [
    "protocol",
    "port"
  ],
  "title": "Values",
  "type": "object"
}
```

### Nested Values Structure

```yaml
# Flat structure (preferred)
serverName: nginx
serverPort: 80

# Nested structure
server:
  name: nginx
  port: 80

# Map of maps (easy to override with --set)
servers:
  foo:
    port: 80
  bar:
    port: 81
```

## Templates

### Basic Template Syntax

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name }}-configmap
data:
  myvalue: "Hello World"
  drink: {{ .Values.favoriteDrink }}
  food: {{ .Values.favoriteFood }}
```

### Built-in Objects

```yaml
# Release information
{{ .Release.Name }}      # Release name
{{ .Release.Namespace }} # Release namespace
{{ .Release.Service }}   # Service conducting the release
{{ .Release.IsUpgrade }} # True if this is an upgrade
{{ .Release.IsInstall }} # True if this is an install

# Chart information
{{ .Chart.Name }}        # Chart name
{{ .Chart.Version }}     # Chart version
{{ .Chart.Description }} # Chart description

# Values
{{ .Values.image.repository }} # Values from values.yaml
{{ .Values.replicaCount }}     # Values from values.yaml

# Files
{{ .Files.Get "config.yaml" }}           # Get file contents
{{ .Files.GetBytes "config.yaml" }}      # Get file as bytes
{{ .Files.Glob "configs/*.yaml" }}       # Glob files
{{ .Files.Lines "config.yaml" }}         # Get file as lines
```

### Template Functions

```yaml
# String functions
{{ .Values.name | quote }}                    # Quote string
{{ .Values.name | upper }}                    # Convert to uppercase
{{ .Values.name | lower }}                    # Convert to lowercase
{{ .Values.name | title }}                    # Title case
{{ .Values.name | trunc 63 }}                 # Truncate to 63 chars
{{ .Values.name | trimSuffix "-" }}           # Trim suffix
{{ .Values.name | trimPrefix "-" }}           # Trim prefix

# List functions
{{ .Values.items | join ", " }}               # Join list with separator
{{ .Values.items | first }}                   # Get first item
{{ .Values.items | last }}                    # Get last item
{{ .Values.items | sortAlpha }}               # Sort alphabetically
{{ .Values.items | uniq }}                    # Remove duplicates

# Math functions
{{ add .Values.port 1 }}                      # Add numbers
{{ sub .Values.port 1 }}                      # Subtract numbers
{{ mul .Values.port 2 }}                      # Multiply numbers
{{ div .Values.port 2 }}                      # Divide numbers
{{ mod .Values.port 2 }}                      # Modulo

# Default values
{{ .Values.port | default 80 }}               # Default value
{{ .Values.enabled | default false }}         # Default boolean
{{ .Values.image | default "nginx:latest" }}  # Default string
```

### Conditional Logic

```yaml
{{- if .Values.ingress.enabled }}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Release.Name }}-ingress
{{- end }}

{{- if eq .Values.environment "production" }}
  replicas: 3
{{- else }}
  replicas: 1
{{- end }}

{{- if and .Values.ingress.enabled .Values.ingress.hosts }}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{ .Release.Name }}-ingress
spec:
  rules:
  {{- range .Values.ingress.hosts }}
  - host: {{ .host }}
    http:
      paths:
      {{- range .paths }}
      - path: {{ .path }}
        pathType: {{ .pathType }}
        backend:
          service:
            name: {{ $.Release.Name }}-service
            port:
              number: {{ $.Values.service.port }}
      {{- end }}
  {{- end }}
{{- end }}
```

### Loops and Ranges

```yaml
# Range over list
{{- range .Values.items }}
- name: {{ .name }}
  value: {{ .value }}
{{- end }}

# Range over map
{{- range $key, $value := .Values.config }}
{{ $key }}: {{ $value }}
{{- end }}

# Range with index
{{- range $index, $item := .Values.items }}
{{ $index }}: {{ $item }}
{{- end }}
```

## Template Files

### Deployment Template

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "mychart.fullname" . }}
  labels:
    {{- include "mychart.labels" . | nindent 4 }}
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: {{ .Values.replicaCount }}
  {{- end }}
  selector:
    matchLabels:
      {{- include "mychart.selectorLabels" . | nindent 6 }}
  template:
    metadata:
      labels:
        {{- include "mychart.selectorLabels" . | nindent 8 }}
    spec:
      {{- with .Values.imagePullSecrets }}
      imagePullSecrets:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      serviceAccountName: {{ include "mychart.serviceAccountName" . }}
      securityContext:
        {{- toYaml .Values.podSecurityContext | nindent 8 }}
      containers:
        - name: {{ .Chart.Name }}
          securityContext:
            {{- toYaml .Values.securityContext | nindent 12 }}
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag | default .Chart.AppVersion }}"
          imagePullPolicy: {{ .Values.image.pullPolicy }}
          ports:
            - name: http
              containerPort: 80
              protocol: TCP
          livenessProbe:
            httpGet:
              path: /
              port: http
          readinessProbe:
            httpGet:
              path: /
              port: http
          resources:
            {{- toYaml .Values.resources | nindent 12 }}
      {{- with .Values.nodeSelector }}
      nodeSelector:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.affinity }}
      affinity:
        {{- toYaml . | nindent 8 }}
      {{- end }}
      {{- with .Values.tolerations }}
      tolerations:
        {{- toYaml . | nindent 8 }}
      {{- end }}
```

### Service Template

```yaml
apiVersion: v1
kind: Service
metadata:
  name: {{ include "mychart.fullname" . }}
  labels:
    {{- include "mychart.labels" . | nindent 4 }}
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: http
      protocol: TCP
      name: http
  selector:
    {{- include "mychart.selectorLabels" . | nindent 4 }}
```

### Ingress Template

```yaml
{{- if .Values.ingress.enabled -}}
{{- $fullName := include "mychart.fullname" . -}}
{{- $svcPort := .Values.service.port -}}
{{- if and .Values.ingress.className (not (hasKey .Values.ingress.annotations "kubernetes.io/ingress.class")) }}
  {{- $_ := set .Values.ingress.annotations "kubernetes.io/ingress.class" .Values.ingress.className}}
{{- end }}
{{- if semverCompare ">=1.19-0" .Capabilities.KubeVersion.GitVersion -}}
apiVersion: networking.k8s.io/v1
{{- else if semverCompare ">=1.14-0" .Capabilities.KubeVersion.GitVersion -}}
apiVersion: networking.k8s.io/v1beta1
{{- else -}}
apiVersion: extensions/v1beta1
{{- end }}
kind: Ingress
metadata:
  name: {{ $fullName }}
  labels:
    {{- include "mychart.labels" . | nindent 4 }}
  {{- with .Values.ingress.annotations }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
spec:
  {{- if and .Values.ingress.className (semverCompare ">=1.19-0" .Capabilities.KubeVersion.GitVersion) }}
  ingressClassName: {{ .Values.ingress.className }}
  {{- end }}
  {{- if .Values.ingress.tls }}
  tls:
    {{- range .Values.ingress.tls }}
    - hosts:
        {{- range .hosts }}
        - {{ . | quote }}
        {{- end }}
      secretName: {{ .secretName }}
    {{- end }}
  {{- end }}
  rules:
    {{- range .Values.ingress.hosts }}
    - host: {{ .host | quote }}
      http:
        paths:
          {{- range .paths }}
          - path: {{ .path }}
            {{- if and .pathType (semverCompare ">=1.18-0" $.Capabilities.KubeVersion.GitVersion) }}
            pathType: {{ .pathType }}
            {{- end }}
            backend:
              {{- if semverCompare ">=1.19-0" $.Capabilities.KubeVersion.GitVersion }}
              service:
                name: {{ $fullName }}
                port:
                  number: {{ $svcPort }}
              {{- else }}
              serviceName: {{ $fullName }}
              servicePort: {{ $svcPort }}
              {{- end }}
          {{- end }}
    {{- end }}
{{- end }}
```

## Named Templates

### _helpers.tpl

```yaml
{{/*
Expand the name of the chart.
*/}}
{{- define "mychart.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "mychart.fullname" -}}
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
{{- define "mychart.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "mychart.labels" -}}
helm.sh/chart: {{ include "mychart.chart" . }}
{{ include "mychart.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "mychart.selectorLabels" -}}
app.kubernetes.io/name: {{ include "mychart.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "mychart.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "mychart.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
```

## Dependencies

### Chart.yaml with Dependencies

```yaml
apiVersion: v2
name: mychart
description: A Helm chart for Kubernetes
type: application
version: 0.1.0
appVersion: "1.16.0"

dependencies:
  - name: mariadb
    version: 5.x.x
    repository: https://charts.helm.sh/stable
    condition: mariadb.enabled
    tags:
      - database
  - name: redis
    version: 1.x.x
    repository: https://charts.helm.sh/stable
    condition: redis.enabled
    tags:
      - cache
```

### Managing Dependencies

```bash
# Add dependency
helm dependency add mariadb https://charts.helm.sh/stable

# Update dependencies
helm dependency update

# Build dependencies
helm dependency build

# List dependencies
helm dependency list
```

## Hooks

### Pre-install Hook

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "{{ .Release.Name }}-pre-install"
  annotations:
    "helm.sh/hook": pre-install
    "helm.sh/hook-weight": "-5"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    metadata:
      name: "{{ .Release.Name }}-pre-install"
    spec:
      restartPolicy: Never
      containers:
      - name: pre-install-job
        image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
        command: ["/bin/sh", "-c", "echo 'Pre-install hook running'"]
```

### Post-install Hook

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: "{{ .Release.Name }}-post-install"
  annotations:
    "helm.sh/hook": post-install
    "helm.sh/hook-weight": "-5"
    "helm.sh/hook-delete-policy": hook-succeeded
spec:
  template:
    metadata:
      name: "{{ .Release.Name }}-post-install"
    spec:
      restartPolicy: Never
      containers:
      - name: post-install-job
        image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
        command: ["/bin/sh", "-c", "echo 'Post-install hook running'"]
```

## Library Charts

### Library Chart Structure

```yaml
# Chart.yaml
apiVersion: v2
name: mylibchart
description: A Helm chart for Kubernetes
type: library
version: 0.1.0
```

### Library Chart Template

```yaml
{{- define "mylibchart.configmap.tpl" -}}
apiVersion: v1
kind: ConfigMap
metadata:
  name: {{ .Release.Name | printf "%s-%s" .Chart.Name }}
data: {}
{{- end -}}

{{- define "mylibchart.configmap" -}}
{{- include "mylibchart.util.merge" (append . "mylibchart.configmap.tpl") -}}
{{- end -}}
```

### Using Library Chart

```yaml
{{- template "mylibchart.configmap" (list . "mylibchart.configmap") -}}
{{- define "mylibchart.configmap" -}}
## Define overrides for your ConfigMap resource here, e.g.
# data:
#   my-key: my-value
{{- end -}}
```

## Values Import/Export

### Exporting Values

```yaml
# Child chart values.yaml
exports:
  data:
    myint: 99
```

### Importing Values

```yaml
# Parent chart Chart.yaml
dependencies:
  - name: subchart
    repository: http://localhost:10191
    version: 0.1.0
    import-values:
      - data
```

### Import with Mapping

```yaml
# Parent chart Chart.yaml
dependencies:
  - name: subchart1
    repository: http://localhost:10191
    version: 0.1.0
    import-values:
      - child: default.data
        parent: myimports
```

## Complete Example

### Full Chart Structure

```shell
mychart/
├── Chart.yaml
├── values.yaml
├── values.schema.json
├── charts/
│   ├── mariadb-5.x.x.tgz
│   └── redis-1.x.x.tgz
└── templates/
    ├── _helpers.tpl
    ├── deployment.yaml
    ├── service.yaml
    ├── ingress.yaml
    ├── serviceaccount.yaml
    ├── hpa.yaml
    ├── NOTES.txt
    └── tests/
        └── test-connection.yaml
```

### Values.yaml

```yaml
# Default values for mychart
replicaCount: 1

image:
  repository: nginx
  pullPolicy: IfNotPresent
  tag: ""

imagePullSecrets: []
nameOverride: ""
fullnameOverride: ""

serviceAccount:
  create: true
  annotations: {}
  name: ""

rbac:
  create: true

service:
  type: ClusterIP
  port: 80

ingress:
  enabled: false
  className: ""
  annotations: {}
  hosts:
    - host: chart-example.local
      paths:
        - path: /
          pathType: ImplementationSpecific
  tls: []

autoscaling:
  enabled: false
  minReplicas: 1
  maxReplicas: 100
  targetCPUUtilizationPercentage: 80

resources: {}
nodeSelector: {}
tolerations: []
affinity: {}

mariadb:
  enabled: true
  auth:
    rootPassword: "rootpassword"
    database: "mydb"
    username: "user"
    password: "password"

redis:
  enabled: true
  auth:
    enabled: false
```

### NOTES.txt

```txt
1. Get the application URL by running these commands:
{{- if .Values.ingress.enabled }}
{{- range $host := .Values.ingress.hosts }}
  {{- range .paths }}
  http{{ if $.Values.ingress.tls }}s{{ end }}://{{ $host.host }}{{ .path }}
  {{- end }}
{{- end }}
{{- else if contains "NodePort" .Values.service.type }}
  export NODE_PORT=$(kubectl get --namespace {{ .Release.Namespace }} -o jsonpath="{.spec.ports[0].nodePort}" services {{ include "mychart.fullname" . }})
  export NODE_IP=$(kubectl get nodes --namespace {{ .Release.Namespace }} -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT
{{- else if contains "LoadBalancer" .Values.service.type }}
     NOTE: It may take a few minutes for the LoadBalancer IP to be available.
           You can watch the status of by running 'kubectl get --namespace {{ .Release.Namespace }} svc -w {{ include "mychart.fullname" . }}'
  export SERVICE_IP=$(kubectl get svc --namespace {{ .Release.Namespace }} {{ include "mychart.fullname" . }} --template "{{"{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}"}}")
  echo http://$SERVICE_IP:{{ .Values.service.port }}
{{- else if contains "ClusterIP" .Values.service.type }}
  export POD_NAME=$(kubectl get pods --namespace {{ .Release.Namespace }} -l "app.kubernetes.io/name={{ include "mychart.name" . }},app.kubernetes.io/instance={{ .Release.Name }}" -o jsonpath="{.items[0].metadata.name}")
  export CONTAINER_PORT=$(kubectl get pod --namespace {{ .Release.Namespace }} $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace {{ .Release.Namespace }} port-forward $POD_NAME 8080:$CONTAINER_PORT
{{- end }}

2. Get the application URL by running these commands:
{{- if contains "NodePort" .Values.service.type }}
  export NODE_PORT=$(kubectl get --namespace {{ .Release.Namespace }} -o jsonpath="{.spec.ports[0].nodePort}" services {{ include "mychart.fullname" . }})
  export NODE_IP=$(kubectl get nodes --namespace {{ .Release.Namespace }} -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT
{{- else if contains "LoadBalancer" .Values.service.type }}
     NOTE: It may take a few minutes for the LoadBalancer IP to be available.
           You can watch the status of by running 'kubectl get --namespace {{ .Release.Namespace }} svc -w {{ include "mychart.fullname" . }}'
  export SERVICE_IP=$(kubectl get svc --namespace {{ .Release.Namespace }} {{ include "mychart.fullname" . }} --template "{{"{{ range (index .status.loadBalancer.ingress 0) }}{{.}}{{ end }}"}}")
  echo http://$SERVICE_IP:{{ .Values.service.port }}
{{- else if contains "ClusterIP" .Values.service.type }}
  export POD_NAME=$(kubectl get pods --namespace {{ .Release.Namespace }} -l "app.kubernetes.io/name={{ include "mychart.name" . }},app.kubernetes.io/instance={{ .Release.Name }}" -o jsonpath="{.items[0].metadata.name}")
  export CONTAINER_PORT=$(kubectl get pod --namespace {{ .Release.Namespace }} $POD_NAME -o jsonpath="{.spec.containers[0].ports[0].containerPort}")
  echo "Visit http://127.0.0.1:8080 to use your application"
  kubectl --namespace {{ .Release.Namespace }} port-forward $POD_NAME 8080:$CONTAINER_PORT
{{- end }}
```

## Best Practices

### Chart Organization

- Use descriptive chart names
- Include proper documentation
- Use semantic versioning
- Include comprehensive values.yaml

### Template Development

- Use named templates for reusability
- Implement proper error handling
- Use conditional logic appropriately
- Follow Kubernetes naming conventions

### Values Management

- Use flat structure when possible
- Provide sensible defaults
- Document all values
- Use validation schemas

### Security

- Use secrets for sensitive data
- Implement proper RBAC
- Use security contexts
- Validate input values

This reference covers the essential aspects of Helm chart development. For more detailed information, refer to the official Helm documentation.
