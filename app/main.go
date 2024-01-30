package main

import (
	"log"
	"lowcode-lambda/clients/dynamodb"
	"lowcode-lambda/clients/rules"
	"lowcode-lambda/handlers"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	err     error
	version string = "beta"
)

func init() {
	rules.FilePath(os.Getenv("FILENAME")).Load(handlers.InputData)

	handlers.Client, err = dynamodb.NewDynamoDBClient()
	if err != nil {
		log.Fatalf("erro ao iniciar cliente dynamodb: %v", err)
	}
}

func main() {
	// fmt.Printf("Versão: %s\n", version)
	lambda.Start(handlers.HandleLambdaEvent)
}
