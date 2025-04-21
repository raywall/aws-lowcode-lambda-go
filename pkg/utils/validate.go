package utils

import (
	"fmt"
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
