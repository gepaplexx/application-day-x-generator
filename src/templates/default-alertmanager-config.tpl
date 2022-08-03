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
  - name: Default
  - name: Watchdog
  - name: SlackCritical
    slack_configs:
      - channel: critical
        api_url: >-
          {{ .SlackChannelCritical}}
        title: Alert
        text: |-
          {{`{{ range .Alerts }}`}}
          *Summary:* {{`{{ .Annotations.summary }}`}}
          *Description:* {{`{{ .Annotations.description }}`}}
          *Details:*
            {{`{{ range .Labels.SortedPairs }} • *{{ .Name }}:* `{{ .Value }}``}}
            {{`{{ end }}`}}
          {{`{{ end }}`}}
  - name: SlackMonitoringInternalApplications
    slack_configs:
      - channel: monitoring-internal-applications
        api_url: >-
          {{ .SlackChannelMonitoringInternalApplications}}
        title: Alert
        text: |-
          {{`{{ range .Alerts }}`}}
          *Summary:* {{`{{ .Annotations.summary }}`}}
          *Description:* {{`{{ .Annotations.description }}`}}
          *Details:*
            {{`{{ range .Labels.SortedPairs }} • *{{ .Name }}:* `{{ .Value }}``}}
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
    # * Alle Critical Alerts aus OpenShift Namespaces kommen in den Critical Channel.
    # Zusätzliche Namespaces, wenn gute Prometheusrules vorhanden sind, wie beispielsweise bei rook-ceph
    # Alle critical Alerts aus namespaces mit label openshift.io/cluster-monitoring: "true"
    - receiver: SlackCritical
      matchers:
        - severity = critical
        - openshift_io_alert_source = platform
    # * Critical Alerts von uns  werden über ein zusätzliches label (vorschlag: type=internal) in den critical channel
    # Alle critical Alerts mit label type: "internal"
    - receiver: SlackCritical
      matchers:
        - type = internal
        - severity = critical
    # * Warning Alerts von uns, kommen entweder in die Cluster Notifications, oder wenn mit label versehen in “monitoring-internal-applications”
    # Warning Alerts mit label
    - receiver: SlackMonitoringInternalApplications
      matchers:
        - type = internal
        - severity = warning
    # * Interessante Alerts, kommen in den “monitoring-internal-applications” Slack Channel, wenn diese wichtig sind.
    # Definition über alertname
    # Warning Alerts mit bestimmten Namen aus namespaces mit label openshift.io/cluster-monitoring: "true"
    - receiver: SlackMonitoringInternalApplications
      matchers:
        - alertname =~ "ClusterNotUpgradeable|APIRemovedInNextEUSReleaseInUse"
