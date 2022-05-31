package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type ClusterConfigValueBuilder struct {
}

const WORKFLOW_REPO_SECRET_TEMPLATE string = `
apiVersion: v1
kind: Secret
metadata:
  name: workflow-repository
  namespace: openshift-gitops
data:
  username: "{{ .ArgocdWorkflowrepositoryUsername }}"
  sshPrivateKey: "{{ .ArgocdWorkflowrepositorySshPrivateKey }}"
`

func (gen *ClusterConfigValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	secretVals := make(map[string]string, 3)
	secretVals["ArgocdWorkflowrepositoryUsername"] = utils.Base64(config["ArgocdWorkflowrepositoryUsername"])
	secretVals["ArgocdWorkflowrepositorySshPrivateKey"] = config["ArgocdWorkflowrepositorySshPrivateKeyBase64"].String()

	secretAsByte, err := utils.ReplaceTemplate(secretVals, WORKFLOW_REPO_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "username", "sshPrivateKey")
	if err != nil {
		return nil, err
	}

	values["ArgocdWorkflowrepositoryUsername"] = encryptedValues["username"]
	values["ArgocdWorkflowrepositorySshPrivateKey"] = encryptedValues["sshPrivateKey"]

	return values, nil
}
