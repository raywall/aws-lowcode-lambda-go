package receiver

import (
	"github.com/aws/aws-lambda-go/events"
)

type SQS interface {
	HandleSQSEvent(event events.SQSEvent) string
}

func (s *Settings) HandleSQSEvent(event events.SQSEvent) string {
	return "SQS event received"
}
