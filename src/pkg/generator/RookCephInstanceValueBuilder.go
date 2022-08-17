package generator

type RookCephInstanceValueBuilder struct{}

func (gen *RookCephInstanceValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)
	values["RookCephInstanceEnabled"] = config["RookCephInstanceEnabled"]
	return values, nil
}
