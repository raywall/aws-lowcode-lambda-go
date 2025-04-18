package loader

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"
	"gopkg.in/yaml.v2"
)

type (
	SSMLoader struct {
		Path   string
		Client SSMClient
	}

	SSMClient interface {
		GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	}
)

func (l *SSMLoader) Load() (*config.LambdaConfig, error) {
	var (
		lambda         config.LambdaConfig
		withDecryption = true
	)

	output, err := l.Client.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           &l.Path,
		WithDecryption: &withDecryption,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting parameter %s: %v", l.Path, err)
	}

	if output.Parameter == nil || output.Parameter.Value == nil {
		return nil, fmt.Errorf("invalid parameter: %v", err)
	}

	err = yaml.Unmarshal([]byte(*output.Parameter.Value), &lambda)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize ssm template file: %v", err)
	}

	return &lambda, nil
}
