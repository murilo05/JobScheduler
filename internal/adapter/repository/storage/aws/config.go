package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/zap"
)

type AWS struct {
	sess   *session.Session
	logger *zap.SugaredLogger
}

func NewAWS(ctx context.Context, logger *zap.SugaredLogger) *AWS {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &AWS{
		sess:   sess,
		logger: logger,
	}
}
