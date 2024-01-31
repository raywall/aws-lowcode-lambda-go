package main

import (
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/raywall/aws-lowcode-lambda-go/config"
	"github.com/raywall/aws-lowcode-lambda-go/server/clients/dynamodb"
	"github.com/raywall/aws-lowcode-lambda-go/server/handlers"

	"github.com/aws/aws-lambda-go/lambda"
)

//go:embed sample.*
var fs embed.FS
var version string = "beta"

func init() {
	conf := &config.Global

	data, err := fs.ReadFile(os.Getenv("CONFIG_SAMPLE"))
	if err != nil {
		log.Fatalf("failed to read lowcode role file: %v", err)
	}

	err = conf.Load(data)
	if err != nil {
		log.Fatalf("failed to load settings: %v", err)
	}

	handlers.Client, err = dynamodb.NewDynamoDBClient()
	if err != nil {
		log.Fatalf("failed to start a dynamodb client: %v", err)
	}
}

func main() {
	fmt.Printf("Versão: %s\n", version)
	lambda.Start(handlers.HandleLambdaEvent)
}
