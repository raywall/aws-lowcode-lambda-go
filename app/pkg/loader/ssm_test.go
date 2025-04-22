package loader_test

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var parameter string = "my_parameter_store"

// MockSSMClient é um mock para o cliente SSM
type MockSSMClient struct {
	mock.Mock
}

// GetParameter mocks a chamada para GetParameter do cliente SSM
func (m *MockSSMClient) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*ssm.GetParameterOutput), args.Error(1)
}

func TestLambdaSSMConfigLoader(t *testing.T) {
	var (
		withDecryption bool   = true
		value          string = string(content)
	)

	t.Run("Deve carregar o conteúdo do arquivo de template em ssm://template_data", func(t *testing.T) {
		mockSSMClient := new(MockSSMClient)
		mockSSMClient.On("GetParameter", mock.Anything, &ssm.GetParameterInput{
			Name:           &parameter,
			WithDecryption: &withDecryption,
		}).Return(&ssm.GetParameterOutput{
			Parameter: &types.Parameter{
				Value: &value,
			},
		}, nil)

		loader := &loader.SSMLoader{
			Path:   parameter,
			Client: mockSSMClient,
		}
		data, err := loader.Load()

		assert.NoError(t, err)
		assert.NotNil(t, data)
	})
}
