package main

import (
	"log"
	"lowcode-lambda/clients/dynamodb"
	"lowcode-lambda/clients/rules"
	"lowcode-lambda/handlers"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var err error

func init() {
	rules.FilePath(os.Getenv("FILENAME")).Load(handlers.InputData)

	handlers.Client, err = dynamodb.NewDynamoDBClient()
	if err != nil {
		log.Fatalf("erro ao iniciar cliente dynamodb: %v", err)
	}
}

func main() {
	lambda.Start(handlers.HandleLambdaEvent)
}
