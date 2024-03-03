package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/0x726f6f6b6965/web3-auth/internal/api/services"
	"github.com/0x726f6f6b6965/web3-auth/internal/proto"
	"github.com/0x726f6f6b6965/web3-auth/internal/utils"
	"github.com/0x726f6f6b6965/web3-auth/pkg/dynamo"
	"github.com/gin-gonic/gin"
)

type auth struct{}

var (
	authAPI      *auth
	onceInitAuth sync.Once
)

func NewAuthAPI() *auth {
	onceInitAuth.Do(func() {
		authAPI = new(auth)

	})
	return authAPI
}

func GetAuthAPI() *auth {
	return authAPI
}

func (s *auth) GetNonce(ctx *gin.Context) {
	var param proto.GetNonceRequest
	if err := ctx.ShouldBindJSON(&param); err != nil {
		utils.InvalidParamErr.Message = err.Error()
		utils.Response(ctx, http.StatusOK, utils.InvalidParamErr, nil)
		return
	}

	if utils.Empty(param.PublicAddress) || !utils.IsValidAddress(param.PublicAddress) {
		utils.InvalidParamErr.Message = "Please enter correct address."
		utils.Response(ctx, http.StatusOK, utils.InvalidParamErr, nil)
		return
	}

	nonce, err := services.GetNonce(ctx, param.PublicAddress)
	if err != nil {
		if err == dynamo.ErrNotFound {
			utils.ErrorCodeNotFoundErr.Message = "user not found, please register first"
			utils.Response(ctx, utils.SuccessCode, utils.ErrorCodeNotFoundErr, nil)
		} else {
			utils.InternalServerError.Message = fmt.Sprintf("Operation failed, %s.", err.Error())
			utils.Response(ctx, utils.SuccessCode, utils.InternalServerError, nil)
		}
		return
	}
	utils.Response(ctx, utils.SuccessCode, utils.Success, nonce)
}

func (s *auth) VerifySignature(ctx *gin.Context) {
	var param proto.VerifySignatureRequest
	if err := ctx.ShouldBindJSON(&param); err != nil {
		utils.InvalidParamErr.Message = err.Error()
		utils.Response(ctx, http.StatusOK, utils.InvalidParamErr, nil)
		return
	}

	if utils.Empty(param.PublicAddress) || !utils.IsValidAddress(param.PublicAddress) {
		utils.InvalidParamErr.Message = "Please enter correct address."
		utils.Response(ctx, http.StatusOK, utils.InvalidParamErr, nil)
		return
	}

	if utils.Empty(param.Signature) {
		utils.InvalidParamErr.Message = "Please enter signature."
		utils.Response(ctx, http.StatusOK, utils.InvalidParamErr, nil)
		return
	}

	if err := services.VerifySignature(ctx, param.PublicAddress, param.Signature); err != nil {
		if err == utils.ErrInvalidNonce {
			utils.InvalidParamErr.Message = "Please enter signature."
			utils.Response(ctx, http.StatusOK, utils.InvalidParamErr, nil)
		} else if err == dynamo.ErrNotFound {
			utils.ErrorCodeNotFoundErr.Message = "user not found, please register first"
			utils.Response(ctx, utils.SuccessCode, utils.ErrorCodeNotFoundErr, nil)
		} else {
			utils.InternalServerError.Message = fmt.Sprintf("Operation failed, %s.", err.Error())
			utils.Response(ctx, utils.SuccessCode, utils.InternalServerError, nil)
		}
		return
	}

	info, err := services.UpdateNonce(ctx, param.PublicAddress)
	if err != nil {
		utils.InternalServerError.Message = fmt.Sprintf("Operation failed, %s.", err.Error())
		utils.Response(ctx, utils.SuccessCode, utils.InternalServerError, nil)
		return
	}

	// generate jwt token
	token, err := utils.GenerateNewAccessToken(info, time.Minute*5)
	if err != nil {
		utils.InternalServerError.Message = fmt.Sprintf("Operation failed, %s.", err.Error())
		utils.Response(ctx, utils.SuccessCode, utils.InternalServerError, nil)
		return
	}
	resp := proto.VerifySignatureResponse{
		PublicAddress: param.PublicAddress,
		Token:         token,
	}
	utils.Response(ctx, utils.SuccessCode, utils.Success, resp)
}
