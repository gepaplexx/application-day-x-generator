package util

import (
	"bytes"
	"fmt"
	"strings"

	b64 "encoding/base64"

	"text/template"

	"reflect"

	"gopkg.in/yaml.v3"
)

func Base64(val Value) string {
	return b64.StdEncoding.EncodeToString([]byte(val.String()))
}

func ReplaceTemplate(config map[string]string, templ string) ([]byte, error) {
	secretTemplate, err := template.New("secret").Parse(templ)
	if err != nil {
		return nil, err
	}

	var secret bytes.Buffer
	err = secretTemplate.Execute(&secret, config)
	if err != nil {
		return nil, err
	}

	return secret.Bytes(), nil
}

// TODO FindValues anpassen/zusammenführen

// data 	=> YAML formatted String
// ks... 	=> YAML path, eg. metadata-name
func FindValues(data []byte, keys ...string) (map[string]Value, error) {
	yamlAsMap := make(map[string]any)
	err := yaml.Unmarshal(data, &yamlAsMap)
	if err != nil {
		return nil, err
	}

	values := make(map[string]Value)
	for _, val := range keys {
		res, err := findValue(yamlAsMap, strings.Split(val, ":")...)
		valWithOutPrefix := strings.TrimPrefix(val, "spec:encryptedData:")
		if err != nil {
			values[valWithOutPrefix] = Value{"not found"}
		} else {
			values[valWithOutPrefix] = Value{fmt.Sprintf("%v", res)}
		}
	}

	return values, nil
}

func FindValue(data []byte, key string) (any, error) {
	yamlAsMap := make(map[string]any)
	err := yaml.Unmarshal(data, &yamlAsMap)
	if err != nil {
		return nil, err
	}

	res, err := findValue(yamlAsMap, strings.Split(key, ".")...)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func FindValuesFlatMap(data []byte, keys ...string) (map[string]Value, error) {
	yamlAsMap := make(map[string]any)
	err := yaml.Unmarshal(data, &yamlAsMap)
	if err != nil {
		return nil, err
	}

	values := make(map[string]Value)
	for _, val := range keys {
		res, err := findValue(yamlAsMap, strings.Split(val, ".")...)
		if err != nil {
			continue
		}

		v := reflect.ValueOf(res)
		switch v.Kind() {
		case reflect.String:
			values[val] = Value{v.String()}
		case reflect.Map:
			for _, key := range v.MapKeys() {
				strct := v.MapIndex(key)
				values[key.String()] = Value{strct.Interface()}
			}
		default:
			values[val] = Value{"not found"}
		}
	}

	return values, nil
}

func findValue(m map[string]any, keys ...string) (rval any, err error) {
	var ok bool
	if len(keys) == 0 { // degenerate input
		return nil, fmt.Errorf("NestedMapLookup needs at least one key")
	}
	if rval, ok = m[keys[0]]; !ok {
		return nil, fmt.Errorf("key not found; remaining keys: %v", keys)
	} else if len(keys) == 1 { // we've reached the final key
		return rval, nil
	} else if m, ok = rval.(map[string]any); !ok {
		return nil, fmt.Errorf("malformed structure at %#v", rval)
	} else { // 1+ more keys
		return findValue(m, keys[1:]...)
	}
}
