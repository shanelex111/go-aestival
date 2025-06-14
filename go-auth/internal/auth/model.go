package auth

type signinRequest struct {
	SigninType string `json:"signin_type" binding:"required"`
	CheckType  string `json:"check_type" binding:"required"`

	Email string `json:"email"`

	PhoneCountryCode string `json:"phone_country_code"`
	PhoneNumber      string `json:"phone_number"`

	VerificationCode string `json:"verification_code"`
	Password         string `json:"password"`
}
