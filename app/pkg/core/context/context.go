package context

import "github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"

type ExecutionContext struct {
	Config    *config.LambdaConfig
	Input     map[string]interface{}
	Variables map[string]interface{}
	Output    interface{}
	// Registry  *plugin.Registry // Adicionado aqui
}
