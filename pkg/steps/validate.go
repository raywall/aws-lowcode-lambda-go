package steps

import (
	"github.com/raywall/aws-lowcode-lambda-go/pkg/core"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
)

var (
	validator = utils.NewJSONValidator()
)

func ExecuteValidate(config map[string]interface{}, ctx *core.ExecutionContext) error {
	dataPath := config["data"].(string)
	data := utils.ResolveTemplates(dataPath, ctx.Variables)

	schemaRef := config["schemaRef"].(string)
	schema := ctx.Config.GetSchema(schemaRef) // Agora acessível via ctx.Config

	if err := validator.Validate(data, schema); err != nil {
		ctx.Variables["step."+config["id"].(string)+".errors"] = err
		return err
	}
	return nil
}
