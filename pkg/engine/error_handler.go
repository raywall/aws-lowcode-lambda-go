package engine

import "github.com/raywall/aws-lowcode-lambda-go/pkg/models"

func BuildErrorResponse(source string, err error) models.ResponseInterface {
	switch source {
	case "apiGatewayProxy", "alb":
		return models.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"error":"` + err.Error() + `"}`,
		}
	case "sqs":
		return models.SQSResponse{
			BatchItemFailures: []struct {
				ItemIdentifier string `json:"itemIdentifier"`
			}{ /* IDs das mensagens falhas */ },
		}
	default:
		return nil
	}
}
