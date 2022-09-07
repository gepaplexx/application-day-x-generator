package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type GepaplexxCicdToolsValueBuilder struct {
}

const GP_CICD_DB_SECRET_TEMPLATE string = `
kind: Secret
apiVersion: v1
metadata:
  name: gepaplexx-cicd-tools-postgresql
  namespace: gepaplexx-cicd-tools
data:
  password: {{ .PostgresqlPassword }}
  postgres-password: {{ .PostgresqlPostgresPassword }}
  username: {{ .PostgresqlUsername }}
type: Opaque

`

const GP_CICD_SSO_SECRET_TEMPLATE string = `
apiVersion: v1
kind: Secret
metadata:
  name: argo-workflows-sso
  namespace: gepaplexx-cicd-tools
data:
  client-secret: {{ .ArgoWorkflowsSsoClientSecret }}

`

func (gen *GepaplexxCicdToolsValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	secretValsDb := make(map[string]string, 3)
	secretValsDb["PostgresqlPassword"] = utils.Base64(config["PostgresqlPassword"])
	secretValsDb["PostgresqlPostgresPassword"] = utils.Base64(config["PostgresqlPostgresPassword"])
	secretValsDb["PostgresqlUsername"] = utils.Base64(config["PostgresqlUsername"])

	secretAsByteDb, err := utils.ReplaceTemplate(secretValsDb, GP_CICD_DB_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}
	encryptedValuesDb, err := seal.SealValues(secretAsByteDb, config["env"], "password", "postgres-password", "username")
	if err != nil {
		return nil, err
	}

	secretValsSso := make(map[string]string, 1)
	secretValsSso["ArgoWorkflowsSsoClientSecret"] = utils.Base64(config["ArgoWorkflowsSsoClientSecret"])
	secretAsByteSso, err := utils.ReplaceTemplate(secretValsSso, GP_CICD_SSO_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}
	encryptedValuesSso, err := seal.SealValues(secretAsByteSso, config["env"], "client-secret")
	if err != nil {
		return nil, err
	}

	values["PostgresqlPassword"] = encryptedValuesDb["password"]
	values["PostgresqlPostgresPassword"] = encryptedValuesDb["postgres-password"]
	values["PostgresqlUsername"] = encryptedValuesDb["username"]
	values["ArgoWorkflowsSsoClientSecret"] = encryptedValuesSso["client-secret"]
	values["ArgoWorkflowsSsoIssuer"] = config["ArgoWorkflowsSsoIssuer"]
	values["ArgoWorkflowsClusterScopedGroupEnabled"] = config["ArgoWorkflowsClusterScopedGroupEnabled"]

	return values, nil
}
