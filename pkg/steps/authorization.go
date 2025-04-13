package steps

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/core"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/utils"
)

func ExecuteAuthorization(config map[string]interface{}, ctx *core.ExecutionContext) error {
	// Similar ao callApi mas com tratamento específico para respostas de autenticação
	resourceName := config["resource"].(string)
	resource := ctx.Config.Resources[resourceName]

	endpoint := utils.ResolveTemplates(resource.Config["BaseURL"].(string), ctx.Variables)

	req := buildAuthRequest(config, endpoint.(string), ctx)

	client := &http.Client{
		Timeout: time.Duration(resource.Config["TimeoutSeconds"].(int)) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	responseData := processAuthResponse(resp, config["outputVar"].(string))
	ctx.Variables[config["outputVar"].(string)] = responseData
	return nil
}

func buildAuthRequest(config map[string]interface{}, endpoint string, ctx *core.ExecutionContext) *http.Request {
	method := config["method"].(string)
	path := utils.ResolveTemplates(config["path"].(string), ctx.Variables).(string)

	fullURL := endpoint + path
	req, _ := http.NewRequest(method, fullURL, nil)

	// Adicionar corpo para fluxos client_credentials
	if body, ok := config["body"].(map[string]interface{}); ok {
		jsonBody := utils.ResolveTemplates(body, ctx.Variables)
		jsonBytes, _ := json.Marshal(jsonBody)
		req.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

func processAuthResponse(resp *http.Response, outputVar string) interface{} {
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Extrair access_token da resposta padrão OAuth
	if token, ok := result["access_token"].(string); ok {
		return token
	}
	return nil
}
