package config

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type SSMLoader struct {
	Path string
}

func (l *SSMLoader) Load() ([]byte, error) {
	sess := session.Must(session.NewSession())
	svc := ssm.New(sess)

	paramInput := &ssm.GetParameterInput{
		Name:           &l.Path,
		WithDecryption: aws.Bool(true),
	}

	result, err := svc.GetParameter(paramInput)
	if err != nil {
		return nil, err
	}

	return []byte(*result.Parameter.Value), nil
}
