package steps

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/core"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
)

func ExecuteCallAPI(config map[string]interface{}, ctx *core.ExecutionContext) error {
	resourceName := config["resource"].(string)
	resource := ctx.Config.Resources[resourceName]

	endpoint := utils.ResolveTemplates(resource.Config["BaseURL"].(string), ctx.Variables)

	req := buildRequest(config, endpoint.(string), ctx)

	client := &http.Client{
		Timeout: time.Duration(resource.Config["TimeoutSeconds"].(int)) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	outputVar := config["outputVar"].(string)
	ctx.Variables[outputVar] = processResponse(resp)
	return nil
}

func buildRequest(config map[string]interface{}, endpoint string, ctx *core.ExecutionContext) *http.Request {
	method := config["method"].(string)
	path := utils.ResolveTemplates(config["path"].(string), ctx.Variables).(string)

	fullURL := endpoint + path
	req, _ := http.NewRequest(method, fullURL, nil)

	if headers, ok := config["headers"].(map[string]interface{}); ok {
		for k, v := range headers {
			req.Header.Add(k, utils.ResolveTemplates(v.(string), ctx.Variables).(string))
		}
	}

	return req
}

func processResponse(resp *http.Response) map[string]interface{} {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return map[string]interface{}{
		"statusCode": resp.StatusCode,
		"headers":    resp.Header,
		"body":       parseBody(resp.Header.Get("Content-Type"), body),
	}
}

func parseBody(contentType string, body []byte) interface{} {
	if strings.Contains(contentType, "application/json") {
		var result interface{}

		_ = json.Unmarshal(body, &result)
		return result
	}
	return string(body)
}
