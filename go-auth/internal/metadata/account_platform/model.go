package account_platform

import "go-auth/internal/base"

type AccountPlatformEntity struct {
	base.BaseModelEntity
}

func (AccountPlatformEntity) TableName() string {
	return cfg.ConfigEntity.TableName
}
