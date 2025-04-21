package engine_test

import (
	"testing"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/core/engine"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"
	"github.com/stretchr/testify/assert"
)

func TestNewExecutionContext(t *testing.T) {
	cfg := &config.LambdaConfig{
		TemplateFormatVersion: "2025-04-19",
		Transform:             "LC::Lambda",
		Metadata: config.Metadata{
			Name: "TestConfig",
		},
		Input: config.Input{
			Source: "test",
			Schema: map[string]any{
				"id":   map[string]any{"type": "string"},
				"name": map[string]any{"type": "string"},
			},
		},
		Steps: []config.Step{
			{ID: "step1", Type: "test"},
		},
	}

	t.Run("Cenário feliz: evento com body string JSON", func(t *testing.T) {
		event := map[string]interface{}{
			"body": `{"id": "123", "name": "Test Event"}`,
			"headers": map[string]interface{}{
				"Content-Type": "application/json",
			},
		}

		expectedInput := map[string]interface{}{
			"body": map[string]interface{}{
				"id":   "123",
				"name": "Test Event",
			},
			"headers": map[string]interface{}{
				"Content-Type": "application/json",
			},
		}

		ctx := engine.NewExecutionContext(cfg, event)

		assert.NotNil(t, ctx)
		assert.Equal(t, cfg, ctx.Config)
		assert.Equal(t, expectedInput, ctx.Input)
		assert.NotNil(t, ctx.Variables)
		assert.Empty(t, ctx.Variables)
		assert.NotNil(t, ctx.Output)
		assert.Empty(t, ctx.Output)
	})

	t.Run("Cenário: evento com body já como map", func(t *testing.T) {
		event := map[string]interface{}{
			"body": map[string]interface{}{
				"id":   "456",
				"name": "Another Event",
			},
			"headers": map[string]interface{}{
				"Content-Type": "application/json",
			},
		}

		ctx := engine.NewExecutionContext(cfg, event)

		assert.NotNil(t, ctx)
		assert.Equal(t, cfg, ctx.Config)
		assert.Equal(t, event, ctx.Input)
		assert.NotNil(t, ctx.Variables)
		assert.Empty(t, ctx.Variables)
		assert.NotNil(t, ctx.Output)
		assert.Empty(t, ctx.Output)
	})

	t.Run("Cenário: evento sem body", func(t *testing.T) {
		event := map[string]interface{}{
			"headers": map[string]interface{}{
				"Content-Type": "application/json",
			},
		}

		ctx := engine.NewExecutionContext(cfg, event)

		assert.NotNil(t, ctx)
		assert.Equal(t, cfg, ctx.Config)
		assert.Equal(t, event, ctx.Input)
		assert.NotNil(t, ctx.Variables)
		assert.Empty(t, ctx.Variables)
		assert.NotNil(t, ctx.Output)
		assert.Empty(t, ctx.Output)
	})
}
