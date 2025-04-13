package core

import (
	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
)

func BuildResponse(source string, config map[string]interface{}, ctx *ExecutionContext) ResponseInterface {
	// Garante que o body seja um mapa antes de resolver templates
	bodyData, ok := config["body"]
	if !ok {
		bodyData = map[string]interface{}{}
	}

	// Converte para map[string]interface{} se necessário
	bodyMap, ok := bodyData.(map[string]interface{})
	if !ok {
		if converted, ok := utils.ConvertToMapStringInterface(bodyData); ok {
			bodyMap = converted
		} else {
			bodyMap = map[string]interface{}{
				"value": bodyData,
			}
		}
	}

	// Resolve os templates no body
	resolvedBody := utils.ResolveTemplates(bodyMap, ctx.Variables)

	// Converte para JSON
	body, err := utils.ConvertAnyMapToJSON(resolvedBody)
	if err != nil {
		return APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"Failed to serialize response"}`,
		}
	}

	switch source {
	case "apiGatewayProxy", "alb":
		return APIGatewayProxyResponse{
			StatusCode: config["statusCode"].(int),
			Headers:    convertHeaders(config["headers"]),
			Body:       body,
		}
	// case "sqs":
	// 	return SQSResponse{}
	// case "sns":
	// 	return SNSResponse{}
	default:
		return APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"Unknown event source"}`,
		}
	}
}

type APIGatewayProxyResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

func (r APIGatewayProxyResponse) ToResponseFormat() map[string]interface{} {
	return map[string]interface{}{
		"statusCode": r.StatusCode,
		"headers":    r.Headers,
		"body":       r.Body,
	}
}
