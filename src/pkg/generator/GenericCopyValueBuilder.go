package generator

import utils "gepaplexx/day-x-generator/pkg/util"

type GenericCopyValueBuilder struct{}

func (gen *GenericCopyValueBuilder) GetValues(config map[string]utils.Value) (map[string]utils.Value, error) {
	values := make(map[string]utils.Value)
	for key, val := range config {
		values[key] = val
	}
	return values, nil
}
