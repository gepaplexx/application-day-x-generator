package generator

import (
	"io/ioutil"

	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type ClusterMonitoringValueBuilder struct{}

const DEFAULT_ALERTMANAGER_CONFIG string = "./templates/default-alertmanager-config.tpl"
const ALERTMANAGER_SECRET_TEMPLATE string = `
apiVersion: v1
data:
  alertmanager.yaml: {{ .AlertmanagerYaml }}
kind: Secret
metadata:
  creationTimestamp: null
  name: alertmanager-main
  namespace: openshift-monitoring
`

func (gen *ClusterMonitoringValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	var alertmanagerConfigPath string
	if config["AlertmanagerConfig"].IsEmpty() {
		alertmanagerConfigPath = DEFAULT_ALERTMANAGER_CONFIG
	} else {
		alertmanagerConfigPath = config["AlertmanagerConfig"].String()
	}

	alertmanagerYaml, err := ioutil.ReadFile(alertmanagerConfigPath)
	if err != nil {
		return nil, err
	}

	alertmanagerVals := make(map[string]string, 1)
	alertmanagerVals["env"] = config["env"].String()
	alertmanagerVals["SlackChannel"] = config["SlackChannel"].String()
	alertManagerConfigByte, err := utils.ReplaceTemplate(alertmanagerVals, string(alertmanagerYaml))
	if err != nil {
		return nil, err
	}

	secretVals := make(map[string]string, 1)
	secretVals["AlertmanagerYaml"] = utils.Base64(utils.Value{Val: string(alertManagerConfigByte)})

	secretAsByte, err := utils.ReplaceTemplate(secretVals, ALERTMANAGER_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "alertmanager.yaml")
	if err != nil {
		return nil, err
	}

	values["AlertmanagerYaml"] = encryptedValues["alertmanager.yaml"]
	values["InfranodesEnabled"] = config["InfranodesEnabled"]

	return values, nil
}
