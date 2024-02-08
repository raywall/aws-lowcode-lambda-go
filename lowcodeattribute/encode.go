package lowcodeattribute

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
)

func SerializeLocalRequest(src map[string]interface{}, dest interface{}) error {
	if src == nil || dest == nil {
		return errors.New("src and dest cannot be nil")
	}

	jsonData, err := json.Marshal(src)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, dest)
	if err != nil {
		return err
	}

	return nil
}

func MarshalAttributeNames(raw map[string]interface{}, separator string) (map[string]*string, error) {
	data := make(map[string]*string)
	for key := range raw {
		data[fmt.Sprintf("%s%s", separator, key)] = aws.String(key)
	}

	return data, nil
}
