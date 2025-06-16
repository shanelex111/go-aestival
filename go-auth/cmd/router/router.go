package router

import (
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

	authGroup := r.Group("/auth/v1")
	{
		authGroup.POST("/signin", auth.Signin)
		authGroup.PUT("/refresh-token", auth.RefreshToken)
		authGroup.DELETE("/signout", auth.Signout)
		authGroup.DELETE("/account", auth.DeleteAccount)

		authGroup.PUT("/password", auth.ResetPassword)
		authGroup.PUT("/avatar", auth.UpdateAvatar)
		authGroup.PUT("/nickname", auth.UpdateNickname)

		authGroup.POST("/send-code", auth.SendCode)
		authGroup.POST("/verify-code", auth.VerifyCode)

	}

	r.Run(":" + engine.GetPort())
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
