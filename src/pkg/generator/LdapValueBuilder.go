package generator

import (
	seal "gepaplexx/day-x-generator/pkg/sealedSecrets"
	utils "gepaplexx/day-x-generator/pkg/util"
)

type LdapValueBuilder struct{}

const LDAP_SECRET_TEMPLATE string = `
apiVersion: v1
kind: Secret
metadata:
  name: ldap-secret
  namespace: openshift-config
data:
  bindPassword: "{{ .LdapBindPassword }}"
  bindDn: "{{ .LdapBindDn }}"
  ldapUrl: "{{ .LdapUrl }}"
  usersQuery: "{{ .LdapUsersQuery }}"
`

func (gen *LdapValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)
	values["LdapEnabled"] = config["LdapEnabled"]
	if config["LdapEnabled"].Equal(false) {
		return values, nil
	}

	secretVals := make(map[string]string, 3)
	secretVals["LdapBindPassword"] = utils.Base64(config["LdapBindPassword"])
	secretVals["LdapBindDn"] = utils.Base64(config["LdapBindDn"])
	secretVals["LdapUrl"] = utils.Base64(config["LdapUrl"])
	secretVals["LdapUsersQuery"] = utils.Base64(config["LdapUsersQuery"])

	secretAsByte, err := utils.ReplaceTemplate(secretVals, LDAP_SECRET_TEMPLATE)
	if err != nil {
		return nil, err
	}

	encryptedValues, err := seal.SealValues(secretAsByte, config["env"], "bindPassword", "bindDn", "ldapUrl", "usersQuery")
	if err != nil {
		return nil, err
	}

	values["LdapBindPassword"] = encryptedValues["bindPassword"]
	values["LdapBindDn"] = encryptedValues["bindDn"]
	values["LdapUrl"] = encryptedValues["ldapUrl"]
	values["LdapUsersQuery"] = encryptedValues["usersQuery"]

	return values, nil
}