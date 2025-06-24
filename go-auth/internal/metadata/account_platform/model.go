package account_platform

import "go-auth/internal/base"

type Entity struct {
	base.BaseModelEntity
}

func (*Entity) TableName() string {
	return cfg.ConfigEntity.TableName
}
