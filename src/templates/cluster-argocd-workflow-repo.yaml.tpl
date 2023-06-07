apiVersion: external-secrets.io/v1beta1
kind: ExternalSecret
metadata:
  name: workflow-repository
spec:
  dataFrom:
    - extract:
        conversionStrategy: Default
        decodingStrategy: None
        key: argocd
  refreshInterval: 12h
  secretStoreRef:
    kind: ClusterSecretStore
    name: internal-cluster-store
  target:
    creationPolicy: Owner
    deletionPolicy: Retain
    name: workflow-repository
    template:
      data:
        insecure: "false"
        name: gepaplexx-cicd
        sshPrivateKey: '{{ `{{ .cicd_repo_ssh_private_key }}` }}'
        type: git
        url: git@github.com:gepaplexx/gepaplexx-cicd.git
        username: '{{ `{{ .cicd_repo_username }}` }}'
      engineVersion: v2
      mergePolicy: Replace
      metadata:
        labels:
          app.kubernetes.io/managed-by: workflow-repository
          argocd.argoproj.io/secret-type: repository
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: admin-config-reader
---