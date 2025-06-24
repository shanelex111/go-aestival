package account

import (
	"errors"
	"go-auth/internal/base"
	"time"

	"github.com/shanelex111/go-common/pkg/db/mysql"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func (e *Entity) SaveInEntity(signinType, checkType string) error {
	var (
		entity    *Entity
		condition *Entity
	)

	// 手机 & 验证码
	if signinType == base.SigninTypePhone && checkType == base.CheckTypeVerificationCode {
		condition = &Entity{
			PhoneCountryCode: e.PhoneCountryCode,
			PhoneNumber:      e.PhoneNumber,
			Status:           StatusEnable,
		}
	}

	// 手机 & 密码
	if signinType == base.SigninTypePhone && checkType == base.CheckTypePassword {
		condition = &Entity{
			PhoneCountryCode: e.PhoneCountryCode,
			PhoneNumber:      e.PhoneNumber,
			Password:         e.Password,
			Status:           StatusEnable,
		}
	}

	// 邮箱 & 验证码
	if signinType == base.SigninTypeEmail && checkType == base.CheckTypeVerificationCode {
		condition = &Entity{
			Email:  e.Email,
			Status: StatusEnable,
		}
	}

	// 邮箱 & 密码
	if signinType == base.SigninTypeEmail && checkType == base.CheckTypePassword {
		condition = &Entity{
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

	*e = *entity
	return mysql.DB.Save(entity).Error
}

func (e *Entity) Update() error {
	return mysql.DB.Save(e).Error
}
func (e *Entity) SetPassword(pw string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	e.Password = string(hashed)
	return nil
}

func (e *Entity) CheckPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(e.Password), []byte(pw)) == nil
}

func FindByEmailInEntity(email string) (*Entity, error) {
	var entity Entity
	if err := mysql.DB.
		Where(&Entity{
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

func FindByPhoneInEntity(phoneCountryCode, phoneNumber string) (*Entity, error) {
	var entity Entity
	if err := mysql.DB.
		Where(&Entity{
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
func FindByAccountID(accountID uint) (*Entity, error) {
	var entity Entity
	if err := mysql.DB.
		Where(&Entity{
			BaseModelEntity: base.BaseModelEntity{
				ID: accountID,
			},
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

func DelAllByAccountID(accountID uint) error {
	if err := mysql.DB.Model(&Entity{}).
		Where(&Entity{
			BaseModelEntity: base.BaseModelEntity{
				ID: accountID,
			},
		}).
		Where("deleted_at = 0 and status != ?", StatusDeleted).
		Updates(&Entity{
			Status: StatusDeleted,
			BaseModelEntity: base.BaseModelEntity{
				DeletedAt: time.Now().UnixMilli(),
			},
		}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}
