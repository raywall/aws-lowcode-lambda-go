package lowcodeattribute

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
)

type ExecutionResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    interface{} `json:"message"`
	Error      error       `json:"error"`
}

func (response *ExecutionResponse) ToGatewayResponse() (events.APIGatewayProxyResponse, error) {
	content := ""
	data, _ := json.Marshal(response.Message)

	if response.Message != nil {
		content = string(data)
	}

	log.Println(content)

	return events.APIGatewayProxyResponse{
		StatusCode: response.StatusCode,
		Body:       content,
	}, response.Error
}
