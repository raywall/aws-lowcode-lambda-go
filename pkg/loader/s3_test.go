package loader_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var bucket string = "my_s3_bucket"

// MockS3Client é um mock para o cliente S3
type MockS3Client struct {
	mock.Mock
}

// GetObject mocks a chamada para GetObject do cliente S3
func (m *MockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func TestLambdaConfigS3Loader(t *testing.T) {
	t.Run("Deve carregar o conteúdo do arquivo de template em s3://meu_bucket/template.yaml", func(t *testing.T) {
		mockS3Client := new(MockS3Client)
		mockS3Client.On("GetObject", mock.Anything, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &filename,
		}).Return(&s3.GetObjectOutput{
			Body: io.NopCloser(strings.NewReader(string(content))),
		}, nil)

		loader := &loader.S3Loader{
			Path:   fmt.Sprintf("s3://%s/%s", bucket, filename),
			Client: mockS3Client,
		}
		data, err := loader.Load()

		assert.NoError(t, err)
		assert.NotNil(t, data)
	})
}
