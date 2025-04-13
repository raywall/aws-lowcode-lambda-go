# AWS Low-Code Lambda Framework

![AWS Lambda](https://img.shields.io/badge/AWS_Lambda-FF9900?style=for-the-badge&logo=amazonaws&logoColor=white)
![YAML](https://img.shields.io/badge/YAML-FF6C37?style=for-the-badge&logo=yaml&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)


## Índice

1. [Visão Geral](#📄-visão-geral)
2. [Estrutura Básica do YAML](#📋-estrutura-básica-do-yaml)
3. [Path Parameters](#🛣️-path-parameters)
	- [Acessando em Responses](#cessando-em-responses)
	- [Validação Recomendada](validação-recomendada)
4. [Acessando Objetos Aninhados](#🔍-acessando-objetos-aninhados)
	- [Exemplos Práticos](#exemplos-praticos)
5. [Validação de Dados](#🛡️-validação-de-dados)
	- [Schema Completo](#schema-completo)
6. [Origens de Entrada](#🔗-origens-de-entrada)
	- [API Gateway](#api-gateway)
	- [Application Load Balancer](#application-load-balancer)
	- [SQS](#sqs)
	- [SNS](#sns)
7. [Recursos Externos](#➕-recursos-externos)
	- [REST API](#rest-api)
	- [GraphQL](#graphql)
8. [Recursos Avançados](#🛠️-recursos-avançados)
	- [Variáveis de Ambiente](#variáveis-de-ambiente)
	- [AWS Systems Manager](#aws-systems-manager)
	- [AWS Secrets Manager](#aws-secrets-manager)
10. [Exemplos Completos](#🧪-exemplos-completos)
	- [Exemplo com API Gateway e Path Parameters](#exemplo-com-api-gateway-e-path-parameters)
	- [Exemplo com SQS e Secrets Manager](#exemplo-com-sqs-e-secrets-manager)
	- [Exemplo Completo com Path Parameters](#exemplo-completo-com-path-parameters)
11. [Implementação em Go](#🖥️ -implementação-em-go)
12. [Configuração Recomendada](#💡-configuração-recomendada)
	- [Variáveis de Ambiente](#variáveis-de-ambiente)
	- [Permissões AWS](#permissões-aws)
	- [Dependências](#dependências)
13. [Boas Práticas](#📝-boas-práticas)
14. [Debugging](#🔎-debugging)
15. [Considerações Finais](#🏁-considerações-finais)



## 📄 Visão Geral

O AWS Low-Code Lambda Framework é uma solução inovadora para criação de funções Lambda usando uma abordagem declarativa com arquivos YAML. Ele permite:

- Configuração simplificada de funções Lambda
- Integração nativa com serviços AWS
- Validação automática de entradas
- Gerenciamento seguro de segredos
- Suporte a múltiplos tipos de eventos



## 📋 Estrutura Básica do YAML

Um arquivo YAML típico contém:

```yaml
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda

# Metadados da função
Metadata:
	Name: minha-funcao
	Version: 1.0.0
	Description: Descrição da função

# Definição da entrada/recepção de requisições pela função lambda
Input:
	# Origem do evento: 'apiGatewayProxy', 'alb', 'sqs' ou 'sns'
	Source: apiGatewayProxy
	# (Opcional) Esquema para validação inicial dos dados de entrada 
	# ex: query parameters, path parameters, body
	# Poderia usar JSON Schema ou uma definição mais simples	
	Schema:
		Body:
			Type: object
			Properties:
				userId: 
					Type: string
			Required:
				- userId

# Definição dos recursos externos a serem utilizados
Resources:
	MyService:
		Type: Rest
		Config:
			BaseURL: https://api.example.com

# Lógica da função: uma sequência de passos
Steps:
- id: firstStep
  type: validate
  config:
	data: ${input.body}

```



## 🛣️ Path Parameters

### Acessando em Responses

```yaml
Steps:
  - id: mockResponse
    type: return
    config:
      body:
        message: "Sucesso!"
        pathParam: ${input.pathParameters.id}  # Acessa parâmetro de rota
        fullPath: "/api/${input.pathParameters.version}/items/${input.pathParameters.id}"
```

### Validação Recomendada

```yaml
Input:
  Schema:
    PathParameters:
      Type: object
      Properties:
        id: 
          Type: string
          Pattern: '^[0-9a-f]{24}$'  # Valida formato MongoDB ID
```



## 🔍 Acessando Objetos Aninhados

### Exemplos Práticos

```yaml
body:
  userProfile: ${input.body.user.profile}
  firstItem: ${input.body.items[0].sku}
  shippingZip: ${input.body.shipping?.address?.zipCode || '00000'}
```



## 🛡️ Validação de Dados

### Schema Completo

```yaml
Input:
  Schema:
    Body:
      Type: object
      Properties:
        user:
          Type: object
          Properties:
            name: {Type: string, MinLength: 2}
            age: {Type: integer, Minimum: 18}
    PathParameters:
      Type: object
      Properties:
        orgId: {Type: string}
```

  

## 🔗 Origens de Entrada

### API Gateway

```yaml
Input:
  Source: apiGatewayProxy
  Schema:
    Body:
      Type: object
      Properties:
        id: {Type: string}
    PathParameters:
      Type: object
      Properties:
        userId: {Type: string}
    QueryStringParameters:
      Type: object
      Properties:
        page: {Type: integer}
```

### Application Load Balancer

```yaml
Input:
  Source: alb
  Schema:
    QueryParameters:
      Type: object
      Properties:
        limit: {Type: integer, Maximum: 100}
```

### SQS

```yaml
Input:
  Source: sqs
  Schema:
    Records:
      Type: array
      Items:
        Type: object
        Properties:
          body: {Type: string}

```

### SNS

```yaml
Input:
  Source: sns
  Schema:
    Message:
      Type: string
    MessageAttributes:
      Type: object
```



## ➕ Recursos Externos

### REST API

```yaml
Resources:
  UserService:
    Type: Rest
    Config:
      BaseURL: https://api.users.com/v1
      TimeoutSeconds: 5
      Headers:
        Authorization: "Bearer ${env:API_TOKEN}"
      Auth:
        Authorizer: "#/Resources/AuthService/Token"

```

### GraphQL

```yaml
Resources:
  ProductService:
    Type: Graphql
    Config:
      Endpoint: https://api.products.com/graphql
      Headers:
        Authorization: "Bearer ${secrets:GRAPHQL_TOKEN}"
      TimeoutSeconds: 10
```



## 🛠️ Recursos Avançados

### Variáveis de Ambiente

```yaml
Resources:
  AuthService:
    Type: Rest
    Config:
      Headers:
        X-Api-Key: "${env:API_KEY}"
```


### AWS Systems Manager

```yaml
Steps:
  - id: getConfig
    type: config
    config:
      dbUrl: "${ssm:/prod/database/url}"
      apiKey: "${ssm:/prod/auth/key}"
```


### AWS Secrets Manager

```yaml
Steps:
  - id: dbConnect
    type: database
    config:
      password: "${secrets:PROD_DB_PASSWORD}"
      apiSecret: "${secrets:PROD_API_SECRET}"
```




## 🧪 Exemplos Completos

### Exemplo com API Gateway e Path Parameters

```yaml
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda

Metadata:
  Name: user-profile
  Description: Obtém perfil do usuário

Input:
  Source: apiGatewayProxy
  Schema:
    PathParameters:
      Type: object
      Properties:
        userId: {Type: string}

Resources:
  UserService:
    Type: Rest
    Config:
      BaseURL: https://api.users.com
      Headers:
        Authorization: "${env:USER_API_KEY}"

Steps:
  - id: getProfile
    type: callApi
    config:
      resource: UserService
      method: GET
      path: "/users/${input.pathParameters.userId}"
      outputVar: profile

  - id: returnResponse
    type: return
    config:
      statusCode: 200
      body:
        userId: ${input.pathParameters.userId}
        profile: ${step.getProfile.outputVar.body}
```


### Exemplo com SQS e Secrets Manager

```yaml
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda

Input:
  Source: sqs

Resources:
  Database:
    Type: SQL
    Config:
      ConnectionString: "${secrets:DB_CONNECTION_STRING}"

Steps:
  - id: processMessages
    type: transform
    config:
      outputVar: processed
      script: |
        return event.Records.map(record => ({
          ...JSON.parse(record.body),
          processedAt: new Date().toISOString()
        }))

  - id: saveToDB
    type: database
    config:
      operation: insert
      table: processed_messages
      items: ${step.processMessages.outputVar}
```


### Exemplo Completo com Path Parameters

```yaml
TemplateFormatVersion: 2025-04-19
Transform: LC::Lambda

Metadata:
  Name: user-profile
  Description: Obtém perfil do usuário

Input:
  Source: apiGatewayProxy
  Schema:
    PathParameters:
      Type: object
      Properties:
        userId: {Type: string}

Steps:
  - id: validateInput
    type: validate
    config:
      data: ${input.pathParameters}
      schemaRef: "#/input/schema/PathParameters"

  - id: getProfile
    type: return
    config:
      statusCode: 200
      body:
        userId: ${input.pathParameters.userId}
        message: "Perfil obtido com sucesso"
```

  

## 🖥️  Implementação em Go  

```go
package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/engine"
)

var (
	lambdaEngine *engine.LambdaEngine
)

func init() {
	config := engine.LoadConfig("function.yaml")
	lambdaEngine = engine.NewEngine(config)
}

func handler(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	return lambdaEngine.Execute(event)
}

func main() {
	lambda.Start(handler)
}
```

  

## 💡 Configuração Recomendada

1. **Variáveis de Ambiente**:

```bash
export API_KEY=my-secret-key
```


2. **Permissões AWS**:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ssm:GetParameter",
                "secretsmanager:GetSecretValue"
            ],
            "Resource": "*"
        }
    ]
}
```

  
3. **Dependências**:

```go
require (
    github.com/aws/aws-lambda-go v1.28.0
    github.com/raywall/aws-lowcode-lambda-go v0.1.0
)
```



## 📝 Boas Práticas

1. **Sempre valide** path parameters e inputs
2. **Use schemas** para documentação e validação automática
3. **Prefira `${input.pathParameters.x}`** sobre acesso direto ao evento
4. **Adicione tratamento de erros** para parâmetros ausentes



## 🔎 Debugging

```yaml
Steps:
- id: debugStep
  type: debug
  config:
	message: |
	  Path Params: ${input.pathParameters}
	  Full Event: ${input}
```



## 🏁 Considerações Finais

- **Variáveis de Ambiente**: Configure-as antes da execução
- **Permissões AWS**: Garanta permissões para SSM/Secrets Manager quando necessário
- **Logging**: Use `context.Logger` para logs consistentes
- **Monitoramento**: Integre com X-Ray para tracing distribuído


Este framework oferece uma maneira poderosa e flexível de criar funções Lambda com mínimo código, mantendo toda a potência do Go e da AWS.