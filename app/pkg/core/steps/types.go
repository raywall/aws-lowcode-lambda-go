package steps

import "github.com/raywall/aws-lowcode-lambda-go/pkg/core/context"

type StepHandler interface {
	Handle(config map[string]interface{}, ctx *context.ExecutionContext) error
}
