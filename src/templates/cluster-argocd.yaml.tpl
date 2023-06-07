apiVersion: argoproj.io/v1alpha1
kind: ArgoCD
metadata:
  name: gepardec-run-argocd
spec:
  server:
    autoscale:
      enabled: false
    grpc:
      ingress:
        enabled: false
    host: openshift-gitops.{{ .env }}.run.gepardec.com
    ingress:
      enabled: false
    insecure: true
    resources:
      limits:
        cpu: 1000m
        memory: 526Mi
      requests:
        cpu: 250m
        memory: 128Mi
    route:
      enabled: false
    service:
      type: ''
  grafana:
    enabled: false
  notifications:
    enabled: true
  prometheus:
    enabled: false
  initialSSHKnownHosts: {}
  sso:
    dex:
      openShiftOAuth: true
      resources:
        limits:
          cpu: 500m
          memory: 256Mi
        requests:
          cpu: 250m
          memory: 128Mi
    provider: dex
  applicationSet:
    resources:
      limits:
        cpu: 1
        memory: 1Gi
      requests:
        cpu: 250m
        memory: 512Mi
  rbac:
    defaultPolicy: ''
    policy: |
      p, role:gepardec-admin, applications, *, *, allow
      p, role:gepardec-admin, repositories, *, *, allow
      p, role:gepardec-admin, projects, *, *, allow
    scopes: '[groups]'
  repo:
    resources:
      limits:
        cpu: 1000m
        memory: 512Mi
      requests:
        cpu: 250m
        memory: 128Mi
  resourceExclusions: |
    - apiGroups:
      - tekton.dev
      clusters:
      - '*'
      kinds:
      - TaskRun
      - PipelineRun
    - apiGroups:
        - cilium.io
      kinds:
        - CiliumIdentity
      clusters:
        - "*"
  ha:
    enabled: false
  tls:
    ca: {}
  redis:
    resources:
      limits:
        cpu: 500m
        memory: 1Gi
      requests:
        cpu: 250m
        memory: 128Mi
  controller:
    processors: {}
    resources:
      limits:
        cpu: 1
        memory: 2Gi
      requests:
        cpu: 250m
        memory: 512Mi
    sharding: {}