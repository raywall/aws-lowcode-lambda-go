package core

import "github.com/raywall/aws-lowcode-lambda-go/pkg/models"

type (
	ExecutionContext struct {
		Config    *models.LambdaConfig
		Input     map[string]interface{}
		Variables map[string]interface{}
		Output    interface{}
		// Registry  *plugin.Registry // Adicionado aqui
	}

	ResponseInterface interface {
		ToResponseFormat() map[string]interface{}
	}

	ResponseBuilder interface {
		Build(source string, config map[string]interface{}, ctx *ExecutionContext) interface{}
	}

	StepHandler interface {
		Handle(config map[string]interface{}, ctx *ExecutionContext) error
	}
)
