package resources

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/raywall/aws-lowcode-lambda-go/config"
	"github.com/raywall/aws-lowcode-lambda-go/lambda/attributes"
)

type ActionRequested string

// Conf refers to the lambda function's configuration, containing all the necessary information
// about the request, database, and response parameters that the client uses to orchestrate requests.
var conf = &config.Global

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
func HandleAPIGatewayEvent(event events.APIGatewayProxyRequest) (any, error) {
	data, err := attributes.DeserializeAvro([]byte(event.Body), "/opt/user.avsc")
	if err != nil {
		return "", err
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
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       fmt.Sprintf("method unsupported: %s", event.HTTPMethod),
		}, nil
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
func saveToDynamoDB(data map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	config := aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}
	config.Endpoint = aws.String(os.Getenv("DYNAMO_ENDPOINT"))

	sess, _ := session.NewSession(&config)
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(data)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("failed marshal data: %v", err),
		}, err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(conf.Resources.Database.TableName),
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("failed input new item: %v", err),
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
	}, nil
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
func readFromDynamoDB(data map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	config := aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}
	config.Endpoint = aws.String(os.Getenv("DYNAMO_ENDPOINT"))

	sess, _ := session.NewSession(&config)
	svc := dynamodb.New(sess)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(conf.Resources.Database.TableName),
		KeyConditionExpression: aws.String("#UserID = :UserID"),
	}

	keyNames, _ := attributes.MarshalAttributeNames(data)
	input.SetExpressionAttributeNames(keyNames)

	keyValues, _ := attributes.MarshalAttributeValues(data)
	input.SetExpressionAttributeValues(keyValues)

	result, err := svc.Query(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("failed input new item: %v", err),
		}, err
	}

	var jsonMap []map[string]interface{}
	if conf.Resources.Response.DataStruct != "" {
		err := json.Unmarshal([]byte(conf.Resources.Response.DataStruct), &jsonMap)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, fmt.Errorf("failed unmarshal data struct config: %v", err)
		}

		err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &jsonMap)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, fmt.Errorf("failed unmarshal record: %v", err)
		}

		jsonResponse, err := json.Marshal(data)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
			}, fmt.Errorf("failed marshal mapped response: %v", err)
		}

		return events.APIGatewayProxyResponse{
			Body:       string(jsonResponse),
			StatusCode: 200,
		}, nil
	}

	response, err := json.Marshal(result.Items)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, fmt.Errorf("failed marshal response: %v", err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(response),
		StatusCode: 200,
	}, nil
}

// updateOnDynamoDB é uma função interna responsável por atualizar, remover ou adicionar os atributos
// de uma tabela do DynamoDB previamente especificada nas configurações da função. Se o ítem for atualizado
// com sucesso, a função retornará um código de status 200 em resposta a sua requisição, entretant,
// se algo der errado, ela retornará o status 500 jutamente com a descrição do erro.
//
// Os atributos do ítem que serão modificados, juntamente com os atributos da chave primária precisam ser
// enviados no corpo da requisição para que a atualização seja efetuada.
func updateOnDynamoDB(data map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

// delete
func deleteOnDynamoDB(data map[string]interface{}) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}
