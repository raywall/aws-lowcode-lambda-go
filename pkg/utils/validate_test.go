package utils_test

import (
	"testing"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestJSONValidator_Validate(t *testing.T) {
	validator := utils.NewJSONValidator()

	t.Run("Cenário feliz: dados válidos contra o schema", func(t *testing.T) {
		schema := map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{"type": "string"},
				"age":  map[string]interface{}{"type": "integer"},
			},
			"required": []string{"name"},
		}

		data := map[string]interface{}{
			"name": "John Doe",
			"age":  30,
		}

		err := validator.Validate(data, schema)
		assert.NoError(t, err)
	})

	t.Run("Cenário de erro: dados inválidos contra o schema (campo ausente)", func(t *testing.T) {
		schema := map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{"type": "string"},
				"age":  map[string]interface{}{"type": "integer"},
			},
			"required": []string{"name"},
		}

		data := map[string]interface{}{
			"age": 30,
		}

		err := validator.Validate(data, schema)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation errors")
	})
}

func TestResolveTemplate(t *testing.T) {
	variables := map[string]interface{}{
		"name":  "World",
		"value": 123,
	}

	t.Run("Cenário feliz: template com variável existente", func(t *testing.T) {
		template := "Hello, ${name}!"
		expected := "Hello, World!"
		result := utils.ResolveTemplate(template, variables)
		assert.Equal(t, expected, result)
	})

	t.Run("Cenário: template sem variáveis para substituir", func(t *testing.T) {
		template := "This is a static string."
		expected := "This is a static string."
		result := utils.ResolveTemplate(template, variables)
		assert.Equal(t, expected, result)
	})

	t.Run("Cenário: template com variável inexistente", func(t *testing.T) {
		template := "The value is ${unknown}."
		expected := "The value is ." // getSafeValueFromPath retorna "" para inexistente
		result := utils.ResolveTemplate(template, variables)
		assert.Equal(t, expected, result)
	})

	t.Run("Cenário: template com múltiplos placeholders", func(t *testing.T) {
		template := "Name: ${name}, Value: ${value}"
		expected := "Name: World, Value: 123"
		result := utils.ResolveTemplate(template, variables)
		assert.Equal(t, expected, result)
	})
}
