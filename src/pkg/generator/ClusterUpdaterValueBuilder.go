package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type ClusterUpdaterValueBuilder struct {
}

const CLUSTER_UPDATE_SECRET_TEMPLATE string = `
kind: Secret
apiVersion: v1
metadata:
  name: cluster-update-secret
  namespace: gp-infrastructure
data:
  slackTargetChannel: {{ .SlackChannel }}
type: Opaque

`

func (gen *ClusterUpdaterValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	secretVals := make(map[string]string, 3)
	secretVals["SlackChannel"] = utils.Base64(config["SlackChannel"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, CLUSTER_UPDATE_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "slackTargetChannel")
	if err != nil {
		return nil, err
	}

	values["env"] = config["env"]
	values["SlackChannel"] = encryptedValues["slackTargetChannel"]

	return values, nil
}
