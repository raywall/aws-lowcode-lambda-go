package config

import (
	"os"
	"strings"
)

type ConfigLoader interface {
	Load() ([]byte, error)
}

type LocalLoader struct {
	Path string
}

func (l *LocalLoader) Load() ([]byte, error) {
	return os.ReadFile(l.Path)
}

func NewLoader(source string) ConfigLoader {
	switch {
	case strings.HasPrefix(source, "s3://"):
		return &S3Loader{Path: source}
	case strings.HasPrefix(source, "ssm://"):
		return &SSMLoader{Path: source}
	default:
		return &LocalLoader{Path: source}
	}
}
