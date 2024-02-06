package attributes

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

type ExecutionResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    interface{} `json:"message"`
	Error      error       `json:"error"`
}

func (response *ExecutionResponse) ToGatewayResponse() (events.APIGatewayProxyResponse, error) {
	data, _ := json.Marshal(response.Message)

	return events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Body:       string(data),
	}, response.Error
}
