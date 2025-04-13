package config

import (
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Loader struct {
	Path string
}

func (l *S3Loader) Load() ([]byte, error) {
	bucket, key := parseS3Path(l.Path)

	sess := session.Must(session.NewSession())
	svc := s3.New(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	result, err := svc.GetObject(input)
	if err != nil {
		return nil, err
	}

	defer result.Body.Close()
	return ioutil.ReadAll(result.Body)
}

func parseS3Path(path string) (bucket, key string) {
	path = strings.TrimPrefix(path, "s3://")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) < 2 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
