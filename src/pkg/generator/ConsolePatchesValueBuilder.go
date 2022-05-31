package generator

type ConsolePatchesValueBuilder struct{}

func (gen *ConsolePatchesValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)
	values["env"] = config["env"]
	values["RouteNameOverride"] = config["RouteNameOverride"]
	return values, nil
}
