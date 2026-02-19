package dynamo

import (
	"context"

	"github.com/andreparelho/order-api/pkg/config"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	aws_dynamo "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoAPI interface {
	GetItem(ctx context.Context, params *aws_dynamo.GetItemInput, optFns ...func(*aws_dynamo.Options)) (*aws_dynamo.GetItemOutput, error)
	PutItem(ctx context.Context, params *aws_dynamo.PutItemInput, optFns ...func(*aws_dynamo.Options)) (*aws_dynamo.PutItemOutput, error)
}

type dynamoAPI struct {
	api DynamoAPI
}

type DynamoClient interface {
	GetItem(ctx context.Context)
	PutItem(ctx context.Context)
}

type client struct {
	client *dynamoAPI
}

func NewDynamoClient(cfg config.Configuration) (DynamoClient, error) {
	dynamoConfig, err := aws_config.LoadDefaultConfig(context.TODO(),
		aws_config.WithRegion(cfg.AWS.Region),
		aws_config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AWS.Key,
				cfg.AWS.Secret,
				cfg.AWS.Session,
			),
		),
	)
	if err != nil {
		return nil, err
	}

	dynamoClient := aws_dynamo.NewFromConfig(dynamoConfig)

	dynamo := &dynamoAPI{
		api: dynamoClient,
	}

	return &client{
		client: dynamo,
	}, nil
}

func (c *client) GetItem(ctx context.Context) {}

func (c *client) PutItem(ctx context.Context) {}
