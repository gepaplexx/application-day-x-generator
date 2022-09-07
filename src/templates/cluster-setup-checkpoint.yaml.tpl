project:
  create: false
  name: gepaplexx

applications:
########################## ROOK/CEPH INSTANCE ##########################
  rookCephInstance:
    name: rook-ceph-instance
    enabled: {{ .RookCephInstanceEnabled }}
    argoProject: gepaplexx
    destination:
      namespace: rook-ceph
      create: false
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-storage-cephcluster
      targetRevision: "*"
    ignoreDifferences:
      - group: ceph.rook.io
        kind: CephCluster
        name: rook-ceph
        jsonPointers:
          - /spec/monitoring/createPrometheusRules
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

################### CLUSTER-LOGGING OPERATOR ##########################
  clusterLoggingOperator:
    name: cluster-logging-operator
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: openshift-logging
      create: false
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-cluster-logging-operator
      targetRevision: "*"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ####################### NFS-STORAGE-PROVISIONER ######################
  nfs-provisioner:
    name: nfs
    enabled: false # TODO: enable sobald FW-Freischaltung da ist.
    argoProject: gepaplexx
    destination:
      namespace: nfs-storage-provisioner
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "*"
      chart: "gp-nfs-provisioner"
      helm:
        parameters:
          - name: "provisioner.nfs.path"
            value: "/data/col1/{{ .env }}"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ####################### Grafana Operator ######################
  grafana-operator:
    name: grafana-operator
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: grafana-operator-system
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      targetRevision: "*"
      chart: "gp-grafana-operator"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### KEYCLOAK-INSTANCE ######################
  keycloak-instance:
    name: keycloak-instance
    enabled: true
    argoProject: gepaplexx
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
      targetRevision: "*"
      helm:
        parameters:
        - name: "ingress.hostname"
          value: "sso.apps.{{ .env }}.gepaplexx.com"
        - name: "persistence.auth.password"
          value: "{{ .KeycloakDbPassword }}"
        - name: "provider.openshift.clientSecret"
          value: "{{ .KeycloakOcpClientSecret }}"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true