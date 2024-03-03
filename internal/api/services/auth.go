package services

import (
	"context"
	"time"

	"github.com/0x726f6f6b6965/web3-auth/internal/proto"
	"github.com/0x726f6f6b6965/web3-auth/internal/utils"
	"github.com/0x726f6f6b6965/web3-auth/pkg/dynamo"
)

func GetNonce(ctx context.Context, address string) (string, error) {
	client := dynamo.GetDynamoClient()
	if client == nil {
		return "", utils.ErrDynamoDBClientNotFound
	}

	_, err := dynamo.GetUserInfo[proto.User](ctx, client, address)
	if err != nil {
		if err != dynamo.ErrNotFound {
			return "", err
		}
		// register
		info := new(proto.User)
		info.PublicAddress = address
		info.CreatedAt = time.Now().Unix()
		info.UpdatedAt = time.Now().Unix()

		nonce, err := utils.GenerateNonce(24)
		if err != nil {
			return "", err
		}
		info.Nonce = nonce
		err = dynamo.PutUserInfo(ctx, client, info.PublicAddress, info)
		if err != nil {
			return "", err
		}
		return info.Nonce, nil
	}

	// update new nonce
	if info, err := UpdateNonce(ctx, address); err != nil {
		return "", err
	} else {
		return info.Nonce, nil
	}
}

func VerifySignature(ctx context.Context, address string, signature string) error {
	client := dynamo.GetDynamoClient()
	if client == nil {
		return utils.ErrDynamoDBClientNotFound
	}

	info, err := dynamo.GetUserInfo[proto.User](ctx, client, address)
	if err != nil {
		return err
	}

	if valid := utils.VerifySignature(address, signature, info.Nonce); valid != nil {
		return utils.ErrInvalidNonce
	}
	return nil
}

func UpdateNonce(ctx context.Context, address string) (*proto.User, error) {
	client := dynamo.GetDynamoClient()
	if client == nil {
		return nil, utils.ErrDynamoDBClientNotFound
	}

	nonce, err := utils.GenerateNonce(24)
	if err != nil {
		return nil, err
	}

	info := new(proto.User)
	info.Nonce = nonce
	info.UpdatedAt = time.Now().Unix()

	newInfo, err := dynamo.UpdateUserInfo(ctx, client, "address", address, *info, []string{"nonce", "updated_at"})
	if err != nil {
		return nil, err
	}
	return newInfo, nil
}
