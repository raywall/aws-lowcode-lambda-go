package loader

import (
	"context"
	"fmt"
	"os"
	"strings"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"
	"gopkg.in/yaml.v2"
)

type (
	ConfigLoader interface {
		Load() (*config.LambdaConfig, error)
	}

	LocalLoader struct {
		Path string
	}
)

func NewLoader(source string) (ConfigLoader, error) {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	switch {
	case strings.HasPrefix(source, "s3://"):
		return &S3Loader{
			Path:   source,
			Client: s3.NewFromConfig(cfg),
		}, nil
	case strings.HasPrefix(source, "ssm://"):
		return &SSMLoader{
			Path:   source,
			Client: ssm.NewFromConfig(cfg),
		}, nil
	default:
		return &LocalLoader{
			Path: source,
		}, nil
	}
}

func (l *LocalLoader) Load() (*config.LambdaConfig, error) {
	var lambda config.LambdaConfig

	data, err := os.ReadFile(l.Path)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	err = yaml.Unmarshal(data, &lambda)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize local template file: %v", err)
	}

	return &lambda, nil
}
