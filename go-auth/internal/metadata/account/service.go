package account

import (
	"errors"
	"go-auth/internal/base"

	"github.com/shanelex111/go-common/pkg/db/mysql"
	"gorm.io/gorm"
)

func FindByEmailInEntity(email string) (*AccountEntity, error) {
	var entity AccountEntity
	if err := mysql.DB.
		Where(&AccountEntity{
			Email:  email,
			Status: statusEnable,
			BaseModelEntity: base.BaseModelEntity{
				DeletedAt: 0,
			},
		}).
		Last(&entity).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func FindByIDInEntity(id uint) (*AccountEntity, error) {
	var entity AccountEntity
	if err := mysql.DB.
		Where(&AccountEntity{
			BaseModelEntity: base.BaseModelEntity{
				ID:        id,
				DeletedAt: 0,
			},
			Status: statusEnable,
		}).
		Last(&entity).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil

}
