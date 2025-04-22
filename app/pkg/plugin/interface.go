package plugin

import "github.com/raywall/aws-lowcode-lambda-go/pkg/core/context"

type Plugin interface {
	Name() string
	Register(registry *Registry)
	Validate(config map[string]interface{}) error
}

type StepExecutor interface {
	Execute(config map[string]interface{}, ctx *context.ExecutionContext) error
}

type ResourceHandler interface {
	Create(config map[string]interface{}) (interface{}, error)
}
