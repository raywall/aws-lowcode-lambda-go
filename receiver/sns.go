package receiver

import (
	"github.com/aws/aws-lambda-go/events"
)

type SNS interface {
	HandleSNSEvent(event events.SNSEvent) string
}

func (s *Settings) HandleSNSEvent(event events.SNSEvent) string {
	return "SNS event received"
}
