package dynamo

import (
	"context"

	"github.com/andreparelho/order-api/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	aws_dynamo "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoAPI interface {
	GetItem(ctx context.Context, params *aws_dynamo.GetItemInput, optFns ...func(*aws_dynamo.Options)) (*aws_dynamo.GetItemOutput, error)
	PutItem(ctx context.Context, params *aws_dynamo.PutItemInput, optFns ...func(*aws_dynamo.Options)) (*aws_dynamo.PutItemOutput, error)
}

type dynamoAPI struct {
	api DynamoAPI
}

type DynamoClient interface {
	GetItem(ctx context.Context, tableName, pk, sk string) (*aws_dynamo.GetItemOutput, error)
	PutItem(ctx context.Context, items map[string]types.AttributeValue, tableName string) error
}

type client struct {
	client *dynamoAPI
}

func NewDynamoClient(ctx context.Context, cfg config.Configuration) (DynamoClient, error) {
	dynamoConfig, err := aws_config.LoadDefaultConfig(context.TODO(),
		aws_config.WithRegion("us-east-1"),
		aws_config.WithEndpointResolver(
			aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           cfg.AWS.Endpoint,
					SigningRegion: cfg.AWS.Region,
				}, nil
			}),
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

func (c *client) GetItem(ctx context.Context, tableName, pk, sk string) (*aws_dynamo.GetItemOutput, error) {
	item, err := c.client.api.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"pk": &types.AttributeValueMemberS{Value: pk},
			"sk": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return &aws_dynamo.GetItemOutput{}, nil
	}

	return item, nil
}

func (c *client) PutItem(ctx context.Context, items map[string]types.AttributeValue, tableName string) error {
	_, err := c.client.api.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      items,
	})
	if err != nil {
		return err
	}

	return nil
}
