package device

import (
	"errors"

	"github.com/shanelex111/go-common/pkg/db/mysql"
	"gorm.io/gorm"
)

func SaveInEntity(e *DeviceEntity) error {
	var entity *DeviceEntity
	if err := mysql.DB.Where(
		&DeviceEntity{
			AccountID:   e.AccountID,
			DeviceID:    e.DeviceID,
			DeviceType:  e.DeviceType,
			DeviceModel: e.DeviceModel,
			AppVersion:  e.AppVersion,
		}).Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			e.SigninTimes++
			e.CreatedIP = e.UpdatedIP
			e.CreatedIPCountryCode = e.UpdatedIPCountryCode
			e.CreatedIPSubdivisionCode = e.UpdatedIPSubdivisionCode
			return mysql.DB.Save(e).Error
		}
		return err
	}
	entity.SigninTimes++
	entity.UpdatedIP = e.UpdatedIP
	entity.UpdatedIPCountryCode = e.UpdatedIPCountryCode
	entity.UpdatedIPSubdivisionCode = e.UpdatedIPSubdivisionCode
	return mysql.DB.Save(entity).Error
}
