package sqs

import (
	"context"

	"github.com/andreparelho/order-api/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	aws_sqs "github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SQSApi interface {
	SendMessage(ctx context.Context, params *aws_sqs.SendMessageInput, optFns ...func(*aws_sqs.Options)) (*aws_sqs.SendMessageOutput, error)
	ReceiveMessage(ctx context.Context, params *aws_sqs.ReceiveMessageInput, optFns ...func(*aws_sqs.Options)) (*aws_sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *aws_sqs.DeleteMessageInput, optFns ...func(*aws_sqs.Options)) (*aws_sqs.DeleteMessageOutput, error)
}

type sqsApi struct {
	sqsApi SQSApi
}

type SQSClient interface {
	SendMessage(ctx context.Context, queueUrl, body string) error
	ReceiveMessage(ctx context.Context, queueUrl string) (*aws_sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, queueUrl string, msg *string) error
}

type client struct {
	client *sqsApi
}

func NewSQSClient(ctx context.Context, config config.Configuration) (SQSClient, error) {

	cfg, err := aws_config.LoadDefaultConfig(context.TODO(),
		aws_config.WithRegion(config.AWS.Region),
		aws_config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				config.AWS.Key, config.AWS.Secret, config.AWS.Session,
			),
		),
		aws_config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				if service == sqs.ServiceID {
					return aws.Endpoint{
						URL:           config.AWS.Endpoint,
						SigningRegion: config.AWS.Region,
					}, nil
				}
				return aws.Endpoint{}, &aws.EndpointNotFoundError{}
			}),
		))
	if err != nil {
		return nil, err
	}

	awsClient := aws_sqs.NewFromConfig(cfg)

	sqs := &sqsApi{
		sqsApi: awsClient,
	}

	return &client{
		client: sqs,
	}, nil
}

func (c *client) SendMessage(ctx context.Context, queueUrl, body string) error {
	_, err := c.client.sqsApi.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &queueUrl,
		MessageBody: &body,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *client) ReceiveMessage(ctx context.Context, queueUrl string) (*aws_sqs.ReceiveMessageOutput, error) {
	msgs, err := c.client.sqsApi.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            &queueUrl,
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (c *client) DeleteMessage(ctx context.Context, queueUrl string, msg *string) error {
	_, err := c.client.sqsApi.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      &queueUrl,
		ReceiptHandle: msg,
	})
	if err != nil {
		return err
	}

	return nil
}
