package config

type Config struct {
	TemplateFormatVersion string       `yaml:"TemplateFormatVersion"`
	Description           string       `yaml:"Description"`
	Receiver              ResourceItem `yaml:"Receiver"`
	Connector             ResourceItem `yaml:"Connector"`
}

type ResourceItem struct {
	ObjectPathSchema string `yaml:"objectPathSchema"`
	ConnectorType    string `yaml:"connectorType"`
	ReceiverType     string `yaml:"receiverType"`

	Properties struct {
		// ApiGateway Receiver
		AllowedMethods []string `yaml:"allowedMethods"`
		AllowedPath    []string `yaml:"allowedPath"`

		// DynamoDB Connector
		TableName     string                 `yaml:"tableName"`
		Keys          map[string]string      `yaml:"keys"`
		Filter        string                 `yaml:"filter"`
		FilterValues  map[string]interface{} `yaml:"filterValues"`
		OutputColumns []string               `yaml:"outputColumns"`
	} `yaml:"properties"`
}
