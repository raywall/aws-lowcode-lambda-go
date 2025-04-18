package loader

// import (
// 	"context"
// 	"fmt"
// 	"io"
// 	"os"
// 	"strings"
// 	"testing"

// 	"github.com/aws/aws-sdk-go-v2/config"
// 	"github.com/aws/aws-sdk-go-v2/service/s3"
// 	"github.com/aws/aws-sdk-go-v2/service/ssm"
// 	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
// 	"github.com/aws/aws-sdk-go/aws"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// var (
// 	local     string = "template.yaml"
// 	bucket    string = "my_s3_bucket"
// 	parameter string = "my_parameter_store"
// )

// // MockS3Client é um mock para o cliente S3
// type MockS3Client struct {
// 	mock.Mock
// }

// // GetObject mocks a chamada para GetObject do cliente S3
// func (m *MockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
// 	args := m.Called(ctx, params)
// 	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
// }

// // MockSSMClient é um mock para o cliente SSM
// type MockSSMClient struct {
// 	mock.Mock
// }

// // GetParameter mocks a chamada para GetParameter do cliente SSM
// func (m *MockSSMClient) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
// 	args := m.Called(ctx, params)
// 	return args.Get(0).(*ssm.GetParameterOutput), args.Error(1)
// }

// // MockConfigLoader é uma interface mock para ConfigLoader
// type MockConfigLoader struct {
// 	mock.Mock
// }

// func (m *MockConfigLoader) Load() (*LambdaConfig, error) {
// 	args := m.Called()
// 	if args.Get(0) == nil {
// 		return nil, args.Error(1)
// 	}
// 	return args.Get(0).(*LambdaConfig), args.Error(1)
// }

// // MockConfigLoadDefaultConfig é uma função mock para config.LoadDefaultConfig
// var MockConfigLoadDefaultConfig = func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
// 	// Será substituída em cada teste conforme necessário
// 	return aws.Config{}, nil
// }

// var content []byte

// func TestLambdaConfigNewLoader(t *testing.T) {
// 	content, _ = os.ReadFile(local)
// }

// func TestLambdaConfigLoader(t *testing.T) {
// 	t.Run("Deve carregar o conteúdo do arquivo de template local template.yaml", func(t *testing.T) {
// 		loader := LocalLoader{
// 			Path: local,
// 		}
// 		data, err := loader.Load()

// 		assert.NoError(t, err)
// 		assert.NotNil(t, data)
// 	})
// }

// func TestLambdaConfigS3Loader(t *testing.T) {
// 	t.Run("Deve carregar o conteúdo do arquivo de template em s3://meu_bucket/template.yaml", func(t *testing.T) {
// 		mockS3Client := new(MockS3Client)
// 		mockS3Client.On("GetObject", mock.Anything, &s3.GetObjectInput{
// 			Bucket: &bucket,
// 			Key:    &local,
// 		}).Return(&s3.GetObjectOutput{
// 			Body: io.NopCloser(strings.NewReader(string(content))),
// 		}, nil)

// 		loader := &S3Loader{
// 			Path:   fmt.Sprintf("s3://%s/%s", bucket, local),
// 			Client: mockS3Client,
// 		}
// 		data, err := loader.Load()

// 		assert.NoError(t, err)
// 		assert.NotNil(t, data)
// 	})
// }

// func TestLambdaSSMConfigLoader(t *testing.T) {
// 	var (
// 		withDecryption bool   = true
// 		value          string = string(content)
// 	)

// 	t.Run("Deve carregar o conteúdo do arquivo de template em ssm://template_data", func(t *testing.T) {
// 		mockSSMClient := new(MockSSMClient)
// 		mockSSMClient.On("GetParameter", mock.Anything, &ssm.GetParameterInput{
// 			Name:           &parameter,
// 			WithDecryption: &withDecryption,
// 		}).Return(&ssm.GetParameterOutput{
// 			Parameter: &types.Parameter{
// 				Value: &value,
// 			},
// 		}, nil)

// 		loader := &SSMLoader{
// 			Path:   parameter,
// 			Client: mockSSMClient,
// 		}
// 		data, err := loader.Load()

// 		assert.NoError(t, err)
// 		assert.NotNil(t, data)
// 	})
// }
