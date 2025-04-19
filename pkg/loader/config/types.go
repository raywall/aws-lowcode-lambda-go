package config

type (
	LambdaConfig struct {
		TemplateFormatVersion string    `yaml:"TemplateFormatVersion"`
		Transform             string    `yaml:"Transform"`
		Metadata              Metadata  `yaml:"Metadata"`
		Input                 Input     `yaml:"Input"`
		Resources             Resources `yaml:"Resources"`
		Steps                 []Step    `yaml:"Steps"`
	}

	Metadata struct {
		Name        string `yaml:"Name"`
		Version     string `yaml:"Version,omitempty"`
		Description string `yaml:"Description,omitempty"`
	}

	Input struct {
		Source string         `yaml:"Source"`
		Schema map[string]any `yaml:"Schema"`
	}

	Step struct {
		ID      string         `yaml:"id"`
		Type    string         `yaml:"type"`
		Config  map[string]any `yaml:"config"`
		OnError ErrorHandler   `yaml:"onError"`
	}

	ErrorHandler struct {
		Action   string `yaml:"action"`
		Response struct {
			StatusCode int                    `yaml:"statusCode"`
			Body       interface{}            `yaml:"body"`
			Headers    map[string]interface{} `yaml:"headers"`
		} `yaml:"response"`
	}

	Resources map[string]Resource

	Resource struct {
		Type   string         `yaml:"Type"`
		Config map[string]any `yaml:"Config"`
	}
)
