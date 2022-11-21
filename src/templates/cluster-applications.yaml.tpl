project:
  create: false
  name: gepaplexx

applications:
####################### CLUSTER-UPDATER #######################
  clusterUpdater:
    name: cluster-updater
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gp-infrastructure
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-cluster-update-checker
      targetRevision: "*"
      helm:
        parameters:
          - name: "clustername"
            value: "{{ .env }}"
          - name: "consoleUrl"
            value: "console.apps.{{ .env }}.gepaplexx.com"
          - name: "slack.channel"
            value: "{{ .SlackChannel }}"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

####################### CLUSTER-LOGGING #######################
  clusterLogging:
    name: cluster-logging-instance
    enabled: {{ .ClusterLoggingEnabled }}
    destination:
      namespace: openshift-logging
      create: false
    argoProject: gepaplexx
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-cluster-logging-instance
      targetRevision: "*"
      helm:
        parameters:
          - name: "lokistack.backend.secretkey"
            value: "{{ .ClusterLoggingS3SecretKey }}"
          - name: "lokistack.minio.enabled"
            value: "true" #TODO: löschen, sobald NetApp S3 Endpoint fertig konfiguriert ist.
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

####################### CLUSTER-LOGGING-EVENTROUTER #######################
  clusterLoggingEventrouter:
    name: cluster-logging-eventrouter
    enabled: true
    destination:
      namespace: openshift-logging
      create: false
    argoProject: gepaplexx
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-cluster-logging-eventrouter
      targetRevision: "*"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

##################### CLUSTER-MONITORING ######################
  clusterMonitoring:
    name: cluster-monitoring
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: openshift-monitoring
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-cluster-monitoring-config
      targetRevision: "*"
      helm:
        parameters:
          - name: "alertmanager.config"
            value: "{{ .AlertmanagerYaml }}"
          - name: "infranodes.enabled"
            value: "{{ .InfranodesEnabled }}"
          - name: "clusterMonitoring.prometheusK8s.clusterName"
            value: "{{ .env }}"
          - name: "clusterMonitoring.prometheusK8s.remoteWrite.password"
            value: "{{ .PrometheusRemoteWritePassword }}"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

##################### OPENSHIFT IMAGE REGISTRY ######################
  internalRegistry:
    name: openshift-registry
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gp-infrastructure
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-ocp-internal-registry
      targetRevision: "*"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### GEPAPLEXX-CICD-TOOLS ######################
  gepaplexx-cicd-tools:
    name: gepaplexx-cicd-tools
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gepaplexx-cicd-tools
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-cicd-tools
      targetRevision: "*"
      helm:
        parameters:
          - name: "argocd.route.hostname"
            value: "argocd.apps.{{ .env }}.gepaplexx.com"
          - name: "argo_workflows.server.ingress.hosts[0]"
            value: "workflows.apps.{{ .env }}.gepaplexx.com"
          - name: "argo_workflows.server.ingress.tls[0].hosts[0]"
            value: "workflows.apps.{{ .env }}.gepaplexx.com"
          - name: "argo_workflows.server.ingress.tls[0].secretName"
            value: "workflows.apps.{{ .env }}.gepaplexx.com-tls"
          - name: "argo_workflows.server.sso.issuer"
            value: "{{ .ArgoWorkflowsSsoIssuer }}"
          - name: "argo_workflows.server.sso.redirectUrl"
            value: "https://workflows.apps.{{ .env }}.gepaplexx.com/oauth2/callback"
          - name: "argo_workflows.rbac.clusterscoped.enabled"
            value: "{{ .ArgoWorkflowsClusterScopedGroupEnabled }}"
          - name: "argo_workflows.rbac.clientSecret"
            value: "{{ .ArgoWorkflowsSsoClientSecret }}"
          - name: "argo_workflows.archive.secretkey"
            value: "{{ .ArgoWorkflowsMinioSecretkey }}"
          - name: "sealedSecret.postgresql.password"
            value: "{{ .PostgresqlPassword }}"
          - name: "sealedSecret.postgresql.postgresPassword"
            value: "{{ .PostgresqlPostgresPassword }}"
          - name: "sealedSecret.postgresql.username"
            value: "{{ .PostgresqlUsername }}"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### GEPAPLEXX-CICD-EVENTBUS ######################
  gepaplexx-cicd-eventbus:
    name: gepaplexx-cicd-eventbus
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gepaplexx-cicd-eventbus
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gepaplexx-cicd-eventbus
      targetRevision: "*"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### GEPAPLEXX-CICD ######################
  gepaplexx-cicd:
    name: gepaplexx-cicd
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: openshift-gitops
      create: false
    source:
      repoURL: "git@github.com:gepaplexx/gepaplexx-cicd.git"
      targetRevision: "main"
      path: "clusters/{{ .env }}/applications"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ####################### Grafana Instance ######################
  grafana-instance:
    name: grafana-instance
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: grafana-operator-system
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "*"
      chart: "gp-grafana-instance"
      helm:
        parameters:
          - name: "ingress.hostname"
            value: "grafana.apps.{{ .env }}.gepaplexx.com"
          - name: "sso.keycloak.clientSecret"
            value: "{{ .KeycloakClientSecret }}"
          - name: "sso.keycloak.realmUrl"
            value: "{{ .KeycloakRealmUrl }}"
    ignoreDifferences:
      - jsonPointers:
          - /spec/config/auth.generic_oauth/client_secret
        kind: Grafana
        group: integreatly.org
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### VAULT ######################
    vault:
      name: vault
      enabled: true
      argoProject: gepaplexx
      destination:
        namespace: gp-vault
        create: true
      source:
        repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
        targetRevision: "*"
        chart: gp-hashicorp-vault
        helm:
          parameters:
            - name: "autoUnseal.creds"
              value: {{ .AutoUnsealCreds }}
      syncPolicy:
        automated:
          prune: true
          selfHeal: true
