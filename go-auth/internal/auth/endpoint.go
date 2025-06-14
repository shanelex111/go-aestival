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

	// 1. 手机号&验证码登录 | 邮箱&验证码登录| 手机号&密码登录 | 邮箱&密码登录

	// 1.1 手机号&验证码校验
	if req.SigninType == signinTypePhone && req.CheckType == checkTypeVerificationCode {
		if req.PhoneCountryCode == "" || req.PhoneNumber == "" || req.VerificationCode == "" {
			response.Failed(c, error_code.AuthBadRequest)
			return
		}
		valid, err := verifyPhoneCode(req.PhoneCountryCode, req.PhoneNumber, req.VerificationCode)
		if err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
		if !valid {
			response.Failed(c, error_code.AuthInvalidVerificationCode)
			return
		}
	}

	// 1.2 邮箱&验证码校验
	if req.SigninType == signinTypeEmail && req.CheckType == checkTypeVerificationCode {
		if req.Email == "" || req.VerificationCode == "" {
			response.Failed(c, error_code.AuthBadRequest)
			return
		}
		valid, err := verifyEmailCode(req.Email, req.VerificationCode)
		if err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
		if !valid {
			response.Failed(c, error_code.AuthInvalidVerificationCode)
			return
		}
	}

	// 2. 查询ip信息

	// 3. 记录账户信息

	// 4. 记录设备信息

	// 5. 生成token

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
