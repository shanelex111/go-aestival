package auth

import (
	"go-auth/internal/base"
	"go-auth/internal/error_code"
	"go-auth/internal/metadata/account"
	"go-auth/internal/metadata/device"
	"go-auth/internal/metadata/verification_code"
	"go-auth/internal/token"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shanelex111/go-common/pkg/request"
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
	var accountEntity = &account.Entity{
		PhoneCountryCode: req.PhoneCountryCode,
		PhoneNumber:      req.PhoneNumber,
		Email:            req.Email,
		Password:         req.Password,
		Status:           account.StatusEnable,
	}

	// 1.1 手机号&验证码校验
	if req.SigninType == base.SigninTypePhone && req.CheckType == base.CheckTypeVerificationCode {
		if req.PhoneCountryCode == "" || req.PhoneNumber == "" || req.VerificationCode == "" {
			response.Failed(c, error_code.AuthBadRequest)
			return
		}
		valid, err := verifyPhoneCode(req.PhoneCountryCode, req.PhoneNumber, req.VerificationCode, verification_code.SceneSignin)
		if err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
		if !valid {
			response.Failed(c, error_code.AuthVerificationCodeUnmatched)
			return
		}
		accountEntity = &account.Entity{
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
		valid, err := verifyEmailCode(req.Email, req.VerificationCode, verification_code.SceneSignin)
		if err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
		if !valid {
			response.Failed(c, error_code.AuthVerificationCodeUnmatched)
			return
		}

		accountEntity = &account.Entity{
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
	deviceEntity := device.Entity{
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
	newToken := &token.CacheToken{
		Account: &token.CacheTokenAccount{
			ID:               accountEntity.ID,
			Email:            accountEntity.Email,
			PhoneCountryCode: accountEntity.PhoneCountryCode,
			PhoneNumber:      accountEntity.PhoneNumber,
		},

		Device: &token.CacheTokenDevice{
			DeviceID:    deviceEntity.DeviceID,
			DeviceType:  deviceEntity.DeviceType,
			DeviceModel: deviceEntity.DeviceModel,
			AppVersion:  deviceEntity.AppVersion,
			CreatedAt:   time.Now().UnixMilli(),
		},
		Geo: geoCity,
	}
	if err = newToken.Create(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	//返回token
	c.AbortWithStatusJSON(http.StatusOK, signinResponse{
		Access: newToken.Access,
	})

}

func Signout(c *gin.Context) {
	requestTokenInfo := request.GetTokenInfo(c)
	if requestTokenInfo == nil {
		response.Failed(c, error_code.AuthUnauthorized)
		return
	}

	cacheToken := &token.CacheToken{
		Account: &token.CacheTokenAccount{
			ID: requestTokenInfo.Account.ID,
		},
		Access: &token.CacheTokenAccess{
			Token:   requestTokenInfo.Access.Token,
			Refresh: requestTokenInfo.Access.Refresh,
		},
	}

	if err := cacheToken.Delete(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}
	c.AbortWithStatus(http.StatusOK)

}

func RefreshToken(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Failed(c, error_code.AuthBadRequest)
		return
	}
	existRefresh, err := token.GetRefresh(req.Refresh)
	if err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}
	if existRefresh == nil {
		response.Failed(c, error_code.AuthUnauthorized)
		return
	}
	// 3. 查询ip信息
	geoCity, err := geo.GetCity(util.GetIP(c))
	if err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 4. 记录设备信息
	deviceEntity := device.Entity{
		AccountID:   existRefresh.Account.ID,
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
	newToken := &token.CacheToken{
		Account: &token.CacheTokenAccount{
			ID: existRefresh.Account.ID,
		},
		Device: &token.CacheTokenDevice{
			DeviceID:    deviceEntity.DeviceID,
			DeviceType:  deviceEntity.DeviceType,
			DeviceModel: deviceEntity.DeviceModel,
			AppVersion:  deviceEntity.AppVersion,
			CreatedAt:   time.Now().UnixMilli(),
		},
		Geo: geoCity,
	}
	if err = newToken.Create(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 6. 删除旧token
	if err := existRefresh.Delete(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	//返回token
	c.AbortWithStatusJSON(http.StatusOK, refreshTokenResponse{
		Access: newToken.Access,
	})
}
func ResetPassword(c *gin.Context) {
	var req resetPasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Failed(c, error_code.AuthBadRequest)
		return
	}
	tokenInfo, exists := c.Get(request.TokenInfoKey)
	if !exists {
		response.Failed(c, error_code.AuthUnauthorized)
		return
	}

	// 1. 手机号&验证码登录 | 邮箱&验证码登录
	var (
		foundEntity      *account.Entity
		requestTokenInfo = tokenInfo.(*request.TokenInfo)
	)

	// 1.1 手机号&验证码校验
	if req.SigninType == base.SigninTypePhone && req.CheckType == base.CheckTypeVerificationCode {
		if req.PhoneCountryCode == "" || req.PhoneNumber == "" || req.VerificationCode == "" {
			response.Failed(c, error_code.AuthBadRequest)
			return
		}
		valid, err := verifyPhoneCode(req.PhoneCountryCode, req.PhoneNumber, req.VerificationCode, verification_code.SceneResetPassword)
		if err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
		if !valid {
			response.Failed(c, error_code.AuthVerificationCodeUnmatched)
			return
		}

		foundEntity, err = account.FindByPhoneInEntity(req.PhoneCountryCode, req.PhoneNumber)
		if err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
		if foundEntity == nil {
			response.Failed(c, error_code.AuthUnauthorized)
			return
		}
		if foundEntity.ID != requestTokenInfo.Account.ID {
			response.Failed(c, error_code.AuthUnauthorized)
			return
		}

	}

	// 1.2 邮箱&验证码校验
	if req.SigninType == base.SigninTypeEmail && req.CheckType == base.CheckTypeVerificationCode {
		if req.Email == "" || req.VerificationCode == "" {
			response.Failed(c, error_code.AuthBadRequest)
			return
		}
		valid, err := verifyEmailCode(req.Email, req.VerificationCode, verification_code.SceneResetPassword)
		if err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
		if !valid {
			response.Failed(c, error_code.AuthVerificationCodeUnmatched)
			return
		}

		foundEntity, err = account.FindByEmailInEntity(req.Email)
		if err != nil {
			response.Failed(c, error_code.AuthInternalServerError)
			return
		}
		if foundEntity == nil {
			response.Failed(c, error_code.AuthUnauthorized)
			return
		}
		if foundEntity.ID != requestTokenInfo.Account.ID {
			response.Failed(c, error_code.AuthUnauthorized)
			return
		}
	}

	if foundEntity == nil {
		response.Failed(c, error_code.AuthUnauthorized)
		return
	}
	if foundEntity.CheckPassword(req.NewPassword) {
		response.Failed(c, error_code.AuthSetTheSamePassword)
		return
	}

	// 2. 修改密码
	if err := foundEntity.SetPassword(req.NewPassword); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}
	if err := foundEntity.Update(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 3. 删除所有token
	if err := token.DelAllByAccountID(foundEntity.ID); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}
}

func SendCode(c *gin.Context) {
	var req sendCodeRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Failed(c, error_code.AuthBadRequest)
		return
	}

	queryEntity := &verification_code.Entity{
		Scene:  req.Scene,
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
	result, err := queryEntity.FindInCache()
	if err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}
	if result != "" {
		response.Failed(c, error_code.AuthVerificationCodeFrequency)
		return
	}

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
	count, err := queryEntity.CountTodayInEntity()
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

	// TODO:调用第三方发送验证码

	// 数据库存储
	queryEntity.Code = code
	queryEntity.Status = verification_code.StatusPending
	queryEntity.ExpiredAt = time.Now().UnixMilli() + verification_code.GetPeriod()

	if err := queryEntity.SaveInEntity(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 存入redis，方便校验
	if err := queryEntity.SaveInCache(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

func VerifyCode(c *gin.Context) {
	var req verifyCodeRequest
	if err := c.ShouldBind(&req); err != nil {
		response.Failed(c, error_code.AuthBadRequest)
		return
	}

	queryEntity := &verification_code.Entity{
		Scene:  req.Scene,
		Type:   req.Type,
		Status: verification_code.StatusPending,
		Code:   req.Code,
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

	// 查询验证码
	result, err := queryEntity.FindInCache()
	if err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}
	if result != "" && result != req.Code {
		response.Failed(c, error_code.AuthVerificationCodeUnmatched)
		return
	}

	foundEntity, err := queryEntity.FindInEntity()
	if err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}
	if foundEntity == nil {
		response.Failed(c, error_code.AuthVerificationCodeUnmatched)
		return
	}

	if time.Now().UnixMilli() > foundEntity.ExpiredAt {
		response.Failed(c, error_code.AuthVerificationCodeExpired)
		return
	}

	// 更新状态
	foundEntity.Status = verification_code.StatusUsed
	if err := foundEntity.SaveInEntity(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	// 删除缓存
	if err := queryEntity.DeleteInCache(); err != nil {
		response.Failed(c, error_code.AuthInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}
