package dynamo

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	onceInit  sync.Once
	daoClient *DaoClient
)

type DaoClient struct {
	dynamoClient *dynamodb.Client
	table        string
}

func NewDynamoClient(ctx context.Context, cfg aws.Config, table string) {
	var (
		client *DaoClient
	)
	// singleton
	onceInit.Do(func() {
		client.table = table
		client.dynamoClient = dynamodb.NewFromConfig(cfg)
		daoClient = client
	})
}

func NewDevLocalClient(host, table string) {
	cfg, _ := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: fmt.Sprintf("http://%s:8000", host)}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
	dynamo := dynamodb.NewFromConfig(cfg)
	client := &DaoClient{
		table:        table,
		dynamoClient: dynamo,
	}
	daoClient = client
}

func GetDynamoClient() *DaoClient {
	return daoClient
}
