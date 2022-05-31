package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type GitOAuthValueBuilder struct{}

const GIT_SECRET_TEMPLATE string = `
apiVersion: v1
kind: Secret
metadata:
  name: github-oauth-secret
  namespace: openshift-config
data:
  clientId: "{{ .GitClientId }}"
  clientSecret: "{{ .GitClientSecret }}"
  restrOrgs: "{{ .GitRestrOrgs }}"
`

func (gen *GitOAuthValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)
	values["GitEnabled"] = config["GitEnabled"]
	if config["GitEnabled"].Equal(false) {
		return values, nil
	}

	secretVals := make(map[string]string, 3)
	secretVals["GitClientId"] = utils.Base64(config["GitClientId"])
	secretVals["GitClientSecret"] = utils.Base64(config["GitClientSecret"])
	secretVals["GitRestrOrgs"] = utils.Base64(config["GitRestrOrgs"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, GIT_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "clientId", "clientSecret", "restrOrgs")
	if err != nil {
		return nil, err
	}

	values["GitClientId"] = encryptedValues["clientId"]
	values["GitClientSecret"] = encryptedValues["clientSecret"]
	values["GitRestrOrgs"] = encryptedValues["restrOrgs"]

	return values, nil
}
