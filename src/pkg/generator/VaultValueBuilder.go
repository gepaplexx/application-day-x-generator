package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type VaultValueBuilder struct {
}

const VAULT_UNSEAL_SECRET_TEMPLATE string = `
kind: Secret
apiVersion: v1
metadata:
  name: gcp-creds
  namespace: gp-vault
data:
  creds.json: {{ .AutoUnsealCreds }}
type: Opaque

`

func (gen *VaultValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	secretVals := make(map[string]string, 1)
	secretVals["AutoUnsealCreds"] = utils.Base64(config["AutoUnsealCreds"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, VAULT_UNSEAL_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "creds.json")
	if err != nil {
		return nil, err
	}

	values["AutoUnsealCreds"] = encryptedValues["creds.json"]

	return values, nil
}
