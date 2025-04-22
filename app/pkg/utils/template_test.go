package utils_test

import (
	"testing"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestConvertToMapStringInterface(t *testing.T) {
	t.Run("Cenário feliz: map[interface{}]interface{}", func(t *testing.T) {
		input := map[interface{}]interface{}{
			1:      "um",
			"dois": 2,
		}
		expected := map[string]interface{}{
			"1":    "um",
			"dois": 2,
		}
		result, ok := utils.ConvertToMapStringInterface(input)
		assert.True(t, ok)
		assert.Equal(t, expected, result)
	})

	t.Run("Cenário: input já é map[string]interface{}", func(t *testing.T) {
		input := map[string]interface{}{
			"chave1": "valor1",
			"chave2": 123,
		}
		result, ok := utils.ConvertToMapStringInterface(input)
		assert.True(t, ok)
		assert.Equal(t, input, result)
	})
}

func TestResolveTemplates(t *testing.T) {
	vars := map[string]interface{}{
		"name": "World",
		"age":  30,
	}

	t.Run("Cenário: string com um template", func(t *testing.T) {
		input := "Hello, ${name}!"
		expected := "Hello, World!"
		result := utils.ResolveTemplates(input, vars)
		assert.Equal(t, expected, result)
	})

	t.Run("Cenário: map com um valor de template", func(t *testing.T) {
		input := map[string]interface{}{
			"greeting": "Hello, ${name}!",
			"age":      30,
		}
		expected := map[string]interface{}{
			"greeting": "Hello, World!",
			"age":      30,
		}
		result := utils.ResolveTemplates(input, vars)
		assert.Equal(t, expected, result)
	})
}
