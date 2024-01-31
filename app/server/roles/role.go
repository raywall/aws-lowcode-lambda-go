package roles

import (
	"encoding/json"
	"fmt"

	"github.com/raywall/aws-lowcode-lambda-go/server/clients/dynamodb"
)

type RoleData []byte

func (roleData RoleData) Load(response *dynamodb.RequestInput) error {

	err = json.Unmarshal(roleData, &response)
	if err != nil {
		return fmt.Errorf("failed to unmarshal lowcode lambda role: %v", err)
	}

	return nil
}
