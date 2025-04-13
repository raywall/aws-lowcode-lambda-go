package utils

import (
	"fmt"
	"reflect"
	"strings"
)

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

func resolveStringTemplates(s string, vars map[string]interface{}) interface{} {
	if !strings.Contains(s, "${") {
		return s
	}

	templateVars := extractTemplateVars(s)
	if len(templateVars) == 0 {
		return s
	}

	// Se a string é exatamente um template (ex: "${input.body}")
	if len(templateVars) == 1 && strings.TrimSpace(s) == "${"+templateVars[0]+"}" {
		return getDeepValue(vars, templateVars[0])
	}

	// Caso contrário, substitui os templates na string
	result := s
	for _, varPath := range templateVars {
		value := getDeepValue(vars, varPath)
		result = strings.ReplaceAll(result, "${"+varPath+"}", fmt.Sprintf("%v", value))
	}
	return result
}

func resolveMapTemplates(m map[string]interface{}, vars map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range m {
		result[k] = ResolveTemplates(v, vars)
	}
	return result
}

func resolveSliceTemplates(s []interface{}, vars map[string]interface{}) []interface{} {
	result := make([]interface{}, len(s))
	for i, v := range s {
		result[i] = ResolveTemplates(v, vars)
	}
	return result
}

func extractTemplateVars(s string) []string {
	var vars []string
	start := 0
	for {
		begin := strings.Index(s[start:], "${")
		if begin == -1 {
			break
		}
		begin += start
		end := strings.Index(s[begin:], "}")
		if end == -1 {
			break
		}
		end += begin
		vars = append(vars, s[begin+2:end])
		start = end + 1
	}
	return vars
}

func getDeepValue(data interface{}, path string) interface{} {
	parts := strings.Split(path, ".")
	var current interface{} = data

	for _, part := range parts {
		if current == nil {
			return nil
		}

		val := reflect.ValueOf(current)

		// Trata maps genéricos
		if val.Kind() == reflect.Map {
			key := reflect.ValueOf(part)
			if !key.IsValid() {
				return nil
			}

			mapVal := val.MapIndex(key)
			if !mapVal.IsValid() {
				return nil
			}

			current = mapVal.Interface()
		} else {
			return nil
		}
	}

	return current
}
