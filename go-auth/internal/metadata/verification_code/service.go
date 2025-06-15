package verification_code

import (
	"errors"

	"github.com/shanelex111/go-common/pkg/db/mysql"
	"gorm.io/gorm"
)

func FindByEmailInEntity(email, code string) (*VerificationCodeEntity, error) {
	var entity VerificationCodeEntity
	if err := mysql.DB.Where(&VerificationCodeEntity{
		Target: email,
		Code:   code,
		Status: statusUsed,
	}).Where("deleted_at = 0").Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func FindByPhoneInEntity(phoneCountryCode, phoneNumber, code string) (*VerificationCodeEntity, error) {
	var entity VerificationCodeEntity
	if err := mysql.DB.Where(&VerificationCodeEntity{
		CountryCode: phoneCountryCode,
		Target:      phoneNumber,
		Code:        code,
		Status:      statusUsed,
	}).Where("deleted_at = 0").Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}
