package sqs

import (
	"github.com/murilo05/JobScheduler/internal/core/ports"
)

type SQS struct {
	repo ports.Repository
}

func newEmailSender(aws string) *SQS {
	return &SQS{}
}
