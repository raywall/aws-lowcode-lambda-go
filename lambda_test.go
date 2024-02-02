package lambda

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	lowcodeFunction := LowcodeFunction{}
	lowcodeFunction.FromConfigFile("sample.yaml", "APIGATEWAY", "DYNAMODB")
}
