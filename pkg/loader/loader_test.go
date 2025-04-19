package loader_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	local    string = "/Users/macmini/Documents/workspace/go/aws-lowcode-lambda-go/examples"
	filename string = "template.yaml"
	content  []byte
)

// MockConfigLoader é uma interface mock para ConfigLoader
type MockConfigLoader struct {
	mock.Mock
}

func (m *MockConfigLoader) Load() (*config.LambdaConfig, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*config.LambdaConfig), args.Error(1)
}

// MockConfigLoadDefaultConfig é uma função mock para config.LoadDefaultConfig
var MockConfigLoadDefaultConfig = func(ctx context.Context, optFns ...func(*awsConfig.LoadOptions) error) (aws.Config, error) {
	return aws.Config{}, nil
}

func TestLambdaConfigNewLoader(t *testing.T) {
	content, _ = os.ReadFile(fmt.Sprintf("%s/%s", local, filename))
}

func TestLambdaConfigLoader(t *testing.T) {
	t.Run("Deve carregar o conteúdo do arquivo de template local template.yaml", func(t *testing.T) {
		loader := loader.LocalLoader{
			Path: fmt.Sprintf("%s/%s", local, filename),
		}
		data, err := loader.Load()

		assert.NoError(t, err)
		assert.NotNil(t, data)
	})
}
