package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
	"github.com/google/uuid"
)

type KeycloakInstanceValueBuilder struct{}

const KEYCLOAK_SECRET_TEMPLATE string = `
apiVersion: v1
kind: Secret
metadata:
  name: keycloak-config
  namespace: gp-sso
data:
  password:  {{ .KeycloakDbPassword }}
  clientSecret:  {{ .KeycloakOcpClientSecret }}
type: Opaque

`

func (gen *KeycloakInstanceValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)
	cs := config["KeycloakOcpClientSecret"].String()
	if cs == utils.NOT_VALID {
		config["KeycloakOcpClientSecret"] = Value{Val: uuid.New().String()}
	}

	secretVals := make(map[string]string, 1)
	secretVals["KeycloakOcpClientSecret"] = utils.Base64(config["KeycloakOcpClientSecret"])
	secretVals["KeycloakDbPassword"] = utils.Base64(config["KeycloakDbPassword"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, KEYCLOAK_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "password", "clientSecret")
	if err != nil {
		return nil, err
	}

	values["KeycloakOcpClientSecret"] = encryptedValues["clientSecret"]
	values["KeycloakDbPassword"] = encryptedValues["password"]

	return values, nil
}
