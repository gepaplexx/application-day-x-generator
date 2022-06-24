package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type RookCephInstanceValueBuilder struct{}

const GP_ROOK_CEPH_INSTANCE_SECRET_TEMPLATE string = `
kind: Secret
apiVersion: v1
metadata:
  name: slack-url
  namespace: rook-ceph
data:
  url: {{ .SlackChannel }}
type: Opaque

`

func (gen *RookCephInstanceValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)
	values["RookCephInstanceEnabled"] = config["RookCephInstanceEnabled"]

	secretVals := make(map[string]string, 1)
	secretVals["SlackChannel"] = utils.Base64(config["SlackChannel"])
	secretAsByte, err := utils.ReplaceTemplate(secretVals, GP_ROOK_CEPH_INSTANCE_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "url")
	if err != nil {
		return nil, err
	}
	values["SlackChannel"] = encryptedValues["url"]

	return values, nil
}
