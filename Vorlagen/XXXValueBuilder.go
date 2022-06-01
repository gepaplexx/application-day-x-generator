package generator

import (
	"gepaplexx/day-x-generator/pkg/util"
)

type Value = util.Value // TODO: REMOVE ME (im package generator bereits vorhanden)

type XXXValueBuilder struct{}

func (gen *XXXValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)

	// TODO

	return values, nil
}
