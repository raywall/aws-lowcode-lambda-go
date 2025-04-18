package response

import "github.com/raywall/aws-lowcode-lambda-go/pkg/core/context"

// ResponseInterface define o contrato para todos os tipos de resposta
// type ResponseInterface interface {
// 	ToResponseFormat() map[string]interface{}
// }

type ResponseInterface interface {
	ToLambdaResponse() interface{}
}

type ResponseBuilder interface {
	Build(source string, config map[string]interface{}, ctx *context.ExecutionContext) interface{}
}

// Respostas específicas para cada fonte de evento
type (
	APIGatewayProxyResponse struct {
		StatusCode        int                 `json:"statusCode"`
		Headers           map[string]string   `json:"headers"`
		MultiValueHeaders map[string][]string `json:"multiValueHeaders,omitempty"`
		Body              string              `json:"body"`
		IsBase64Encoded   bool                `json:"isBase64Encoded,omitempty"`
	}

	ALBResponse struct {
		StatusCode        int                 `json:"statusCode"`
		Headers           map[string]string   `json:"headers"`
		MultiValueHeaders map[string][]string `json:"multiValueHeaders,omitempty"`
		Body              string              `json:"body"`
		IsBase64Encoded   bool                `json:"isBase64Encoded"`
	}

	SQSResponse struct {
		BatchItemFailures []struct {
			ItemIdentifier string `json:"itemIdentifier"`
		} `json:"batchItemFailures"`
	}

	SNSResponse struct {
		// SNS não espera resposta específica
		Message string `json:"message,omitempty"`
	}
)

// Implementações de ToLambdaResponse
func (r APIGatewayProxyResponse) ToLambdaResponse() interface{} { return r }
func (r ALBResponse) ToLambdaResponse() interface{}             { return r }
func (r SQSResponse) ToLambdaResponse() interface{}             { return r }
func (r SNSResponse) ToLambdaResponse() interface{}             { return nil } // SNS não precisa de resposta
