package generator

var GENERATORS = []Generator{
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "keycloak-operator",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "external-secret",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "grafana-operator",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "vault",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "kyverno",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "keycloak-instance",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "external-secrets-config",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "cicd-tools",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "cicd-eventbus",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "gerpardec-run-cicd",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "grafana-instance",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "grafana-dashboards",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "multena",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "altermanager",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "one-time-secret",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        VaultSetupScript,
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ArgoBootstrap,
		Name:         "initial-cluster-setup",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ArgoBootstrap,
		Name:         "cluster-applications",
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterArgoCD,
	},
	{
		ValueBuilder: &InitialSetupScriptValueBuilder{},
		Stage:        InitialSetupScript,
	},
	{
		ValueBuilder: &GenericCopyValueBuilder{},
		Stage:        ClusterArgoCDRepo,
	},
}
