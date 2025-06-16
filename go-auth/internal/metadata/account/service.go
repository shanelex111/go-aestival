package account

import (
	"errors"
	"go-auth/internal/base"

	"github.com/shanelex111/go-common/pkg/db/mysql"
	"gorm.io/gorm"
)

func SaveInEntity(e *AccountEntity, signinType, checkType string) error {
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

	return mysql.DB.Save(entity).Error
}
