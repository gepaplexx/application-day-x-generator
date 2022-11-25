package generator

import (
	"io/ioutil"

	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type ClusterMonitoringValueBuilder struct{}

const DEFAULT_ALERTMANAGER_CONFIG string = "./templates/default-alertmanager-config.yaml.tpl"
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
const REMOTE_WRITE_SECRET_TEMPLATE string = `
apiVersion: v1
kind: Secret
metadata:
  name: hub-remote-write-authentication
  namespace: openshift-monitoring
data:
  password: {{ .PrometheusRemoteWritePassword }}
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
	alertmanagerVals["SlackChannelCritical"] = config["SlackChannelCritical"].String()
	alertmanagerVals["SlackChannelMonitoringInternalApplications"] = config["SlackChannelMonitoringInternalApplications"].String()
	alertManagerConfigByte, err := utils.ReplaceTemplate(alertmanagerVals, string(alertmanagerYaml))
	if err != nil {
		return nil, err
	}

	secretVals := make(map[string]string, 2)
	secretVals["AlertmanagerYaml"] = utils.Base64(utils.Value{Val: string(alertManagerConfigByte)})
	secretVals["PrometheusRemoteWritePassword"] = utils.Base64(config["PrometheusRemoteWritePassword"])

	secretAsByteAlertmanager, err := utils.ReplaceTemplate(secretVals, ALERTMANAGER_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}
	encryptedValuesAlertmanager, err := seal.SealValues(secretAsByteAlertmanager, config["env"], "alertmanager.yaml")
	if err != nil {
		return nil, err
	}

	secretAsByteRemoteWrite, err := utils.ReplaceTemplate(secretVals, REMOTE_WRITE_SECRET_TEMPLATE)
	encryptedValuesRemoteWrite, err := seal.SealValues(secretAsByteRemoteWrite, config["env"], "password")

	values["AlertmanagerYaml"] = encryptedValuesAlertmanager["alertmanager.yaml"]
	values["PrometheusRemoteWritePassword"] = encryptedValuesRemoteWrite["password"]
	values["InfranodesEnabled"] = config["InfranodesEnabled"]

	return values, nil
}
