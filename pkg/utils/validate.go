package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type JSONValidator struct{}

var templateRegex = regexp.MustCompile(`\${([^}]+)}`)

func NewJSONValidator() *JSONValidator {
	return &JSONValidator{}
}

func (v *JSONValidator) Validate(data interface{}, schema interface{}) error {
	// Converter data para map[string]interface{} se necessário
	dataMap, err := convertToMap(data)
	if err != nil {
		return fmt.Errorf("invalid data format: %v", err)
	}

	schemaLoader := gojsonschema.NewGoLoader(schema)
	documentLoader := gojsonschema.NewGoLoader(dataMap)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		return fmt.Errorf("validation errors: %v", result.Errors())
	}

	return nil
}

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

func ResolveTemplate(template string, variables map[string]interface{}) interface{} {
	if !strings.Contains(template, "${") {
		return template
	}

	return templateRegex.ReplaceAllStringFunc(template, func(match string) string {
		path := strings.TrimPrefix(match, "${")
		path = strings.TrimSuffix(path, "}")
		return getSafeValueFromPath(path, variables)
	})
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
