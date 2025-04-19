package config_test

import (
	"testing"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func load(content string) (*config.LambdaConfig, error) {
	var lambda config.LambdaConfig

	err := yaml.Unmarshal([]byte(content), &lambda)
	if err != nil {
		return nil, err
	}

	return &lambda, nil
}

func TestLambdaConfig(t *testing.T) {
	t.Run("Cenário: validar versão do formato e transform do template", func(t *testing.T) {
		want := config.LambdaConfig{
			TemplateFormatVersion: "2025-04-19",
			Transform:             "LC::Lambda",
		}

		data, err := load(`
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda`)

		assert.Equal(t, *data, want)
		assert.NoError(t, err)
	})

	t.Run("Cenário: configuração mínima", func(t *testing.T) {
		want := config.LambdaConfig{
			TemplateFormatVersion: "2025-04-19",
			Transform:             "LC::Lambda",
			Metadata: config.Metadata{
				Name: "MinimalConfig",
			},
			Input: config.Input{
				Source: "apiGatewayProxy",
			},
			Steps: []config.Step{
				{
					ID:   "health",
					Type: "return",
					Config: map[string]any{
						"statusCode": 200,
						"body":       "Ok",
					},
				},
			},
		}

		data, err := load(`
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda
Metadata:
  Name: MinimalConfig
Input:
  Source: apiGatewayProxy
Steps:
- id: health
  type: return
  config:
    statusCode: 200
    body: Ok
`)

		assert.Equal(t, *data, want)
		assert.NoError(t, err)
	})

	t.Run("Cenário feliz: configuração completa", func(t *testing.T) {
		data, err := load(`
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda
Metadata:
  Name: MyLambdaConfig
  Version: 1.0.0
  Description: Configuração de exemplo do AWS lowCode Lambda
Input:
  Source: apiGatewayProxy
  Schema:
    Body:
      Type: object
      Properties:
        convenioId:
          Type: string
      Required:
      - convenioId
Resources:
  ConvenioServiceGraphQL:
    Type: GraphQL
    Config:
      BaseURL: https://api.convenios.com/graphql
      Headers:
        Authorization: ${env:CONVENIO_SERVICE_API_KEY}
      TimeoutSeconds: 10
Steps:
- id: getConvenioInfo
  type: callApi
  config:
    resource: convenioServiceGraphQL
    operationType: query
    query: |
      query: dataSources(codigoConvenio:$codigoConvenio) {
        convenio {
          codigoConvenio
          nomeConvenio
        }
      }
    variables:
      codigoConvenio: ${input.body.convenioId}
    outputVar: ${this.response}
- id: validateConvenioResponse
  type: validate
  config:
    data: ${step.getConvenioInfo.outputVar}
    rules:
    - field: statusCode
      assertion: equals
      value: 200
    - field: body.errors
      assertion: notExists
    - field: body.data.product
      assertion: exists
  onError:
    action: return
    response:
      statusCode: 404
      body:
        error: Convênio não encontrado!
        details: ${step.getConvenioInfo.outputVar.body.errors}
- id: finalResponse
  type: return
  config:
    statusCode: 200
    headers:
      Content-Type: application/json
    body:
      convenioInfo:
        id: ${step.getConvenioInfo.outputVar.body.codigoConvenio}
        name: ${step.getConvenioInfo.outputVar.body.nomeConvenio}
`)

		want := config.LambdaConfig{
			TemplateFormatVersion: "2025-04-19",
			Transform:             "LC::Lambda",
			Metadata: config.Metadata{
				Name:        "MyLambdaConfig",
				Version:     "1.0.0",
				Description: "Configuração de exemplo do AWS lowCode Lambda",
			},
			Input: config.Input{
				Source: "apiGatewayProxy",
				Schema: map[string]any{
					"Body": map[any]any{
						"Type": "object",
						"Properties": map[any]any{
							"convenioId": map[any]any{
								"Type": "string",
							},
						},
						"Required": []any{"convenioId"},
					},
				},
			},
			Resources: config.Resources{
				"ConvenioServiceGraphQL": config.Resource{
					Type: "GraphQL",
					Config: map[string]any{
						"BaseURL": "https://api.convenios.com/graphql",
						"Headers": map[any]any{
							"Authorization": "${env:CONVENIO_SERVICE_API_KEY}",
						},
						"TimeoutSeconds": 10,
					},
				},
			},
			Steps: []config.Step{
				{
					ID:   "getConvenioInfo",
					Type: "callApi",
					Config: map[string]any{
						"resource":      "convenioServiceGraphQL",
						"operationType": "query",
						"query": `query: dataSources(codigoConvenio:$codigoConvenio) {
  convenio {
    codigoConvenio
    nomeConvenio
  }
}
`,
						"variables": map[any]any{
							"codigoConvenio": "${input.body.convenioId}",
						},
						"outputVar": "${this.response}",
					},
				},
				{
					ID:   "validateConvenioResponse",
					Type: "validate",
					Config: map[string]any{
						"data": "${step.getConvenioInfo.outputVar}",
						"rules": []any{
							map[any]any{
								"field":     "statusCode",
								"assertion": "equals",
								"value":     200,
							},
							map[any]any{
								"field":     "body.errors",
								"assertion": "notExists",
							},
							map[any]any{
								"field":     "body.data.product",
								"assertion": "exists",
							},
						},
					},
					OnError: config.ErrorHandler{
						Action: "return",
						Response: struct {
							StatusCode int                    `yaml:"statusCode"`
							Body       interface{}            `yaml:"body"`
							Headers    map[string]interface{} `yaml:"headers"`
						}{
							StatusCode: 404,
							Body: map[any]any{
								"error":   "Convênio não encontrado!",
								"details": "${step.getConvenioInfo.outputVar.body.errors}",
							},
							Headers: nil,
						},
					},
				},
				{
					ID:   "finalResponse",
					Type: "return",
					Config: map[string]any{
						"statusCode": 200,
						"headers": map[any]any{
							"Content-Type": "application/json",
						},
						"body": map[any]any{
							"convenioInfo": map[any]any{
								"id":   "${step.getConvenioInfo.outputVar.body.codigoConvenio}",
								"name": "${step.getConvenioInfo.outputVar.body.nomeConvenio}",
							},
						},
					},
				},
			},
		}

		assert.Equal(t, *data, want)
		assert.NoError(t, err)
	})

	t.Run("Cenário: metadados opcionais ausentes", func(t *testing.T) {
		want := config.LambdaConfig{
			TemplateFormatVersion: "2025-04-19",
			Transform:             "LC::Lambda",
			Metadata: config.Metadata{
				Name: "MetadadosOpcionaisAusentes",
			},
			Input: config.Input{
				Source: "apiGatewayProxy",
			},
			Resources: map[string]config.Resource{},
			Steps:     []config.Step{},
		}

		data, err := load(`
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda
Metadata:
  Name: MetadadosOpcionaisAusentes
Input:
  Source: apiGatewayProxy
Resources: {}
Steps: []
`)

		assert.Equal(t, *data, want)
		assert.NoError(t, err)
	})

	t.Run("Cenário: onError com campos ausentes", func(t *testing.T) {
		want := config.LambdaConfig{
			TemplateFormatVersion: "2025-04-19",
			Transform:             "LC::Lambda",
			Metadata: config.Metadata{
				Name: "onError com campos ausentes",
			},
			Input: config.Input{
				Source: "apiGatewayProxy",
			},
			Resources: map[string]config.Resource{},
			Steps: []config.Step{
				{
					ID:   "etapa",
					Type: "return",
					Config: map[string]any{
						"statusCode": 200,
					},
					OnError: config.ErrorHandler{
						Action: "return",
						Response: struct {
							StatusCode int                    `yaml:"statusCode"`
							Body       interface{}            `yaml:"body"`
							Headers    map[string]interface{} `yaml:"headers"`
						}{
							StatusCode: 500,
						},
					},
				},
			},
		}

		data, err := load(`
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda
Metadata:
  Name: onError com campos ausentes
Input:
  Source: apiGatewayProxy
Resources: {}
Steps:
- id: etapa
  type: return
  config:
    statusCode: 200
  onError:
    action: return
    response:
      statusCode: 500
`)

		assert.Equal(t, *data, want)
		assert.NoError(t, err)
	})

	t.Run("Cenário de erro: YAML mal formatado", func(t *testing.T) {
		data, err := load(`
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda
Metadata:
  Name: InvalidConfig
Input:
  Source: apiGatewayProxy
Resources
  RestService:
  Type: Rest`)

		assert.Nil(t, data)
		assert.Error(t, err)
	})

	t.Run("Cenário de erro: tipos incorretos", func(t *testing.T) {
		data, err := load(`
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda
Metadata:
  Name: WrongTypes
Input:
  Source: 123
Resources:
  RestService:
    Type: 123
Steps:
- id: 123
  type: return
  config: "string"
  onError: "error"`)

		assert.Nil(t, data)
		assert.Error(t, err)
	})

	t.Run("Cenário: lista de steps vazia", func(t *testing.T) {
		want := config.LambdaConfig{
			TemplateFormatVersion: "2025-04-19",
			Transform:             "LC::Lambda",
			Metadata: config.Metadata{
				Name: "EmptySteps",
			},
			Input: config.Input{
				Source: "apiGatewayProxy",
			},
		}

		data, err := load(`
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda
Metadata:
  Name: EmptySteps
Input:
  Source: apiGatewayProxy
`)

		assert.Equal(t, *data, want)
		assert.NoError(t, err)
	})

	t.Run("Cenário: resources vazio", func(t *testing.T) {
		want := config.LambdaConfig{
			TemplateFormatVersion: "2025-04-19",
			Transform:             "LC::Lambda",
			Metadata: config.Metadata{
				Name: "EmptyResources",
			},
			Input: config.Input{
				Source: "apiGatewayProxy",
			},
		}

		data, err := load(`
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda
Metadata:
  Name: EmptyResources
Input:
  Source: apiGatewayProxy
`)

		assert.Equal(t, *data, want)
		assert.NoError(t, err)
	})
}
