package verification_code

import "go-auth/internal/base"

type VerificationCodeEntity struct {
	base.BaseModelEntity
}

func (VerificationCodeEntity) TableName() string {
	return cfg.ConfigEntity.TableName
}
