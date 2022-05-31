package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type GepaplexxCicdToolsValueBuilder struct {
}

const GP_CICD_SECRET_TEMPLATE string = `
kind: Secret
apiVersion: v1
metadata:
  name: gepaplexx-cicd-tools-postgresql
  namespace: gepaplexx-cicd-tools
data:
  password: {{ .PostgresqlPassword }}
  postgres-password: {{ .PostgresqlPostgresPassword }}
type: Opaque

`

func (gen *GepaplexxCicdToolsValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	secretVals := make(map[string]string, 3)
	secretVals["PostgresqlPassword"] = utils.Base64(config["PostgresqlPassword"])
	secretVals["PostgresqlPostgresPassword"] = utils.Base64(config["PostgresqlPostgresPassword"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, GP_CICD_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "password", "postgres-password")
	if err != nil {
		return nil, err
	}

	values["PostgresqlPassword"] = encryptedValues["password"]
	values["PostgresqlPostgresPassword"] = encryptedValues["postgres-password"]

	return values, nil
}
