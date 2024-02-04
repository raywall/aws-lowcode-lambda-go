package lambda

import (
	"os"

	goavro "github.com/linkedin/goavro/v2"
)

func deserializeAvro(avroData []byte, avroSchemaPath string) (map[string]interface{}, error) {
	avroSchema, err := os.ReadFile(avroSchemaPath)
	if err != nil {
		return nil, err
	}

	codec, err := goavro.NewCodec(string(avroSchema))
	if err != nil {
		return nil, err
	}

	native, _, err := codec.NativeFromTextual(avroData)
	if err != nil {
		return nil, err
	}

	return native.(map[string]interface{}), nil
}
