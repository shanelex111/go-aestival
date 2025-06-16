package verification_code

import (
	"errors"
	"go-auth/internal/base"

	"github.com/shanelex111/go-common/pkg/db/mysql"
	"gorm.io/gorm"
)

func FindByEmailInEntity(email, code string) (*VerificationCodeEntity, error) {
	var entity VerificationCodeEntity
	if err := mysql.DB.Where(&VerificationCodeEntity{
		Target: email,
		Code:   code,
		Status: StatusUsed,
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
		Status:      StatusUsed,
	}).Where("deleted_at = 0").Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (e *VerificationCodeEntity) FindLastInEntity() (*VerificationCodeEntity, error) {
	var (
		entity    VerificationCodeEntity
		condition = &VerificationCodeEntity{
			Type:   e.Type,
			Target: e.Target,
			Status: e.Status,
		}
	)
	if e.Type == base.SendCodeTypePhone {
		condition.CountryCode = e.CountryCode
	}

	if err := mysql.DB.Where(condition).Where("deleted_at = 0").Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}
