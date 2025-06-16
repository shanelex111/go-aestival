package auth

import (
	"fmt"
	"go-auth/internal/base"
	"go-auth/internal/error_code"
	"go-auth/internal/metadata/account"
	"go-auth/internal/metadata/device"
	"go-auth/internal/metadata/verification_code"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shanelex111/go-common/pkg/response"
	"github.com/shanelex111/go-common/pkg/util"
	"github.com/shanelex111/go-common/third_party/geo"
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
	var accountEntity *account.AccountEntity

	// 1.1 手机号&验证码校验
	if req.SigninType == base.SigninTypePhone && req.CheckType == base.CheckTypeVerificationCode {
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
			response.Failed(c, error_code.AuthVerificationCodeUnmatched)
			return
		}
		accountEntity = &account.AccountEntity{
			PhoneCountryCode: req.PhoneCountryCode,
			PhoneNumber:      req.PhoneNumber,
			Status:           account.StatusEnable,
		}
	}

	// 1.2 邮箱&验证码校验
	if req.SigninType == base.SigninTypeEmail && req.CheckType == base.CheckTypeVerificationCode {
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
			response.Failed(c, error_code.AuthVerificationCodeUnmatched)
			return
		}

		accountEntity = &account.AccountEntity{
			Email:  req.Email,
			Status: account.StatusEnable,
		}
	}

	// 1.3 手机号&密码校验
	if req.SigninType == base.SigninTypePhone && req.CheckType == base.CheckTypePassword {
		if req.PhoneCountryCode == "" || req.PhoneNumber == "" || req.Password == "" {
			response.Failed(c, error_code.AuthBadRequest)
			return
		}
		// 查询账户是否存在
		foundEntity, err := account.FindByPhoneInEntity(req.PhoneCountryCode, req.PhoneNumber)
		if err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
		// 校验密码
		if foundEntity != nil {
			if !foundEntity.CheckPassword(req.Password) {
				response.Failed(c, error_code.AuthInvalidPassword)
				return
			}
			accountEntity = foundEntity
		}

	}

	// 1.4 邮箱&密码校验
	if req.SigninType == base.SigninTypeEmail && req.CheckType == base.CheckTypePassword {
		if req.Email == "" || req.Password == "" {
			response.Failed(c, error_code.AuthBadRequest)
			return
		}

		// 查询账户是否存在
		foundEntity, err := account.FindByEmailInEntity(req.Email)
		if err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}

		// 校验密码
		if foundEntity != nil {
			if !foundEntity.CheckPassword(req.Password) {
				response.Failed(c, error_code.AuthInvalidPassword)
				return
			}
			accountEntity = foundEntity
		}
	}

	// 2. 记录账户信息
	if err := accountEntity.SaveInEntity(req.SigninType, req.CheckType); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 3. 查询ip信息
	geoCity, err := geo.GetCity(util.GetIP(c))
	if err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 4. 记录设备信息
	deviceEntity := device.DeviceEntity{
		AccountID:   accountEntity.ID,
		DeviceID:    req.Device.ID,
		DeviceType:  req.Device.Type,
		DeviceModel: req.Device.Model,
		AppVersion:  req.Device.AppVersion,
	}

	if geoCity != nil {
		deviceEntity.UpdatedIP = geoCity.IP
		deviceEntity.UpdatedIPContinentCode = geoCity.ContinentCode
		deviceEntity.UpdatedIPCountryCode = geoCity.CountryCode
		deviceEntity.UpdatedIPSubdivisionCode = geoCity.SubvisionCode
		deviceEntity.UpdatedIPCityName = geoCity.CityName
	}

	if err := deviceEntity.SaveInEntity(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 5. 生成token
	fmt.Println(deviceEntity)

}

func Signout(c *gin.Context) {

}

func RefreshToken(c *gin.Context) {

}
func ResetPassword(c *gin.Context) {

}
func UpdateAvatar(c *gin.Context) {

}
func UpdateNickname(c *gin.Context) {

}
func DeleteAccount(c *gin.Context) {

}
func SendCode(c *gin.Context) {
	var req sendCodeRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Failed(c, error_code.AuthBadRequest)
		return
	}

	queryEntity := &verification_code.VerificationCodeEntity{
		Type:   req.Type,
		Status: verification_code.StatusPending,
	}

	switch req.Type {
	case base.SendCodeTypeEmail:
		if req.Email == "" {
			response.Failed(c, error_code.AuthBadRequest)
			return
		}
		queryEntity.Target = req.Email
	case base.SendCodeTypePhone:
		if req.PhoneCountryCode == "" || req.PhoneNumber == "" {
			response.Failed(c, error_code.AuthBadRequest)
			return
		}

		queryEntity.Target = req.PhoneNumber
		queryEntity.CountryCode = req.PhoneCountryCode
	default:
		response.Failed(c, error_code.AuthBadRequest)
		return
	}

	// 查询出最近一次pending验证码
	foundEntity, err := queryEntity.FindLastInEntity()
	if err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}
	if foundEntity != nil {
		if time.Now().UnixMilli() <= foundEntity.ExpiredAt {
			response.Failed(c, error_code.AuthVerificationCodeFrequency)
			return
		}
	}

	// 是否超出限制
	count, err := queryEntity.CountInEntity()
	if err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}
	if count >= int64(verification_code.GetLimited()) {
		response.Failed(c, error_code.AuthVerificationCodeLimited)
		return
	}

	// 过期所有pending
	if err := queryEntity.ExpiredAllInEntity(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 创建验证码
	code := util.GetRandomNumber(verification_code.GetNumber())

	queryEntity.Code = code
	queryEntity.Status = verification_code.StatusPending
	queryEntity.ExpiredAt = time.Now().UnixMilli() + verification_code.GetPerid()

	if err := queryEntity.SaveInEntity(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

}

func VerifyCode(c *gin.Context) {

}
