package auth

import (
	"go-auth/internal/error_code"

	"github.com/gin-gonic/gin"
	"github.com/shanelex111/go-common/pkg/response"
)

func Signin(c *gin.Context) {
	var (
		req signinRequest
	)
	if err := c.ShouldBind(&req); err != nil {
		response.Failed(c, error_code.AuthBadRequest)
		return
	}

	// 1. 手机号&验证码登录

	// 2. 手机号&密码登录

	// 3. 邮箱&密码登录

	// 4. 邮箱&验证码登录

}

func Signout(c *gin.Context) {

}

func RefreshToken(c *gin.Context) {

}
func ResetPassword(c *gin.Context) {

}
func SendCode(c *gin.Context) {

}

func VerifyCode(c *gin.Context) {

}
