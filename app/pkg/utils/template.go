package utils

import (
	"fmt"
	"reflect"
)

// Função auxiliar para converter map[interface{}]interface{} para map[string]interface{}
func ConvertToMapStringInterface(data interface{}) (map[string]interface{}, bool) {
	val := reflect.ValueOf(data)
	if val.Kind() != reflect.Map {
		return nil, false
	}

	result := make(map[string]interface{})
	for _, key := range val.MapKeys() {
		keyStr := fmt.Sprintf("%v", key.Interface())
		result[keyStr] = val.MapIndex(key).Interface()
	}
	return result, true
}

func ResolveTemplates(data interface{}, vars map[string]interface{}) interface{} {
	// Converter para map[string]interface{} se for map[interface{}]interface{}
	if m, ok := ConvertToMapStringInterface(data); ok {
		return resolveMapTemplates(m, vars)
	}

	switch v := data.(type) {
	case string:
		return resolveStringTemplates(v, vars)
	case map[string]interface{}:
		return resolveMapTemplates(v, vars)
	case []interface{}:
		return resolveSliceTemplates(v, vars)
	default:
		return v
	}
}
