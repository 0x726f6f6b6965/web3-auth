package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/0x726f6f6b6965/web3-auth/internal/api/services"
	"github.com/0x726f6f6b6965/web3-auth/internal/proto"
	"github.com/0x726f6f6b6965/web3-auth/internal/utils"
	"github.com/gin-gonic/gin"
)

type user struct{}

var (
	userAPI      *user
	onceInitUser sync.Once
)

func NewUserAPI() *user {
	onceInitUser.Do(func() {
		userAPI = new(user)

	})
	return userAPI
}

func GetUserAPI() *user {
	return userAPI
}

func (s *user) GetUserInfo(ctx *gin.Context) {
	var userToken = new(proto.UserToken)
	if info, ok := ctx.Get("access_token"); !ok {
		utils.InvalidParamErr.Message = "Please carry token."
		utils.Response(ctx, http.StatusOK, utils.InvalidParamErr, nil)
		return
	} else {
		userToken = info.(*proto.UserToken)
	}

	if utils.Empty(userToken.PublicAddress) || !utils.IsValidAddress(userToken.PublicAddress) {
		utils.InvalidParamErr.Message = "Please enter correct address."
		utils.Response(ctx, http.StatusOK, utils.InvalidParamErr, nil)
		return
	}

	info, err := services.GetUserInfo(ctx, userToken.PublicAddress, userToken.Nonce)
	if err != nil {
		utils.InternalServerError.Message = fmt.Sprintf("Operation failed, %s.", err.Error())
		utils.Response(ctx, utils.SuccessCode, utils.InternalServerError, nil)
		return
	}
	utils.Response(ctx, utils.SuccessCode, utils.Success, info)
}

func (s *user) UpdateUserInfo(ctx *gin.Context) {
	var param proto.UpdateUserRequest
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

	if len(param.UpdateMask) == 0 {
		utils.InvalidParamErr.Message = "Please enter update mask."
		utils.Response(ctx, http.StatusOK, utils.InvalidParamErr, nil)
		return
	}

	info, err := services.UpdateUserInfo(ctx, param.PublicAddress, &param.User, param.UpdateMask)
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
	resp := proto.UpdateUserResponse{
		PublicAddress: param.PublicAddress,
		Token:         token,
		User:          *info,
	}
	utils.Response(ctx, utils.SuccessCode, utils.Success, resp)
}
