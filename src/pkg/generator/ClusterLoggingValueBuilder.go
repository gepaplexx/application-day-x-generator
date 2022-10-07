package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type ClusterLoggingValueBuilder struct{}

const LOGGING_BACKEND_SECRET_TEMPLATE string = `
apiVersion: v1
data:
  secretkey: {{ .ClusterLoggingS3SecretKey }}
kind: Secret
metadata:
  creationTimestamp: null
  name: loki-backend-secret
  namespace: openshift-logging
`

func (gen *ClusterLoggingValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)
	values["ClusterLoggingEnabled"] = config["ClusterLoggingEnabled"]

	secretVals := make(map[string]string, 1)
	secretVals["ClusterLoggingS3SecretKey"] = utils.Base64(config["ClusterLoggingS3SecretKey"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, LOGGING_BACKEND_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "secretkey")
	if err != nil {
		return nil, err
	}

	values["ClusterLoggingS3SecretKey"] = encryptedValues["secretkey"]

	return values, nil
}
