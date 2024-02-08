package config

import (
	"fmt"
	"os"

	"github.com/linkedin/goavro"
)

func (res *ResourceItem) EncodeJSON(data map[string]interface{}) (interface{}, error) {
	jsonSchema, err := os.ReadFile(res.ObjectPathSchema)
	if err != nil {
		return nil, fmt.Errorf("error when creating avro codec: %v", err)
	}

	// Create an avro codec using the scheme
	codec, err := goavro.NewCodec(string(jsonSchema))
	if err != nil {
		return nil, fmt.Errorf("error when serialize data for avro: %v", err)
	}

	// Serialize data to Avro using the codec
	binaryData, err := codec.BinaryFromNative(nil, data)
	if err != nil {
		return nil, fmt.Errorf("error when serialize data for avro: %v", err)
	}

	// Develop Avro's data back to a Go Map in GO
	nativeData, _, err := codec.NativeFromBinary(binaryData)
	if err != nil {
		return nil, fmt.Errorf("error when deserving avro data: %v", err)
	}

	return nativeData, nil
}
