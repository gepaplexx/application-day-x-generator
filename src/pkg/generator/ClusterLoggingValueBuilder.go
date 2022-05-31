package generator

type ClusterLoggingValueBuilder struct{}

func (gen *ClusterLoggingValueBuilder) GetValues(config map[string]Value) (map[string]Value, error) {
	values := make(map[string]Value)
	values["ClusterLoggingEnabled"] = config["ClusterLoggingEnabled"]
	values["ClusterLoggingRequestMemory"] = config["ClusterLoggingRequestMemory"]
	values["ClusterLoggingLimitsMemory"] = config["ClusterLoggingLimitsMemory"]
	return values, nil
}
