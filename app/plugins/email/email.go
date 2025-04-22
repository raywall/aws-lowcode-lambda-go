package email

import (
	"github.com/raywall/aws-lowcode-lambda-go/pkg/core/context"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/plugin"
)

type EmailPlugin struct{}

func (p *EmailPlugin) Name() string {
	return "email"
}

func (p *EmailPlugin) Register(registry *plugin.Registry) {
	registry.RegisterStep("sendEmail", &EmailStep{})
}

type EmailStep struct{}

func (s *EmailStep) Execute(config map[string]interface{}, ctx *context.ExecutionContext) error {
	// Implementação do envio de email
	return nil
}
