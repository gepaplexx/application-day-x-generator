package generator

import (
	utils "gepaplexx/day-x-generator/pkg/util"
	"strings"
)

type InitialSetupScriptValueBuilder struct{}

func (gen *InitialSetupScriptValueBuilder) GetValues(config map[string]utils.Value) (map[string]utils.Value, error) {
	values := make(map[string]utils.Value)
	env := config["env"].String()
	initialClusterSetupPath := buildFilename(ArgoBootstrap, env, "initial-cluster-setup", stages[ArgoBootstrap])
	values["BootstrapInitialClusterSetupYaml"] = utils.Value{
		Val: strings.TrimPrefix(initialClusterSetupPath, utils.TARGET_DIR+"/"),
	}

	clusterApplicationsPath := buildFilename(ArgoBootstrap, env, "cluster-applications", stages[ArgoBootstrap])
	values["BootstrapClusterApplicationsYaml"] = utils.Value{
		Val: strings.TrimPrefix(clusterApplicationsPath, utils.TARGET_DIR+"/"),
	}

	clusterArgoCDPath := buildFilename(ClusterArgoCD, env, "", stages[ClusterArgoCD])
	values["ClusterArgocdYaml"] = utils.Value{
		Val: strings.TrimPrefix(clusterArgoCDPath, utils.TARGET_DIR+"/"),
	}

	clusterArgoCDWorkflowRepoPath := buildFilename(ClusterArgoCDRepo, env, "", stages[ClusterArgoCDRepo])
	values["ClusterArgocdWorkflowRepoSecret"] = utils.Value{
		Val: strings.TrimPrefix(clusterArgoCDWorkflowRepoPath, utils.TARGET_DIR+"/"),
	}

	for key, val := range config {
		values[key] = val
	}

	return values, nil
}
