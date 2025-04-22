package loader

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/raywall/aws-lowcode-lambda-go/pkg/loader/config"
	"gopkg.in/yaml.v2"
)

type (
	S3Client interface {
		GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	}

	S3Loader struct {
		Path   string
		Client S3Client
	}
)

func (l *S3Loader) Load() (*config.LambdaConfig, error) {
	var lambda config.LambdaConfig

	bucket, key := parseS3Path(l.Path)
	output, err := l.Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		return nil, fmt.Errorf("error getting object %s from bucket %s: %v", key, bucket, err)
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	err = yaml.Unmarshal(data, &lambda)
	if err != nil {
		return nil, fmt.Errorf("unable to serialize s3 template file: %v", err)
	}

	return &lambda, nil
}

func parseS3Path(path string) (bucket, key string) {
	path = strings.TrimPrefix(path, "s3://")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
