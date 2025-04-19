package context_test

import (
	"testing"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/core/context"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"
	"github.com/stretchr/testify/assert"
)

func TestExecutionContext(t *testing.T) {
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

	inputData := map[string]interface{}{
		"body": map[string]interface{}{
			"id":   "123",
			"name": "Test Input",
		},
		"headers": map[string]interface{}{
			"Content-Type": "application/json",
		},
	}

	variablesData := map[string]interface{}{
		"executionId": "abc-123",
		"startTime":   "2025-04-19T18:00:00Z",
	}

	outputData := map[string]interface{}{
		"result": "success",
		"code":   200,
	}

	tests := []struct {
		name   string
		config *config.LambdaConfig
		input  map[string]interface{}
		vars   map[string]interface{}
		output interface{}
		want   context.ExecutionContext
	}{
		{
			name:   "cenário feliz: todos os campos preenchidos",
			config: cfg,
			input:  inputData,
			vars:   variablesData,
			output: outputData,
			want: context.ExecutionContext{
				Config:    cfg,
				Input:     inputData,
				Variables: variablesData,
				Output:    outputData,
			},
		},
		{
			name:   "cenário: config nulo",
			config: nil,
			input:  inputData,
			vars:   variablesData,
			output: outputData,
			want: context.ExecutionContext{
				Config:    nil,
				Input:     inputData,
				Variables: variablesData,
				Output:    outputData,
			},
		},
		{
			name:   "cenário: input nulo",
			config: cfg,
			input:  nil,
			vars:   variablesData,
			output: outputData,
			want: context.ExecutionContext{
				Config:    cfg,
				Input:     nil,
				Variables: variablesData,
				Output:    outputData,
			},
		},
		{
			name:   "cenário: variables nulo",
			config: cfg,
			input:  inputData,
			vars:   nil,
			output: outputData,
			want: context.ExecutionContext{
				Config:    cfg,
				Input:     inputData,
				Variables: nil,
				Output:    outputData,
			},
		},
		{
			name:   "cenário: output nulo",
			config: cfg,
			input:  inputData,
			vars:   variablesData,
			output: nil,
			want: context.ExecutionContext{
				Config:    cfg,
				Input:     inputData,
				Variables: variablesData,
				Output:    nil,
			},
		},
		{
			name:   "cenário: todos os campos nulos",
			config: nil,
			input:  nil,
			vars:   nil,
			output: nil,
			want: context.ExecutionContext{
				Config:    nil,
				Input:     nil,
				Variables: nil,
				Output:    nil,
			},
		},
		{
			name:   "cenário: input e variables vazios",
			config: cfg,
			input:  map[string]interface{}{},
			vars:   map[string]interface{}{},
			output: outputData,
			want: context.ExecutionContext{
				Config:    cfg,
				Input:     map[string]interface{}{},
				Variables: map[string]interface{}{},
				Output:    outputData,
			},
		},
		{
			name:   "cenário: output com tipo diferente",
			config: cfg,
			input:  inputData,
			vars:   variablesData,
			output: "string output",
			want: context.ExecutionContext{
				Config:    cfg,
				Input:     inputData,
				Variables: variablesData,
				Output:    "string output",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := context.ExecutionContext{
				Config:    tt.config,
				Input:     tt.input,
				Variables: tt.vars,
				Output:    tt.output,
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
