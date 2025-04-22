package response

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/core/context"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
)

func BuildResponse(source string, config map[string]interface{}, ctx *context.ExecutionContext) ResponseInterface {
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

func (r APIGatewayProxyResponse) ToResponseFormat() map[string]interface{} {
	return map[string]interface{}{
		"statusCode": r.StatusCode,
		"headers":    r.Headers,
		"body":       r.Body,
	}
}

// convertHeaders transforma o mapa de headers do YAML para o formato de resposta
// Suporta:
// - Valores string diretos
// - Números convertidos para string
// - Booleanos convertidos para string
// - Slices unidos com vírgula (para multi-value)
func convertHeaders(headersConfig interface{}) map[string]string {
	headers := make(map[string]string)

	if headersMap, ok := headersConfig.(map[string]interface{}); ok {
		for key, value := range headersMap {
			switch v := value.(type) {
			case string:
				headers[key] = v
			case bool:
				headers[key] = strconv.FormatBool(v)
			case int, int32, int64, float32, float64:
				headers[key] = fmt.Sprintf("%v", v)
			case []interface{}: // Para valores múltiplos (ex: Set-Cookie)
				var strValues []string
				for _, item := range v {
					strValues = append(strValues, fmt.Sprintf("%v", item))
				}
				headers[key] = strings.Join(strValues, ", ")
			default:
				headers[key] = fmt.Sprintf("%v", v)
			}
		}
	}

	// Garante Content-Type padrão se não especificado
	if _, exists := headers["Content-Type"]; !exists {
		headers["Content-Type"] = "application/json"
	}

	return headers
}
