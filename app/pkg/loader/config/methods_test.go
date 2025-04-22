package config_test

import (
	"testing"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"
	"github.com/stretchr/testify/assert"
)

func TestGetSchema(t *testing.T) {
	t.Run("Cenário: referência válida ", func(t *testing.T) {
		var (
			ref  = "#/input/schema/name"
			want = map[string]any{"type": "string"}

			data = config.LambdaConfig{
				Input: config.Input{
					Schema: map[string]any{
						"name": map[string]any{"type": "string"},
						"age":  map[string]any{"type": "integer"},
					},
				},
			}
		)

		assert.Equal(t, data.GetSchema(ref), want)
	})

	t.Run("Cennário: referência válida para outro campo", func(t *testing.T) {
		var (
			ref  = "#/input/schema/age"
			want = map[string]any{"type": "integer"}

			data = config.LambdaConfig{
				Input: config.Input{
					Schema: map[string]any{
						"name": map[string]any{"type": "string"},
						"age":  map[string]any{"type": "integer"},
					},
				},
			}
		)

		assert.Equal(t, data.GetSchema(ref), want)
	})

	t.Run("Cennário: referência inválida - prefixo incorreto", func(t *testing.T) {
		var (
			ref = "input/schema/name"

			data = config.LambdaConfig{
				Input: config.Input{
					Schema: map[string]any{
						"name": map[string]any{"type": "string"},
					},
				},
			}
		)

		assert.Nil(t, data.GetSchema(ref))
	})

	t.Run("Cennário: referência inválida - chave não existente", func(t *testing.T) {
		var (
			ref = "#/input/schema/email"

			data = config.LambdaConfig{
				Input: config.Input{
					Schema: map[string]any{
						"name": map[string]any{"type": "string"},
					},
				},
			}
		)

		assert.Nil(t, data.GetSchema(ref))
	})

	t.Run("Cennário: schema vazio", func(t *testing.T) {
		var (
			ref = "#/input/schema/any"

			data = config.LambdaConfig{
				Input: config.Input{
					Schema: map[string]any{},
				},
			}
		)

		assert.Nil(t, data.GetSchema(ref))
	})

	t.Run("Cennário: input vazio", func(t *testing.T) {
		var (
			ref = "#/input/schema/any"

			data = config.LambdaConfig{
				Input: config.Input{},
			}
		)

		assert.Nil(t, data.GetSchema(ref))
	})

	t.Run("Cennário: referência com espaços extras", func(t *testing.T) {
		var (
			ref  = "#/input/schema/field_with_spaces"
			want = "value"

			data = config.LambdaConfig{
				Input: config.Input{
					Schema: map[string]any{
						"field_with_spaces": "value",
					},
				},
			}
		)

		assert.Equal(t, data.GetSchema(ref), want)
	})
}
