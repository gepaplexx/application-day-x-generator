project:
  create: false
  name: gepardec-run

applications:
#################### KEYCLOAK-INSTANCE ######################
  keycloak-instance:
    name: keycloak-instance
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-sso
      create: false
    ignoreDifferences:
      - jsonPointers:
          - /secret
        kind: OAuthClient
        group: oauth.openshift.io
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-keycloak-instance
      targetRevision: "{{ or .KeyCloakInstanceChartVersion "1.0.*" }}"
      helm:
        parameters:
          - name: "ingress.hostname"
            value: "sso.{{ .env }}.run.gepardec.com"
          - name: "keycloakConfigCli.enabled"
            value: "{{ .KeycloakInstanceConfigCliEnabled }}"
          - name: "keycloakConfigCli.cluster"
            value: "{{ .env }}"
          - name: "keycloakConfigCli.identityProvider.openshift.baseUrl"
            value: "sso.{{ .env }}.run.gepardec.com"
          - name: "postgresql.backup.external.bucket"
            value: "c-gepa-{{ .env }}-keycloak-backup"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

############## EXTERNAL SECRETS CONFIGURATION ###############
  external-secrets:
    name: external-secrets-configuration
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-external-secrets
      create: false
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "{{ or .ExternalSecretsConfigChartVersion "1.0.*" }}"
      chart: gp-external-secrets-configuration
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### CICD-TOOLS ######################
  cicd-tools:
    name: cicd-tools
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-cicd-tools
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-cicd-tools
      targetRevision: "{{ or .CicdToolsChartVersion "1.0.*" }}"
      helm:
        parameters:
          - name: "clustername"
            value: "{{ .env }}"
          - name: "argo_workflows.rbac.clusterscoped.enabled"
            value: "{{ .ArgoWorkflowsClusterScopedGroupEnabled }}"
          - name: "argo_workflows.server.sso.issuer"
            value: "https://sso.{{ .env }}.run.gepardec.com/realms/internal"
          - name: "argo_workflows.server.sso.redirectUrl"
            value: "https://workflows.{{ .env }}.run.gepardec.com/oauth2/callback"
          - name: "argo_workflows.artifactRepository.s3.bucket"
            value: "c-gepa-{{ .env }-argo-workflows"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### CICD-EVENTBUS ######################
  cicd-eventbus:
    name: cicd-eventbus
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-cicd-eventbus
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-cicd-eventbus
      targetRevision: "{{ or .CicdEventbusChartVersion "1.0.*" }}"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

##################### GEPARDEC-RUN-CICD ####################
  gepardec-run-cicd:
    name: gepardec-run-cicd
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gepardec-run-gitops
      create: false
    source:
      repoURL: "git@github.com:gepaplexx/gepaplexx-cicd.git"
      targetRevision: "{{ or .GepardecRunCicdChartVersion "main" }}"
      path: "clusters/{{ .env }}/applications"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

######################### Grafana Instance ######################
  grafana-instance:
    name: grafana-instance
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-grafana
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "{{ or .GrafanaInstanceChartVersion "1.0.*" }}"
      chart: "gp-grafana-instance"
      helm:
        parameters:
          - name: "ingress.hostname"
            value: "grafana.{{ .env }}.run.gepardec.com"
          - name: "sso.keycloak.realmUrl"
            value: "https://sso.{{ .env }}.run.gepardec.com/realms/internal"
    ignoreDifferences:
      - jsonPointers:
          - /spec/config/auth.generic_oauth/client_secret
        kind: Grafana
        group: integreatly.org
    syncPolicy:
      automated:
        prune: true
        selfHeal: true
#
#  ####################### Grafana Dashboards ######################
  grafana-dashboards:
    name: grafana-dashboards
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-grafana
      create: true
    source:
      repoURL: "https://github.com/gepaplexx/gepardec-run-cluster-configuration"
      targetRevision: "{{ or .GrafanaDashboardsChartVersion "main" }}"
      path: "observability/dashboards/"
      directory:
        recurse: true
        include: "{all/*,{{ .env }}/*}"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true
#
#  ####################### Multena ######################
  multena:
    name: multena
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-grafana
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "{{ or .MultenaChartVersion "1.0.*" }}"
      chart: "gp-multena"
      helm:
        parameters:
          - name: "multena.jwksCertUrl"
            value: https://sso.{{ .env }}.run.gepardec.com/realms/internal/protocol/openid-connect/certs
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

##################### ALERTMANAGER ######################
  alertmanager:
    name: alertmanager
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: openshift-user-workload-monitoring
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "{{ or .AlertmanagerChartVersion "0.1.*" }}"
      chart: gp-alertmanager
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

##################### ONE-TIME-SECRET ######################
  one-time-secret:
    name: one-time-secret
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-one-time-secret
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "{{ or .OneTimeSecretChartVersion "1.1.*" }}"
      chart: gp-one-time-secret
      helm:
        parameters:
          - name: "ingress.hostname"
            value: "secret.{{ .env }}.run.gepardec.com"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true
