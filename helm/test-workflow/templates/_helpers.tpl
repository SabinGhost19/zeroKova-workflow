{{/*
=============================================================================
helper templates for test-workflow helm chart
these templates provide reusable functions across all chart templates
=============================================================================
*/}}

{{/*
expand the name of the chart
uses nameOverride if set, otherwise uses chart name
*/}}
{{- define "test-workflow.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
create a default fully qualified app name
uses fullnameOverride if set, otherwise combines release name and chart name
truncated to 63 chars for kubernetes name compliance
*/}}
{{- define "test-workflow.fullname" -}}
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
create chart name and version for chart label
*/}}
{{- define "test-workflow.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
common labels applied to all resources
includes standard kubernetes and helm labels
*/}}
{{- define "test-workflow.labels" -}}
helm.sh/chart: {{ include "test-workflow.chart" . }}
{{ include "test-workflow.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- with .Values.global.commonLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
selector labels used for pod selection
these labels must match between deployment and service
*/}}
{{- define "test-workflow.selectorLabels" -}}
app.kubernetes.io/name: {{ include "test-workflow.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
service-specific labels
includes component label for the specific service
*/}}
{{- define "test-workflow.serviceLabels" -}}
{{ include "test-workflow.labels" . }}
app.kubernetes.io/component: {{ .component }}
{{- end }}

{{/*
service-specific selector labels
*/}}
{{- define "test-workflow.serviceSelectorLabels" -}}
{{ include "test-workflow.selectorLabels" . }}
app.kubernetes.io/component: {{ .component }}
{{- end }}

{{/*
create the name of the service account to use
*/}}
{{- define "test-workflow.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "test-workflow.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
return the proper image name for a service
combines global registry, service repository, and tag
*/}}
{{- define "test-workflow.image" -}}
{{- $registry := .global.customImageRegistry -}}
{{- $repository := .service.image.repository -}}
{{- $tag := .service.image.tag | default .global.imageTag -}}
{{- printf "%s/%s:%s" $registry $repository $tag -}}
{{- end }}

{{/*
return the image pull policy
uses service-specific policy if set, otherwise uses global
*/}}
{{- define "test-workflow.imagePullPolicy" -}}
{{- .service.image.pullPolicy | default .global.imagePullPolicy -}}
{{- end }}

{{/*
return image pull secrets
*/}}
{{- define "test-workflow.imagePullSecrets" -}}
{{- with .Values.global.imagePullSecrets }}
imagePullSecrets:
  {{- toYaml . | nindent 2 }}
{{- end }}
{{- end }}

{{/*
return the namespace
uses global namespace if set, otherwise uses release namespace
*/}}
{{- define "test-workflow.namespace" -}}
{{- .Values.global.namespace | default .Release.Namespace -}}
{{- end }}

{{/*
postgresql host - returns internal or external host
*/}}
{{- define "test-workflow.postgresql.host" -}}
{{- if .Values.postgresql.enabled }}
{{- printf "%s-postgresql" .Release.Name -}}
{{- else }}
{{- .Values.externalDatabase.host -}}
{{- end }}
{{- end }}

{{/*
postgresql port
*/}}
{{- define "test-workflow.postgresql.port" -}}
{{- if .Values.postgresql.enabled }}
{{- 5432 -}}
{{- else }}
{{- .Values.externalDatabase.port -}}
{{- end }}
{{- end }}

{{/*
postgresql database name
*/}}
{{- define "test-workflow.postgresql.database" -}}
{{- if .Values.postgresql.enabled }}
{{- .Values.postgresql.auth.database -}}
{{- else }}
{{- .Values.externalDatabase.database -}}
{{- end }}
{{- end }}

{{/*
postgresql username
*/}}
{{- define "test-workflow.postgresql.username" -}}
{{- if .Values.postgresql.enabled }}
{{- .Values.postgresql.auth.username -}}
{{- else }}
{{- .Values.externalDatabase.user -}}
{{- end }}
{{- end }}

{{/*
postgresql secret name for password
*/}}
{{- define "test-workflow.postgresql.secretName" -}}
{{- if .Values.postgresql.enabled }}
{{- printf "%s-postgresql" .Release.Name -}}
{{- else }}
{{- .Values.externalDatabase.existingSecret -}}
{{- end }}
{{- end }}

{{/*
kafka bootstrap servers - returns internal or external
*/}}
{{- define "test-workflow.kafka.brokers" -}}
{{- if .Values.kafka.enabled }}
{{- printf "%s-kafka:9092" .Release.Name -}}
{{- else }}
{{- .Values.externalKafka.brokers -}}
{{- end }}
{{- end }}

{{/*
service addresses for grpc communication
*/}}
{{- define "test-workflow.orderService.address" -}}
{{- printf "%s-order-service:%d" (include "test-workflow.fullname" .) (int .Values.orderService.service.grpcPort) -}}
{{- end }}

{{- define "test-workflow.inventoryService.address" -}}
{{- printf "%s-inventory-service:%d" (include "test-workflow.fullname" .) (int .Values.inventoryService.service.grpcPort) -}}
{{- end }}

{{- define "test-workflow.notificationService.address" -}}
{{- printf "%s-notification-service:%d" (include "test-workflow.fullname" .) (int .Values.notificationService.service.grpcPort) -}}
{{- end }}

{{/*
common annotations including global ones
*/}}
{{- define "test-workflow.annotations" -}}
{{- with .Values.global.commonAnnotations }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
pod labels - combination of selector labels and any additional pod labels
*/}}
{{- define "test-workflow.podLabels" -}}
{{ include "test-workflow.serviceSelectorLabels" . }}
tier: {{ .tier | default "backend" }}
{{- end }}

{{/*
checksum annotation for config changes
forces pod restart when configmap changes
*/}}
{{- define "test-workflow.configChecksum" -}}
checksum/config: {{ include (print $.Template.BasePath "/configmap.yaml") . | sha256sum }}
{{- end }}
