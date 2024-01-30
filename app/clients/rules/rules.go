package rules

import (
	"embed"
	"encoding/json"
	"fmt"
	"lowcode-lambda/clients/dynamodb"
)

var (
	//go:embed assets/*
	fs embed.FS
)

type FilePath string

func (r FilePath) Load(response *dynamodb.RequestInput) error {
	rules, err := fs.ReadFile("assets/sample.json")
	if err != nil {
		return fmt.Errorf("failed to read lowcode policy file: %v", err)
	}

	err = json.Unmarshal(rules, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal lowcode rules: %v", err)
	}

	return nil
}
