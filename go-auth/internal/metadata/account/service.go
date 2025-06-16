package account

import (
	"errors"
	"go-auth/internal/base"

	"github.com/shanelex111/go-common/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (e *AccountEntity) SaveInEntity(signinType, checkType string) error {
	var (
		entity    *AccountEntity
		condition *AccountEntity
	)

	// 手机 & 验证码
	if signinType == base.SigninTypePhone && checkType == base.CheckTypeVerificationCode {
		condition = &AccountEntity{
			PhoneCountryCode: e.PhoneCountryCode,
			PhoneNumber:      e.PhoneNumber,
			Status:           StatusEnable,
		}
	}

	// 手机 & 密码
	if signinType == base.SigninTypePhone && checkType == base.CheckTypePassword {
		condition = &AccountEntity{
			PhoneCountryCode: e.PhoneCountryCode,
			PhoneNumber:      e.PhoneNumber,
			Password:         e.Password,
			Status:           StatusEnable,
		}
	}

	// 邮箱 & 验证码
	if signinType == base.SigninTypeEmail && checkType == base.CheckTypeVerificationCode {
		condition = &AccountEntity{
			Email:  e.Email,
			Status: StatusEnable,
		}
	}

	// 邮箱 & 密码
	if signinType == base.SigninTypeEmail && checkType == base.CheckTypePassword {
		condition = &AccountEntity{
			Email:    e.Email,
			Password: e.Password,
			Status:   StatusEnable,
		}
	}

	if err := mysql.DB.Where(
		condition,
	).Where("deleted_at = 0").Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return mysql.DB.Save(e).Error
		}
		return err
	}

	e = entity
	return mysql.DB.Save(entity).Error
}

func (e *AccountEntity) SetPassword(pw string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.Password = string(hashed)
	return nil
}

func (e *AccountEntity) CheckPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(pw)) == nil
}

func FindByEmailInEntity(email string) (*AccountEntity, error) {
	var entity AccountEntity
	if err := mysql.DB.
		Where(&AccountEntity{
			Email:  email,
			Status: StatusEnable,
		}).
		Where("deleted_at = 0").
		Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func FindByPhoneInEntity(phoneCountryCode, phoneNumber string) (*AccountEntity, error) {
	var entity AccountEntity
	if err := mysql.DB.
		Where(&AccountEntity{
			PhoneCountryCode: phoneCountryCode,
			PhoneNumber:      phoneNumber,
			Status:           StatusEnable,
		}).
		Where("deleted_at = 0").
		Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil

}
