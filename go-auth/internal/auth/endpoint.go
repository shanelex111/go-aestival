package auth

import "github.com/gin-gonic/gin"

func Signin(c *gin.Context) {
	var (
		req signinRequest
	)
	if err := c.ShouldBind(&req); err != nil {

	}

	// 1. 手机号&验证码登录

	// 2. 手机号&密码登录

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
