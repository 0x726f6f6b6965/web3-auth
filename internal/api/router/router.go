package router

import (
	"github.com/0x726f6f6b6965/web3-auth/internal/api"
	"github.com/0x726f6f6b6965/web3-auth/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	RegisterAuthRouter(server.Group("/auth/"))
	RegisterUserRouter(server.Group("/user/"))
}

func RegisterAuthRouter(router *gin.RouterGroup) {
	auth := api.GetAuthAPI()
	router.POST("/nonce", auth.GetNonce)
	router.POST("/verify", auth.VerifySignature)
}

func RegisterUserRouter(router *gin.RouterGroup) {
	user := api.GetUserAPI()
	router.Use(middleware.UserAuthorization())
	router.GET("/info", user.GetUserInfo)
	router.PATCH("/info", user.UpdateUserInfo)
}
