package engine

import (
	"github.com/raywall/aws-lowcode-lambda-go/pkg/core/response"
)

func BuildErrorResponse(source string, err error) response.ResponseInterface {
	switch source {
	case "apiGatewayProxy", "alb":
		return response.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"` + err.Error() + `"}`,
		}
	case "sqs":
		return response.SQSResponse{
			BatchItemFailures: []struct {
				ItemIdentifier string `json:"itemIdentifier"`
			}{ /* IDs das mensagens falhas */ },
		}
	default:
		return nil
	}
}
