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
  # this rule inhibits ALL alerts during a clusterupdate,
  # except the MCDDrainError that can indicate a blocked update
  - source_matchers: [alertname = ClusterIsUpdating]
    target_matchers: [alertname != MCDDrainError]
receivers:
  - name: Default
  - name: Watchdog
  - name: SlackCritical
    slack_configs:
      - channel: critical
        api_url: >-
          {{ .SlackChannelCritical }}
        title: Alert
        text: |-
          {{`{{ range .Alerts }}`}}
          *Summary:* {{`{{ .Annotations.summary }}`}}
          *Description:* {{`{{ .Annotations.description }}`}}
          *Details:*
            {{ printf " {{ range .Labels.SortedPairs }} • *{{ .Name }}:* `{{ .Value }}` "}}
            {{`{{ end }}`}}
          {{`{{ end }}`}}
  - name: SlackMonitoringInternalApplications
    slack_configs:
      - channel: monitoring-internal-applications
        api_url: >-
          {{ .SlackChannelMonitoringInternalApplications }}
        title: Alert
        text: |-
          {{`{{ range .Alerts }}`}}
          *Summary:* {{`{{ .Annotations.summary }}`}}
          *Description:* {{`{{ .Annotations.description }}`}}
          *Details:*
            {{ printf " {{ range .Labels.SortedPairs }} • *{{ .Name }}:* `{{ .Value }}` "}}
            {{`{{ end }}`}}
          {{`{{ end }}`}}

route:
  group_by:
    - namespace
  group_interval: 5m
  group_wait: 30s
  receiver: Default
  repeat_interval: 12h
  routes:
    - receiver: Watchdog
      matchers:
        - alertname = Watchdog
    # Alle critical Alerts aus namespaces mit label openshift.io/cluster-monitoring: "true"
    - receiver: SlackCritical
      matchers:
        - severity = critical
        - openshift_io_alert_source = platform
    # Critical Alerts von uns werden über ein zusätzliches label type: "internal" in den critical channel gesendet
    - receiver: SlackCritical
      matchers:
        - type = internal
        - severity = critical
    # Warning Alerts von uns werden über ein zusätzliches label type: "internal" in den monitoring-internal-applications channel gesendet
    - receiver: SlackMonitoringInternalApplications
      matchers:
        - type = internal
        - severity = warning
    # Interessante Alerts, kommen in den monitoring-internal-applications Slack Channel.
    # Definition über alertname
    - receiver: SlackMonitoringInternalApplications
      matchers:
        - alertname =~ "ClusterNotUpgradeable|APIRemovedInNextEUSReleaseInUse|MCDDrainError"
