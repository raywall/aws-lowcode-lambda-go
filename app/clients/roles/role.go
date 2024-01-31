package roles

import (
	"encoding/json"
	"fmt"
	"lowcode-lambda/clients/dynamodb"
)

type RoleData []byte

func (roleData RoleData) Load(response *dynamodb.RequestInput) error {

	err = json.Unmarshal(roleData, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal lowcode lambda role: %v", err)
	}

	return nil
}
