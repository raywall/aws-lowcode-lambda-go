package steps

import (
	"github.com/raywall/aws-lowcode-lambda-go/pkg/core"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
)

func ExecuteReturn(config map[string]interface{}, ctx *core.ExecutionContext) error {
	ctx.Variables["input"] = map[string]interface{}{
		"body":    ctx.Input["body"],
		"headers": ctx.Input["headers"],
		"path":    ctx.Input["path"],
	}

	// Converte o config para map[string]interface{} se necessário
	convertedConfig, _ := utils.ConvertToMapStringInterface(config)
	if convertedConfig == nil {
		convertedConfig = config
	}

	// Resolve todos os templates
	resolvedConfig := utils.ResolveTemplates(convertedConfig, ctx.Variables)

	// Converte para o tipo esperado
	finalConfig, ok := resolvedConfig.(map[string]interface{})
	if !ok {
		finalConfig = map[string]interface{}{
			"statusCode": 200,
			"headers":    map[string]string{"Content-Type": "application/json"},
			"body":       resolvedConfig,
		}
	}

	// Constrói a resposta
	ctx.Output = core.BuildResponse(ctx.Config.Input.Source, finalConfig, ctx).ToResponseFormat()
	return nil
}
