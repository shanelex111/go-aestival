package auth

import "go-auth/internal/metadata/verification_code"

func verifyEmailCode(email, code string) (bool, error) {
	entity, err := verification_code.FindByEmailInEntity(email, code)
	if err != nil {
		return false, err
	}
	if entity == nil {
		return false, nil
	}
	return true, nil
}

func verifyPhoneCode(phoneCountryCode, phoneNumber, code string) (bool, error) {
	entity, err := verification_code.FindByPhoneInEntity(phoneCountryCode, phoneNumber, code)
	if err != nil {
		return false, err
	}
	if entity == nil {
		return false, nil
	}
	return true, nil
}
