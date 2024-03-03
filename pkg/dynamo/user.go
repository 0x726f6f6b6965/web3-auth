package dynamo

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// GetUserInfo - get user information
// PK: USER#<public address>
// SK: #PROFILE#<public address>
func GetUserInfo[Info any](ctx context.Context, client *DaoClient, publicAddress string) (*Info, error) {
	info := new(Info)

	data, err := client.dynamoClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(client.table),
		Key:       getUserInfoKey(publicAddress),
	})

	if err != nil {
		return info, err
	}

	if data.Item == nil {
		return info, ErrNotFound
	}

	if err := attributevalue.UnmarshalMap(data.Item, info); err != nil {
		return info, err
	}

	return info, nil
}

// PutUserInfo - put user information
// PK: USER#<public address>
// SK: #PROFILE#<public address>
func PutUserInfo[Info any](ctx context.Context, client *DaoClient, publicAddress string, info Info) error {
	data, err := attributevalue.MarshalMap(info)
	if err != nil {
		return err
	}

	data[pk] = &types.AttributeValueMemberS{
		Value: fmt.Sprintf(UserKey, publicAddress),
	}
	data[sk] = &types.AttributeValueMemberS{
		Value: fmt.Sprintf(ProfileKey, publicAddress),
	}

	_, err = client.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(client.table),
		Item:                data,
		ConditionExpression: aws.String(pkNotExists),
	})

	return err
}

// DeleteUserInfo - delete user information
// PK: USER#<public address>
// SK: #PROFILE#<public address>
func DeleteUserInfo(ctx context.Context, client *DaoClient, publicAddress string) error {
	_, err := client.dynamoClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:           aws.String(client.table),
		Key:                 getUserInfoKey(publicAddress),
		ConditionExpression: aws.String(pkExists),
	})

	return err
}

// UpdateUserInfo - update user information
// PK: USER#<public address>
// SK: #PROFILE#<public address>
func UpdateUserInfo[Info any](ctx context.Context, client *DaoClient, pk, publicAddress string, info Info, updateMask []string) (*Info, error) {
	newInfo := new(Info)
	expr, err := getUpdateExpression(info, pk, "", updateMask)
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
		return newInfo, err
	}

	resp, err := client.dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(client.table),
		Key:                       getUserInfoKey(publicAddress),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueAllNew,
		ConditionExpression:       aws.String(pkExists),
	})

	if err != nil {
		log.Printf("Couldn't update user %v. Here's why: %v\n", publicAddress, err)
		return newInfo, err
	}

	err = attributevalue.UnmarshalMap(resp.Attributes, newInfo)
	if err != nil {
		log.Printf("Couldn't unmarshall update response. Here's why: %v\n", err)
		return newInfo, err
	}

	return newInfo, nil
}
