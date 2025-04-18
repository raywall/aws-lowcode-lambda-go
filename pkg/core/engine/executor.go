package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/core/context"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/core/steps"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/plugin"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
)

func NewExecutionContext(config *config.LambdaConfig, event map[string]interface{}) *context.ExecutionContext {
	if bodyStr, ok := event["body"].(string); ok {
		var parsedBody map[string]interface{}
		if err := json.Unmarshal([]byte(bodyStr), &parsedBody); err == nil {
			event["body"] = parsedBody
		}
	}

	return &context.ExecutionContext{
		Input:     event,
		Config:    config,
		Variables: make(map[string]interface{}),
		Output:    make(map[string]interface{}), // Inicializa o output
	}
}

func Execute(ctx *context.ExecutionContext) error {
	for _, step := range ctx.Config.Steps {
		err := executeStep(step, ctx, nil)
		if err != nil {
			return handleError(step, err, ctx)
		}
	}

	return nil
}

func executeStep(step config.Step, ctx *context.ExecutionContext, registry *plugin.Registry) error {
	// Tenta executar como passo customizado primeiro
	if registry != nil {
		if executor, exists := registry.GetStep(step.Type); exists {
			return executor.Execute(step.Config, ctx)
		}
	}

	switch step.Type {
	case "validate":
		return steps.ExecuteValidate(step.Config, ctx)
	case "callApi":
		return steps.ExecuteCallAPI(step.Config, ctx)
	case "authorization":
		return steps.ExecuteAuthorization(step.Config, ctx)
	case "return":
		return steps.ExecuteReturn(step.Config, ctx)
	default:
		return fmt.Errorf("tipo de passo não suportado: %s", step.Type)
	}
}

func handleError(step config.Step, err error, ctx *context.ExecutionContext) error {
	if step.OnError.Action == "return" {
		// Registrar o erro original no contexto
		ctx.Variables["error"] = map[string]interface{}{
			"message": err.Error(),
			"step":    step.ID,
		}

		// Resolver templates
		body := utils.ResolveTemplates(step.OnError.Response.Body, ctx.Variables)

		// Garantir que o body seja serializável
		if s, ok := body.(string); ok {
			body = strings.ReplaceAll(s, "\\n", "\n")
		}

		return nil
	}
	return err
}
