package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type GrafanaValueBuilder struct{}

const GRAFANA_CLIENTSECRET_TEMPLATE string = `
apiVersion: v1
kind: Secret
metadata:
  name: grafana-sso-config
  namespace: grafana-operator-system
data:
  clientsecret: "{{ .KeycloakClientSecret }}"
`

func (gen *GrafanaValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	secretVals := make(map[string]string, 1)
	secretVals["KeycloakClientSecret"] = utils.Base64(config["KeycloakClientSecret"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, GRAFANA_CLIENTSECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "clientsecret")
	if err != nil {
		return nil, err
	}

	values["KeycloakClientSecret"] = encryptedValues["clientsecret"]
	values["KeycloakRealmUrl"] = config["KeycloakRealmUrl"]
	values["ClusterAdminBearerToken"] = config["ClusterAdminBearerToken"]

	return values, nil
}
