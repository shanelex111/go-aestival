package account

import "go-auth/internal/base"

type AccountEntity struct {
	base.BaseModelEntity
}

func (AccountEntity) TableName() string {
	return cfg.ConfigEntity.TableName
}
