package models

import "strings"

type LambdaConfig struct {
	TemplateFormatVersion string    `yaml:"TemplateFormatVersion"`
	Metadata              Metadata  `yaml:"Metadata"`
	Input                 Input     `yaml:"Input"`
	Resources             Resources `yaml:"Resources"`
	Steps                 []Step    `yaml:"Steps"`
}

type Metadata struct {
	Name        string `yaml:"Name"`
	Version     string `yaml:"Version"`
	Description string `yaml:"Description"`
}

type Input struct {
	Source string         `yaml:"Source"`
	Schema map[string]any `yaml:"Schema"`
}

type Step struct {
	ID      string         `yaml:"id"`
	Type    string         `yaml:"type"`
	Config  map[string]any `yaml:"config"`
	OnError ErrorHandler   `yaml:"onError"`
}

type ErrorHandler struct {
	Action   string `yaml:"action"`
	Response struct {
		StatusCode int                    `yaml:"statusCode"`
		Body       interface{}            `yaml:"body"`
		Headers    map[string]interface{} `yaml:"headers"`
	} `yaml:"response"`
}

type Resources map[string]Resource

type Resource struct {
	Type   string         `yaml:"Type"`
	Config map[string]any `yaml:"Config"`
}

// Método auxiliar para acessar schemas
func (c *LambdaConfig) GetSchema(ref string) interface{} {
	// Implementar lógica de resolução de JSON Pointer
	// Exemplo simplificado:
	if strings.HasPrefix(ref, "#/input/schema/") {
		key := strings.TrimPrefix(ref, "#/input/schema/")
		return c.Input.Schema[key]
	}
	return nil
}
