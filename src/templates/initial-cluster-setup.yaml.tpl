project:
  create: true
  name: gepardec-run

applications:

  ##################### KEYCLOAK-OPERATOR ######################
  keycloak-operator:
    name: keycloak-operator
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-sso
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-keycloak-operator
      targetRevision: "{{ or .KeycloakOperatorChartVersion "1.0.*" }}"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### EXTERNAL SECRETS ######################
  external-secrets:
    name: external-secrets
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-external-secrets
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "{{ or .ExternalSecretsChartVersion "1.0.*" }}"
      chart: gp-external-secrets-operator
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ####################### Grafana Operator ######################
  grafana-operator:
    name: grafana-operator
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-grafana
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "{{ or .GrafanaOperatorChartVersion "1.0.*" }}"
      chart: "gp-grafana-operator"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### VAULT ######################
  vault:
    name: vault
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-vault
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "{{ or .VaultChartVersion "1.0.*" }}"
      chart: gp-hashicorp-vault
      helm:
        parameters:
          - name: "ingress.hostname"
            value: "vault.{{ .env }}.run.gepardec.com"
          - name: "vault.server.ha.disruptionBudget.enabled"
            value: "false"
          - name: "backup.external.bucket"
            value: "c-gepa-{{ .env }}-vault-backup"
    ignoreDifferences:
      - group: admissionregistration.k8s.io
        jsonPointers:
          - /webhooks/0/clientConfig/caBundle
        kind: MutatingWebhookConfiguration
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### KYVERNO/Cluster-Policies ######################
  kyverno:
    name: kyverno
    enabled: true
    argoProject: gepardec-run
    destination:
      namespace: gp-kyverno
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-kyverno
      targetRevision: "{{ or .KyvernoChartVersion "1.0.*" }}"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true
      replace: true