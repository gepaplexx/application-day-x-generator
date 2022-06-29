global:
  resolve_timeout: 5m
inhibit_rules:
  - equal:
      - namespace
      - alertname
    source_match:
      severity: critical
    target_match_re:
      severity: warning|info
  - equal:
      - namespace
      - alertname
    source_match:
      severity: warning
    target_match_re:
      severity: info
receivers:
  - name: Critical
  - name: Default
  - name: Watchdog
  - name: ClusterNotUpgradeable
    slack_configs:
      - channel: critical
        api_url: >-
          {{ .SlackChannel }}
        title: ClusterNotUpgradeable
        text: "Environment: {{ .env }}.gepaplexx.com \n{{ `{{ .CommonAnnotations.message }}` }}"
  - name: APIRemovedInNextEUSReleaseInUse
    slack_configs:
      - channel: critical
        api_url: >-
          {{ .SlackChannel }}
        text: "Environment: {{ .env }}.gepaplexx.com \n{{ `{{ range .Alerts }}` }} {{ `{{ .Annotations.message }}` }} \n ---------------------------------------------------------------------------------------------------- \n{{ `{{ end }}` }}"
  - name: PvcMonitoringAlerts
    slack_configs:
      - channel: critical
        api_url: >-
        {{ .SlackChannel }}
        title: PersistentVolumeClaim Alert
        text: |-
          Themengebiet: STORAGE
          Summary: {{`{{ .CommonAnnotations.summary }}`}}
          Environment: {{ .env }}.gepaplexx.com
          Alerts:
          {{`{{- range .Alerts -}}`}}
          - {{`{{ .Annotations.description }}`}}\n
          {{`{{- end -}}`}}
route:
  group_by:
    - namespace
  group_interval: 5m
  group_wait: 30s
  receiver: Default
  repeat_interval: 12h
  routes:
    - match:
        alertname: Watchdog
      receiver: Watchdog
    - match:
        severity: critical
      receiver: Critical
    - receiver: ClusterNotUpgradeable
      match:
        alertname: ClusterNotUpgradeable
    - receiver: APIRemovedInNextEUSReleaseInUse
      match:
        alertname: APIRemovedInNextEUSReleaseInUse
    - receiver: PvcMonitoringAlerts
      matchers:
      - monitoring=pvc