project:
  create: true
  name: gepaplexx

applications:
  ################## PATCH OPERATOR #############################
  patchOperator:
    name: patch-operator
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gp-infrastructure
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-patch-operator
      targetRevision: "*"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

################## SEALED SECRETS OPERATOR ####################
  sealedSecretsOperator:
    name: sealed-secrets-operator
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gp-infrastructure
      create: false
    source:
      chart: sealed-secrets
      repoURL: https://bitnami-labs.github.io/sealed-secrets
      targetRevision: 2.*
      helm:
        parameters:
          - name: containerSecurityContext.enabled
            value: 'false'
          - name: podSecurityContext.enabled
            value: 'false'
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

###################### IDENTITY PROVIDER ######################
  identityProvider:
    name: identity-provider
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gp-infrastructure
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-chart-development/"
      chart: gp-identity-provider
      targetRevision: "*"
      helm:
        parameters:
        {{- if .GoogleEnabled.Val }}
        - name: "google.clientSecret"
          value: "{{ .GoogleClientSecret }}"
        - name: "google.clientId"
          value: "{{ .GoogleClientId }}"
        - name: "google.restrDomain"
          value: "{{ .GoogleRestrictedDomain }}"
        {{- end }}
        - name: "google.enable"
          value: "{{ .GoogleEnabled }}"
        {{- if .GitEnabled.Val }}
        - name: "git.clientSecret"
          value: "{{ .GitClientSecret }}"
        - name: "git.clientId"
          value: "{{ .GitClientId }}"
        - name: "git.restrOrgs"
          value: "{{ .GitRestrOrgs }}"
        {{- end }}
        - name: "git.enable"
          value: "{{ .GitEnabled }}"
        - name: "htpasswd.enable"
          value: "{{ .HtpasswdEnabled }}"
        {{- if .HtpasswdEnabled.Val }}          
        - name: "htpasswd.data"
          value: "{{ .HtpasswdData }}"
        {{- end }}
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

###################### OAUTH GROUP SYNC  ######################
  oauthGroupSync:
    name: oauth-group-sync
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gp-infrastructure
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-chart-development/"
      chart: gp-oauth-group-sync
      targetRevision: "*"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true


####################### CLUSTER-CONFIG ########################
  clusterConfig:
    name: cluster-config
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gp-infrastructure
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-chart-development/"
      chart: gp-cluster-config
      targetRevision: "*"
      helm:
        parameters:
          - name: "argocd.workflowrepository.username"
            value: "{{ .ArgocdWorkflowrepositoryUsername }}"
          - name: "argocd.workflowrepository.sshPrivateKey"
            value: "{{ .ArgocdWorkflowrepositorySshPrivateKey }}"
          - name: "alerts.slackUrl.kube"
            value: "{{ .AlertsSlackUrlKube }}"
          - name: "alerts.slackUrl.apiserver"
            value: "{{ .AlertsSlackUrlApiServer }}"
          - name: "alerts.environment"
            value: "{{ .env }}"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

########################## ROOK/CEPH OPERATOR ##########################
  rookCephOperator:
    name: rook-ceph-operator
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: rook-ceph
      create: false
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-chart-development/"
      chart: gp-rook-ceph-operator
      targetRevision: "*"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

##################### IMAGE REGISTRY CACHE ######################
  internalRegistryMirror:
    name: openshift-registry-mirror
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gp-infrastructure
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-chart-development/"
      chart: gp-pull-through-cache
      targetRevision: "*"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

  ##################### CERT-MANAGER ######################
  certManager:
    name: cert-manager
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: cert-manager
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-charts/"
      chart: gp-cert-manager
      targetRevision: "*"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

##################### CERTIFICATES-PATCHES ######################
  certificatesPatches:
    name: certificate-patches
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gp-infrastructure
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-chart-development/"
      chart: gp-certificates-patches
      targetRevision: "*"
      helm:
        parameters:
          - name: "apiserver.customApiUrl"
            value: "api.{{ .env }}.gepaplexx.com"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

##################### CLUSTER-ISSUER ######################
  clusterIssuer:
    name: cluster-issuer
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: cert-manager
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-chart-development/"
      chart: gp-cluster-issuer
      targetRevision: "*"
      helm:
        parameters:
          - name: "solvers.dnsZones[0]"
            value: "{{ .env }}.gepaplexx.com"
          - name: "solvers.accessKeyId"
            value: "{{ .SolversAccesKeyId }}"
          - name: "solvers.secretName"
            value: "{{ .env }}-route53-credentials-secret"
          - name: "solvers.secretAccessKey"
            value: "{{ .SolversSecretAccessKey }}"
          - name: "certificates.defaultIngress"
            value: "*.apps.{{ .env }}.gepaplexx.com"
          - name: "certificates.console"
            value: "console.apps.{{ .env }}.gepaplexx.com"
          - name: "certificates.api"
            value: "api.{{ .env }}.gepaplexx.com"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true

##################### CONSOLE-PATCHES ######################
  consolePatches:
    name: console-patches
    enabled: true
    argoProject: gepaplexx
    destination:
      namespace: gp-infrastructure
      create: true
    source:
      repoURL: "https://gepaplexx.github.io/gp-helm-chart-development/"
      chart: gp-console-patches
      targetRevision: "*"
      helm:
        parameters:
          - name: "route.nameOverride"
            value: "{{ .RouteNameOverride }}"
          - name: "route.hostname"
            value: "apps.{{ .env }}.gepaplexx.com"
    syncPolicy:
      automated:
        prune: true
        selfHeal: true