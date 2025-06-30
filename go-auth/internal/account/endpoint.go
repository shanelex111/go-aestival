package account

import (
	"go-auth/internal/error_code"
	"go-auth/internal/metadata/account"
	"go-auth/internal/metadata/device"
	"go-auth/internal/metadata/verification_code"
	"go-auth/internal/token"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shanelex111/go-common/pkg/request"
	"github.com/shanelex111/go-common/pkg/response"
)

func GetInfo(c *gin.Context) {

}
func UpdateAvatar(c *gin.Context) {

}
func UpdateNickname(c *gin.Context) {

}
func DeleteAccount(c *gin.Context) {
	tokenInfo, exists := c.Get(request.TokenInfoKey)
	if !exists {
		response.Failed(c, error_code.AuthUnauthorized)
		return
	}
	requestTokenInfo := tokenInfo.(*request.TokenInfo)

	var (
		accountID        = requestTokenInfo.Account.ID
		email            = requestTokenInfo.Account.Email
		phoneCountryCode = requestTokenInfo.Account.PhoneCountryCode
		phoneNumber      = requestTokenInfo.Account.PhoneNumber
	)

	// 1. 删除账户
	if err := account.DelAllByAccountID(accountID); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 2. 删除所有token
	if err := token.DelAllByAccountID(accountID); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 3. 删除所有设备
	if err := device.DelAllByAccountID(accountID); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 4. 删除所有verify codes
	if email != "" {
		if err := verification_code.DelAllByEmail(email); err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
	}
	if phoneCountryCode != "" && phoneNumber != "" {
		if err := verification_code.DelAllByPhone(phoneCountryCode, phoneNumber); err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
	}

	c.AbortWithStatus(http.StatusOK)

}
