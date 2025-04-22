package utils_test

import (
	"testing"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestToJsonString(t *testing.T) {
	t.Run("Cenário feliz: conversão de mapa para JSON string", func(t *testing.T) {
		data := map[string]interface{}{"name": "John", "age": 30}
		expected := `{"name":"John","age":30}`
		result := utils.ToJsonString(data)
		assert.JSONEq(t, expected, result)
	})

	t.Run("Cenário: input já é string", func(t *testing.T) {
		data := `{"status": "ok"}`
		expected := `{"status": "ok"}`
		result := utils.ToJsonString(data)
		assert.Equal(t, expected, result)
	})
}

func TestConvertMapToJSON(t *testing.T) {
	t.Run("Cenário feliz: conversão de map[interface{}]interface{} para JSON string", func(t *testing.T) {
		input := map[interface{}]interface{}{1: "one", "two": 2.0, true: "yes"}
		expected := `{"1":"one","two":2,"true":"yes"}`
		result, err := utils.ConvertMapToJSON(input)
		assert.NoError(t, err)
		assert.JSONEq(t, expected, result)
	})

	t.Run("Cenário: mapa vazio", func(t *testing.T) {
		input := map[interface{}]interface{}{}
		expected := `{}`
		result, err := utils.ConvertMapToJSON(input)
		assert.NoError(t, err)
		assert.JSONEq(t, expected, result)
	})
}

func TestConvertAnyMapToJSON(t *testing.T) {
	t.Run("Cenário feliz: conversão de map[string]interface{} para JSON string", func(t *testing.T) {
		input := map[string]interface{}{"key": "value", "number": 1}
		expected := `{"key":"value","number":1}`
		result, err := utils.ConvertAnyMapToJSON(input)
		assert.NoError(t, err)
		assert.JSONEq(t, expected, result)
	})

	t.Run("Cenário: input é string JSON válida", func(t *testing.T) {
		input := `{"status":"pending","id":123}`
		expected := `{"status":"pending","id":123}`
		result, err := utils.ConvertAnyMapToJSON(input)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestSetNestedValue(t *testing.T) {
	t.Run("Cenário feliz: definir valor em caminho aninhado existente", func(t *testing.T) {
		data := map[string]interface{}{"user": map[string]interface{}{"address": map[string]interface{}{"city": "Old City"}}}
		path := "user.address.city"
		value := "New City"
		err := utils.SetNestedValue(data, path, value)
		assert.NoError(t, err)
		assert.Equal(t, value, data["user"].(map[string]interface{})["address"].(map[string]interface{})["city"])
	})

	t.Run("Cenário: definir valor em caminho aninhado inexistente (criação de mapas)", func(t *testing.T) {
		data := map[string]interface{}{}
		path := "user.address.city"
		value := "New City"
		err := utils.SetNestedValue(data, path, value)
		assert.NoError(t, err)
		assert.Equal(t, value, data["user"].(map[string]interface{})["address"].(map[string]interface{})["city"])
	})
}

func TestGetNestedValue(t *testing.T) {
	data := map[string]interface{}{"user": map[string]interface{}{"profile": map[string]interface{}{"name": "John"}}}

	t.Run("Cenário feliz: obter valor de caminho aninhado existente", func(t *testing.T) {
		path := "user.profile.name"
		expected := "John"
		result, ok := utils.GetNestedValue(data, path)
		assert.True(t, ok)
		assert.Equal(t, expected, result)
	})

	t.Run("Cenário: caminho inexistente", func(t *testing.T) {
		path := "user.profile.age"
		result, ok := utils.GetNestedValue(data, path)
		assert.False(t, ok)
		assert.Nil(t, result)
	})
}

func TestToStringMap(t *testing.T) {
	t.Run("Cenário feliz: converter map[string]interface{} para map[string]string", func(t *testing.T) {
		data := map[string]interface{}{"name": "Alice", "age": 25}
		expected := map[string]string{"name": "Alice", "age": "25"}
		result := utils.ToStringMap(data)
		assert.Equal(t, expected, result)
	})

	t.Run("Cenário: input nulo", func(t *testing.T) {
		var data interface{} = nil
		expected := map[string]string{}
		result := utils.ToStringMap(data)
		assert.Equal(t, expected, result)
	})
}

func TestGetEnvVarsMap(t *testing.T) {
	t.Run("Cenário: retorna um mapa de variáveis de ambiente (placeholder)", func(t *testing.T) {
		result := utils.GetEnvVarsMap()
		assert.Contains(t, result, "EXAMPLE_API_KEY")
	})

	t.Run("Cenário: mapa de variáveis de ambiente não está vazio (placeholder)", func(t *testing.T) {
		result := utils.GetEnvVarsMap()
		assert.NotEmpty(t, result)
	})
}

func TestGetSecretsMap(t *testing.T) {
	t.Run("Cenário: retorna um mapa de segredos (placeholder)", func(t *testing.T) {
		result := utils.GetSecretsMap()
		assert.Contains(t, result, "DB_PASSWORD")
	})

	t.Run("Cenário: mapa de segredos não está vazio (placeholder)", func(t *testing.T) {
		result := utils.GetSecretsMap()
		assert.NotEmpty(t, result)
	})
}
