package utils

import "errors"

const (
	SuccessCode                    = 200
	ErrorCode                      = -1
	ErrorCodeOfInternalServerError = 100        // internal server error, please check server log
	ErrorCodeOfInvalidParams       = 100 + iota // param error

	ErrorCodeLogin    = 401
	ErrorCodeNotFound = 404
)

var (
	Success              = ErrorString{SuccessCode, "success"}
	InvalidParamErr      = ErrorString{ErrorCodeOfInvalidParams, "Wrong request parameter"}
	InternalServerError  = ErrorString{ErrorCodeOfInternalServerError, "Service internal exception"}
	ErrorCodeNotFoundErr = ErrorString{ErrorCodeNotFound, "The resource is not found"}

	ErrDynamoDBClientNotFound = errors.New("the client of dynamodb not found")
	ErrInvalidNonce           = errors.New("invalid nonce")
	ErrTokenExpire            = errors.New("the token expired")
)

type ErrorString struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
