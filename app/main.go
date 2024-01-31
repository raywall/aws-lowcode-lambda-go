package main

import (
	"embed"
	"fmt"
	"log"
	"lowcode-lambda/clients/dynamodb"
	"lowcode-lambda/clients/roles"
	"lowcode-lambda/handlers"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

var version string = "beta"

func init() {
	//go:embed sample.json
	var fs embed.FS

	roleData, err := fs.ReadFile(os.Getenv("FILENAME"))
	if err != nil {
		return fmt.Errorf("failed to read lowcode role file: %v", err)
	}

	roles.RoleData(roleData).Load(handlers.InputData)

	handlers.Client, err := dynamodb.NewDynamoDBClient()
	if err != nil {
		log.Fatalf("failed to start a dynamodb client: %v", err)
	}
}

func main() {
	// fmt.Printf("Versão: %s\n", version)
	lambda.Start(handlers.HandleLambdaEvent)
}
