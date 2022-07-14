package generator

var GENERATORS = []Generator{
	{
		ValueBuilder: &GoogleOAuthValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "google-oauth",
	},
	{
		ValueBuilder: &GitOAuthValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "github-oauth",
	},
	{
		ValueBuilder: &HtpasswdValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "htpasswd",
	},
	{
		ValueBuilder: &ClusterConfigValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "cluster-config",
	},
	{
		ValueBuilder: &ClusterIssuerValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "cluster-issuer",
	},
	{
		ValueBuilder: &ConsolePatchesValueBuilder{},
		Stage:        InitialClusterSetup,
		Name:         "console-patches",
	},
	{
		ValueBuilder: &RookCephInstanceValueBuilder{},
		Stage:        ClusterSetupCheckpoint,
		Name:         "rook-ceph-instance",
	},
	{
		ValueBuilder: &ClusterLoggingValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "cluster-logging",
	},
	{
		ValueBuilder: &ClusterMonitoringValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "cluster-monitoring",
	},
	{
		ValueBuilder: &ClusterUpdaterValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "cluster-updater",
	},
	{
		ValueBuilder: &GepaplexxCicdToolsValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "gepaplexx-cicd-tools",
	},
	{
		ValueBuilder: &VaultValueBuilder{},
		Stage:        ClusterApplications,
		Name:         "vault",
	},
}
