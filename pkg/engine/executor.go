package engine

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/core"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/models"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/plugin"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/steps"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
)

func NewExecutionContext(config *models.LambdaConfig, event map[string]interface{}) *core.ExecutionContext {
	// Verifica se existe um body na requisição e se é string
	if bodyStr, ok := event["body"].(string); ok {
		var parsedBody map[string]interface{}
		if err := json.Unmarshal([]byte(bodyStr), &parsedBody); err == nil {
			event["body"] = parsedBody
		}
	}

	return &core.ExecutionContext{
		Input:     event,
		Config:    config,
		Variables: make(map[string]interface{}),
		Output:    make(map[string]interface{}), // Inicializa o output
	}
}

func Execute(ctx *core.ExecutionContext) error {
	for _, step := range ctx.Config.Steps {
		err := executeStep(step, ctx, nil)
		if err != nil {
			return handleError(step, err, ctx)
		}
	}

	return nil
}

func executeStep(step models.Step, ctx *core.ExecutionContext, registry *plugin.Registry) error {
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

func handleError(step models.Step, err error, ctx *core.ExecutionContext) error {
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
