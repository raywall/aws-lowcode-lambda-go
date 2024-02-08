package receiver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/raywall/aws-lowcode-lambda-go/config"
	"github.com/raywall/aws-lowcode-lambda-go/lowcodeattribute"
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
func HandleAPIGatewayEvent(event events.APIGatewayProxyRequest, client *dynamodb.DynamoDB) *lowcodeattribute.ExecutionResponse {
	svc = client

	var data map[string]interface{}
	err := json.Unmarshal([]byte(event.Body), &data)
	if err != nil {
		return &lowcodeattribute.ExecutionResponse{
			StatusCode: 500,
			Error:      err,
		}
	}

	jsonMap, err := conf.Resources.Receiver.EncodeJSON(data)
	if err != nil {
		return &lowcodeattribute.ExecutionResponse{
			StatusCode: 500,
			Error:      err,
		}
	}

	switch ActionRequested(event.HTTPMethod) {
	case Create:
		return saveToDynamoDB(jsonMap)
	case Read:
		return readFromDynamoDB(jsonMap)
	case Update:
		return updateOnDynamoDB(jsonMap)
	case Delete:
		return deleteOnDynamoDB(jsonMap)
	default:
		return &lowcodeattribute.ExecutionResponse{
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
func saveToDynamoDB(data interface{}) *lowcodeattribute.ExecutionResponse {
	item, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return &lowcodeattribute.ExecutionResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("failed marshal data: %v", err),
			Error:      err,
		}
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String(conf.Resources.Connector.Properties.TableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return &lowcodeattribute.ExecutionResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("failed input new item: %v", err),
			Error:      err,
		}
	}

	return &lowcodeattribute.ExecutionResponse{
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
func readFromDynamoDB(data interface{}) *lowcodeattribute.ExecutionResponse {
	names, _ := conf.Resources.Connector.GetKeyAttributeNames(data)
	values, _ := conf.Resources.Connector.GetKeyAttributeValues(data)
	conditions, _ := conf.Resources.Connector.GetKeyConditions(data)

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(conf.Resources.Connector.Properties.TableName),
		KeyConditionExpression:    aws.String(conditions),
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
	}
	log.Println("queryInput:", queryInput) // remover
	result, err := svc.Query(queryInput)
	if err != nil {
		return &lowcodeattribute.ExecutionResponse{
			StatusCode: 500,
			Message:    fmt.Sprintf("failed to execute a table query: %v", err),
			Error:      err,
		}
	}

	var jsonMap []map[string]interface{}
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &jsonMap)
	if err != nil {
		return &lowcodeattribute.ExecutionResponse{
			StatusCode: 500,
			Error:      fmt.Errorf("failed to deserialize response: %v", err),
		}
	}

	jsonResponse, err := json.Marshal(jsonMap)
	if err != nil {
		return &lowcodeattribute.ExecutionResponse{
			StatusCode: 500,
			Error:      fmt.Errorf("failed to serialize query result: %v", err),
		}
	}

	return &lowcodeattribute.ExecutionResponse{
		StatusCode: 200,
		Message:    string(jsonResponse),
	}
}

// updateOnDynamoDB é uma função interna responsável por atualizar, remover ou adicionar os atributos
// de uma tabela do DynamoDB previamente especificada nas configurações da função. Se o ítem for atualizado
// com sucesso, a função retornará um código de status 200 em resposta a sua requisição, entretant,
// se algo der errado, ela retornará o status 500 jutamente com a descrição do erro.
//
// Os atributos do ítem que serão modificados, juntamente com os atributos da chave primária precisam ser
// enviados no corpo da requisição para que a atualização seja efetuada.
func updateOnDynamoDB(data interface{}) *lowcodeattribute.ExecutionResponse {
	keys, _ := conf.Resources.Connector.GetPrimaryKeyAttributeValue(data)
	names, _ := conf.Resources.Connector.GetAttributeNames(data)
	values, _ := conf.Resources.Connector.GetAttributeValues(data)
	updateExpr, _ := conf.Resources.Connector.GetUpdateExpression(data)

	updateInput := &dynamodb.UpdateItemInput{
		TableName:                 aws.String(conf.Resources.Connector.Properties.TableName),
		UpdateExpression:          aws.String(updateExpr),
		ExpressionAttributeNames:  names,
		ExpressionAttributeValues: values,
		Key:                       keys,
	}

	_, err := svc.UpdateItem(updateInput)
	if err != nil {
		return &lowcodeattribute.ExecutionResponse{
			StatusCode: 500,
			Error:      err,
		}
	}

	return &lowcodeattribute.ExecutionResponse{
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
func deleteOnDynamoDB(data interface{}) *lowcodeattribute.ExecutionResponse {
	keys, err := conf.Resources.Connector.GetPrimaryKeyAttributeValue(data.(map[string]interface{}))
	if err != nil {
		return &lowcodeattribute.ExecutionResponse{
			StatusCode: 500,
			Error:      fmt.Errorf("failed to get primary key: %v", err),
		}
	}

	deleteInput := dynamodb.DeleteItemInput{
		TableName: aws.String(conf.Resources.Connector.Properties.TableName),
		Key:       keys,
	}

	_, err = svc.DeleteItem(&deleteInput)
	if err != nil {
		return &lowcodeattribute.ExecutionResponse{
			StatusCode: 500,
			Error:      fmt.Errorf("failed to remove table item: %v", err),
		}
	}

	return &lowcodeattribute.ExecutionResponse{
		StatusCode: 200,
	}
}
