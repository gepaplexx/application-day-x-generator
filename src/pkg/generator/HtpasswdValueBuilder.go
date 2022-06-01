package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type HtpasswdValueBuilder struct{}

const HTPASSWD_SECRET_TEMPLATE string = `
kind: Secret
apiVersion: v1
metadata:
  name: htpass-secret
  namespace: openshift-config
data:
  htpasswd: >-
    {{ .HtpasswdData }}
type: Opaque
`

func (gen *HtpasswdValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)
	values["HtpasswdEnabled"] = config["HtpasswdEnabled"]
	if config["HtpasswdEnabled"].Equal(false) {
		return values, nil
	}

	secretVals := make(map[string]string, 1)
	secretVals["HtpasswdData"] = utils.Base64(config["HtpasswdData"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, HTPASSWD_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "htpasswd")
	if err != nil {
		return nil, err
	}

	values["HtpasswdData"] = encryptedValues["htpasswd"]

	return values, nil
}
