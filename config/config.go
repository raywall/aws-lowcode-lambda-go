package config

import (
	"gopkg.in/yaml.v2"
)

// Global is the object that contains your configuration
var Global = Config{}

// Load is the function responsible for unmarshaling the configuration YAML file into an object that can be
// used by the DynamoDBClient.
func (config *Config) Load(data []byte) error {
	err := yaml.Unmarshal(data, config)
	if err != nil {
		return err
	}

	return nil
}
