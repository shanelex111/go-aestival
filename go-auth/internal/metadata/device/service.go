package device

import (
	"errors"
	"go-auth/internal/base"
	"time"

	"github.com/shanelex111/go-common/pkg/db/mysql"
	"gorm.io/gorm"
)

func (e *Entity) SaveInEntity() error {
	var entity *Entity
	if err := mysql.DB.Where(
		&Entity{
			AccountID:   e.AccountID,
			DeviceID:    e.DeviceID,
			DeviceType:  e.DeviceType,
			DeviceModel: e.DeviceModel,
			AppVersion:  e.AppVersion,
		}).Where("deleted_at = 0").Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.SigninTimes++
			e.CreatedIP = e.UpdatedIP
			e.CreatedIPContinentCode = e.UpdatedIPContinentCode
			e.CreatedIPCountryCode = e.UpdatedIPCountryCode
			e.CreatedIPSubdivisionCode = e.UpdatedIPSubdivisionCode
			e.CreatedIPCityName = e.UpdatedIPCityName
			return mysql.DB.Save(e).Error
		}
		return err
	}
	entity.SigninTimes++
	entity.UpdatedIP = e.UpdatedIP
	entity.UpdatedIPContinentCode = e.UpdatedIPContinentCode
	entity.UpdatedIPCountryCode = e.UpdatedIPCountryCode
	entity.UpdatedIPSubdivisionCode = e.UpdatedIPSubdivisionCode
	entity.UpdatedIPCityName = e.UpdatedIPCityName
	*e = *entity
	return mysql.DB.Save(entity).Error
}

func DelAllByAccountID(accountID uint) error {
	if err := mysql.DB.Model(&Entity{}).
		Where(&Entity{
			AccountID: accountID,
		}).
		Where("deleted_at = 0").
		Updates(&Entity{
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
