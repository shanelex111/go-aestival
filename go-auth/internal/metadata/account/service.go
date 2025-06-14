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
			Status: StatusEnable,
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
				ID: id,
			},
			Status: StatusEnable,
		}).
		Last(&entity).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil

}
