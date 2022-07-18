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

const ALERTS_SLACK_HOOK_SECRET_KUBE_TEMPLATE string = `
kind: Secret
apiVersion: v1
metadata:
  name: slack-hook
  namespace: openshift-kube-controller-manager-operator
data:
  url: >-
    {{ .AlertsSlackUrl }}
`

const ALERTS_SLACK_HOOK_SECRET_APISERVER_TEMPLATE string = `
kind: Secret
apiVersion: v1
metadata:
  name: slack-hook
  namespace: openshift-apiserver
data:
  url: >-
    {{ .AlertsSlackUrl }}
`

func (gen *ClusterConfigValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	secretValsArgo := make(map[string]string, 3)
	secretValsArgo["ArgocdWorkflowrepositoryUsername"] = utils.Base64(config["ArgocdWorkflowrepositoryUsername"])
	secretValsArgo["ArgocdWorkflowrepositorySshPrivateKey"] = config["ArgocdWorkflowrepositorySshPrivateKeyBase64"].String()
	encryptedValues, err := encrypt(WORKFLOW_REPO_SECRET_TEMPLATE, secretValsArgo, config["env"], "username", "sshPrivateKey")
	if err != nil {
	    return nil, err
	}
	values["ArgocdWorkflowrepositoryUsername"] = encryptedValues["username"]
	values["ArgocdWorkflowrepositorySshPrivateKey"] = encryptedValues["sshPrivateKey"]

	secretValsAlert := make(map[string]string, 1)
	secretValsAlert["AlertsSlackUrl"] = utils.Base64(config["AlertsSlackUrl"])

	encryptedValues, err = encrypt(ALERTS_SLACK_HOOK_SECRET_KUBE_TEMPLATE, secretValsAlert, config["env"], "url")
	if err != nil {
	    return nil, err
	}
	values["AlertsSlackUrlKube"] = encryptedValues["url"]

    encryptedValues, err = encrypt(ALERTS_SLACK_HOOK_SECRET_APISERVER_TEMPLATE, secretValsAlert, config["env"], "url")
    if err != nil {
        return nil, err
    }
    values["AlertsSlackUrlApiServer"] = encryptedValues["url"]

	return values, nil
}

func encrypt(template string, values map[string]string, env Value, keys ...string) (map[string]Value, error) {

	secretAsByte, err := utils.ReplaceTemplate(values, template)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, env, keys...)
	if err != nil {
		return nil, err
	}

	return encryptedValues, nil
}
