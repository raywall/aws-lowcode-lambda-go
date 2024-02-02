# aws-lowcode-lambda-go


# Creating a config file

```yaml
# config.yaml
TemplateFormatVersion: 2024-01-31
Description: lowcode-lambda configuration

Resources:
  Request:
    AllowedMethods:
    - GET
    - POST
    - PUT

  Database:
    TableName: users
    Keys:
      email:
        Operator: "="
      order:
        Operator: "="
    Filter: "#age > :age"
    FilterValues:
      age: "5"
    ProjectionCols:
    - email
    - username
    - age
  
  Response:
    DataStruct: '[{"username": "", "age": ""}]'
```

# Building a lowcode lambda function

``` Go
// main.go
package main

import (
  "os"
  "fmt"

  "github.com/raywall/aws-lowcode-lambda-go/config"
	"github.com/raywall/aws-lowcode-lambda-go/server/clients/dynamodb"
	"github.com/raywall/aws-lowcode-lambda-go/server/handlers"
)

func init() {
  conf := &config.Global

	data, err := os.ReadFile(os.Getenv("CONFIG_SAMPLE"))
	if err != nil {
		log.Fatalf("failed reading lowcode role file: %v", err)
	}

	err = conf.Load(data)
	if err != nil {
		log.Fatalf("failed loading settings: %v", err)
	}

  // create a handler for integration between an api gateway and a dynamodb table
	handlers.Client, err = dynamodb.NewDynamoDBClient(conf)
	if err != nil {
		log.Fatalf("failed starting a dynamodb client: %v", err)
	}
}

func main() {
    // make the handler available for remote procedure call by aws lambda
    lambda.Start(handlers.HandleLambdaEvent)
}
```

# Testing your function locally with SAM
```shell 

# building your project
sam build

# running your function
sam local start-api
```