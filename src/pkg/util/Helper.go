package util

import (
	"fmt"
	"strings"

	"reflect"

	"gopkg.in/yaml.v3"
)

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
