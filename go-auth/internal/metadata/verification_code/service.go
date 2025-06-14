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
		Status: statusUsed,
		BaseModelEntity: base.BaseModelEntity{
			DeletedAt: 0,
		},
	}).Last(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}
