package main

import (
	"context"
	"os"

	"github.com/raywall/aws-lowcode-lambda-go/pkg/core/engine"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"

	"github.com/aws/aws-lambda-go/lambda"
)

var lambdaConfig *config.LambdaConfig

func init() {
	loader, err := loader.NewLoader(getConfigSource())
	if err != nil {
		panic(err)
	}

	lambdaConfig, err = loader.Load()
	if err != nil {
		panic(err)
	}
}

func handler(ctx context.Context, event map[string]interface{}) (interface{}, error) {
	lambdaContext := engine.NewExecutionContext(lambdaConfig, event)
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
