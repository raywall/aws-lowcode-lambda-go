package resources

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/raywall/aws-lowcode-lambda-go/config"
	"github.com/raywall/aws-lowcode-lambda-go/lambda/attributes"
)

type ActionRequested string

// Conf refers to the lambda function's configuration, containing all the necessary information
// about the request, database, and response parameters that the client uses to orchestrate requests.
var conf = &config.Global
var svc *dynamodb.DynamoDB

const (
	Create ActionRequested = "POST"
	Read   ActionRequested = "GET"
	Update ActionRequested = "PUT"
	Delete ActionRequested = "DELETE"
)

// handleAPIGatewayEvent é uma função interna que valida a requisição recebida do gateway e
// direciona a ação solicitada de acordo com o método http recebido.
//
// É importante destacar, que apenas os métodos http indicados como permitido na configuração
// da função serão permitidos.
//
// Caso o método enviado não seja suportado pela função ainda, ela responderá com um código 400
func HandleAPIGatewayEvent(event events.APIGatewayProxyRequest, client *dynamodb.DynamoDB) *attributes.ExecutionResponse {
	svc = client

	data, err := attributes.DeserializeAvro([]byte(event.Body), "/opt/user.avsc")
	if err != nil {
		return &attributes.ExecutionResponse{
			StatusCode: 500,
			Error:      err,
		}
	}

	switch ActionRequested(event.HTTPMethod) {
	case Create:
		return saveToDynamoDB(data)
	case Read:
		return readFromDynamoDB(data)
	case Update:
		return updateOnDynamoDB(data)
	case Delete:
		return deleteOnDynamoDB(data)
	default:
		return &attributes.ExecutionResponse{
			StatusCode: 404,
			Message:    fmt.Sprintf("method unsupported: %s", event.HTTPMethod),
		}
	}
}

// saveToDynamoDB é uma função interna responsável por inserir um novo registro do DynamoDB na tabela
// que foi previamente indicado na configuração da função.
//
// Se o novo ítem for inserido com sucesso na tabela, a função retornará um código de status 201, indicando
// que o registro foi criado com sucesso, no entanto, se algo der errado, ele deverá retornar um erro 500 e
// indicar uma mensagem com o erro ocorrido.
//
// Os dados a serem registrados na tabela devem ser indicados no corpo da requisição e o caminho do
// arquivo avsc (avro) com a estrutura do objeto deve ter sido especificado na configuração da função
//
// Para usar esta função, você também precisa especificar o Nome da Tabela do DynamoDB e as chaves que
// compõem a chave primária da tabela.
func saveToDynamoDB(data map[string]interface{}) *attributes.ExecutionResponse {
	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return &attributes.ExecutionResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("failed marshal data: %v", err),
			Error:      err,
		}
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(conf.Resources.Database.TableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return &attributes.ExecutionResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("failed input new item: %v", err),
			Error:      err,
		}
	}

	return &attributes.ExecutionResponse{
		StatusCode: 201,
	}
}

// readFromDynamoDB é uma função interna que possibilita realizar consultas em uma tabela do DynamoDB
// previamente indicada nas configurações da função.
//
// Se você setar o atributo 'ProjectionCols', a consulta irá retornar apenas as colunas que foram
// préviamente indicadas.
//
// Se você setar os atributos 'Filter' e 'FilterValues', este filtro será aplicado a query, realizando
// uma consulta mais específica mediante as regras indicadas.
//
// Para usar esta função, você também precisa especificar o Nome da Tabela do DynamoDB e as chaves que
// compõem a chave primária da tabela.
func readFromDynamoDB(data map[string]interface{}) *attributes.ExecutionResponse {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(conf.Resources.Database.TableName),
		KeyConditionExpression: aws.String("#UserID = :UserID"),
	}

	keyNames, _ := attributes.MarshalAttributeNames(data, "#")
	input.SetExpressionAttributeNames(keyNames)

	keyValues, _ := attributes.MarshalAttributeValues(data, ":")
	input.SetExpressionAttributeValues(keyValues)

	result, err := svc.Query(input)
	if err != nil {
		return &attributes.ExecutionResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("failed input new item: %v", err),
			Error:      err,
		}
	}

	var jsonMap []map[string]interface{}
	if conf.Resources.Response.DataStruct != "" {
		err := json.Unmarshal([]byte(conf.Resources.Response.DataStruct), &jsonMap)
		if err != nil {
			return &attributes.ExecutionResponse{
				StatusCode: 500,
				Error:      fmt.Errorf("failed unmarshal data struct config: %v", err),
			}
		}

		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &jsonMap)
		if err != nil {
			return &attributes.ExecutionResponse{
				StatusCode: 500,
				Error:      fmt.Errorf("failed unmarshal record: %v", err),
			}
		}

		jsonResponse, err := json.Marshal(data)
		if err != nil {
			return &attributes.ExecutionResponse{
				StatusCode: 500,
				Error:      fmt.Errorf("failed marshal mapped response: %v", err),
			}
		}

		return &attributes.ExecutionResponse{
			StatusCode: 200,
			Message:    string(jsonResponse),
		}
	}

	response, err := json.Marshal(result.Items)
	if err != nil {
		return &attributes.ExecutionResponse{
			StatusCode: 500,
			Error:      fmt.Errorf("failed marshal response: %v", err),
		}
	}

	return &attributes.ExecutionResponse{
		StatusCode: 200,
		Message:    string(response),
	}
}

// updateOnDynamoDB é uma função interna responsável por atualizar, remover ou adicionar os atributos
// de uma tabela do DynamoDB previamente especificada nas configurações da função. Se o ítem for atualizado
// com sucesso, a função retornará um código de status 200 em resposta a sua requisição, entretant,
// se algo der errado, ela retornará o status 500 jutamente com a descrição do erro.
//
// Os atributos do ítem que serão modificados, juntamente com os atributos da chave primária precisam ser
// enviados no corpo da requisição para que a atualização seja efetuada.
func updateOnDynamoDB(data map[string]interface{}) *attributes.ExecutionResponse {
	updateInput := &dynamodb.UpdateItemInput{
		TableName: aws.String(conf.Resources.Database.TableName),
	}

	keyNames, _ := attributes.MarshalAttributeNames(data, "#")
	updateInput.SetExpressionAttributeNames(keyNames)

	keyValues, _ := attributes.MarshalAttributeValues(data, ":")
	updateInput.SetExpressionAttributeValues(keyValues)

	cols := []string{}
	updateMode := "SET"

	keys := make(map[string]interface{})
	for key, _ := range conf.Resources.Database.Keys {
		cols = append(cols, fmt.Sprintf("#%s = :%s", key, key))
		if value, ok := data[key]; ok {
			keys[key] = value
		}
	}

	primaryKeys, err := dynamodbattribute.MarshalMap(keys)
	if err != nil {
		return &attributes.ExecutionResponse{
			StatusCode: 500,
			Error:      err,
		}
	}

	updateInput.SetKey(primaryKeys)
	updateInput.SetUpdateExpression(fmt.Sprintf("%s %s", updateMode, strings.Join(cols, ",")))

	_, err = svc.UpdateItem(updateInput)
	if err != nil {
		return &attributes.ExecutionResponse{
			StatusCode: 500,
			Error:      err,
		}
	}

	return &attributes.ExecutionResponse{
		StatusCode: 200,
	}
}

// delete is an internal function responsible for remove an item of the DynamoDB table using the settings
// specified in your configuration file. If the item is removed successfully, you will receive a 200 (Ok)
// status code in response of your request. However, if something goes wrong, you will receive a 500 status
// code and an error specifying the problem
//
// you need to send the values of the keys in your request to properly remove the item
//
// To use this function, you need to specify the 'TableName' and 'Keys' in your configuration file.
func deleteOnDynamoDB(data map[string]interface{}) *attributes.ExecutionResponse {
	keys, err := attributes.MarshalAttributeValues(data, "")
	if err != nil {
		return &attributes.ExecutionResponse{
			StatusCode: 500,
			Error:      err,
		}
	}

	deleteInput := dynamodb.DeleteItemInput{
		TableName: aws.String(conf.Resources.Database.TableName),
		Key:       keys,
	}

	_, err = svc.DeleteItem(&deleteInput)
	if err != nil {
		return &attributes.ExecutionResponse{
			StatusCode: 500,
			Error:      fmt.Errorf("failed to remove table item: %v", err),
		}
	}

	return &attributes.ExecutionResponse{
		StatusCode: 200,
	}
}
