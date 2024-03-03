package services

import (
	"context"

	"github.com/0x726f6f6b6965/web3-auth/internal/proto"
	"github.com/0x726f6f6b6965/web3-auth/internal/utils"
	"github.com/0x726f6f6b6965/web3-auth/pkg/dynamo"
)

func UpdateUserInfo(ctx context.Context, address string, info *proto.User, updateMask []string) (*proto.User, error) {
	client := dynamo.GetDynamoClient()
	if client == nil {
		return nil, utils.ErrDynamoDBClientNotFound
	}

	newInfo, err := dynamo.UpdateUserInfo(ctx, client, "address", address, *info, updateMask)
	if err != nil {
		return nil, err
	}
	return newInfo, nil
}

func GetUserInfo(ctx context.Context, address string, nonce string) (*proto.User, error) {
	client := dynamo.GetDynamoClient()
	if client == nil {
		return nil, utils.ErrDynamoDBClientNotFound
	}

	info, err := dynamo.GetUserInfo[proto.User](ctx, client, address)
	if err != nil {
		return nil, err
	}
	if info.Nonce != nonce {
		return nil, utils.ErrInvalidNonce
	}
	return info, nil
}
