package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// Função auxiliar para conversão segura
func convertToMap(data interface{}) (map[string]interface{}, error) {
	switch v := data.(type) {
	case map[string]interface{}:
		return v, nil
	case string:
		var result map[string]interface{}
		if err := json.Unmarshal([]byte(v), &result); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unsupported data type: %T", data)
	}
}

// extractTemplateVars extrai variáveis no formato ${...} de uma string.
// Ele ignora construções aninhadas (${...${...}...}), variáveis vazias (${}),
// e remove espaços em branco ao redor do nome da variável extraída.
func extractTemplateVars(s string) []string {
	var vars []string
	currentIndex := 0

	for {
		begin := strings.Index(s[currentIndex:], "${")
		if begin == -1 {
			break
		}
		begin += currentIndex

		end := strings.Index(s[begin+2:], "}")
		if end == -1 {
			break
		}
		end += begin + 2

		varName := s[begin+2 : end]

		if strings.Contains(varName, "${") {
			currentIndex = begin + 2
			continue
		}

		varName = strings.TrimSpace(varName)

		if len(varName) == 0 {
			currentIndex = end + 1
			continue
		}

		vars = append(vars, varName)
		currentIndex = end + 1
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

func getSafeValueFromPath(path string, vars map[string]interface{}) string {
	parts := strings.Split(path, ".")
	var current interface{} = vars

	// Caso especial para pathParameters
	if len(parts) > 1 && parts[0] == "input" && parts[1] == "pathParameters" {
		if _, exists := vars["input"]; !exists {
			return ""
		}
		inputMap, ok := vars["input"].(map[string]interface{})
		if !ok {
			return ""
		}
		if _, exists := inputMap["pathParameters"]; !exists {
			return ""
		}
	}

	for _, part := range parts {
		val := reflect.ValueOf(current)
		if val.Kind() != reflect.Map {
			return ""
		}

		key := reflect.ValueOf(part)
		if !key.IsValid() {
			return ""
		}

		mapVal := val.MapIndex(key)
		if !mapVal.IsValid() {
			return ""
		}
		current = mapVal.Interface()
	}

	if current == nil {
		return ""
	}

	switch v := current.(type) {
	case string:
		return v
	case map[string]interface{}, []interface{}:
		jsonBytes, _ := json.Marshal(v)
		return string(jsonBytes)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// findFieldByJSONTag busca um campo em uma struct pela sua tag JSON.
func findFieldByJSONTag(structValue reflect.Value, jsonTag string) reflect.Value {
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		tag := field.Tag.Get("json")
		// Considera a tag exata ou a parte antes da vírgula (ex: "name,omitempty")
		tagParts := strings.Split(tag, ",")
		if len(tagParts) > 0 && tagParts[0] == jsonTag {
			return structValue.Field(i)
		}
	}
	return reflect.Value{} // Retorna valor inválido se não encontrado
}
