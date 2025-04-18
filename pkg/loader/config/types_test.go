package config

import (
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
	"gopkg.in/yaml.v2"
)

var filepath = "/Users/macmini/Documents/workspace/go/aws-lowcode-lambda-go/examples/template.yaml"

func load() *LambdaConfig {
	var lambda LambdaConfig

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil
	}

	err = yaml.Unmarshal(data, &lambda)
	if err != nil {
		return nil
	}

	return &lambda
}

func TestLambdaConfigTemplate(t *testing.T) {
	data := load()

	t.Run("Deve retornar a versão de template 2025-04-19", func(t *testing.T) {
		assert.Equal(t, data.TemplateFormatVersion, "2025-04-19")
	})

	t.Run("Deve retornar o transform do lowcode lambda: LC::Lambda", func(t *testing.T) {
		assert.Equal(t, data.Transform, "LC::Lambda")
	})
}

func TestLambdaConfigMetadata(t *testing.T) {
	data := load()

	t.Run("Deve retornar o nome do serviço", func(t *testing.T) {
		assert.Equal(t, data.Metadata.Name, "user-service")
	})
}
