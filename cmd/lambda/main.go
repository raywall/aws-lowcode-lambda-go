package main

import (
	"context"
	"os"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/config"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/engine"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/models"

	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/yaml.v2"
)

var lambdaConfig models.LambdaConfig

func init() {
	loader := config.NewLoader(getConfigSource())
	data, err := loader.Load()
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(data, &lambdaConfig)
	if err != nil {
		panic(err)
	}
}

func handler(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	lambdaContext := engine.NewExecutionContext(&lambdaConfig, event)
	err := engine.Execute(lambdaContext)
	if err != nil {
		return engine.BuildErrorResponse(lambdaContext.Config.Input.Source, err), nil
	}
	return lambdaContext.Output, nil
}

func main() {
	lambda.Start(handler)
}

func getConfigSource() string {
	if source := os.Getenv("CONFIG_SOURCE"); source != "" {
		return source
	}
	return "config.yaml"
}
