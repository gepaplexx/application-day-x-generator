package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type GoogleOAuthValueBuilder struct{}

const GOOGLE_SECRET_TEMPLATE string = `
apiVersion: v1
kind: Secret
metadata:
  name: google-oauth-secret
  namespace: openshift-config
data:
  clientId: "{{ .GoogleClientId }}"
  clientSecret: "{{ .GoogleClientSecret }}"
  restrDomain: "{{ .GoogleRestrictedDomain }}"
`

func (gen *GoogleOAuthValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)
	values["GoogleEnabled"] = config["GoogleEnabled"]
	if config["GoogleEnabled"].Equal(false) {
		return values, nil
	}

	secretVals := make(map[string]string, 3)
	secretVals["GoogleClientId"] = utils.Base64(config["GoogleClientId"])
	secretVals["GoogleClientSecret"] = utils.Base64(config["GoogleClientSecret"])
	secretVals["GoogleRestrictedDomain"] = utils.Base64(config["GoogleRestrictedDomain"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, GOOGLE_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "clientId", "clientSecret", "restrDomain")
	if err != nil {
		return nil, err
	}

	values["GoogleClientId"] = encryptedValues["clientId"]
	values["GoogleClientSecret"] = encryptedValues["clientSecret"]
	values["GoogleRestrictedDomain"] = encryptedValues["restrDomain"]

	return values, nil
}
