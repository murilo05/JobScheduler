package aws

import (
	"context"
	"errors"
	"flag"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/murilo05/JobScheduler/internal/adapter/repository/storage"
)

var _ storage.AWS = &AWS{}

func (a *AWS) SendEmailValidationToSQS(ctx context.Context, email, token string) error {
	queue := flag.String("q", "", "The name of the queue")
	flag.Parse()

	if *queue == "" {
		err := errors.New("you must supply the name of a queue (-q QUEUE)")
		a.logger.Error("Invalid queue name: %s", err)
		return err
	}

	body, err := createEmailPayload(email, token)
	if err != nil {
		return err
	}

	svc := sqs.New(a.sess)
	_, err = svc.SendMessage(&sqs.SendMessageInput{
		DelaySeconds: aws.Int64(10),
		MessageBody:  aws.String(string(body)),
		QueueUrl:     queue,
	})
	if err != nil {
		a.logger.Error("failed to send message to queue: %s", err)
		return err
	}

	return nil
}
