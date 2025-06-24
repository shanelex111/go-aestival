package auth

import "go-auth/internal/token"

type signinRequest struct {
	SigninType string `json:"signin_type" binding:"required"`
	CheckType  string `json:"check_type" binding:"required"`

	Email string `json:"email"`

	PhoneCountryCode string `json:"phone_country_code"`
	PhoneNumber      string `json:"phone_number"`

	VerificationCode string `json:"verification_code"`
	Password         string `json:"password"`

	Device *deviceRequest `json:"device" binding:"required"`
}
type signinResponse struct {
	Access *token.CacheTokenAccess `json:"access"`
}
type refreshTokenRequest struct {
	Device  *deviceRequest `json:"device" binding:"required"`
	Refresh string         `json:"refresh" binding:"required"`
}
type refreshTokenResponse struct {
	Access *token.CacheTokenAccess `json:"access"`
}
type deviceRequest struct {
	ID         string `json:"id" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Model      string `json:"model" binding:"required"`
	AppVersion int    `json:"app_version" binding:"required"`
}

type sendCodeRequest struct {
	Scene            string `json:"scene" binding:"required,oneof=signin reset_password"`
	Type             string `json:"type" binding:"required,oneof=email phone"`
	Email            string `json:"email"`
	PhoneCountryCode string `json:"phone_country_code"`
	PhoneNumber      string `json:"phone_number"`
}

type verifyCodeRequest struct {
	sendCodeRequest
	Code string `json:"code" binding:"required"`
}
