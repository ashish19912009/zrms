package logger

type LogKey string

func BaseLogContext(layer string, layer_name string, method string, method_name string) map[string]interface{} {
	return map[string]interface{}{
		string(layer):  layer_name,
		string(method): method_name,
	}
}
