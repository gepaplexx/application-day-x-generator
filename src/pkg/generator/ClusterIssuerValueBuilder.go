package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type ClusterIssuerValueBuilder struct{}

const ROUTE53_CREDS_SECRET_TEMPLATE string = `
apiVersion: v1
kind: Secret
metadata:
  name: {{ .env }}-route53-credentials-secret
  namespace: cert-manager
data:
  secret-access-key: "{{ .SolversSecretAccessKey }}"
`

func (gen *ClusterIssuerValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	secretVals := make(map[string]string, 2)
	secretVals["env"] = config["env"].String()
	secretVals["SolversSecretAccessKey"] = utils.Base64(config["SolversSecretAccessKey"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, ROUTE53_CREDS_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "secret-access-key")
	if err != nil {
		return nil, err
	}

	values["SolversSecretAccessKey"] = encryptedValues["secret-access-key"]
	values["SolversAccesKeyId"] = config["SolversAccesKeyId"]
	values["env"] = config["env"]

	return values, nil
}
