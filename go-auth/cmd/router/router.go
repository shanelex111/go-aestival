package router

import (
	"go-auth/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/shanelex111/go-common/pkg/engine"
)

func Run() {
	gin.SetMode(getGinMode())
	r := gin.New()
	r.Use(gin.Recovery())

	authGroup := r.Group("/auth/v1")
	{
		authGroup.POST("/signup", auth.Signup)
		authGroup.POST("/signin", auth.Signin)
		authGroup.PUT("/refresh-token", auth.RefreshToken)
		authGroup.DELETE("/signout", auth.Signout)

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
