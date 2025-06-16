package auth

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
type deviceRequest struct {
	ID         string `json:"id" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Model      string `json:"model" binding:"required"`
	AppVersion int    `json:"app_version" binding:"required"`
}

type signinResponse struct {
}

type sendCodeRequest struct {
}

type sendCodeResponse struct {
}
