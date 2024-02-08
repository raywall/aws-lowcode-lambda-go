#!/bash/sh
set -e

aws dynamodb create-table --endpoint-url http://localhost:8000 --table-name UserTable --attribute-definitions AttributeName="UserID",AttributeType="S" --key-schema AttributeName="UserID",KeyType=HASH --provisioned-throughput ReadCapacityUnits=2,WriteCapacityUnits=2