package receiver

import "github.com/raywall/aws-lowcode-lambda-go/config"

// Conf refers to the lambda function's configuration, containing all the necessary information
// about the request, database, and response parameters that the client uses to orchestrate requests.
type Settings config.Settings
