package auth

import (
	"go-auth/internal/base"
	"go-auth/internal/metadata/verification_code"
	"time"
)

func verifyEmailCode(email, code, scene string) (bool, error) {
	queryEntity := &verification_code.Entity{
		Scene:  scene,
		Type:   base.SendCodeTypeEmail,
		Code:   code,
		Target: email,
		Status: verification_code.StatusUsed,
	}
	entity, err := queryEntity.FindInEntity()
	if err != nil {
		return false, err
	}
	if entity == nil {
		return false, nil
	}
	if entity.ExpiredAt < time.Now().UnixMilli() {
		return false, nil
	}
	return true, nil
}

func verifyPhoneCode(phoneCountryCode, phoneNumber, code, scene string) (bool, error) {
	queryEntity := &verification_code.Entity{
		Scene:       scene,
		Type:        base.SendCodeTypeEmail,
		Code:        code,
		Target:      phoneNumber,
		CountryCode: phoneCountryCode,
		Status:      verification_code.StatusUsed,
	}
	entity, err := queryEntity.FindInEntity()
	if err != nil {
		return false, err
	}
	if entity == nil {
		return false, nil
	}
	if entity.ExpiredAt < time.Now().UnixMilli() {
		return false, nil
	}
	return true, nil
}
