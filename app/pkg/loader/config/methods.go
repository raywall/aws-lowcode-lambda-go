package config

import (
	"strings"
)

func (c *LambdaConfig) GetSchema(ref string) interface{} {
	if strings.HasPrefix(ref, "#/input/schema/") {
		key := strings.TrimPrefix(ref, "#/input/schema/")
		return c.Input.Schema[key]
	}
	return nil
}
