package router

import (
	"go-auth/internal/account"
	"go-auth/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/shanelex111/go-common/pkg/engine"
	"github.com/shanelex111/go-common/pkg/request"
)

func Run() {
	gin.SetMode(getGinMode())
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(request.SetUUID())
	r.Use(request.SetLogger())

	// auth
	{

		authNoTokenGroup := r.Group("/auth/v1")
		{
			authNoTokenGroup.POST("/signin", auth.Signin)
			authNoTokenGroup.POST("/refresh-token", auth.RefreshToken)
			authNoTokenGroup.POST("/send-code", auth.SendCode)
			authNoTokenGroup.POST("/verify-code", auth.VerifyCode)
			authNoTokenGroup.PUT("/password", auth.ResetPassword)

		}

		authTokenGroup := r.Group("/auth/v1")
		authTokenGroup.Use(request.AuthTokenInfo())
		{
			authTokenGroup.DELETE("/signout", auth.Signout)
		}
	}

	// account
	{

		accountGroup := r.Group("/account/v1")
		accountGroup.Use(request.AuthTokenInfo())
		{
			accountGroup.GET("/info", account.GetInfo)
			accountGroup.PUT("/avatar", account.UpdateAvatar)
			accountGroup.PUT("/nickname", account.UpdateNickname)
			accountGroup.DELETE("/account", account.DeleteAccount)
		}
	}

	if err := r.Run(":" + engine.GetPort()); err != nil {
		panic(err)
	}
}

func getGinMode() string {
	switch engine.GetMode() {
	case engine.ModePreProd, engine.ModeRelease:
		return gin.ReleaseMode
	case engine.ModeTest:
		return gin.TestMode
	default:
		return gin.DebugMode
	}
}
