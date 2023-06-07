apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: {{ .Name }}
  namespace: {{ or .Namespace "gepardec-run-gitops" }}
spec:
  destination:
    namespace: {{ or .TargetNs "gepardec-run-gitops" }}
    server: 'https://kubernetes.default.svc'
  project: {{ or .ArgoProject "bootstrap-argoproject" }}
  source:
    helm:
      valueFiles:
        - {{ .ValuesFile }}
    path: {{ .Path }}
    repoURL: '{{ or .Repo "https://github.com/gepaplexx/gepardec-run-cluster-configuration/" }}'
    targetRevision: {{ or .Revision "main" }}
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
    syncOptions:
      - CreateNamespace=true
      - ServerSideApply=true